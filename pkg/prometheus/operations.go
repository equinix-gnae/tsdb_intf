package prometheus

import (
	"github.com/equinix-gnae/tsdb_intf/pkg/ts"
)

type add ts.Add

func (r add) Apply(queryStr *string) (err error) {
	applyOperation(queryStr, "+", r.Left, r.Right)
	return nil
}

type subtract ts.Subtract

func (r subtract) Apply(queryStr *string) (err error) {
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

func applyOperation(queryStr *string, op, left, right string) {
	if left == "$" {
		*queryStr = "(" + *queryStr + " " + op + " " + right + ") "
	} else {
		*queryStr = "(" + left + " " + op + " " + *queryStr + ") "
	}
}
