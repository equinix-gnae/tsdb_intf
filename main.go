package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/grafana/mimir/integration/e2emimir"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/kylelemons/godebug/pretty"
	"github.com/prometheus/common/model"
)

// *** return type Data Structures ***
type TimeValue struct {
	Time  int64   `json:"time"`
	Value float64 `json:"value"`
}

type TimeSeries struct {
	Name            string            `json:"name"`
	Labels          map[string]string `json:"labels"`
	TimeValueSeries []TimeValue       `json:"timeValueSeries"`
}

type TSDBQueryResult []TimeSeries

// ***********************************

// timestamp time.Time, timeout time.Duration,
type TSDBStore interface {
	Query(ctx context.Context, query TSQuery, opts map[string]any) TSDBQueryResult
}

type TSQuery struct {
	Table     string            `json:"table"`
	StartTime time.Time         `json:"startTime"`
	EndTime   time.Time         `json:"endTime"`
	Step      time.Duration     `json:"step"`
	Filters   map[string]string `json:"filters"`
}

// *** Mimir ************************
type MimirDBStore struct {
	Client *e2emimir.Client
}

func NewMimirDBStore(url string, id string) MimirDBStore {
	mimirE2eClient, err := e2emimir.NewClient("", url, "", "", id)

	if err != nil {
		log.Fatalln(err)
	}

	return MimirDBStore{Client: mimirE2eClient}
}

/*
	 Query String Example:
		query := `bits{index_num="bb1-ngn.gv51.1001"}`
*/
func (r MimirDBStore) Query(ctx context.Context, query TSQuery, opts map[string]any) TSDBQueryResult {
	if timeout, ok := opts["timeout"]; ok {
		r.Client.SetTimeout(timeout.(time.Duration))
	}

	// build query string
	var queryBuilder strings.Builder
	queryBuilder.WriteString(query.Table)
	queryBuilder.WriteString("{")
	for k, v := range query.Filters {
		queryBuilder.WriteString(k)
		queryBuilder.WriteString("=")
		queryBuilder.WriteString("\"")
		queryBuilder.WriteString(v)
		queryBuilder.WriteString("\"")
	}
	queryBuilder.WriteString("}")

	resp, err := r.Client.QueryRange(queryBuilder.String(), query.StartTime, query.EndTime, query.Step)

	if err != nil {
		log.Fatalln(err)
	}

	matrix, ok := resp.(model.Matrix)

	if !ok {
		log.Fatalf("unable to convert resp to vector")
	}

	if len(matrix) == 0 {
		log.Fatalf("empty response is returned for query: %q", query)
	}

	//pretty.Print(matrix)

	result := make(TSDBQueryResult, 0, len(matrix))

	for _, sampleStream := range matrix {
		labels := make(map[string]string)
		for key, val := range sampleStream.Metric {
			labels[string(key)] = string(val)
		}

		timeValueSeries := make([]TimeValue, 0, len(sampleStream.Values))
		for _, sample := range sampleStream.Values {
			timeValueSeries = append(timeValueSeries, TimeValue{Time: int64(sample.Timestamp), Value: float64(sample.Value)})
		}

		result = append(result, TimeSeries{
			Labels:          labels,
			TimeValueSeries: timeValueSeries,
		})
	}
	return result
}

// *** InfluxDB *******************

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

func main() {
	location, err := time.LoadLocation("UTC")

	if err != nil {
		log.Fatal(err)
	}

	query := TSQuery{
		Table:     "bits",
		StartTime: time.Date(2024, time.February, 5, 0, 0, 0, 0, location),
		EndTime:   time.Date(2024, time.February, 9, 0, 0, 0, 0, location),
		Filters:   map[string]string{"index_num": "bb1-ngn.gv51.1001"},
		Step:      time.Minute * 30,
	}

	// var tsdbStore1 TSDBStore = NewMimirDBStore("sv5-edn-mimir-stg.lab.equinix.com", "eot-telemetry")
	// pretty.Print(tsdbStore1.Query(context.Background(), query, map[string]any{}))

	var tsdbStore2 TSDBStore = NewInfluxDBStore("http://devsv3ednmgmt09.lab.equinix.com:30320", "mytoken")
	pretty.Print(tsdbStore2.Query(context.Background(), query, map[string]any{"org": "primary", "bucket": "testing_script"}))
}
