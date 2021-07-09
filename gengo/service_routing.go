package gengo

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
	"strings"
)

func generateRouting(version *spec.Version, packageName string, generatePath string) *gen.TextFile {
	w := NewGoWriter()
	w.Line("package %s", packageName)
	w.EmptyLine()
	w.Line("import (")
	w.Line(`  "encoding/json"`)
	w.Line(`  "fmt"`)
	w.Line(`  "github.com/husobee/vestigo"`)
	w.Line(`  "net/http"`)
	w.Line(`)`)
	generateCheckErrors(w)
	generateCheckOperationErrors(w)
	w.EmptyLine()

	for _, api := range version.Http.Apis {
		generateApiRouter(w, api)
	}

	return &gen.TextFile{
		Path:    filepath.Join(generatePath, "routing.go"),
		Content: w.String(),
	}
}

func generateCheckErrors(w *gen.Writer) {
	w.Line(`
func checkErrors(params *ParamsParser, w http.ResponseWriter) bool {
	if len(params.Errors) > 0 {
		w.WriteHeader(400)
		fmt.Println(params.Errors)
		return false
		}
	return true
}`)
}

func generateCheckOperationErrors(w *gen.Writer) {
	w.Line(`
func checkOperationErrors(err error, w http.ResponseWriter) bool {
	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err.Error())
		return false
	}
	return true
}`)
}

func getVestigoUrl(operation spec.NamedOperation) []string {
	reminder := operation.FullUrl()
	urlParams := []string{}
	if operation.Endpoint.UrlParams != nil && len(operation.Endpoint.UrlParams) > 0 {
		for _, param := range operation.Endpoint.UrlParams {
			parts := strings.Split(reminder, spec.UrlParamStr(param.Name.Source))
			urlParams = append(urlParams, fmt.Sprintf("%s:%s", parts[0], param.Name.Source))
			reminder = parts[1]
		}
	}
	urlParams = append(urlParams, fmt.Sprintf("%s", reminder))
	return urlParams
}

func addSetCors(w *gen.Writer, operation spec.NamedOperation) {
	w.Line(`router.SetCors("%s", &vestigo.CorsAccessControl{`, JoinParams(getVestigoUrl(operation)))
	params := []string{}
	for _, param := range operation.HeaderParams {
		params = append(params, fmt.Sprintf("%s", param.Name.Source))
	}
	w.Line(`  AllowHeaders: []string{"%s"},`, JoinDelimParams(params))
	w.Line(`})`)
}

func generateApiRouter(w *gen.Writer, api spec.Api) {
	w.Line(`func Add%sRoutes(router *vestigo.Router, %s %s) {`, api.Name.PascalCase(), serviceInterfaceTypeVar(&api), serviceInterfaceTypeName(&api))
	for _, operation := range api.Operations {
		w.Indent()
		w.Line(`router.%s("%s", func(w http.ResponseWriter, r *http.Request) {`, ToPascalCase(operation.Endpoint.Method), JoinParams(getVestigoUrl(operation)))
		generateOperationMethod(w, api, operation)
		w.Line(`})`)
		if operation.HeaderParams != nil && len(operation.HeaderParams) > 0 {
			addSetCors(w, operation)
		}
		w.Unindent()
	}
	w.Line(`}`)
	w.EmptyLine()
}

func parserParameterCall(operation spec.NamedOperation, param *spec.NamedParam, paramsParserName string) string {
	paramNameSource := param.Name.Source
	if operation.Endpoint.UrlParams != nil && len(operation.Endpoint.UrlParams) > 0 {
		paramNameSource = ":" + paramNameSource
	}
	parserParams := []string{fmt.Sprintf(`"%s"`, paramNameSource)}
	methodName, defaultParam := parserDefaultName(param)
	isEnum := param.Type.Definition.Info.Model != nil && param.Type.Definition.Info.Model.IsEnum()
	enumModel := param.Type.Definition.Info.Model
	if isEnum {
		parserParams = append(parserParams, enumModel.Name.PascalCase()+`ValuesStrings`)
	}
	if defaultParam != nil {
		parserParams = append(parserParams, *defaultParam)
	}
	call := fmt.Sprintf(`%s.%s(%s)`, paramsParserName, methodName, JoinDelimParams(parserParams))
	if isEnum {
		call = fmt.Sprintf(`%s(%s)`, enumModel.Name.PascalCase(), call)
	}
	return call
}

func addParamsParser(w *gen.Writer, operation spec.NamedOperation, namedParams []spec.NamedParam, paramsParserName string, paramName string) {
	if namedParams != nil && len(namedParams) > 0 {
		w.Line(`  %s := NewParamsParser(r.%s)`, paramsParserName, paramName)
		for _, param := range namedParams {
			w.Line(`  %s := %s`, param.Name.CamelCase(), parserParameterCall(operation, &param, paramsParserName))
		}
		w.Line(`  if !checkErrors(%s, w) { return }`, paramsParserName)
	}
}

func generateOperationMethod(w *gen.Writer, api spec.Api, operation spec.NamedOperation) {
	if operation.Body != nil {
		w.Line(`  var body %s`, GoType(&operation.Body.Type.Definition))
		w.Line(`  json.NewDecoder(r.Body).Decode(&body)`)
	}
	addParamsParser(w, operation, operation.QueryParams, "queryParams", "URL.Query()")
	addParamsParser(w, operation, operation.HeaderParams, "headerParams", "Header")
	addParamsParser(w, operation, operation.Endpoint.UrlParams, "urlParams", "URL.Query()")

	w.Line(`  response, err := %s.%s(%s)`, serviceInterfaceTypeVar(&api), operation.Name.PascalCase(), JoinDelimParams(addOperationMethodParams(operation)))
	w.Line(`  if !checkOperationErrors(err, w) { return }`)
	for _, response := range operation.Responses {
		w.Line(`  if response.%s != nil {`, response.Name.PascalCase())
		w.Line(`    w.WriteHeader(%s)`, spec.HttpStatusCode(response.Name))
		if spec.HttpStatusCode(response.Name) == "200" {
			w.Line(`    json.NewEncoder(w).Encode(response.%s)`, response.Name.PascalCase())
		}
		w.Line(`    return`)
		w.Line(`  }`)
	}
}

func addOperationMethodParams(operation spec.NamedOperation) []string {
	urlParams := []string{}
	if operation.Body != nil {
		urlParams = append(urlParams, fmt.Sprintf("%s", "&body"))
	}
	for _, param := range operation.QueryParams {
		urlParams = append(urlParams, fmt.Sprintf("%s", param.Name.CamelCase()))
	}
	for _, param := range operation.HeaderParams {
		urlParams = append(urlParams, fmt.Sprintf("%s", param.Name.CamelCase()))
	}
	for _, param := range operation.Endpoint.UrlParams {
		urlParams = append(urlParams, fmt.Sprintf("%s", param.Name.CamelCase()))
	}
	return urlParams
}
