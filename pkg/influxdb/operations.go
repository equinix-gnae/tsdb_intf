package influxdb

import "github.com/equinix-gnae/tsdb_intf/pkg/ts"

type add ts.Add

func (r add) Apply(queryStr *string) (err error) {
	panic("Apply should be implemented by TSDB driver")
}

type sub ts.Sub

func (r sub) Apply(queryStr *string) (err error) {
	panic("Apply should be implemented by TSDB driver")
}

type multiply ts.Multiply

func (r multiply) Apply(queryStr *string) (err error) {
	panic("Apply should be implemented by TSDB driver")
}

type divide ts.Divide

func (r divide) Apply(queryStr *string) (err error) {
	panic("Apply should be implemented by TSDB driver")
}
