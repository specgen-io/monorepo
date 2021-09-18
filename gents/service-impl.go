package gents

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"strings"
)

func generateServicesImplementations(specification *spec.Spec, servicesPath string, generatePath string) []gen.TextFile {
	files := []gen.TextFile{}
	for _, version := range specification.Versions {
		for _, api := range version.Http.Apis {
			files = append(files, *generateServiceImplementation(&api, servicesPath, generatePath))
		}
	}
	return files
}

func generateServiceImplementation(api *spec.Api, servicesPath string, generatePath string) *gen.TextFile {
	path := versionedPath(servicesPath, api.Apis.Version, api.Name.SnakeCase()+"_service.ts")
	w := NewTsWriter()

	w.Line("import * as service from '%s'", importPath(serviceApiPath(generatePath, api), path))
	w.EmptyLine()
	w.Line("export let %sService = (): service.%s => {", api.Name.CamelCase(), serviceInterfaceName(api)) //TODO: remove services

	operations := []string{}
	for _, operation := range api.Operations {
		operations = append(operations, operation.Name.CamelCase())
		params := ""
		if operation.Body != nil || operation.HasParams() {
			params = fmt.Sprintf(`params: service.%s`, operationParamsTypeName(&operation))
		}
		w.Line("  let %s = async (%s): Promise<service.%s> => {", operation.Name.CamelCase(), params, responseTypeName(&operation))
		w.Line("    throw new Error('Not Implemented')")
		w.Line("  }")
		w.EmptyLine()
	}
	w.Line("  return {%s}", strings.Join(operations, ", "))
	w.Line("}")
	return &gen.TextFile{path, w.String()}
}
