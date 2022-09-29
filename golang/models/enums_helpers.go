package models

import (
	"generator"
	"golang/module"
	"strings"
)

func GenerateEnumsHelperFunctions(module module.Module) *generator.CodeFile {
	code := `
package [[.PackageName]]

import (
	"encoding/json"
	"errors"
	"fmt"
)

func contains(lookFor string, arr []string) bool {
	for _, value := range arr {
		if lookFor == value {
			return true
		}
	}
	return false
}

func ReadStringValue(b []byte, values []string) (string, error) {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return "", err
	}
	if !contains(str, values) {
		return "", errors.New(fmt.Sprintf("Unknown enum value: %s", str))
	}
	return str, nil
}
`
	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{module.Name})
	return &generator.CodeFile{module.GetPath("helpers.go"), strings.TrimSpace(code)}
}
