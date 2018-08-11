package main

import (
	"log"
)

func buildActionList(t *TestDef) ([]Action, bool) {
	var valid = true
	actions := make([]Action, len(t.Actions), len(t.Actions))
	for _, element := range t.Actions {
		for key, value := range element {
			var action Action
			actionMap := value.(map[interface{}]interface{})
			switch key {
			case "sleep":
				action = NewSleepAction(actionMap)
				break
			case "http":
				action = NewHttpAction(actionMap)
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
