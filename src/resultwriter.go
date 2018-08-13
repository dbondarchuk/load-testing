package main

import (
	"bufio"
	"encoding/json"
	"os"
)

var w *bufio.Writer
var f *os.File
var err error

var opened bool = false

func OpenResultsFile(fileName string) {
	if !opened {
		opened = true
	} else {
		return
	}
	f, err = os.Create(fileName)
	if err != nil {
		os.Mkdir("results", 0777)
		os.Mkdir("results/log", 0777)
		f, err = os.Create(fileName)
		if err != nil {
			panic(err)
		}
	}
	w = bufio.NewWriter(f)
	_, err = w.WriteString(string("["))
}

func CloseResultsFile() {
	if opened {
		_, err = w.WriteString(string("]"))
		w.Flush()
		f.Close()
	}
	// Do nothing if not opened
}

func writeResult(httpResult *HttpReqResult) {
	jsonString, err := json.Marshal(httpResult)
	if err != nil {
		panic(err)
	}
	_, err = w.WriteString(string(jsonString))
	_, err = w.WriteString(",")

	if err != nil {
		panic(err)
	}
	w.Flush()

}
