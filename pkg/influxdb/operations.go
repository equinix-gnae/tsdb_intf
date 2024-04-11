package influxdb

import (
	"fmt"

	"github.com/equinix-gnae/tsdb_intf/pkg/ts"
)

type add ts.Add

func (r add) Apply(queryStr *string) (err error) {
	applyOperation(queryStr, "+", r.Left, r.Right)
	return nil
}

type sub ts.Sub

func (r sub) Apply(queryStr *string) (err error) {
	applyOperation(queryStr, "-", r.Left, r.Right)
	return nil
}

type multiply ts.Multiply

func (r multiply) Apply(queryStr *string) (err error) {
	applyOperation(queryStr, "*", r.Left, r.Right)
	return nil
}

type divide ts.Divide

func (r divide) Apply(queryStr *string) (err error) {
	applyOperation(queryStr, "/", r.Left, r.Right)
	return nil
}

// _value column contains the value of the mertic
func applyOperation(queryStr *string, op, left, right string) {
	if left == "$" {
		*queryStr += fmt.Sprintf("|> map(fn: (r) => ({ r with _value: r._value %s %s }))\n", op, right)
	} else {
		*queryStr += fmt.Sprintf("|> map(fn: (r) => ({ r with _value: %s %s r._value }))\n", left, op)
	}
}
