package main

import (
	"context"
	"fmt"
	"log"
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
	Query(ctx context.Context, query string, opts map[string]any) TSDBQueryResult
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

func (r MimirDBStore) Query(ctx context.Context, query string, opts map[string]any) TSDBQueryResult {
	// TODO: check errors
	startTime := opts["start_time"].(time.Time)
	endTime := opts["end_time"].(time.Time)
	step := opts["step"].(time.Duration)

	if timeout, ok := opts["timeout"]; ok {
		r.Client.SetTimeout(timeout.(time.Duration))
	}

	resp, err := r.Client.QueryRange(query, startTime, endTime, step)

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

func (r InfluxDBStore) Query(ctx context.Context, query string, opts map[string]any) TSDBQueryResult {
	queryAPI := r.Client.QueryAPI(opts["org"].(string))

	result, err := queryAPI.Query(ctx, query)

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
		//fmt.Printf("value: %v, measurement: %v, field: %v, time: %v\n", result.Record().Values(), result.Record().Measurement(), result.Record().Field(), result.Record().Time())
	}
	// check for an error
	if result.Err() != nil {
		fmt.Printf("query parsing error: %s\n", result.Err().Error())
	}

	return returnReturn
}

func main() {
	// var tsdbStore1 TSDBStore = NewMimirDBStore("sv5-edn-mimir-stg.lab.equinix.com", "eot-telemetry")
	// pretty.Print(tsdbStore1.Query(context.Background(), `bits{index_num="bb1-ngn.gv51.1001"}`, map[string]any{
	// 	"start_time": time.Now().Add(3 * 24 * -time.Hour),
	// 	"end_time":   time.Now().Add(2 * 24 * -time.Hour),
	// 	"step":       time.Minute * 30,
	// }))

	var tsdbStore2 TSDBStore = NewInfluxDBStore("http://devsv3ednmgmt09.lab.equinix.com:30320", "mytoken")

	query := `from(bucket: "testing_script")
	|> range( start: -3d, stop: -2d)
	|> filter(fn: (r) => r["_field"] == "bits")
	|> filter(fn: (r) => r["index_num"] == "bb1-ngn.gv51.1001")`

	pretty.Print(tsdbStore2.Query(context.Background(), query, map[string]any{"org": "primary"}))
}
