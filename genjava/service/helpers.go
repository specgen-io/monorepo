package service

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/spec"
)

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

func controllerName(api *spec.Api) string {
	return fmt.Sprintf(`%sController`, api.Name.PascalCase())
}

func versionControllerName(controllerName string, version *spec.Version) string {
	return fmt.Sprintf(`%s%s`, controllerName, version.Version.PascalCase())
}

func controllerMethodName(operation *spec.NamedOperation) string {
	return fmt.Sprintf(`%sController`, operation.Name.CamelCase())
}
