package genjava

import (
	"fmt"
	"github.com/specgen-io/spec"
)

func getterName(field *spec.NamedDefinition) string {
	return fmt.Sprintf(`get%s`, field.Name.PascalCase())
}

func setterName(field *spec.NamedDefinition) string {
	return fmt.Sprintf(`set%s`, field.Name.PascalCase())
}

func serviceName(api *spec.Api) string {
	return fmt.Sprintf(`%sService`, api.Name.PascalCase())
}

func serviceImplName(api *spec.Api) string {
	return fmt.Sprintf(`%sServiceImpl`, api.Name.PascalCase())
}

func versionServiceName(serviceName string, version *spec.Version) string {
	return fmt.Sprintf(`%s%s`, serviceName, version.Version.PascalCase())
}

func serviceVarName(api *spec.Api) string {
	return fmt.Sprintf(`%sService`, api.Name.CamelCase())
}

func serviceInterfaceName(api *spec.Api) string {
	return fmt.Sprintf(`%sService`, api.Name.PascalCase())
}

func serviceResponseInterfaceName(operation spec.NamedOperation) string {
	return fmt.Sprintf(`%sResponse`, operation.Name.PascalCase())
}

func serviceResponseImplName(operation spec.NamedOperation, response spec.NamedResponse) string {
	return fmt.Sprintf(`%s%s`, serviceResponseInterfaceName(operation), response.Name.PascalCase())
}

func controllerName(api *spec.Api) string {
	return fmt.Sprintf(`%sController`, api.Name.PascalCase())
}

func versionControllerName(controllerName string, version *spec.Version) string {
	return fmt.Sprintf(`%s%s`, controllerName, version.Version.PascalCase())
}

func controllerMethodName(operation spec.NamedOperation) string {
	return fmt.Sprintf(`%sController`, operation.Name.CamelCase())
}

func versionUrl(version *spec.Version, url string) string {
	if version.Version.Source != "" {
		return fmt.Sprintf(`/%s%s`, version.Version.FlatCase(), url)
	}
	return fmt.Sprintf(`%s`, url)
}
