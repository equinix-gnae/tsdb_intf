package ts_test

import (
	"context"
	"testing"
	"time"

	"github.com/equinix-gnae/tsdb_intf/pkg/ts"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/kr/pretty"
)

var BaseQuery = ts.TSQuery{
	Table:     "bits",
	StartTime: time.Now().Add(-1 * time.Hour),
	EndTime:   time.Now().UTC(),
	Filters:   map[string]string{"index_num": "use4-ngn.gv52.4"},
	GroupBy:   []string{"index_num", "_measurement"},
	Step:      time.Minute * 5,
	Timeout:   time.Second * 30,
}

func TestPrometheus(t *testing.T) {
	tsdb := ts.NewPrometheusClient("http://mgmtsrv1.sv11.edn.equinix.com:32090", "", "")

	if result, err := tsdb.Query(context.Background(), BaseQuery); err != nil {
		t.Errorf("got an error: %v", err)
	} else {
		pretty.Print(result)
	}

}
func TestMimir(t *testing.T) {
	tsdb := ts.NewMimirClient("sv5-edn-mimir-stg.lab.equinix.com", "eot-telemetry")

	if result, err := tsdb.Query(context.Background(), BaseQuery); err != nil {
		t.Errorf("got an error: %v", err)
	} else {
		pretty.Print(result)
	}

}

func TestInfluxDB(t *testing.T) {
	options := influxdb2.DefaultOptions()
	options.SetLogLevel(3)
	tsdb := ts.NewInfluxDBClient("http://devsv3ednmgmt09.lab.equinix.com:30320", "mytoken", "testing_script", "primary", options)

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
		{name: "Prometheus", db: ts.NewPrometheusClient("http://mgmtsrv1.sv11.edn.equinix.com:32090", "", "")},
		{name: "Mimir", db: ts.NewMimirClient("sv5-edn-mimir-stg.lab.equinix.com", "eot-telemetry")},
		{name: "InfluxDB", db: ts.NewInfluxDBClient("http://devsv3ednmgmt09.lab.equinix.com:30320", "mytoken", "testing_script", "primary", options)},
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
