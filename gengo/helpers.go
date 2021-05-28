package gengo

import (
	"github.com/specgen-io/specgen/v2/gen"
	"strings"
)

func generateHelperFunctions(packageName string, path string) *gen.TextFile {
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
	code, _ = gen.ExecuteTemplate(code, struct { PackageName string } {packageName })
	return &gen.TextFile{path, strings.TrimSpace(code)}
}
