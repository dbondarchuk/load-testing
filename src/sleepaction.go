package main

import (
	"time"
)

type SleepAction struct {
	TimeOut int           `json:"timeOut"`
	Step    TestStepValue `json:"-"`
}

func (s SleepAction) Execute(httpResultsChannel chan HttpReqResult, variables map[string]interface{}) error {
	time.Sleep(time.Duration(s.TimeOut) * time.Second)
	return nil
}

func (s SleepAction) GetStep() *TestStepValue {
	return &s.Step
}

func NewSleepAction(s TestStepValue) SleepAction {
	return SleepAction{
		s.PropertyValues["timeOut"].(int),
		s,
	}
}
