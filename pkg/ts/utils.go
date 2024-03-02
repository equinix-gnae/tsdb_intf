package ts

import (
	"fmt"
	"strings"

	"github.com/prometheus/common/model"
)

/*
	 Query String Example:
		query := `bits{index_num="bb1-ngn.gv51.1001"}`
*/
func GeneratePromQueryString(query TSQuery) string {
	var queryBuilder strings.Builder
	queryBuilder.WriteString(query.Table)

	queryBuilder.WriteString("{")
	for k, v := range query.Filters {
		queryBuilder.WriteString(fmt.Sprintf("%s=%q, ", k, v))
	}
	queryBuilder.WriteString("}")

	return queryBuilder.String()
}

func PromQueryResultToTS(promQueryResult model.Value, strQuery string) (TSQueryResult, error) {
	matrix, ok := promQueryResult.(model.Matrix)

	if !ok {
		return nil, fmt.Errorf("for query %q, unable to convert resp to vector", strQuery)
	}

	if len(matrix) == 0 {
		return nil, fmt.Errorf("for query %q, empty response is returned", strQuery)
	}

	result := make(TSQueryResult, 0, len(matrix))

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
	return result, nil
}
