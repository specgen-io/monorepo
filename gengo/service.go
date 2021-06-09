package gengo

import (
	"fmt"
	"github.com/pinzolo/casee"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
	"strings"
)

var ToPascalCase = casee.ToPascalCase

func GenerateService(serviceFile string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}
	files := []gen.TextFile{}
	for _, version := range specification.Versions {
		folder := "spec"
		if version.Version.Source != "" {
			folder += "_" + version.Version.FlatCase()
		}

		files = append(files, *generateRouting(&version, "spec", generatePath, folder))

		files = append(files, *generateParamsParser(folder, filepath.Join(generatePath, folder, "parsing.go")))
		files = append(files, *generateEchoService(folder, filepath.Join(generatePath, folder, "echo_service.go")))
		files = append(files, *generateCheckService("spec", filepath.Join(generatePath, "spec", "check_service.go")))
		files = append(files, *generateEchoServiceV2("spec_v2", filepath.Join(generatePath, "spec_v2", "echo_service.go")))
	}
	return gen.WriteFiles(files, true)
}

func generateRouting(version *spec.Version, packageName, generatePath, folder string) *gen.TextFile {
	w := NewGoWriter()
	w.Line("package %s", versionedPackage(version.Version, packageName))
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
		Path:    filepath.Join(generatePath, folder, "routing.go"),
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
	w.Line(`router.SetCors("%s", &vestigo.CorsAccessControl{`, strings.Join(getVestigoUrl(operation), ""))
	params := []string{}
	for _, param := range operation.HeaderParams {
		params = append(params, fmt.Sprintf("%s", param.Name.Source))
	}
	w.Line(`  AllowHeaders: []string{"%s"},`, strings.Join(params, ", "))
	w.Line(`})`)
}

func generateApiRouter(w *gen.Writer, api spec.Api) {
	apiName := api.Name.PascalCase()
	w.Line(`func Add%sRoutes(router *vestigo.Router, %sService I%sService) {`, apiName, api.Name.Source, apiName)
	for _, operation := range api.Operations {
		w.Indent()
		w.Line(`router.%s("%s", func(w http.ResponseWriter, r *http.Request) {`, ToPascalCase(operation.Endpoint.Method), strings.Join(getVestigoUrl(operation), ""))
		generateOperationMethod(w, api.Name.Source, operation)
		w.Line(`})`)
		if operation.HeaderParams != nil && len(operation.HeaderParams) > 0 {
			addSetCors(w, operation)
		}
		w.Unindent()
	}
	w.Line(`}`)
	w.EmptyLine()
}

func parserDefaultName(param *spec.NamedParam) (string, *string) {
	methodName := parserMethodName(&param.Type.Definition)
	if param.Default != nil {
		defaultValue := DefaultValue(&param.Type.Definition, *param.Default)
		return methodName + `Defaulted`, &defaultValue
	} else {
		return methodName, nil
	}
}

func parserMethodName(typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return parserMethodNamePlain(typ)
	case spec.NullableType:
		return parserMethodNamePlain(typ.Child) + "Nullable"
	case spec.ArrayType:
		return parserMethodNamePlain(typ.Child) + "Array"
	default:
		panic(fmt.Sprintf("Unsupported string param type: %v", typ.Plain))
	}
}

func parserMethodNamePlain(typ *spec.TypeDef) string {
	if typ.Info.Model != nil && typ.Info.Model.IsEnum() {
		return "StringEnum"
	}
	switch typ.Plain {
	case spec.TypeInt32:
		return "Int"
	case spec.TypeInt64:
		return "Int64"
	case spec.TypeFloat:
		return "Float32"
	case spec.TypeDouble:
		return "Float64"
	case spec.TypeDecimal:
		return "Decimal"
	case spec.TypeBoolean:
		return "Bool"
	case spec.TypeString:
		return "String"
	case spec.TypeUuid:
		return "Uuid"
	case spec.TypeDate:
		return "Date"
	case spec.TypeDateTime:
		return "DateTime"
	default:
		panic(fmt.Sprintf("Unsupported string param type: %v", typ.Plain))
	}
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
	call := fmt.Sprintf(`%s.%s(%s)`, paramsParserName, methodName, strings.Join(parserParams, ", "))
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

func generateOperationMethod(w *gen.Writer, apiName string, operation spec.NamedOperation) {
	if operation.Body != nil {
		w.Line(`  var body %s`, GoType(&operation.Body.Type.Definition))
		w.Line(`  json.NewDecoder(r.Body).Decode(&body)`)
	}
	addParamsParser(w, operation, operation.QueryParams, "query", "URL.Query()")
	addParamsParser(w, operation, operation.HeaderParams, "header", "Header")
	addParamsParser(w, operation, operation.Endpoint.UrlParams, "query", "URL.Query()")

	w.Line(`  response, err := %sService.%s(%s)`, apiName, operation.Name.PascalCase(), strings.Join(addUrlParams(operation), ", "))
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

func addUrlParams(operation spec.NamedOperation) []string {
	urlParams := []string{}
	if operation.Body == nil {
		for _, param := range operation.QueryParams {
			urlParams = append(urlParams, fmt.Sprintf("%s", param.Name.CamelCase()))
		}
		for _, param := range operation.HeaderParams {
			urlParams = append(urlParams, fmt.Sprintf("%s", param.Name.CamelCase()))
		}
		for _, param := range operation.Endpoint.UrlParams {
			urlParams = append(urlParams, fmt.Sprintf("%s", param.Name.CamelCase()))
		}
	} else {
		urlParams = append(urlParams, fmt.Sprintf("%s", "&body"))
	}

	return urlParams
}
