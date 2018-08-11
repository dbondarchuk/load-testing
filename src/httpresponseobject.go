package main

// HttpResponseObject contains repsonse from http request
type HttpResponseObject struct {
	StatusCode int               `json:"StatusCode"`
	Body       interface{}       `json:"Body"`
	Headers    map[string]string `json:"Headers"`
}
