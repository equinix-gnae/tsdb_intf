package ts

/*
All the query operations can specify operation query result using '$' sign

Example:
========
	Add{Left: '$', Right, '10'} => queryResult + 10


Note: Operations should be applied after the functions like

*/

type QueryOperation interface {
	Apply(queryStr *string) (err error)
}

type Add struct {
	Left  string `json:"left"`
	Right string `json:"right"`
}

func (r Add) Apply(queryStr *string) (err error) {
	panic("Apply should be implemented by TSDB driver")
}

type Sub struct {
	Left  string `json:"left"`
	Right string `json:"right"`
}

func (r Sub) Apply(queryStr *string) (err error) {
	panic("Apply should be implemented by TSDB driver")
}

type Multiply struct {
	Left  string `json:"left"`
	Right string `json:"right"`
}

func (r Multiply) Apply(queryStr *string) (err error) {
	panic("Apply should be implemented by TSDB driver")
}

type Divide struct {
	Left  string `json:"left"`
	Right string `json:"right"`
}

func (r Divide) Apply(queryStr *string) (err error) {
	panic("Apply should be implemented by TSDB driver")
}
