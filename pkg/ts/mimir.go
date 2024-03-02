package ts

import (
	"context"
	"log"

	"github.com/grafana/mimir/integration/e2emimir"
)

type MimirClient struct {
	Client *e2emimir.Client
}

func NewMimirClient(url string, id string) MimirClient {
	mimirE2eClient, err := e2emimir.NewClient("", url, "", "", id)

	if err != nil {
		log.Fatalln(err)
	}

	return MimirClient{Client: mimirE2eClient}
}

func (r MimirClient) Query(ctx context.Context, query TSQuery) (TSQueryResult, error) {
	// XXX: handle the case where client is shared b/w goroutines?
	r.Client.SetTimeout(query.Timeout)
	strQuery := GeneratePromQueryString(query)
	resp, err := r.Client.QueryRange(strQuery, query.StartTime, query.EndTime, query.Step)

	if err != nil {
		log.Fatalln(err)
	}

	return PromQueryResultToTS(resp, strQuery)
}
