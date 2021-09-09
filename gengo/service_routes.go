package gengo

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
)

func generateRoutes(moduleName string, specification *spec.Spec, packageName string, generatePath string) *gen.TextFile {
	w := NewGoWriter()
	w.Line("package %s", packageName)
	w.EmptyLine()
	w.Line("import (")
	w.Line(`  "github.com/husobee/vestigo"`)
	for _, version := range specification.Versions {
		if version.Version.Source != "" {
			w.Line(`  "%s/%s"`, moduleName, versionedFolder(version.Version, packageName))
		}
	}
	w.Line(`)`)
	w.EmptyLine()
	w.Line(`func AddRoutes(router *vestigo.Router, %s) {`, JoinDelimParams(addRoutesParams(specification, packageName)))
	for _, version := range specification.Versions {
		for _, api := range version.Http.Apis {
			w.Line(`  %sAdd%sRoutes(router, %s)`, addVersionedPackage(version.Version, packageName), api.Name.PascalCase(), addVersionedInterfaceVar(&api, version.Version))
		}
	}
	w.Line(`}`)

	return &gen.TextFile{
		Path:    filepath.Join(generatePath, "routes.go"),
		Content: w.String(),
	}
}

func addRoutesParams(specification *spec.Spec, packageName string) []string {
	routesParams := []string{}
	for _, version := range specification.Versions {
		for _, api := range version.Http.Apis {
			routesParams = append(routesParams, fmt.Sprintf(`%s %s`, addVersionedInterfaceVar(&api, version.Version), addVersionedPackage(version.Version, packageName)+serviceInterfaceTypeName(&api)))
		}
	}
	return routesParams
}
