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
