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
	Name            string            `json:"name"`
	Labels          map[string]string `json:"labels"`
	TimeValueSeries []TimeValue       `json:"timeValueSeries"`
}

type TSDBQueryResult []TimeSeries

type TSDBStore interface {
	Query(ctx context.Context, query TSQuery, opts map[string]any) TSDBQueryResult
}

type TSQuery struct {
	Table     string            `json:"table"`
	StartTime time.Time         `json:"startTime"`
	EndTime   time.Time         `json:"endTime"`
	Step      time.Duration     `json:"step"`
	Filters   map[string]string `json:"filters"`
	GroupBy   []string          `jsone:"groupBy"` // used by influx to generate group key which common for TS
}
