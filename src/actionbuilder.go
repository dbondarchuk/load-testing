package main

import (
	"log"
	"strings"
)

func buildActionList(t *TestDef) ([]Action, bool) {
	var valid = true
	actions := make([]Action, 0, len(t.Steps))
	for _, step := range t.Steps {
		if !step.Enabled {
			continue
		}

		var action Action
		switch strings.ToLower(step.TypeName) {
		case "sleep":
			action = NewSleepAction(step)
			break
		case "restrequest":
			action = NewHttpAction(step)
			break
		case "comparenumbervalues":
			action = NewCompareNumberValuesAction(step)
			break
		case "comparevalues":
			action = NewCompareValuesAction(step)
			break
		case "comparevariable":
			action = NewCompareVariableAction(step)
			break
		case "randomnumbervalue":
			action = NewRandomNumberValueAction(step)
			break
		case "randomstringvalue":
			action = NewRandomStringValueAction(step)
			break
		default:
			valid = false
			log.Fatal("Unknown action type encountered: " + step.TypeName)
			break
		}
		if valid {
			actions = append(actions, action)
		}
	}
	return actions, valid
}
