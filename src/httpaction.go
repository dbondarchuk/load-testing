package main

// HttpAction type for info about running http request
type HttpAction struct {
	Method       string            `json:"method"`
	Endpoint     string            `json:"endPoint"`
	BodyType     string            `json:"bodyType"`
	RawData      string            `json:"rawData"`
	FormData     map[string]string `json:"formData"`
	Files        map[string]string `json:"files"`
	Headers      map[string]string `json:"headers"`
	Cookies      map[string]string `json:"cookies"`
	Name         string            `json:"name"`
	VariableName string            `json:"variableName"`
	TimeOut      int               `json:"timeOut"`
	StoreCookie  string            `json:"storeCookie"`
}

// Execute action
func (h HttpAction) Execute(resultsChannel chan HttpReqResult, variables map[string]interface{}) {
	DoHttpRequest(h, resultsChannel, variables)
}

// NewHttpAction - creates new HttpAction
func NewHttpAction(a map[interface{}]interface{}) HttpAction {
	var storeCookie string
	if a["storeCookie"] != nil && a["storeCookie"].(string) != "" {
		storeCookie = a["storeCookie"].(string)
	}

	httpAction := HttpAction{
		a["method"].(string),
		a["endPoint"].(string),
		a["bodyType"].(string),
		a["rawData"].(string),
		a["formData"].(map[string]string),
		a["files"].(map[string]string),
		a["headers"].(map[string]string),
		a["cookies"].(map[string]string),
		a["name"].(string),
		a["variableName"].(string),
		a["timeOut"].(int),
		storeCookie,
	}

	return httpAction
}
