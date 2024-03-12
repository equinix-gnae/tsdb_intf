package ts

import (
	"context"
	"fmt"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/spf13/viper"
)

type TimeValue struct {
	Time  int64   `json:"time"`
	Value float64 `json:"value"`
}

type TimeSeries struct {
	Labels          map[string]string `json:"labels"`
	TimeValueSeries []TimeValue       `json:"timeValueSeries"`
}

type TSQueryResult []TimeSeries

/*
caution: for prometheus/mimir step can't be 0 otherwise we get following error:
bad_data: invalid parameter "step": zero or negative query resolution step widths are not accepted. Try a positive intege
*/
type TSQuery struct {
	Table     string            `json:"table"`   // for prometheus the table is metric name, for influx the table is _field
	Filters   map[string]string `json:"filters"` // for prometheus filters are kv lables, for influx they are `where` clause
	GroupBy   []string          `json:"groupBy"` // used by influx to generate group key which is common for a TS
	StartTime time.Time         `json:"startTime"`
	EndTime   time.Time         `json:"endTime"`
	Timeout   time.Duration     `json:"timeout"`
	Step      time.Duration     `json:"step"`
}

// read-only
type TSDB interface {
	Query(ctx context.Context, query TSQuery) (TSQueryResult, error)
}

func NewTSDBClient(config *Config) (TSDB, error) {
	switch config.Name {
	case "prometheus":
		return NewPrometheusClient(config.URL), nil
	case "mimir":
		return NewMimirClient(config.URL, config.Org), nil
	case "influxDB":
		// TODO: need to pass client speicific options as well
		return NewInfluxDBClient(config.URL, config.Secret, config.Database, config.Org, influxdb2.DefaultOptions()), nil
	default:
		return nil, fmt.Errorf("unsupported config name: %q", config.Name)
	}
}

func init() {
	v := viper.GetViper()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
}
