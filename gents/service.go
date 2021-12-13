package gents

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/genopenapi"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateService(specification *spec.Spec, swaggerPath string, generatePath string, servicesPath string, server string, validation string) error {
	generateModule := Module(generatePath)

	sourcesOverwrite := []gen.TextFile{}
	sourcesScaffold := []gen.TextFile{}

	validationModule := generateModule.Submodule(validation)
	validationFile := generateValidation(validation, validationModule)
	sourcesOverwrite = append(sourcesOverwrite, *validationFile)
	paramsModule := generateModule.Submodule("params")
	sourcesOverwrite = append(sourcesOverwrite, *generateParamsStaticCode(paramsModule))

	for _, version := range specification.Versions {
		versionModule := generateModule.Submodule(version.Version.FlatCase())
		modelsModule := versionModule.Submodule("models")
		sourcesOverwrite = append(sourcesOverwrite, *generateVersionModels(&version, validation, validationModule, modelsModule))
		sourcesOverwrite = append(sourcesOverwrite, generateServiceApis(&version, modelsModule, versionModule)...)
		sourcesOverwrite = append(sourcesOverwrite, *generateVersionRouting(&version, validation, server, validationModule, paramsModule, versionModule))
	}
	sourcesOverwrite = append(sourcesOverwrite, *generateSpecRouter(specification, server, generateModule, generateModule.Submodule("spec_router")))

	if swaggerPath != "" {
		sourcesOverwrite = append(sourcesOverwrite, *genopenapi.GenerateOpenapi(specification, swaggerPath))
	}

	if servicesPath != "" {
		sourcesScaffold = generateServicesImplementations(specification, generateModule, Module(servicesPath))
	}

	err := gen.WriteFiles(sourcesScaffold, false)
	if err != nil {
		return err
	}

	err = gen.WriteFiles(sourcesOverwrite, true)
	if err != nil {
		return err
	}

	return nil
}

func generateVersionRouting(version *spec.Version, validation string, server string, validationModule, paramsModule, module module) *gen.TextFile {
	routingModule := module.Submodule("routing")
	if server == express {
		return generateExpressVersionRouting(version, validation, validationModule, paramsModule, routingModule)
	}
	if server == koa {
		return generateKoaVersionRouting(version, validation, validationModule, paramsModule, routingModule)
	}
	panic(fmt.Sprintf("Unknown server: %s", server))
}

func generateSpecRouter(specification *spec.Spec, server string, rootModule module, module module) *gen.TextFile {
	if server == express {
		return generateExpressSpecRouter(specification, rootModule, module)
	}
	if server == koa {
		return generateKoaSpecRouter(specification, rootModule, module)
	}
	panic(fmt.Sprintf("Unknown server: %s", server))
}

func apiRouterName(api *spec.Api) string {
	return api.Name.CamelCase() + "Router"
}

func apiRouterNameVersioned(api *spec.Api) string {
	result := apiRouterName(api)
	version := api.Apis.Version.Version
	if version.Source != "" {
		result = result + version.PascalCase()
	}
	return result
}

func apiServiceParamName(api *spec.Api) string {
	version := api.Apis.Version
	name := api.Name.CamelCase() + "Service"
	if version.Version.Source != "" {
		name = name + version.Version.PascalCase()
	}
	return name
}

//TODO: Same as above
func apiVersionedRouterName(api *spec.Api) string {
	version := api.Apis.Version
	name := api.Name.CamelCase() + "Router"
	if version.Version.Source != "" {
		name = name + version.Version.PascalCase()
	}
	return name
}
