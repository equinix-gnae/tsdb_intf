package main

import (
	"context"
	"log"
	"time"

	"github.com/equinix-nspa/tsdb_intf/pkg/ts"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/kr/pretty"
	// influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	query := ts.TSQuery{
		Table:     "bits",
		StartTime: time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC),
		Filters:   map[string]string{"index_num": "bb1-ngn.gv51.1001"},
		GroupBy:   []string{"index_num", "_measurement"},
		Step:      time.Hour * 2,
		Timeout:   time.Second * 30,
	}

	// *** prometheus ***
	//var tsdbStore ts.TSSB = ts.NewPrometheusClient("http://mgmtsrv1.sv11.edn.equinix.com:32090")

	// *** mimir ***
	//var tsdbStore ts.TSSB = ts.NewMimirClient("sv5-edn-mimir-stg.lab.equinix.com", "eot-telemetry")

	// *** influx ***

	options := influxdb2.DefaultOptions()
	options.SetFlushInterval(5_000)
	options.SetLogLevel(3)
	var tsdbStore ts.TSDB = ts.NewInfluxDBClient("http://devsv3ednmgmt09.lab.equinix.com:30320", "mytoken", "testing_script", "primary", options)

	// *** query ***

	if result, err := tsdbStore.Query(context.Background(), query); err != nil {
		log.Printf("got an error: %v", err)
	} else {
		pretty.Print(result)
	}

}
