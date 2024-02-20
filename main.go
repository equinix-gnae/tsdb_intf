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
		Filters:   map[string]string{"index_num": "bb1-ngn.gv51.1001", "_measurement": ""},
		GroupBy:   []string{"index_num", "_measurement"},
		Step:      time.Hour * 2,
	}

	// var tsdbStore1 tsdb.TSDBStore = tsdb.NewMimirDBStore("sv5-edn-mimir-stg.lab.equinix.com", "eot-telemetry")
	// pretty.Print(tsdbStore1.Query(context.Background(), query, map[string]any{"timeout": time.Second * 30}))

	var tsdbStore2 tsdb.TSDBStore = tsdb.NewInfluxDBStore("http://devsv3ednmgmt09.lab.equinix.com:30320", "mytoken")
	pretty.Print(tsdbStore2.Query(context.Background(), query, map[string]any{"org": "primary", "bucket": "testing_script"}))
}
