package gents

import (
	"fmt"
	"github.com/specgen-io/spec"
	"specgen/gen"
	"specgen/genopenapi"
)

func GenerateService(serviceFile string, swaggerPath string, generatePath string, servicesPath string, server string, validation string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}

	sourcesOverwrite := []gen.TextFile{}
	sourcesScaffold := []gen.TextFile{}

	for _, version := range specification.Versions {
		sourcesOverwrite = append(sourcesOverwrite, *generateServiceApis(&version, generatePath))
		sourcesOverwrite = append(sourcesOverwrite, *generateVersionRouting(&version, validation, server, generatePath))
	}
	sourcesOverwrite = append(sourcesOverwrite, *generateSpecRouter(specification, server, generatePath))

	modelsFiles := generateModels(specification, validation, generatePath)
	sourcesOverwrite = append(sourcesOverwrite, modelsFiles...)

	if swaggerPath != "" {
		sourcesOverwrite = append(sourcesOverwrite, *genopenapi.GenerateOpenapi(specification, swaggerPath))
	}

	if servicesPath != "" {
		sourcesScaffold = generateServicesImplementations(specification, servicesPath)
	}

	err = gen.WriteFiles(sourcesScaffold, false)
	if err != nil {
		return err
	}

	err = gen.WriteFiles(sourcesOverwrite, true)
	if err != nil {
		return err
	}

	return nil
}

func generateVersionRouting(version *spec.Version, validation string, server string, generatePath string) *gen.TextFile {
	if server == express {
		return generateExpressVersionRouting(version, validation, generatePath)
	}
	if server == koa {
		return generateKoaVersionRouting(version, validation, generatePath)
	}
	panic(fmt.Sprintf("Unknown server: %s", server))
}

func generateSpecRouter(specification *spec.Spec, server string, generatePath string) *gen.TextFile {
	if server == express {
		return generateExpressSpecRouter(specification, generatePath)
	}
	if server == koa {
		return generateKoaSpecRouter(specification, generatePath)
	}
	panic(fmt.Sprintf("Unknown server: %s", server))
}

func apiRouterName(api *spec.Api) string {
	return api.Name.CamelCase() + "Router"
}

func apiServiceParamName(api *spec.Api) string {
	version := api.Apis.Version
	name := api.Name.CamelCase() + "Service"
	if version.Version.Source != "" {
		name = name + version.Version.PascalCase()
	}
	return name
}

func apiVersionedRouterName(api *spec.Api) string {
	version := api.Apis.Version
	name := api.Name.CamelCase() + "Router"
	if version.Version.Source != "" {
		name = name + version.Version.PascalCase()
	}
	return name
}
