package prometheus

import (
	"fmt"
	"strings"

	"github.com/equinix-gnae/tsdb_intf/pkg/ts"
	"github.com/prometheus/common/model"
)

// Query String Example => query := `bits{index_num="bb1-ngn.gv51.1001"}`
func generateInstantVector(query ts.TSQuery) string {
	var queryBuilder strings.Builder

	queryBuilder.WriteString(query.Table)
	queryBuilder.WriteString("{")
	for k, v := range query.Filters {
		queryBuilder.WriteString(fmt.Sprintf("%s=%q, ", k, v))
	}
	queryBuilder.WriteString("}")

	return queryBuilder.String()
}

// Query String Example => query := `rate(bits{index_num="bb1-ngn.gv51.1001"}[5m])`
func applyFunctions(queryStr *string, functions []ts.QueryFunction) error {
	for _, queryFunction := range functions {
		switch t := queryFunction.(type) {
		case ts.Rate:
			rate(t).Apply(queryStr)
		case ts.Sum:
			sum(t).Apply(queryStr)
		default:
			return fmt.Errorf("unsupported Function: %v", queryFunction)
		}
	}
	return nil
}

// Query String Example => query := `rate(bits{index_num="bb1-ngn.gv51.1001"}[5m])`
func applyOperations(queryStr *string, operations []ts.QueryOperation) error {
	for _, operation := range operations {
		switch t := operation.(type) {
		case ts.Add:
			add(t).Apply(queryStr)
		case ts.Substract:
			substract(t).Apply(queryStr)
		case ts.Multiply:
			multiply(t).Apply(queryStr)
		case ts.Divide:
			divide(t).Apply(queryStr)
		default:
			return fmt.Errorf("unsupported operation: %v", operation)
		}
	}
	return nil
}

func GeneratePromQueryString(query ts.TSQuery) (string, error) {
	queryStr := generateInstantVector(query)

	if err := applyFunctions(&queryStr, query.Functions); err != nil {
		return "", err
	}

	if err := applyOperations(&queryStr, query.Operations); err != nil {
		return "", err
	}

	fmt.Printf("query: %q\n", queryStr)
	return queryStr, nil
}

func PromQueryResultToTS(promQueryResult model.Value, strQuery string) (ts.TSQueryResult, error) {
	matrix, ok := promQueryResult.(model.Matrix)

	if !ok {
		return nil, fmt.Errorf("for query %q, unable to convert resp to vector", strQuery)
	}

	if len(matrix) == 0 {
		return nil, fmt.Errorf("for query %q, empty response is returned", strQuery)
	}

	result := make(ts.TSQueryResult, 0, len(matrix))

	for _, sampleStream := range matrix {
		labels := make(map[string]string)
		for key, val := range sampleStream.Metric {
			labels[string(key)] = string(val)
		}

		timeValueSeries := make([]ts.TimeValue, 0, len(sampleStream.Values))
		for _, sample := range sampleStream.Values {
			timeValueSeries = append(timeValueSeries, ts.TimeValue{Time: int64(sample.Timestamp), Value: float64(sample.Value)})
		}

		result = append(result, ts.TimeSeries{
			Labels:          labels,
			TimeValueSeries: timeValueSeries,
		})
	}
	return result, nil
}
