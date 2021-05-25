package genruby

import (
	"fmt"
	spec "github.com/specgen-io/spec.v2"
)

func versionedModule(moduleName string, version spec.Name) string {
	if version.Source != "" {
		return fmt.Sprintf("%s::%s", moduleName, version.PascalCase())
	} else {
		return moduleName
	}
}
