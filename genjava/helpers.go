package genjava

import (
	"fmt"
	"github.com/specgen-io/spec"
	"strings"
)

func versionedFolder(version spec.Name, folder string) string {
	if version.Source != "" {
		return fmt.Sprintf(`%s_%s`, folder, version.FlatCase())
	}
	return folder
}

func versionedPackage(specification *spec.Spec, version spec.Name, packageName string) string {
	if version.Source != "" {
		return fmt.Sprintf("%s.models.%s_%s", validSpecificationName(specification), packageName, version.FlatCase())
	}
	return fmt.Sprintf("%s.models.%s", validSpecificationName(specification), packageName)
}

func validSpecificationName(specification *spec.Spec) string {
	if strings.Contains(specification.Name.Source, "-") {
		return strings.ReplaceAll(specification.Name.Source, "-", "_")
	}
	return specification.Name.Source
}
