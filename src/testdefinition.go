package main

// TestDef - definition of test
type TestDef struct {
	Id         string                              `json:"id"`
	Iterations int                                 `json:"loopCount"`
	Users      int                                 `json:"usersCount"`
	Rampup     int                                 `json:"rampup"`
	Variables  map[string]interface{}              `json:"variables"`
	Actions    []map[string]map[string]interface{} `json:"actions"`
}
