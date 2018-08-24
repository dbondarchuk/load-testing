package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

var SimulationStart time.Time

func main() {

	spec := parseSpecFile()

	// defer profile.Start(profile.CPUProfile).Stop()

	// Start the web socket server, will not block exit until forced

	SimulationStart = time.Now()
	dir, _ := os.Getwd()
	dat, _ := ioutil.ReadFile(dir + "/" + spec)

	var t TestDef
	json.Unmarshal([]byte(dat), &t)

	if !ValidateTestDefinition(&t) {
		return
	}

	actions, isValid := buildActionList(&t)
	if !isValid {
		return
	}

	OpenResultsFiles(dir+"/"+t.Id+".result.log", dir+"/"+t.Id+".http.log")
	spawnUsers(&t, actions)

	FlushResults()

	elapsed := time.Since(SimulationStart)

	time.Sleep(time.Duration(time.Second))

	fmt.Printf("Done in %v\n", elapsed)

	CloseResultsFiles()
}

func parseSpecFile() string {
	if len(os.Args) == 1 {
		fmt.Errorf("No command line arguments, exiting...\n")
		panic("Cannot start simulation, no JSON simulaton specification supplied as command-line argument")
	}
	var s, sep string
	for i := 1; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}
	if s == "" {
		panic(fmt.Sprintf("Specified simulation file '%s' is not a .json file", s))
	}
	return s
}

func spawnUsers(t *TestDef, actions []Action) {
	httpResultsChannel := make(chan HttpReqResult, 100000) // buffer?
	resultsChannel := make(chan Result, 100000)            // buffer?
	go aggregateResultPerSecondHandler(resultsChannel)
	go aggregateHttpResultPerSecondHandler(httpResultsChannel)
	wg := sync.WaitGroup{}
	for i := 0; i < t.Users; i++ {
		wg.Add(1)
		UID := strconv.Itoa(rand.Intn(t.Users+1) + 100000)
		go launchActions(t, httpResultsChannel, resultsChannel, &wg, actions, UID, i)
		var waitDuration float32 = float32(t.Rampup) / float32(t.Users)
		time.Sleep(time.Duration(int(1000*waitDuration)) * time.Millisecond)
	}
	fmt.Println("All users started, waiting at WaitGroup")
	wg.Wait()
}

func launchActions(t *TestDef, httpResultsChannel chan HttpReqResult, resultsChannel chan Result, wg *sync.WaitGroup, actions []Action, UID string, user int) {
	var variables = make(map[string]interface{})

	for i := 0; i < t.Iterations; i++ {
		var startTime = time.Now()

		// Make sure the variables is cleared before each iteration - except for the UID which stays
		resetVariablesAndUID(t.Variables, UID, variables)

		var errs = make(map[string]error)
		isError := false
		// Iterate over the actions. Note the use of the command-pattern like Execute method on the Action interface
		for k, action := range actions {
			step := action.GetStep()
			if action != nil && (!isError || step.RunOnFailure) {
				err := action.(Action).Execute(httpResultsChannel, variables)
				if err != nil {
					errs["#"+strconv.Itoa(k)+" '"+step.Name+"'"] = err
					if !step.IgnoreError {
						isError = true
					}
				}
			}
		}

		var errorlist = make([]string, len(errs))
		for name, err := range errs {
			errorlist = append(errorlist, "Step '"+name+"' has failed with error: '"+err.Error()+"'")
		}

		var totalTime = time.Since(startTime).Nanoseconds() / 1000

		result := Result{
			user + 1,
			i + 1,
			isError,
			errorlist,
			totalTime,
		}

		resultsChannel <- result
	}

	wg.Done()
}

func resetVariablesAndUID(original map[string]interface{}, UID string, variables map[string]interface{}) {
	b, e := json.Marshal(original)
	e = json.Unmarshal(b, &variables)

	if e == nil {
		variables["UID"] = UID
	}
}
