package main

import (
	"errors"
	"regexp"
	"strings"
)

type CompareValuesAction struct {
	Value      string        `json:"value"`
	Method     string        `json:"method"`
	IgnoreCase bool          `json:"ignoreCase"`
	To         string        `json:"to"`
	Step       TestStepValue `json:"-"`
}

func (h CompareValuesAction) GetStep() *TestStepValue {
	return &h.Step
}

// Execute action
func (h CompareValuesAction) Execute(httpResultsChannel chan HttpReqResult, variables map[string]interface{}) error {
	isValid := false
	value, err := SubstParams(variables, h.Value)
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

	case "endswith":
		isValid = strings.HasSuffix(value, to)
		break

	case "contains":
		isValid = strings.Contains(value, to)
		break

	case "notequals":
		isValid = value != to
		break

	case "notstartswith":
		isValid = !strings.HasPrefix(value, to)
		break

	case "notendswith":
		isValid = !strings.HasSuffix(value, to)
		break

	case "notcontains":
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

func NewCompareValuesAction(s TestStepValue) CompareValuesAction {
	compareValuesAction := CompareValuesAction{
		s.PropertyValues["value"].(string),
		s.PropertyValues["method"].(string),
		s.PropertyValues["ignoreCase"].(bool),
		s.PropertyValues["to"].(string),
		s,
	}

	return compareValuesAction
}
