package ts

import (
	"context"
	"strings"
	"time"

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

type TSQueryFilter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Regex bool   `json:"regex"`
	Not   bool   `json:"not"` // for not equal etc
}

/*
caution: for prometheus/mimir step can't be 0 otherwise we get following error:
bad_data: invalid parameter "step": zero or negative query resolution step widths are not accepted. Try a positive intege
*/
type TSQuery struct {
	Table      string           `json:"table"`      // for prometheus the table is metric name, for influx the table is _field
	Filters    []TSQueryFilter  `json:"filters"`    // for prometheus filters are kv lables, for influx they are `where` clause
	Functions  []QueryFunction  `json:"functions"`  // functions like rate, sum
	Operations []QueryOperation `json:"operations"` // operations like add, sub, div
	GroupBy    []string         `json:"groupBy"`    // used by influx to generate group key which is common for a TS
	StartTime  time.Time        `json:"startTime"`
	EndTime    time.Time        `json:"endTime"`
	Timeout    time.Duration    `json:"timeout"`
	Step       time.Duration    `json:"step"`
}

type TSDB interface {
	Query(ctx context.Context, query TSQuery) (TSQueryResult, error)
}

func init() {
	v := viper.GetViper()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
}
