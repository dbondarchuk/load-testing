package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type CompareValuesAction struct {
	Value      string `json:"value"`
	Method     string `json:"method"`
	IgnoreCase bool   `json:"ignoreCase"`
	To         string `json:"to"`
}

// Execute action
func (h CompareValuesAction) Execute(resultsChannel chan HttpReqResult, variables map[string]interface{}) error {
	isValid := false
	value := SubstParams(variables, h.Value)
	to := SubstParams(variables, h.To)

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

func NewCompareValuesAction(a map[string]interface{}) CompareValuesAction {
	var compare = a["to"].(string)

	firstIndex := strings.Index(compare, "|")

	method := compare[0:firstIndex]
	leftOver := compare[firstIndex+1 : len(compare)]

	secondIndex := strings.Index(leftOver, "|")
	ignoreCase, _ := strconv.ParseBool(leftOver[0:secondIndex])

	to := leftOver[secondIndex+1 : len(leftOver)]

	compareValuesAction := CompareValuesAction{
		a["value"].(string),
		method,
		ignoreCase,
		to,
	}

	return compareValuesAction
}
