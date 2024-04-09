package utils

import (
	"fmt"

	"github.com/equinix-gnae/tsdb_intf/pkg/influxdb"
	"github.com/equinix-gnae/tsdb_intf/pkg/prometheus"
	"github.com/equinix-gnae/tsdb_intf/pkg/ts"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func NewTSDBClient(config *ts.Config) (ts.TSDB, error) {
	switch config.Name {
	case "prometheus":
		return prometheus.NewPrometheusClient(config.URL, config.ID, config.Secret), nil
	case "mimir":
		return prometheus.NewMimirClient(config.URL, config.Org), nil
	case "influxDB":
		// TODO: need to pass client speicific options as well
		return influxdb.NewInfluxDBClient(config.URL, config.Secret, config.Database, config.Org, influxdb2.DefaultOptions()), nil
	default:
		return nil, fmt.Errorf("unsupported config name: %q", config.Name)
	}
}
