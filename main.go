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
	if timeout, ok := opts["timeout"]; ok {
		r.Client.SetTimeout(timeout.(time.Duration))
	}

	ts, ok := opts["timestamp"]

	if !ok {
		ts = time.Now()
	}

	resp, err := r.Client.Query(query, ts.(time.Time))

	if err != nil {
		log.Fatalln(err)
	}

	resp_vector := resp.(model.Vector)
	result := make(TSDBQueryResult, len(resp_vector))

	for _, sample := range resp_vector {
		tv := TimeValue{Time: int64(sample.Timestamp), Value: float64(sample.Value)}

		labels := make(map[string]string)

		for key, val := range sample.Metric {
			labels[string(key)] = string(val)
		}

		result = append(result, TimeSeries{
			Labels:          labels,
			TimeValueSeries: []TimeValue{tv},
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

	returnReturn := make(TSDBQueryResult, 10)

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
		fmt.Printf("value: %v, measurement: %v, field: %v, time: %v\n", result.Record().Values(), result.Record().Measurement(), result.Record().Field(), result.Record().Time())
	}
	// check for an error
	if result.Err() != nil {
		fmt.Printf("query parsing error: %s\n", result.Err().Error())
	}

	return returnReturn
}

func main() {
	var tsdbStore1 TSDBStore = NewMimirDBStore("sv5-edn-mimir-stg.lab.equinix.com", "eot-telemetry")
	pretty.Print(tsdbStore1.Query(context.Background(), "bits", map[string]any{}))

	// var tsdbStore2 TSDBStore = NewInfluxDBStore("http://devsv3ednmgmt09.lab.equinix.com:30320", "mytoken")

	// query := `from(bucket: "testing_script")
	// |> range( start: -3d, stop: -2d)
	// |> filter(fn: (r) => r["_measurement"] == "in_traffic" or r["_measurement"] == "out_traffic")
	// |> filter(fn: (r) => r["_field"] == "bits")
	// |> filter(fn: (r) => r["index_num"] == "bb1-ngn.gv51.1001")
	// |> yield(name: "last")`

	// pretty.Print(tsdbStore2.Query(context.Background(), query, map[string]any{"org": "primary"}))
}
