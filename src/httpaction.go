package main

import (
	"encoding/json"
	"strings"
)

// HttpAction type for info about running http request
type HttpAction struct {
	Method       string         `json:"method"`
	Endpoint     string         `json:"endPoint"`
	BodyType     string         `json:"bodyType"`
	RawData      string         `json:"rawData"`
	FormData     []KeyValuePair `json:"formData"`
	Files        []KeyValuePair `json:"files"`
	Headers      []KeyValuePair `json:"headers"`
	Cookies      []KeyValuePair `json:"cookies"`
	VariableName string         `json:"variableName"`
	TimeOut      int            `json:"timeOut"`
	Step         TestStepValue  `json:"-"`
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
	a := s.PropertyValues

	bodyType, _ := a["bodyType"].(string)
	rawData, _ := a["rawData"].(string)
	variableName, _ := a["variableName"].(string)
	timeOut, _ := a["timeOut"].(int)

	var formData []KeyValuePair
	formDataJson, _ := json.Marshal(a["formData"])
	err := json.Unmarshal(formDataJson, &formData)

	var files []KeyValuePair
	filesJson, _ := json.Marshal(a["files"])
	err = json.Unmarshal(filesJson, &files)

	var headers []KeyValuePair
	headersJson, _ := json.Marshal(a["headers"])
	err = json.Unmarshal(headersJson, &headers)

	var cookies []KeyValuePair
	cookiesJson, _ := json.Marshal(a["cookies"])
	err = json.Unmarshal(cookiesJson, &cookies)

	if err != nil {
	}

	httpAction := HttpAction{
		strings.ToUpper(a["method"].(string)),
		a["endPoint"].(string),
		bodyType,
		rawData,
		formData,
		files,
		headers,
		cookies,
		variableName,
		timeOut,
		s,
	}

	return httpAction
}
