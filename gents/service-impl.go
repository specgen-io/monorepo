package gents

import (
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
	"strings"
)

func generateServicesImplementations(specification *spec.Spec, generatePath string) []gen.TextFile {
	files := []gen.TextFile{}
	for _, version := range specification.Versions {
		for _, api := range version.Http.Apis {
			files = append(files, *generateServiceImplementation(&api, generatePath))
		}
	}
	return files
}

func generateServiceImplementation(api *spec.Api, generatePath string) *gen.TextFile {
	w := NewTsWriter()

	versionModule(api.Apis.Version, "service")

	w.Line("import * as services from './spec/%s'", versionModule(api.Apis.Version, "services"))
	w.EmptyLine()
	w.Line("export let %sService = (): services.%s => {", api.Name.CamelCase(), serviceInterfaceName(api))

	operations := []string{}
	for _, operation := range api.Operations {
		operations = append(operations, operation.Name.CamelCase())
		w.Line("  let %s = async (params: services.%s): Promise<services.%s> => {", operation.Name.CamelCase(), operationParamsTypeName(&operation), responseTypeName(&operation))
		w.Line("    throw new Error('Not Implemented')")
		w.Line("  }")
		w.EmptyLine()
	}
	w.Line("  return {%s}", strings.Join(operations, ", "))
	w.Line("}")
	return &gen.TextFile{filepath.Join(generatePath, versionFilename(api.Apis.Version, api.Name.SnakeCase()+"_service", "ts")), w.String()}
}
