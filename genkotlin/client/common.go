package client

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/spec"
)

func serviceResponseInterfaceName(operation *spec.NamedOperation) string {
	return fmt.Sprintf(`%sResponse`, operation.Name.PascalCase())
}

func clientName(api *spec.Api) string {
	return fmt.Sprintf(`%sClient`, api.Name.PascalCase())
}

func addBuilderParam(param *spec.NamedParam) string {
	if param.Type.Definition.IsNullable() {
		return fmt.Sprintf(`%s!!`, param.Name.CamelCase())
	}
	return param.Name.CamelCase()
}
