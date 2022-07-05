package client

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func clientName(api *spec.Api) string {
	return fmt.Sprintf(`%sClient`, api.Name.PascalCase())
}

func addBuilderParam(param *spec.NamedParam) string {
	if param.Type.Definition.IsNullable() {
		return fmt.Sprintf(`%s!!`, param.Name.CamelCase())
	}
	return param.Name.CamelCase()
}

func joinParams(params []string) string {
	return strings.Join(params, ", ")
}
