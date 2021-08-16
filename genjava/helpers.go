package genjava

import (
	"fmt"
	"github.com/specgen-io/spec"
)

func versionedPackage(version spec.Name, packageName string) string {
	if version.Source != "" {
		return fmt.Sprintf(`%s_%s`, packageName, version.FlatCase())
	}
	return packageName
}
