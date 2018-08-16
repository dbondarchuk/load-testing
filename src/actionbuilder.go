package main

import (
	"log"
	"strings"
)

func buildActionList(t *TestDef) ([]Action, bool) {
	var valid = true
	actions := make([]Action, len(t.Actions), len(t.Actions))
	for _, element := range t.Actions {
		for key, value := range element {
			var action Action
			actionMap := value
			switch strings.ToLower(key) {
			case "sleep":
				action = NewSleepAction(actionMap)
				break
			case "http":
				action = NewHttpAction(actionMap)
				break
			case "comparenumbervalues":
				action = NewCompareNumberValuesAction(actionMap)
				break
			case "comparevalues":
				action = NewCompareValuesAction(actionMap)
				break
			case "comparevariable":
				action = NewCompareVariableAction(actionMap)
				break
			default:
				valid = false
				log.Fatal("Unknown action type encountered: " + key)
				break
			}
			if valid {
				actions = append(actions, action)
			}
		}
	}
	return actions, valid
}
