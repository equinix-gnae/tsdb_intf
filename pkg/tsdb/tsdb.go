package tsdb

import (
	"context"
	"time"
)

type TimeValue struct {
	Time  int64   `json:"time"`
	Value float64 `json:"value"`
}

type TimeSeries struct {
	//Name            string            `json:"name"`
	Labels          map[string]string `json:"labels"`
	TimeValueSeries []TimeValue       `json:"timeValueSeries"`
}

type TSDBQueryResult []TimeSeries

type TSQuery struct {
	Table     string            `json:"table"` // for prometheus the table is metric name, for influx the table is _field
	StartTime time.Time         `json:"startTime"`
	EndTime   time.Time         `json:"endTime"`
	Step      time.Duration     `json:"step"`
	Filters   map[string]string `json:"filters"`  // for prometheus filters are kv lables, for influx they are `where` clause
	GroupBy   []string          `jsone:"groupBy"` // used by influx to generate group key which is common for a TS
}

// read-only
type TSDBStore interface {
	Query(ctx context.Context, query TSQuery, opts map[string]any) TSDBQueryResult
}
