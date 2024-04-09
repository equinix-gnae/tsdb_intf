package prometheus

import (
	"github.com/equinix-gnae/tsdb_intf/pkg/ts"
)

type add ts.Add

func (r add) Apply(queryStr *string) (err error) {
	if r.Left == "$" {
		*queryStr = *queryStr + " + " + r.Right
	} else {
		*queryStr = r.Left + " + " + *queryStr
	}

	return nil
}

type sub ts.Sub

func (r sub) Apply(queryStr *string) (err error) {
	if r.Left == "$" {
		*queryStr = *queryStr + " - " + r.Right
	} else {
		*queryStr = r.Left + " - " + *queryStr
	}

	return nil
}

type multiply ts.Multiply

func (r multiply) Apply(queryStr *string) (err error) {
	if r.Left == "$" {
		*queryStr = *queryStr + " * " + r.Right
	} else {
		*queryStr = r.Left + " * " + *queryStr
	}

	return nil
}

type divide ts.Divide

func (r divide) Apply(queryStr *string) (err error) {
	if r.Left == "$" {
		*queryStr = *queryStr + " / " + r.Right
	} else {
		*queryStr = r.Left + " / " + *queryStr
	}

	return nil
}
