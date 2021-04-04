package main

import (
	"fmt"
	"time"
)

var flushResult = false
var flushHttpResult = false
var lastTime int64 = 0

var totalLoopsDone int = 0

func FlushResults() {
	flushResult = true
	flushHttpResult = true
}

func aggregateResultPerSecondHandler(resultChannel chan Result) {
	for {
		var results []*Result

		until := time.Now().UnixNano() + 1000000000

		for time.Now().UnixNano() < until {
			select {
			case msg := <-resultChannel:
				results = append(results, &msg)

			default:
				if flushResult {
					flushResult = false
					break
				} else {
					// Can be trouble. Uses too much CPU if low, limits throughput if too high
					time.Sleep(100 * time.Microsecond)
				}
			}
		}

		// concurrently assemble the result and send it off to the websocket.
		go assembleAndSendResult(results)
	}
}

func aggregateHttpResultPerSecondHandler(httpResultPerSecondChannel chan HttpReqResult) {
	for {
		var totalReq int = 0
		var totalLatency int = 0

		var byStatus = make(map[string]map[int]int)

		until := time.Now().UnixNano() + 1000000000

		for time.Now().UnixNano() < until {
			select {
			case httpMsg := <-httpResultPerSecondChannel:
				if byStatus[httpMsg.Url] == nil {
					byStatus[httpMsg.Url] = make(map[int]int)
				}

				byStatus[httpMsg.Url][httpMsg.Status]++
				totalReq++
				totalLatency += int(httpMsg.Latency / 1000) // measure in microseconds

			default:
				if flushHttpResult {
					flushHttpResult = false
					break
				} else {
					// Can be trouble. Uses too much CPU if low, limits throughput if too high
					time.Sleep(100 * time.Microsecond)
				}
			}
		}

		// concurrently assemble the result and send it off to the websocket.
		go assembleAndSendHttpResult(totalReq, totalLatency, byStatus)
	}
}

func assembleAndSendResult(results []*Result) {
	timePassed := time.Since(SimulationStart).Nanoseconds() / 1000000000
	if timePassed == lastTime {
		timePassed++
	}

	totalLoopsDone += len(results)

	time.Sleep(50 * time.Microsecond)
	fmt.Printf("Second: %d. Loops done: %d. Total loops done: %d\n", timePassed, len(results), totalLoopsDone)
	writeResult(results)
}

func assembleAndSendHttpResult(totalReq int, totalLatency int, byStatus map[string]map[int]int) {
	avgLatency := 0
	if totalReq > 0 {
		avgLatency = totalLatency / totalReq
	}

	time := time.Since(SimulationStart).Nanoseconds() / 1000000000
	if time == lastTime {
		time++
	}

	lastTime = time

	statFrame := StatFrame{
		time,       // seconds
		avgLatency, // microseconds
		totalReq,
		byStatus,
	}
	fmt.Printf("Second: %d. Avg latency: %d μs (%d ms). req/s: %d\n", statFrame.Time, statFrame.Latency, statFrame.Latency/1000, statFrame.Reqs)
	writeHttpResult(&statFrame)
}
