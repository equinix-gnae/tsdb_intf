package ts_test

import (
	"context"
	"testing"
	"time"

	"github.com/equinix-gnae/tsdb_intf/pkg/influxdb"
	"github.com/equinix-gnae/tsdb_intf/pkg/prometheus"
	"github.com/equinix-gnae/tsdb_intf/pkg/ts"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/kr/pretty"
)

var BaseQuery = ts.TSQuery{
	Table: "bits",
	//StartTime: time.Now().Add(-1 * time.Hour),
	//EndTime:   time.Now().UTC(),
	StartTime: time.Date(2024, time.May, 8, 0, 0, 0, 0, time.UTC),
	EndTime:   time.Date(2024, time.May, 9, 0, 0, 0, 0, time.UTC),
	Filters:   []ts.TSQueryFilter{{Key: "index_num", Value: "use4-ngn.gv52.4", Regex: false, Not: false}},
	GroupBy:   []string{"index_num"},
	Step:      time.Minute * 5,
	Timeout:   time.Second * 30,
}

func TestPrometheus(t *testing.T) {
	tsdb := prometheus.NewPrometheusClient("http://mgmtsrv1.sv11.edn.equinix.com:32090", "", "")

	if result, err := tsdb.Query(context.Background(), BaseQuery); err != nil {
		t.Errorf("got an error: %v", err)
	} else {
		pretty.Print(result)
	}

}
func TestMimir(t *testing.T) {
	tsdb := prometheus.NewMimirClient("sv5-edn-mimir-stg.lab.equinix.com", "eot-telemetry")

	if result, err := tsdb.Query(context.Background(), BaseQuery); err != nil {
		t.Errorf("got an error: %v", err)
	} else {
		pretty.Print(result)
	}

}

func TestInfluxDB(t *testing.T) {
	options := influxdb2.DefaultOptions()
	options.SetLogLevel(3)
	tsdb := influxdb.NewInfluxDBClient("http://devsv3ednmgmt09.lab.equinix.com:30320", "mytoken", "testing_script", "primary", options)

	if result, err := tsdb.Query(context.Background(), BaseQuery); err != nil {
		t.Errorf("got an error: %v", err)
	} else {
		pretty.Print(result)
	}

}
func TestAllTSDBs(t *testing.T) {
	options := influxdb2.DefaultOptions()
	options.SetFlushInterval(5_000)
	options.SetLogLevel(3)

	tsdbs := []struct {
		name string
		db   ts.TSDB
	}{
		{name: "Prometheus", db: prometheus.NewPrometheusClient("http://mgmtsrv1.sv11.edn.equinix.com:32090", "", "")},
		{name: "Mimir", db: prometheus.NewMimirClient("sv5-edn-mimir-stg.lab.equinix.com", "eot-telemetry")},
		{name: "InfluxDB", db: influxdb.NewInfluxDBClient("http://devsv3ednmgmt09.lab.equinix.com:30320", "mytoken", "testing_script", "primary", options)},
	}

	for _, tsdb := range tsdbs {
		t.Run(tsdb.name, func(t *testing.T) {
			if result, err := tsdb.db.Query(context.Background(), BaseQuery); err != nil {
				t.Errorf("got an error: %v", err)
			} else {
				pretty.Print(result)
			}
		})

	}
}
