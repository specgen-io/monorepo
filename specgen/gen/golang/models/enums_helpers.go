package models

import (
	"strings"

	"github.com/specgen-io/specgen/generator/v2"
	"github.com/specgen-io/specgen/v2/gen/golang/module"
)

func generateEnumsHelperFunctions(module module.Module) *generator.CodeFile {
	code := `
package [[.PackageName]]

import (
	"encoding/json"
	"errors"
	"fmt"
)

func contains(lookFor string, arr []string) bool {
	for _, value := range arr{
		if lookFor == value {
			return true
		}
	}
	return false
}

func readEnumStringValue(b []byte, values []string) (string, error) {
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
	return &generator.CodeFile{module.GetPath("enums_helpers.go"), strings.TrimSpace(code)}
}
