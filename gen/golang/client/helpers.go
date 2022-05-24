package client

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func logFieldsName(operation *spec.NamedOperation) string {
	return fmt.Sprintf("log%s", operation.Name.PascalCase())
}

func JoinParams(params []string) string {
	return strings.Join(params, ", ")
}
