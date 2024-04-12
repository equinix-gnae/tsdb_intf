package influxdb

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/equinix-gnae/tsdb_intf/pkg/ts"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type InfluxDBClient struct {
	Client influxdb2.Client
	Bucket string
	Org    string
}

func NewInfluxDBClient(url string, token string, bucket string, org string, options *influxdb2.Options) InfluxDBClient {
	influxClient := influxdb2.NewClientWithOptions(url, token, options)
	running, err := influxClient.Ping(context.Background())

	if err != nil {
		log.Fatalf("error running ping test for influx: %v", err)
	}

	if !running {
		log.Fatal("influx is not running")
	}

	return InfluxDBClient{Client: influxClient, Bucket: bucket, Org: org}
}

func (r InfluxDBClient) Query(ctx context.Context, query ts.TSQuery) (ts.TSQueryResult, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, query.Timeout)
	defer cancel()

	strQuery, err := r.GenerateQueryString(query)

	if err != nil {
		return nil, err
	}

	queryAPI := r.Client.QueryAPI(r.Org)
	resp, err := queryAPI.Query(ctxWithTimeout, strQuery)

	if err != nil {
		return nil, err
	}

	// caution: result.TableChanged() is not working for some reason that why we are using
	// preTableId/currentTableId to implement the logic to figure out if table has changed
	preTableId := -1
	returnResult := make(ts.TSQueryResult, 0, 10)

	for resp.Next() {
		// Notice when group key has changed
		if resp.TableChanged() {
			fmt.Printf("table: %s\n", resp.TableMetadata().String())

		}
		record := resp.Record()
		currentTableId := record.Table()

		// new time series
		if preTableId != currentTableId {
			labels := make(map[string]string)
			for key, val := range record.Values() {
				labels[key] = fmt.Sprintf("%v", val)
			}
			returnResult = append(returnResult, ts.TimeSeries{
				Labels:          labels,
				TimeValueSeries: make([]ts.TimeValue, 0, 10),
			})
			preTableId = currentTableId
		}

		// same TS: update the TSVals of the last element in
		returnResult[len(returnResult)-1].TimeValueSeries = append(
			returnResult[len(returnResult)-1].TimeValueSeries,
			ts.TimeValue{Time: record.Time().UnixMilli(), Value: record.Value().(float64)},
		)

	}
	if resp.Err() != nil {
		return nil, fmt.Errorf("for query %q, query parsing error: %s", strQuery, resp.Err().Error())

	}

	return returnResult, nil
}

/*
	 Query String Example:
		query := `from(bucket: "testing_script")
		|> range( start: 2024-02-05, stop: 2024-02-09)
		|> filter(fn: (r) => r["_field"] == "bits")
		|> filter(fn: (r) => r["index_num"] == "bb1-ngn.gv51.1001")
		|> group (columns: ["index_num"])
		|> aggregateWindow(every: 5m, fn: last, createEmpty: false)
		`
*/
func (r InfluxDBClient) GenerateQueryString(query ts.TSQuery) (string, error) {
	var queryBuilder strings.Builder

	queryBuilder.WriteString(fmt.Sprintf("from(bucket: %q)\n", r.Bucket))
	queryBuilder.WriteString(fmt.Sprintf("|> range( start: %s, stop: %s)\n", query.StartTime.Format(time.RFC3339), query.EndTime.Format(time.RFC3339)))
	queryBuilder.WriteString(fmt.Sprintf("|> filter(fn: (r) => r[\"_field\"] == %q)\n", query.Table))

	if query.Filters != nil && len(query.Filters) > 0 {
		queryBuilder.WriteString("// filters\n")
	}

	for k, v := range query.Filters {
		queryBuilder.WriteString(fmt.Sprintf("|> filter(fn: (r) => r[%q] == %q)\n", k, v))
	}

	queryStr := queryBuilder.String()
	if err := applyFunctions(&queryStr, query.Functions); err != nil {
		return "", err
	}

	if err := applyOperations(&queryStr, query.Operations); err != nil {
		return "", err
	}

	queryStr += "// others\n"
	// TODO: add "_measurement" and "_field" by default for groupBy key? or its upto query maker?
	if len(query.GroupBy) > 0 {
		groupKey, _ := json.Marshal(query.GroupBy)
		queryStr += fmt.Sprintf("|> group (columns: %s)\n", groupKey)
	}

	if query.Step != 0 {
		queryStr += fmt.Sprintf("|> aggregateWindow(every: %s, fn: last, createEmpty: false)", query.Step)
	}

	return queryStr, nil
}

func applyFunctions(queryStr *string, functions []ts.QueryFunction) error {
	if len(functions) > 0 {
		*queryStr += "// functions\n"
	}

	for _, queryFunction := range functions {
		switch t := queryFunction.(type) {
		case ts.Rate:
			rate(t).Apply(queryStr)
		case ts.Sum:
			sum(t).Apply(queryStr)
		default:
			return fmt.Errorf("unsupported Function: %v", queryFunction)
		}
	}
	return nil
}

// Query String Example => query := `rate(bits{index_num="bb1-ngn.gv51.1001"}[5m])`
func applyOperations(queryStr *string, operations []ts.QueryOperation) error {
	if len(operations) > 0 {
		*queryStr += "// operations\n"
	}
	for _, operation := range operations {
		switch t := operation.(type) {
		case ts.Add:
			add(t).Apply(queryStr)
		case ts.Subtract:
			subtract(t).Apply(queryStr)
		case ts.Multiply:
			multiply(t).Apply(queryStr)
		case ts.Divide:
			divide(t).Apply(queryStr)
		default:
			return fmt.Errorf("unsupported operation: %v", operation)
		}
	}
	return nil
}
