package ts

import "time"

type QueryFunction interface {
	Apply(queryStr *string) (err error)
}

type Rate struct {
	Range time.Duration `json:"range"`
}

func (r Rate) Apply(queryStr *string) (err error) {
	panic("Apply should be implemented by TSDB driver")
}

type Sum struct {
	Range time.Duration `json:"range"`
}

func (r Sum) Apply(queryStr *string) (err error) {
	panic("Apply should be implemented by TSDB driver")
}
