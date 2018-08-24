package main

// Action basic interface
type Action interface {
	GetStep() *TestStepValue
	Execute(httpResultsChannel chan HttpReqResult, variables map[string]interface{}) error
}
