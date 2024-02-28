package main

import (
	"context"
	"time"

	"github.com/equinix-gnae/tsdb_intf/pkg/tsdb"
	"github.com/kr/pretty"
)

func main() {
	query := tsdb.TSQuery{
		Table:     "bits",
		StartTime: time.Date(2024, time.February, 12, 0, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2024, time.February, 13, 0, 0, 0, 0, time.UTC),
		Filters:   map[string]string{"index_num": "bb1-ngn.gv51.1001"},
		GroupBy:   []string{"index_num", "_measurement"},
		Step:      time.Hour * 2,
		Timeout:   time.Second * 30,
	}
	var tsdbStore tsdb.TSDBStore = tsdb.NewMimirDBStore("sv5-edn-mimir-stg.lab.equinix.com", "eot-telemetry")
	//var tsdbStore tsdb.TSDBStore = tsdb.NewInfluxDBStore("http://devsv3ednmgmt09.lab.equinix.com:30320", "mytoken", "testing_script", "primary")
	pretty.Print(tsdbStore.Query(context.Background(), query, map[string]any{}))
}
