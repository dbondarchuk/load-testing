package main

import (
	"time"
)

type SleepAction struct {
	TimeOut int `json:"timeOut"`
}

func (s SleepAction) Execute(resultsChannel chan HttpReqResult, variables map[string]interface{}) {
	time.Sleep(time.Duration(s.TimeOut) * time.Second)
}

func NewSleepAction(a map[string]interface{}) SleepAction {
	return SleepAction{a["timeOut"].(int)}
}
