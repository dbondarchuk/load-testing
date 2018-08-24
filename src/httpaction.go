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
	Step         TestStepValue          `json:"-"`
}

// Execute action
func (h HttpAction) Execute(httpResultsChannel chan HttpReqResult, variables map[string]interface{}) error {
	return DoHttpRequest(h, httpResultsChannel, variables)
}

func (h HttpAction) GetStep() *TestStepValue {
	return &h.Step
}

// NewHttpAction - creates new HttpAction
func NewHttpAction(s TestStepValue) HttpAction {
	var storeCookie string
	a := s.PropertyValues
	if a["storeCookie"] != nil && a["storeCookie"].(string) != "" {
		storeCookie = a["storeCookie"].(string)
	}

	bodyType, _ := a["bodyType"].(string)
	rawData, _ := a["rawData"].(string)
	variableName, _ := a["variableName"].(string)
	timeOut, _ := a["timeOut"].(int)
	formData, _ := a["formData"].(map[string]interface{})
	files, _ := a["files"].(map[string]interface{})
	headers, _ := a["headers"].(map[string]interface{})
	cookies, _ := a["cookies"].(map[string]interface{})

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
		s,
	}

	return httpAction
}
