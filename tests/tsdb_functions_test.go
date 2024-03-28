package ts_test

import (
	"context"
	"testing"
	"time"

	"github.com/equinix-gnae/tsdb_intf/pkg/ts"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/kr/pretty"
)

var BaseQueryWithFunc = ts.TSQuery{
	Table:     "bits",
	StartTime: time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC),
	EndTime:   time.Date(2024, time.March, 22, 0, 0, 0, 0, time.UTC),
	Filters:   map[string]string{"index_num": "use4-ngn.gv52.4"},
	GroupBy:   []string{"index_num"},
	Step:      time.Minute * 5,
	Timeout:   time.Second * 30,
	Functions: []ts.QueryFunction{
		ts.Rate{Range: BaseQuery.EndTime.Sub(BaseQuery.StartTime)},
		ts.Sum{},
	},
}

func TestPrometheusFuncs(t *testing.T) {
	tsdb := ts.NewPrometheusClient("http://mgmtsrv1.sv11.edn.equinix.com:32090", "", "")

	if result, err := tsdb.Query(context.Background(), BaseQueryWithFunc); err != nil {
		t.Errorf("got an error: %v", err)
	} else {
		pretty.Print(result)
	}

}

func TestMimirFuncs(t *testing.T) {
	tsdb := ts.NewMimirClient("sv5-edn-mimir-stg.lab.equinix.com", "eot-telemetry")

	if result, err := tsdb.Query(context.Background(), BaseQueryWithFunc); err != nil {
		t.Errorf("got an error: %v", err)
	} else {
		pretty.Print(result)
	}

}

func TestInfluxDBFuncs(t *testing.T) {
	options := influxdb2.DefaultOptions()
	options.SetLogLevel(3)
	tsdb := ts.NewInfluxDBClient("http://devsv3ednmgmt09.lab.equinix.com:30320", "mytoken", "testing_script", "primary", options)

	if result, err := tsdb.Query(context.Background(), BaseQueryWithFunc); err != nil {
		t.Errorf("got an error: %v", err)
	} else {
		pretty.Print(result)
	}

}
