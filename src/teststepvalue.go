package main

type TestStepValue struct {
	TypeName       string                 `json:"typeName"`
	Condition      string                 `json:"condition"`
	Name           string                 `json:"name"`
	Enabled        bool                   `json:"enabled"`
	IgnoreError    bool                   `json:"ignoreError"`
	RunOnFailure   bool                   `json:"runOnFailure"`
	PropertyValues map[string]interface{} `json:"propertyValues"`
}
