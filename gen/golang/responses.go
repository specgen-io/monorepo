package golang

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/spec"
)

func NewResponse(response *spec.OperationResponse, body string) string {
	return fmt.Sprintf(`%s{%s: &%s}`, responseTypeName(response.Operation), response.Name.PascalCase(), body)
}
