package models

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func getterName(field *spec.NamedDefinition) string {
	return fmt.Sprintf(`get%s`, field.Name.PascalCase())
}

func setterName(field *spec.NamedDefinition) string {
	return fmt.Sprintf(`set%s`, field.Name.PascalCase())
}

func JoinParams(params []string) string {
	return strings.Join(params, ", ")
}
