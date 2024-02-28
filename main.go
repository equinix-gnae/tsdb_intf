package main

import (
	"context"
	"time"

	"github.com/equinix-gnae/tsdb_intf/pkg/ts"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/kr/pretty"
)

func main() {
	query := ts.TSQuery{
		Table:     "bits",
		StartTime: time.Date(2024, time.February, 12, 0, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2024, time.February, 13, 0, 0, 0, 0, time.UTC),
		Filters:   map[string]string{"index_num": "bb1-ngn.gv51.1001"},
		GroupBy:   []string{"index_num", "_measurement"},
		Step:      time.Hour * 2,
		Timeout:   time.Second * 30,
	}
	//var tsdbStore ts.TSStore = ts.NewMimirDBStore("sv5-edn-mimir-stg.lab.equinix.com", "eot-telemetry")

	options := influxdb2.DefaultOptions()
	options.SetFlushInterval(5_000)
	options.SetLogLevel(3)

	var tsdbStore ts.TSStore = ts.NewInfluxDBStore("http://devsv3ednmgmt09.lab.equinix.com:30320", "mytoken", "testing_script", "primary", options)
	pretty.Print(tsdbStore.Query(context.Background(), query, map[string]any{}))
}
