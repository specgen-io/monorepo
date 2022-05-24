package service

import (
	"fmt"
	imports2 "github.com/specgen-io/specgen/v2/gen/golang/imports"
	"github.com/specgen-io/specgen/v2/gen/golang/module"
	"github.com/specgen-io/specgen/v2/gen/golang/writer"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func generateSpecRouting(specification *spec.Spec, module module.Module) *sources.CodeFile {
	w := writer.NewGoWriter()
	w.Line("package %s", module.Name)

	imports := imports2.Imports()
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
			routesParams = append(routesParams, fmt.Sprintf(`%s %s.%s`, serviceApiNameVersioned(&api), versionedApiImportAlias(&api), serviceInterfaceName))
		}
	}
	w.Line(`func AddRoutes(router *vestigo.Router, %s) {`, strings.Join(routesParams, ", "))
	for _, version := range specification.Versions {
		for _, api := range version.Http.Apis {
			w.Line(`  %s%s`, packageFrom(&version), callAddRouting(&api, serviceApiNameVersioned(&api)))
		}
	}
	w.Line(`}`)

	return &sources.CodeFile{
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

func serviceApiNameVersioned(api *spec.Api) string {
	return fmt.Sprintf(`%sService%s`, api.Name.Source, api.Apis.Version.Version.PascalCase())
}

func packageFrom(version *spec.Version) string {
	if version.Version.Source != "" {
		return fmt.Sprintf(`%s.`, version.Version.FlatCase())
	}
	return ""
}
