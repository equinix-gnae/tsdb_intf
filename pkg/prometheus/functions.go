package prometheus

import "github.com/equinix-gnae/tsdb_intf/pkg/ts"

type rate ts.Rate

func (r rate) Apply(queryStr *string) (err error) {

	*queryStr = "rate(" + *queryStr + "[" + r.Range.String() + "]" + ")"
	return err
}

type sum ts.Sum

func (r sum) Apply(queryStr *string) (err error) {
	*queryStr = "sum (" + *queryStr + ")"
	return err
}
