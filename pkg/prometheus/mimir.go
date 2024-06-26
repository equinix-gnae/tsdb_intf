package prometheus

import (
	"context"
	"log"

	"github.com/equinix-gnae/tsdb_intf/pkg/ts"
	"github.com/grafana/mimir/integration/e2emimir"
)

type MimirClient struct {
	Client *e2emimir.Client
}

func NewMimirClient(address string, id string) MimirClient {
	mimirE2eClient, err := e2emimir.NewClient("", address, "", "", id)

	if err != nil {
		log.Fatalln(err)
	}

	return MimirClient{Client: mimirE2eClient}
}

func (r MimirClient) Query(ctx context.Context, query ts.TSQuery) (ts.TSQueryResult, error) {
	// XXX: handle the case where client is shared b/w goroutines?
	r.Client.SetTimeout(query.Timeout)
	strQuery, err := GeneratePromQueryString(query)

	if err != nil {
		return nil, err
	}

	resp, err := r.Client.QueryRange(strQuery, query.StartTime, query.EndTime, query.Step)

	if err != nil {
		return nil, err
	}

	return PromQueryResultToTS(resp, strQuery)
}
