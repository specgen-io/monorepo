package gengo

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/spec"
	"path/filepath"
	"strings"
)

func serviceInterfaceTypeVar(api *spec.Api) string {
	return fmt.Sprintf(`%sService`, api.Name.Source)
}

func clientTypeName() string {
	return `Client`
}

func responseTypeName(operation *spec.NamedOperation) string {
	return fmt.Sprintf(`%sResponse`, operation.Name.PascalCase())
}

func createPackageName(args ...string) string {
	parts := []string{}
	for _, arg := range args {
		arg = strings.TrimPrefix(arg, "./")
		if arg != "" {
			parts = append(parts, arg)
		}
	}
	return strings.Join(parts, "/")
}

func createPath(args ...string) string {
	parts := []string{}
	for _, arg := range args {
		if arg != "" {
			parts = append(parts, arg)
		}
	}
	return filepath.Join(parts...)
}

func getShortPackageName(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}

func responsesSignature(operation *spec.NamedOperation) string {
	if len(operation.Responses) == 1 {
		response := operation.Responses[0]
		if response.Type.Definition.IsEmpty() {
			return fmt.Sprintf(`error`)
		} else {
			return fmt.Sprintf(`(*%s, error)`, GoType(&response.Type.Definition))
		}
	}
	return fmt.Sprintf(`(*%s, error)`, responseTypeName(operation))
}
