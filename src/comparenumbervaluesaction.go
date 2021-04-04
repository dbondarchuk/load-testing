package main

import (
	"errors"
	"strconv"
)

type CompareNumberValuesAction struct {
	Method string        `json:"method"`
	Value  string        `json:"value"`
	To     string        `json:"to"`
	Step   TestStepValue `json:"-"`
}

func (h CompareNumberValuesAction) GetStep() *TestStepValue {
	return &h.Step
}

// Execute action
func (h CompareNumberValuesAction) Execute(httpResultsChannel chan HttpReqResult, variables map[string]interface{}) error {
	valueStr, err := SubstParams(variables, h.Value)
	if err != nil {
		return nil
	}

	toStr, err := SubstParams(variables, h.To)
	if err != nil {
		return nil
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return nil
	}

	to, err := strconv.Atoi(toStr)
	if err != nil {
		return nil
	}

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

	case "ne":
		isValid = value != to
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

func NewCompareNumberValuesAction(s TestStepValue) CompareNumberValuesAction {
	compareNumberValuesAction := CompareNumberValuesAction{
		s.PropertyValues["method"].(string),
		s.PropertyValues["value"].(string),
		s.PropertyValues["to"].(string),
		s,
	}

	return compareNumberValuesAction
}
