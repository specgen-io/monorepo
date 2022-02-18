package ruby

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/spec"
)

func versionedModule(moduleName string, version spec.Name) string {
	if version.Source != "" {
		return fmt.Sprintf("%s::%s", moduleName, version.PascalCase())
	} else {
		return moduleName
	}
}
