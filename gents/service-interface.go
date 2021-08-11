package gents

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
	"strings"
)

func generateServiceApis(version *spec.Version, generatePath string) *gen.TextFile {
	w := NewTsWriter()

	w.Line("import * as %s from './%s'", modelsPackage, versionFilename(version, "models", ""))

	for _, api := range version.Http.Apis {
		generateApiService(w, &api)
	}

	filename := versionFilename(version, "services", "ts")
	return &gen.TextFile{filepath.Join(generatePath, filename), w.String()}
}

func generateApiService(w *gen.Writer, api *spec.Api) {
	for _, operation := range api.Operations {
		if operation.Body != nil || operation.HasParams() {
			w.EmptyLine()
			generateOperationParams(w, &operation)
		}
		w.EmptyLine()
		generateOperationResponse(w, &operation)
	}
	w.EmptyLine()
	w.Line("export interface %s {", serviceInterfaceName(api))
	for _, operation := range api.Operations {
		params := ""
		if operation.Body != nil || operation.HasParams() {
			params = fmt.Sprintf(`params: %s`, operationParamsTypeName(&operation))
		}
		w.Line("  %s(%s): Promise<%s>", operation.Name.CamelCase(), params, responseTypeName(&operation))
	}
	w.Line("}")
}

func serviceInterfaceName(api *spec.Api) string {
	return api.Name.PascalCase() + "Service"
}

func operationParamsTypeName(operation *spec.NamedOperation) string {
	return operation.Name.PascalCase() + "Params"
}

func addServiceParam(w *gen.Writer, paramName string, typ *spec.TypeDef) {
	if typ.IsNullable() {
		paramName = paramName + "?"
	}
	w.Line("  %s: %s,", paramName, TsType(typ))
}

func generateServiceParams(w *gen.Writer, params []spec.NamedParam, isHeader bool) {
	for _, param := range params {
		paramName := param.Name.Source
		if isHeader {
			paramName = strings.ToLower(param.Name.Source)
		}
		addServiceParam(w, tsIdentifier(paramName), &param.Type.Definition)
	}
}

func tsIdentifier(name string) string {
	if !checkFormat(tsIdentifierFormat, name) {
		return fmt.Sprintf("'%s'", name)
	}
	return name
}

func generateOperationParams(w *gen.Writer, operation *spec.NamedOperation) {
	w.Line("export interface %s {", operationParamsTypeName(operation))
	generateServiceParams(w, operation.HeaderParams, true)
	generateServiceParams(w, operation.Endpoint.UrlParams, false)
	generateServiceParams(w, operation.QueryParams, false)
	if operation.Body != nil {
		addServiceParam(w, "body", &operation.Body.Type.Definition)
	}
	w.Line("}")
}
