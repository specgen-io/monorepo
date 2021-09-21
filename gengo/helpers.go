package gengo

import (
	"fmt"
	"github.com/specgen-io/spec"
	"path/filepath"
	"strings"
)

func serviceInterfaceTypeVar(api *spec.Api) string {
	return fmt.Sprintf(`%sService`, api.Name.Source)
}

func clientTypeName(api *spec.Api) string {
	return fmt.Sprintf(`%sClient`, api.Name.SnakeCase())
}

func responseTypeName(operation *spec.NamedOperation) string {
	return fmt.Sprintf(`%sResponse`, operation.Name.PascalCase())
}

func versionedFolder(version spec.Name, folder string) string {
	if version.Source != "" {
		return fmt.Sprintf(`%s/%s`, folder, version.FlatCase())
	}
	return folder
}

func versionedPackage(version spec.Name, packageName string) string {
	if version.Source != "" {
		return version.FlatCase()
	}
	return packageName
}

func apiPackage(root string, packageName string, api *spec.Api) string {
	return fmt.Sprintf(`"%s/%s/%s"`, root, packageName, api.Name.SnakeCase())
}

func createPackageName(args ...string) string {
	parts := []string{}
	for _, arg := range args {
		if arg != "" {
			parts = append(parts, arg)
		}
	}
	return fmt.Sprintf(`"%s"`, strings.Join(parts, "/"))
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