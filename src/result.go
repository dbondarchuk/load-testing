package main

type Result struct {
	User         int      `json:"user"`
	Loop         int      `json:"loop"`
	IsSuccessful bool     `json:"isSuccessful"`
	ErrorList    []string `json:"errorList"`
	Time         int64    `json:"time"`
}
