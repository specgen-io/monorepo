package gents

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/genopenapi"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateService(specification *spec.Spec, swaggerPath string, generatePath string, servicesPath string, server string, validation string) *sources.Sources {
	sources := sources.NewSources()
	generateModule := Module(generatePath)

	validationModule := generateModule.Submodule(validation)
	sources.AddGenerated(generateValidation(validation, validationModule))
	paramsModule := generateModule.Submodule("params")
	sources.AddGenerated(generateParamsStaticCode(paramsModule))

	for _, version := range specification.Versions {
		versionModule := generateModule.Submodule(version.Version.FlatCase())
		modelsModule := versionModule.Submodule("models")
		sources.AddGenerated(generateVersionModels(&version, validation, validationModule, modelsModule))
		sources.AddGeneratedAll(generateServiceApis(&version, modelsModule, versionModule))
		sources.AddGenerated(generateVersionRouting(&version, validation, server, validationModule, paramsModule, versionModule))
	}
	sources.AddGenerated(generateSpecRouter(specification, server, generateModule, generateModule.Submodule("spec_router")))

	if swaggerPath != "" {
		sources.AddGenerated(genopenapi.GenerateOpenapi(specification, swaggerPath))
	}

	if servicesPath != "" {
		sources.AddScaffoldedAll(generateServicesImplementations(specification, generateModule, Module(servicesPath)))
	}

	return sources
}

func generateVersionRouting(version *spec.Version, validation string, server string, validationModule, paramsModule, module module) *sources.CodeFile {
	routingModule := module.Submodule("routing")
	if server == express {
		return generateExpressVersionRouting(version, validation, validationModule, paramsModule, routingModule)
	}
	if server == koa {
		return generateKoaVersionRouting(version, validation, validationModule, paramsModule, routingModule)
	}
	panic(fmt.Sprintf("Unknown server: %s", server))
}

func generateSpecRouter(specification *spec.Spec, server string, rootModule module, module module) *sources.CodeFile {
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
