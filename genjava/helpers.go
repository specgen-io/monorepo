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

func versionedPackage(version spec.Name, packageName string) string {
	if version.Source != "" {
		return fmt.Sprintf("%s_%s", packageName, version.FlatCase())
	}
	return packageName
}
