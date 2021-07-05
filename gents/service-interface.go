package gents

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
	"strings"
)

func genereateServiceApis(version *spec.Version, generatePath string) *gen.TextFile {
	w := NewTsWriter()

	w.Line("import * as models from './%s'", versionFilename(version, "models", ""))

	for _, api := range version.Http.Apis {
		generateApiService(w, &api)
	}

	filename := versionFilename(version, "services", "ts")
	return &gen.TextFile{filepath.Join(generatePath, filename), w.String()}
}

func generateApiService(w *gen.Writer, api *spec.Api) {
	for _, operation := range api.Operations {
		w.EmptyLine()
		generateOperationParams(w, &operation)
		w.EmptyLine()
		generateResponse(w, &operation, &modelsPackage)
	}
	w.EmptyLine()
	w.Line("export interface %s {", serviceInterfaceName(api))
	for _, operation := range api.Operations {
		w.Line("  %s(params: %s): Promise<%s>", operation.Name.CamelCase(), paramsTypeName(&operation), responseTypeName(&operation))
	}
	w.Line("}")
}

func serviceInterfaceName(api *spec.Api) string {
	return api.Name.PascalCase() + "Service"
}

func paramsTypeName(operation *spec.NamedOperation) string {
	return operation.Name.PascalCase() + "Params"
}

var modelsPackage = "models"

func addServiceParam(w *gen.Writer, paramName string, typ *spec.TypeDef, modelsPackage *string) {
	if typ.IsNullable() {
		paramName = paramName + "?"
	}
	w.Line("  %s: %s,", paramName, PackagedTsType(typ, modelsPackage))
}

func generateServiceParams(w *gen.Writer, params []spec.NamedParam, isHeader bool) {
	for _, param := range params {
		paramName := param.Name.Source
		if isHeader {
			paramName = isTsIdentifier(strings.ToLower(param.Name.Source))
		}
		addServiceParam(w, paramName, &param.Type.Definition, &modelsPackage)
	}
}

func isTsIdentifier(name string) string {
	if strings.Contains(name, "-") {
		return fmt.Sprintf("'%s'", name)
	}
	return name
}

func generateOperationParams(w *gen.Writer, operation *spec.NamedOperation) {
	w.Line("export interface %s {", paramsTypeName(operation))
	generateServiceParams(w, operation.HeaderParams, true)
	generateServiceParams(w, operation.Endpoint.UrlParams, false)
	generateServiceParams(w, operation.QueryParams, false)
	if operation.Body != nil {
		addServiceParam(w, "body", &operation.Body.Type.Definition, &modelsPackage)
	}
	w.Line("}")
}
