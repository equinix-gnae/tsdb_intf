package tsdb

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/grafana/mimir/integration/e2emimir"
	"github.com/prometheus/common/model"
)

type MimirDBStore struct {
	Client *e2emimir.Client
}

func NewMimirDBStore(url string, id string) MimirDBStore {
	mimirE2eClient, err := e2emimir.NewClient("", url, "", "", id)

	if err != nil {
		log.Fatalln(err)
	}

	return MimirDBStore{Client: mimirE2eClient}
}

/*
	 Query String Example:
		query := `bits{index_num="bb1-ngn.gv51.1001"}`
*/
func (r MimirDBStore) Query(ctx context.Context, query TSQuery, opts map[string]any) TSDBQueryResult {
	if timeout, ok := opts["timeout"]; ok {
		r.Client.SetTimeout(timeout.(time.Duration))
	}

	// build query string
	var queryBuilder strings.Builder
	queryBuilder.WriteString(query.Table)
	queryBuilder.WriteString("{")
	for k, v := range query.Filters {
		queryBuilder.WriteString(fmt.Sprintf("%s=%q, ", k, v))
	}
	queryBuilder.WriteString("}")

	// caution: for prometheus/mimir step can't be 0/negative
	// bad_data: invalid parameter "step": zero or negative query resolution step widths are not accepted.
	// Try a positive intege

	resp, err := r.Client.QueryRange(queryBuilder.String(), query.StartTime, query.EndTime, query.Step)

	if err != nil {
		log.Fatalln(err)
	}

	matrix, ok := resp.(model.Matrix)

	if !ok {
		log.Fatalf("unable to convert resp to vector")
	}

	if len(matrix) == 0 {
		log.Fatalf("empty response is returned for query: %q", query)
	}

	//pretty.Print(matrix)

	result := make(TSDBQueryResult, 0, len(matrix))

	for _, sampleStream := range matrix {
		labels := make(map[string]string)
		for key, val := range sampleStream.Metric {
			labels[string(key)] = string(val)
		}

		timeValueSeries := make([]TimeValue, 0, len(sampleStream.Values))
		for _, sample := range sampleStream.Values {
			timeValueSeries = append(timeValueSeries, TimeValue{Time: int64(sample.Timestamp), Value: float64(sample.Value)})
		}

		result = append(result, TimeSeries{
			Labels:          labels,
			TimeValueSeries: timeValueSeries,
		})
	}
	return result
}
