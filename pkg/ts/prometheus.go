package ts

import (
	"context"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

type withBasicAuthRoundTripper struct {
	username string
	password string
	next     http.RoundTripper
}

func (r *withBasicAuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if r.username != "" && r.password != "" {
		req.SetBasicAuth(r.username, r.password)
	}

	return r.next.RoundTrip(req)
}

type PrometheusClient struct {
	Client v1.API
}

func NewPrometheusClient(url string, username string, password string) PrometheusClient {
	roundTripper := withBasicAuthRoundTripper{
		username: username,
		password: password,
		next:     http.DefaultTransport.(*http.Transport).Clone(),
	}

	cfg := api.Config{
		Address:      url,
		RoundTripper: &roundTripper,
	}
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

	resp, warn, err := r.Client.QueryRange(ctx, strQuery, Range, v1.WithTimeout(query.Timeout))

	if err != nil {
		return nil, err
	}

	if warn != nil {
		log.Printf("WARN| query: %q, warning: %v", query, warn)
	}

	return PromQueryResultToTS(resp, strQuery)
}
