package influxdb

import (
	"fmt"

	"github.com/equinix-gnae/tsdb_intf/pkg/ts"
)

type rate ts.Rate

func (r rate) Apply(queryStr *string) (err error) {
	*queryStr += fmt.Sprintf("|> derivative(unit:%v, nonNegative: true)\n", r.Range)
	return nil
}

type sum ts.Sum

func (r sum) Apply(queryStr *string) (err error) {
	*queryStr += "|> sum()\n"
	return nil
}
