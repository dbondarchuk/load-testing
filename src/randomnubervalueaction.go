package main

import (
	"math/rand"
	"strconv"
	"time"
)

type RandomNumberValueAction struct {
	Min          string        `json:"min"`
	Max          string        `json:"max"`
	VariableName string        `json:"variableName"`
	Step         TestStepValue `json:"-"`
}

func (r RandomNumberValueAction) Execute(httpResultsChannel chan HttpReqResult, variables map[string]interface{}) error {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	minStr, err := SubstParams(variables, r.Min)
	if err != nil {
		return err
	}

	maxStr, err := SubstParams(variables, r.Max)
	if err != nil {
		return err
	}

	min, err := strconv.Atoi(minStr)
	if err != nil {
		return err
	}

	max, err := strconv.Atoi(maxStr)
	if err != nil {
		return err
	}

	num := r1.Intn(max+1-min) + min

	variables[r.VariableName] = num

	return nil
}

func (r RandomNumberValueAction) GetStep() *TestStepValue {
	return &r.Step
}

func NewRandomNumberValueAction(s TestStepValue) RandomNumberValueAction {
	return RandomNumberValueAction{
		s.PropertyValues["min"].(string),
		s.PropertyValues["max"].(string),
		s.PropertyValues["variableName"].(string),
		s,
	}
}
