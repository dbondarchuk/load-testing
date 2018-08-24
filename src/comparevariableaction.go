package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type CompareVariableAction struct {
	Value      string        `json:"value"`
	Method     string        `json:"method"`
	IgnoreCase bool          `json:"ignoreCase"`
	To         string        `json:"to"`
	Step       TestStepValue `json:"-"`
}

func (h CompareVariableAction) GetStep() *TestStepValue {
	return &h.Step
}

// Execute action
func (h CompareVariableAction) Execute(httpResultsChannel chan HttpReqResult, variables map[string]interface{}) error {
	isValid := false
	value, err := SubstParams(variables, "$("+h.Value+")") // convert to variable string
	if err != nil {
		return nil
	}

	to, err := SubstParams(variables, h.To)
	if err != nil {
		return nil
	}

	if h.IgnoreCase {
		value = strings.ToLower(value)
		to = strings.ToLower(to)
	}

	switch method := strings.ToLower(h.Method); method {
	case "equals":
		isValid = value == to
		break

	case "startswith":
		isValid = strings.HasPrefix(value, to)
		break

	case "endswitheq":
		isValid = strings.HasSuffix(value, to)
		break

	case "contains":
		isValid = strings.Contains(value, to)
		break

	case "notequal":
		isValid = value != to
		break

	case "notstartwith":
		isValid = !strings.HasPrefix(value, to)
		break

	case "notendwith":
		isValid = !strings.HasSuffix(value, to)
		break

	case "notcontain":
		isValid = !strings.Contains(value, to)
		break

	case "regex":
		re := regexp.MustCompile(to)

		isValid = re.MatchString(value)
		break
	}

	if !isValid {
		return errors.New("Wrong assertion. Expected: '" + to + "' be " + h.Method + " to '" + value + "'")
	}

	return nil
}

func NewCompareVariableAction(s TestStepValue) CompareVariableAction {
	var compare = s.PropertyValues["to"].(string)

	firstIndex := strings.Index(compare, "|")

	method := compare[0:firstIndex]
	leftOver := compare[firstIndex+1 : len(compare)]

	secondIndex := strings.Index(leftOver, "|")
	ignoreCase, _ := strconv.ParseBool(leftOver[0:secondIndex])

	to := leftOver[secondIndex+1 : len(leftOver)]

	compareVariableAction := CompareVariableAction{
		s.PropertyValues["variableName"].(string),
		method,
		ignoreCase,
		to,
		s,
	}

	return compareVariableAction
}
