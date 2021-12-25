package gengo

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/spec"
)

func generateSpecRouting(specification *spec.Spec, module module) *gen.TextFile {
	w := NewGoWriter()
	w.Line("package %s", module.Name)

	imports := Imports()
	imports.Add("github.com/husobee/vestigo")
	for _, version := range specification.Versions {
		versionModule := module.Submodule(version.Version.FlatCase())
		if version.Version.Source != "" {
			imports.Add(versionModule.Package)
		}
		for _, api := range version.Http.Apis {
			apiModule := versionModule.Submodule(api.Name.SnakeCase())
			imports.AddAlias(apiModule.Package, versionedApiImportAlias(&api))
		}
	}
	imports.Write(w)

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
		Path:    module.GetPath("spec_routing.go"),
		Content: w.String(),
	}
}

func versionedApiImportAlias(api *spec.Api) string {
	version := api.Apis.Version.Version
	if version.Source != "" {
		return api.Name.CamelCase() + version.PascalCase()
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
