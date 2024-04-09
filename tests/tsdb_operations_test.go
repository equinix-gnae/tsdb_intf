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

var BaseQueryWithOps = ts.TSQuery{
	Table:     "bits",
	StartTime: time.Date(2024, time.April, 8, 0, 0, 0, 0, time.UTC),
	EndTime:   time.Date(2024, time.April, 9, 0, 0, 0, 0, time.UTC),
	Filters:   map[string]string{"index_num": "use4-ngn.gv52.4"},
	GroupBy:   []string{"index_num"},
	Step:      time.Minute * 5,
	Timeout:   time.Second * 30,
	Functions: []ts.QueryFunction{
		ts.Rate{Range: BaseQuery.EndTime.Sub(BaseQuery.StartTime)},
	},
	Operations: []ts.QueryOperation{
		ts.Add{Left: "$", Right: "100"},
		ts.Multiply{Left: "$", Right: "100"},
		ts.Divide{Left: "$", Right: "100"},
	},
}

func TestPrometheusOps(t *testing.T) {
	tsdb := prometheus.NewPrometheusClient("http://mgmtsrv1.sv11.edn.equinix.com:32090", "", "")

	if result, err := tsdb.Query(context.Background(), BaseQueryWithOps); err != nil {
		t.Errorf("got an error: %v", err)
	} else {
		pretty.Print(result)
	}

}

func TestMimirOps(t *testing.T) {
	tsdb := prometheus.NewMimirClient("sv5-edn-mimir-stg.lab.equinix.com", "eot-telemetry")

	if result, err := tsdb.Query(context.Background(), BaseQueryWithOps); err != nil {
		t.Errorf("got an error: %v", err)
	} else {
		pretty.Print(result)
	}

}

func TestInfluxDBOps(t *testing.T) {
	options := influxdb2.DefaultOptions()
	options.SetLogLevel(3)
	tsdb := influxdb.NewInfluxDBClient("http://devsv3ednmgmt09.lab.equinix.com:30320", "mytoken", "testing_script", "primary", options)

	if result, err := tsdb.Query(context.Background(), BaseQueryWithOps); err != nil {
		t.Errorf("got an error: %v", err)
	} else {
		pretty.Print(result)
	}

}
