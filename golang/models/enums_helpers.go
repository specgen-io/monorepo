package models

import (
	"generator"
	"strings"
)

func (g *EncodingJsonGenerator) GenerateEnumsHelperFunctions() *generator.CodeFile {
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
	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{g.Modules.Enums.Name})
	return &generator.CodeFile{g.Modules.Enums.GetPath("helpers.go"), strings.TrimSpace(code)}
}
