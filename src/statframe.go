package main

type StatFrame struct {
	Time     int64                  `json:"time"`
	Latency  int                    `json:"latency"`
	Reqs     int                    `json:"reqs"`
	ByStatus map[string]map[int]int `json:"byStatus"`
}
