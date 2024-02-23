package tsdb

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type InfluxDBStore struct {
	Client influxdb2.Client
	Bucket string
	Org    string
}

func NewInfluxDBStore(url string, token string, bucket string, org string) InfluxDBStore {
	options := influxdb2.DefaultOptions()
	options.SetPrecision(time.Nanosecond)
	options.SetFlushInterval(5_000)
	options.SetLogLevel(3)

	influxClient := influxdb2.NewClientWithOptions(url, token, options)

	running, err := influxClient.Ping(context.Background())

	if err != nil {
		log.Fatalf("error running ping test for influx: %v", err)
	}

	if !running {
		log.Fatal("influx is not running")
	}

	return InfluxDBStore{Client: influxClient, Bucket: bucket, Org: org}
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
func (r InfluxDBStore) Query(ctx context.Context, query TSQuery, opts map[string]any) TSDBQueryResult {
	queryAPI := r.Client.QueryAPI(r.Org)

	// build query string
	var queryBuilder strings.Builder

	queryBuilder.WriteString(fmt.Sprintf("from(bucket: %q)\n", r.Bucket))
	queryBuilder.WriteString(fmt.Sprintf("|> range( start: %s, stop: %s)\n", query.StartTime.Format(time.RFC3339), query.EndTime.Format(time.RFC3339)))
	queryBuilder.WriteString(fmt.Sprintf("|> filter(fn: (r) => r[\"_field\"] == %q)\n", query.Table))

	for k, v := range query.Filters {
		queryBuilder.WriteString(fmt.Sprintf("|> filter(fn: (r) => r[%q] == %q)\n", k, v))
	}

	// TODO: add "_measurement" and "_field" by default for groupBy key? or its upto query maker?
	if len(query.GroupBy) > 0 {
		groupKey, _ := json.Marshal(query.GroupBy)
		queryBuilder.WriteString(fmt.Sprintf("|> group (columns: %s)\n", groupKey))
	}

	if query.Step != 0 {
		queryBuilder.WriteString(fmt.Sprintf("|> aggregateWindow(every: %s, fn: last, createEmpty: false)", query.Step))
	}

	result, err := queryAPI.Query(ctx, queryBuilder.String())

	if err != nil {
		log.Fatalln(err)
	}

	// caution: result.TableChanged() is not working for some reason that why we are using
	// preTableId/currentTableId to implement the logic to figure out if table has changed
	preTableId := -1
	returnResult := make(TSDBQueryResult, 0, 10)

	for result.Next() {
		// Notice when group key has changed
		if result.TableChanged() {
			fmt.Printf("table: %s\n", result.TableMetadata().String())

		}
		record := result.Record()
		currentTableId := record.Table()

		// new time series
		if preTableId != currentTableId {
			labels := make(map[string]string)
			for key, val_intf := range record.Values() {
				switch val := val_intf.(type) {
				case string:
					labels[key] = val
				default:
					labels[key] = fmt.Sprintf("%v", val)
				}
			}
			returnResult = append(returnResult, TimeSeries{
				Labels:          labels,
				TimeValueSeries: make([]TimeValue, 0, 10),
			})
			preTableId = currentTableId
		}

		// same TS: update the TSVals of the last element in
		returnResult[len(returnResult)-1].TimeValueSeries = append(
			returnResult[len(returnResult)-1].TimeValueSeries,
			TimeValue{Time: record.Time().Unix(), Value: record.Value().(float64)},
		)

	}
	if result.Err() != nil {
		fmt.Printf("query parsing error: %s\n", result.Err().Error())
	}

	return returnResult
}
