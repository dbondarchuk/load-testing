package main

import (
	"time"
)

type SleepAction struct {
	TimeOut int `json:"timeOut"`
}

func (s SleepAction) Execute(resultsChannel chan HttpReqResult, variables map[string]interface{}) error {
	time.Sleep(time.Duration(s.TimeOut) * time.Second)
	return nil
}

func NewSleepAction(a map[string]interface{}) SleepAction {
	return SleepAction{a["timeOut"].(int)}
}
