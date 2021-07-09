package gengo

import (
	"fmt"
	"github.com/specgen-io/spec"
)

func serviceInterfaceTypeName(api *spec.Api) string {
	return fmt.Sprintf(`I%sService`, api.Name.PascalCase())
}

func serviceInterfaceTypeVar(api *spec.Api) string {
	return fmt.Sprintf(`%sService`, api.Name.Source)
}

func serviceTypeName(api *spec.Api) string {
	return fmt.Sprintf(`Service%s`, api.Name.PascalCase())
}

func clientTypeName(api *spec.Api) string {
	return fmt.Sprintf(`%sClient`, api.Name.SnakeCase())
}

func responseTypeName(operation *spec.NamedOperation) string {
	return fmt.Sprintf(`%sResponse`, operation.Name.PascalCase())
}

func versionedFolder(version spec.Name, folder string) string {
	if version.Source != "" {
		return fmt.Sprintf(`%s_%s`, folder, version.FlatCase())
	}
	return folder
}

func versionedPackage(version spec.Name, packageName string) string {
	if version.Source != "" {
		return fmt.Sprintf("%s_%s", packageName, version.FlatCase())
	}
	return packageName
}
