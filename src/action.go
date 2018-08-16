package main

// Action basic interface
type Action interface {
	Execute(resultsChannel chan HttpReqResult, variables map[string]interface{}) error
}
