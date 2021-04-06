package main

import (
	"encoding/json"
	"flag"
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
	inputFilePtr := flag.String("i", "input.json", "Input file with test")
	idPtr := flag.String("id", "", "Id of the test")
	outputDirPtr := flag.String("o", ".", "Output dir for the logs")
	tolerancePercentage := flag.Float64("t", 0, "Maximum percentage of tolerated failed steps")

	flag.Parse()

	if *idPtr == "" {
		fmt.Printf("Please specify ID of the test")
		os.Exit(1)
	}

	// defer profile.Start(profile.CPUProfile).Stop()

	// Start the web socket server, will not block exit until forced

	SimulationStart = time.Now()
	dat, inputError := ioutil.ReadFile(*inputFilePtr)
	if inputError != nil {
		fmt.Print(inputError.Error())
		os.Exit(2)
	}

	var t TestDef
	json.Unmarshal([]byte(dat), &t)

	t.Id = *idPtr

	if t.Variables == nil {
		t.Variables = make(map[string]interface{})
	}

	if !ValidateTestDefinition(&t) {
		os.Exit(3)
	}

	actions, isValid := buildActionList(&t)
	if !isValid {
		os.Exit(4)
	}

	OpenResultsFiles(*outputDirPtr, t.Id+".result.log", t.Id+".http.log")
	var failedLoops = spawnUsers(&t, actions)

	FlushResults()

	elapsed := time.Since(SimulationStart)

	time.Sleep(time.Duration(time.Second))

	totalLoops := t.Users * t.Iterations

	fmt.Printf("Total tests: %d. Failed tests: %d\n", totalLoops, failedLoops)
	fmt.Printf("Done in %v\n", elapsed)

	CloseResultsFiles()

	failPercentage := float64(failedLoops) / float64(totalLoops) * 100.0
	fmt.Printf("Tolerance percentage: %.3f. Fail percentage: %.3f\n", *tolerancePercentage, failPercentage)
	if failPercentage > *tolerancePercentage {
		fmt.Println("Test has failed")
		os.Exit(-1)
	}
}

func spawnUsers(t *TestDef, actions []Action) int {
	httpResultsChannel := make(chan HttpReqResult, 100000) // buffer?
	resultsChannel := make(chan Result, 100000)            // buffer?
	go aggregateResultPerSecondHandler(resultsChannel)
	go aggregateHttpResultPerSecondHandler(httpResultsChannel)
	wg := sync.WaitGroup{}

	totalFailedLoops := 0

	for i := 0; i < t.Users; i++ {
		wg.Add(1)
		UID := strconv.Itoa(rand.Intn(t.Users+1) + 100000)
		go launchActions(t, httpResultsChannel, resultsChannel, &wg, actions, UID, i, &totalFailedLoops)
		var waitDuration float32 = float32(t.Rampup) / float32(t.Users)
		time.Sleep(time.Duration(int(1000*waitDuration)) * time.Millisecond)
	}
	fmt.Println("All users started, waiting at WaitGroup")
	wg.Wait()

	return totalFailedLoops
}

func launchActions(t *TestDef, httpResultsChannel chan HttpReqResult, resultsChannel chan Result, wg *sync.WaitGroup, actions []Action, UID string, user int, totalFailedLoops *int) {
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
					errs["#"+strconv.Itoa(k+1)+" '"+step.Name+"'"] = err
					if !step.IgnoreError {
						isError = true
					}
				}
			}
		}

		errorlist := make([]string, len(errs))
		k := 0

		for name, err := range errs {
			errorlist[k], k = "Step '"+name+"' has failed with error: '"+err.Error()+"'", k+1
		}

		if k > 0 {
			(*totalFailedLoops)++
		}

		var totalTime = time.Since(startTime).Nanoseconds() / 1000

		result := Result{
			user + 1,
			i + 1,
			!isError,
			errorlist,
			totalTime,
		}

		resultsChannel <- result
	}

	wg.Done()
}

func resetVariablesAndUID(original map[string]interface{}, UID string, variables map[string]interface{}) {
	b, _ := json.Marshal(original)
	e := json.Unmarshal(b, &variables)

	if e == nil {
		variables["UID"] = UID
	}
}
