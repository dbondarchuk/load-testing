package main

// HttpAction type for info about running http request
type HttpAction struct {
	Method       string                 `json:"method"`
	Endpoint     string                 `json:"endPoint"`
	BodyType     string                 `json:"bodyType"`
	RawData      string                 `json:"rawData"`
	FormData     map[string]interface{} `json:"formData"`
	Files        map[string]interface{} `json:"files"`
	Headers      map[string]interface{} `json:"headers"`
	Cookies      map[string]interface{} `json:"cookies"`
	Name         string                 `json:"name"`
	VariableName string                 `json:"variableName"`
	TimeOut      int                    `json:"timeOut"`
	StoreCookie  string                 `json:"storeCookie"`
}

// Execute action
func (h HttpAction) Execute(resultsChannel chan HttpReqResult, variables map[string]interface{}) {
	DoHttpRequest(h, resultsChannel, variables)
}

// NewHttpAction - creates new HttpAction
func NewHttpAction(a map[string]interface{}) HttpAction {
	var storeCookie string
	if a["storeCookie"] != nil && a["storeCookie"].(string) != "" {
		storeCookie = a["storeCookie"].(string)
	}

	bodyType, ok := a["bodyType"].(string)
	rawData, ok := a["rawData"].(string)
	variableName, ok := a["variableName"].(string)
	timeOut, ok := a["timeOut"].(int)
	formData, ok := a["formData"].(map[string]interface{})
	files, ok := a["files"].(map[string]interface{})
	headers, ok := a["headers"].(map[string]interface{})
	cookies, ok := a["cookies"].(map[string]interface{})

	if !ok {
	}

	httpAction := HttpAction{
		a["method"].(string),
		a["endPoint"].(string),
		bodyType,
		rawData,
		formData,
		files,
		headers,
		cookies,
		a["name"].(string),
		variableName,
		timeOut,
		storeCookie,
	}

	return httpAction
}
