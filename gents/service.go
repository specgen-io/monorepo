package gents

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/genopenapi"
)

func GenerateService(serviceFile string, swaggerPath string, generatePath string, servicesPath string, server string, validation string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}

	sources := []gen.TextFile{}

	for _, version := range specification.Versions {
		sources = append(sources, *generateServiceApis(&version, generatePath))
		sources = append(sources, *generateVersionRouting(&version, validation, server, generatePath))
	}
	sources = append(sources, *generateSpecRouter(specification, server, generatePath))

	modelsFiles := generateModels(specification, validation, generatePath)
	sources = append(sources, modelsFiles...)

	sources = append(sources, *genopenapi.GenerateOpenapi(specification, swaggerPath))

	if servicesPath != "" {
		services := generateServicesImplementations(specification, servicesPath)
		err = gen.WriteFiles(services, false)
		if err != nil {
			return err
		}
	}

	err = gen.WriteFiles(sources, true)
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
	paramName := api.Name.CamelCase() + "Service"
	if version.Version.Source != "" {
		paramName = paramName + version.Version.PascalCase()
	}
	return paramName
}