package gengo

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
)

func generateRoutes(moduleName string, specification *spec.Spec, rootPackage string, modulePath string) *gen.TextFile {
	w := NewGoWriter()
	w.Line("package %s", getShortPackageName(modulePath))
	w.EmptyLine()

	imports := []string{
		`"github.com/husobee/vestigo"`,
	}
	for _, version := range specification.Versions {
		versionModule := versionedFolder(version.Version, modulePath)
		if version.Version.Source != "" {
			imports = append(imports, fmt.Sprintf(`%s`, createPackageName(moduleName, modulePath, version.Version.Source)))
		}
		for _, api := range version.Http.Apis {
			imports = append(imports, fmt.Sprintf(`%s %s`, versionedApiImportAlias(&api), apiPackage(rootPackage, versionModule, &api)))
		}
	}

	w.Line(`import (`)
	for _, imp := range imports {
		w.Line(`  %s`, imp)
	}
	w.Line(`)`)
	w.EmptyLine()
	routesParams := []string{}
	for _, version := range specification.Versions {
		for _, api := range version.Http.Apis {
			routesParams = append(routesParams, fmt.Sprintf(`%s %s`, addVersionedInterfaceParam(&api), versionedApiInterfaceType(&api)))
		}
	}
	w.Line(`func AddRoutes(router *vestigo.Router, %s) {`, JoinDelimParams(routesParams))
	for _, version := range specification.Versions {
		for _, api := range version.Http.Apis {
			//TODO: Add%sRoutes should be abstracted behind some function
			w.Line(`  %sAdd%sRoutes(router, %s)`, packageFrom(&version), api.Name.PascalCase(), addVersionedInterfaceParam(&api))
		}
	}
	w.Line(`}`)

	return &gen.TextFile{
		Path:    filepath.Join(modulePath, "routes.go"),
		Content: w.String(),
	}
}

func versionedApiImportAlias(api *spec.Api) string {
	version := api.Apis.Version.Version
	if version.Source != "" {
		return api.Name.CamelCase()+version.PascalCase()
	}
	return api.Name.CamelCase()
}

func versionedApiInterfaceType(api *spec.Api) string {
	return fmt.Sprintf("%s.Service", versionedApiImportAlias(api))
}

func addVersionedInterfaceParam(api *spec.Api) string {
	return serviceInterfaceTypeVar(api) + api.Apis.Version.Version.PascalCase()
}

func packageFrom(version *spec.Version) string {
	if version.Version.Source != "" {
		return fmt.Sprintf(`%s.`, version.Version.FlatCase())
	}
	return ""
}