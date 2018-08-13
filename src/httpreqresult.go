package main

type HttpReqResult struct {
	Type    string
	Latency int64
	Size    int
	Url     string
	Status  int
	Name    string
	When    int64
}
