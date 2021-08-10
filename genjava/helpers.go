package genjava

import (
	"fmt"
	"github.com/specgen-io/spec"
)

func versionedFolder(version spec.Name, folder string) string {
	if version.Source != "" {
		return fmt.Sprintf(`%s_%s`, folder, version.FlatCase())
	}
	return folder
}

func versionedPackage(specification *spec.Spec, version spec.Name, packageName string) string {
	if version.Source != "" {
		return fmt.Sprintf("%s.models.%s_%s", specification.Name.SnakeCase(), packageName, version.FlatCase())
	}
	return fmt.Sprintf("%s.models.%s", specification.Name.SnakeCase(), packageName)
}
