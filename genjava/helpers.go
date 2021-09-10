package genjava

import (
	"fmt"
	"github.com/specgen-io/spec"
)

func versionedPackage(version spec.Name, packageName string) string {
	if version.Source != "" {
		return fmt.Sprintf(`%s.%s`, packageName, version.FlatCase())
	}
	return packageName
}

func versionedPath(version spec.Name, generatePath string) string {
	if version.Source != "" {
		return fmt.Sprintf(`%s/%s`, generatePath, version.FlatCase())
	}
	return generatePath
}

func getterName(field *spec.NamedDefinition) string {
	return fmt.Sprintf(`get%s`, field.Name.PascalCase())
}

func setterName(field *spec.NamedDefinition) string {
	return fmt.Sprintf(`set%s`, field.Name.PascalCase())
}
