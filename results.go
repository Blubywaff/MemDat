package main

const NO_STATUS = status(0)
const SUCCESS = status(1)
const FAILURE = status(2)

type status uint8

type Result struct {
	result string
	status status
}

func (r *Result) IsError() bool {
	return r.status == FAILURE
}

func (r *Result) IsNil() bool {
	return r.status == NO_STATUS
}

func (r *Result) Result() string {
	return r.result
}

func newResult(result string, status status) *Result {
	return &Result{
		result: result,
		status: status,
	}
}

func nilResult() *Result {
	return &Result{
		result: "",
		status: NO_STATUS,
	}
}
