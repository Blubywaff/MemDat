package main

const NoStatus = status(0)
const Success = status(1)
const Failure = status(2)

type status uint8

type Result struct {
	result string
	status status
}

func (r *Result) IsError() bool {
	return r.status == Failure
}

func (r *Result) IsNil() bool {
	return r.status == NoStatus
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
		status: NoStatus,
	}
}
