package tsdb

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type InfluxDBStore struct {
	Client influxdb2.Client
}

func NewInfluxDBStore(url string, token string) InfluxDBStore {
	options := influxdb2.DefaultOptions()
	options.SetPrecision(time.Second)
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

	return InfluxDBStore{Client: influxClient}
}

/*
	 Query String Example:
		query := `from(bucket: "testing_script")
		|> range( start: 2024-02-05, stop: 2024-02-09)
		|> filter(fn: (r) => r["_field"] == "bits")
		|> filter(fn: (r) => r["index_num"] == "bb1-ngn.gv51.1001")
		|> aggregateWindow(every: 5m, fn: last, createEmpty: false)
		`
*/
func (r InfluxDBStore) Query(ctx context.Context, query TSQuery, opts map[string]any) TSDBQueryResult {
	// TODO: both bucket and org should moved to config file
	queryAPI := r.Client.QueryAPI(opts["org"].(string))
	bucket := opts["bucket"].(string)

	// build query string
	var queryBuilder strings.Builder
	queryBuilder.WriteString(fmt.Sprintf("from(bucket: %q)\n", bucket))
	queryBuilder.WriteString(fmt.Sprintf("|> range( start: %s, stop: %s)\n", query.StartTime.Format(time.RFC3339), query.EndTime.Format(time.RFC3339)))
	queryBuilder.WriteString(fmt.Sprintf("|> filter(fn: (r) => r[\"_field\"] == %q)\n", query.Table))
	for k, v := range query.Filters {
		queryBuilder.WriteString(fmt.Sprintf("|> filter(fn: (r) => r[%q] == %q)\n", k, v))
	}
	queryBuilder.WriteString(fmt.Sprintf("|> aggregateWindow(every: %s, fn: last, createEmpty: false)", query.Step))

	result, err := queryAPI.Query(ctx, queryBuilder.String())

	if err != nil {
		log.Fatalln(err)
	}

	returnReturn := make(TSDBQueryResult, 0, 10)

	for result.Next() {
		// Notice when group key has changed
		if result.TableChanged() {
			fmt.Printf("table: %s\n", result.TableMetadata().String())
		}
		record := result.Record()

		tv := TimeValue{Time: record.Time().Unix(), Value: record.Value().(float64)}

		labels := make(map[string]string)

		for key, val_intf := range record.Values() {
			switch val := val_intf.(type) {
			case string:
				labels[key] = val
			default:
				labels[key] = fmt.Sprintf("%v", val)
			}
		}

		returnReturn = append(returnReturn, TimeSeries{
			Name:            record.Field() + "_" + record.Measurement(),
			Labels:          labels,
			TimeValueSeries: []TimeValue{tv},
		})
		//gfmt.Printf("value: %v, measurement: %v, field: %v, time: %v\n", result.Record().Values(), result.Record().Measurement(), result.Record().Field(), result.Record().Time())
	}
	// check for an error
	if result.Err() != nil {
		fmt.Printf("query parsing error: %s\n", result.Err().Error())
	}

	return returnReturn
}
