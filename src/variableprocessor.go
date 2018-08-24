package main

import (
	"regexp"
	"strings"

	"github.com/robertkrimen/otto"
)

var re = regexp.MustCompile("\\$\\(([a-zA-Z0-9_\\.\\[\\]]{1,})\\)")

func SubstParams(dictionary map[string]interface{}, textData string) (string, error) {
	if strings.ContainsAny(textData, "$(") {
		res := re.FindAllStringSubmatch(textData, -1)
		vm := otto.New()
		for key, value := range dictionary {
			vm.Set(key, value)
		}

		for _, v := range res {
			value, err := vm.Run(v[1])

			if err != nil {
				return "", err
			}

			textData = strings.Replace(textData, "$("+v[1]+")", value.String(), 1)
		}

		return textData, nil
	}

	return textData, nil
}
