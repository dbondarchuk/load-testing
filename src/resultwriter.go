package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var httpResultWriter *bufio.Writer
var resultWriter *bufio.Writer

var httpResultFile *os.File
var resultFile *os.File
var err error

var results []*Result
var httpResults []*StatFrame

var opened bool = false

func OpenResultsFiles(resultFileName string, httpResultFileName string) {
	if !opened {
		opened = true
	} else {
		return
	}

	openFile(resultFileName, &resultFile, &resultWriter)
	openFile(httpResultFileName, &httpResultFile, &httpResultWriter)
}

func CloseResultsFiles() {
	if opened {
		closeFile(resultFile, resultWriter)

		closeFile(httpResultFile, httpResultWriter)
	}

	// Do nothing if not opened
}

func writeHttpResult(statFrame *StatFrame) {
	byStatus, _ := json.Marshal(statFrame.ByStatus)
	str := fmt.Sprintf("%d|%d|%d|%s", statFrame.Time, statFrame.Reqs, statFrame.Latency, byStatus)

	httpResultWriter.WriteString(str + "\n")
	httpResultWriter.Flush()
}

func writeResult(results []*Result) {
	for _, result := range results {
		errors, _ := json.Marshal(result.ErrorList)
		str := fmt.Sprintf("%d|%d|%t|%d|%s", result.User, result.Loop, result.IsSuccessful, result.Time, filterNewLines(string(errors)))
		resultWriter.WriteString(str + "\n")
	}

	resultWriter.Flush()
}

func openFile(fileName string, file **os.File, writer **bufio.Writer) {
	*file, err = os.Create(fileName)
	if err != nil {
		os.Mkdir("results", 0777)
		os.Mkdir("results/log", 0777)
		*file, err = os.Create(fileName)
		if err != nil {
			panic(err)
		}
	}

	*writer = bufio.NewWriter(*file)
}

func closeFile(file *os.File, writer *bufio.Writer) {
	writer.Flush()
	file.Close()
}

func filterNewLines(s string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case 0x000A, 0x000B, 0x000C, 0x000D, 0x0085, 0x2028, 0x2029:
			return -1
		default:
			return r
		}
	}, s)
}
