package main

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var upperCharacters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
var lowerCharacters = strings.ToLower(upperCharacters)

var digitCharacters = "0123456789"
var specialCharacters = "!@#$%^&*()-_+=[]\\/,.<>;:'\""

type RandomStringValueAction struct {
	MinLength          string        `json:"minLength"`
	MaxLength          string        `json:"maxLength"`
	SpecialCharacters  bool          `json:"specialCharacters"`
	ExcludedCharacters string        `json:"excludedCharacters"`
	VariableName       string        `json:"variableName"`
	Step               TestStepValue `json:"-"`
}

func (r RandomStringValueAction) Execute(httpResultsChannel chan HttpReqResult, variables map[string]interface{}) error {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	symbols := upperCharacters + lowerCharacters + digitCharacters
	if r.SpecialCharacters {
		symbols = symbols + specialCharacters
	}

	if len(r.ExcludedCharacters) > 0 {
		if strings.Contains(r.ExcludedCharacters, ",,") {
			symbols = strings.Replace(symbols, ",", "", -1)
		}

		excludedSymbolsList := strings.Split(r.ExcludedCharacters, ",")
		for _, symbol := range excludedSymbolsList {
			if len(symbol) == 1 {
				symbols = strings.Replace(symbols, symbol, "", -1)
			}
		}
	}

	minLengthStr, err := SubstParams(variables, r.MinLength)
	if err != nil {
		return err
	}

	maxLengthStr, err := SubstParams(variables, r.MaxLength)
	if err != nil {
		return err
	}

	minLength, err := strconv.Atoi(minLengthStr)
	if err != nil {
		return err
	}

	maxLength, err := strconv.Atoi(maxLengthStr)
	if err != nil {
		return err
	}

	length := r1.Intn(maxLength+1-minLength) + minLength

	result := ""
	for index := 0; index < length; index++ {
		result = result + string(symbols[r1.Intn(len(symbols))])
	}

	variables[r.VariableName] = result

	return nil
}

func (r RandomStringValueAction) GetStep() *TestStepValue {
	return &r.Step
}

func NewRandomStringValueAction(s TestStepValue) RandomStringValueAction {
	return RandomStringValueAction{
		s.PropertyValues["minLength"].(string),
		s.PropertyValues["maxLength"].(string),
		s.PropertyValues["specialCharacters"].(bool),
		s.PropertyValues["excludedCharacters"].(string),
		s.PropertyValues["variableName"].(string),
		s,
	}
}
