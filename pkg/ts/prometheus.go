package ts

import (
	"context"
	"log"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

type PrometheusClient struct {
	Client v1.API
}

func NewPrometheusClient(url string) PrometheusClient {
	cfg := api.Config{Address: url}

	httpClient, err := api.NewClient(cfg)
	if err != nil {
		log.Fatalf("get httpClient failed: %+v", err)
	}
	v1API := v1.NewAPI(httpClient)

	if err != nil {
		log.Fatalln(err)
	}

	return PrometheusClient{Client: v1API}
}

func (r PrometheusClient) Query(ctx context.Context, query TSQuery) (TSQueryResult, error) {
	Range := v1.Range{Start: query.StartTime, End: query.EndTime, Step: query.Step}
	strQuery := GeneratePromQueryString(query)

	resp, warn, err := r.Client.QueryRange(ctx, strQuery, Range)

	if err != nil {
		return nil, err
	}

	if warn != nil {
		log.Printf("WARN| query: %q, warning: %v", query, warn)
	}

	return PromQueryResultToTS(resp, strQuery)
}
