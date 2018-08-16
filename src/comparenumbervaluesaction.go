package main

import (
	"errors"
	"strconv"
)

type CompareNumberValuesAction struct {
	Method string `json:"method"`
	Value  string `json:"value"`
	To     string `json:"to"`
}

// Execute action
func (h CompareNumberValuesAction) Execute(resultsChannel chan HttpReqResult, variables map[string]interface{}) error {
	value, _ := strconv.Atoi(SubstParams(variables, h.Value))
	to, _ := strconv.Atoi(SubstParams(variables, h.To))
	var isValid = false

	switch method := h.Method; method {
	case "lt":
		isValid = value < to
		break

	case "le":
		isValid = value <= to
		break

	case "eq":
		isValid = value == to
		break

	case "gt":
		isValid = value > to
		break

	case "ge":
		isValid = value >= to
		break
	}

	if !isValid {
		return errors.New("Wrong assertion. Expected: " + strconv.Itoa(to) + ", was: " + strconv.Itoa(value) + ", method: " + h.Method)
	}

	return nil
}

func NewCompareNumberValuesAction(a map[string]interface{}) CompareNumberValuesAction {
	compareNumberValuesAction := CompareNumberValuesAction{
		a["method"].(string),
		a["value"].(string),
		a["to"].(string),
	}

	return compareNumberValuesAction
}
