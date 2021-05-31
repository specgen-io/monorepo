package gengo

import (
	"fmt"
	"github.com/pinzolo/casee"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
	"strings"
)

func GenerateService(serviceFile string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}
	files := []gen.TextFile{}
	for _, version := range specification.Versions {
		w := NewGoWriter()
		generateRouting(w, &version, "spec")
		folder := "spec"
		if version.Version.Source != "" {
			folder += "_" + version.Version.FlatCase()
		}
		files = append(files, gen.TextFile{Path: filepath.Join(generatePath, folder, "routing.go"), Content: w.String()})
		files = append(files, *generateParamsParser(folder, filepath.Join(generatePath, folder, "parsing.go")))
	}
	return gen.WriteFiles(files, true)
}

func generateRouting(w *gen.Writer, version *spec.Version, packageName string) {
	w.Line("package %s", versionedPackage(version.Version, packageName))
	w.EmptyLine()
	w.Line("import (")
	w.Line(`  "encoding/json"`)
	w.Line(`  "fmt"`)
	w.Line(`  "github.com/husobee/vestigo"`)
	w.Line(`  "net/http"`)
	w.Line(`)`)
	generateCheckErrors(w)
	for _, api := range version.Http.Apis {
		generateApiRouter(w, api)
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
}
`)
}

func generateApiRouter(w *gen.Writer, api spec.Api) {
	apiName := api.Name.PascalCase()
	w.Line(`func Add%sRoutes(router *vestigo.Router, %sService I%sService) {`, apiName, api.Name.Source, apiName)
	for _, operation := range api.Operations {
		reminder := operation.FullUrl()
		urlParams := []string{}
		for _, param := range operation.Endpoint.UrlParams {
			paramName := param.Name.Source
			parts := strings.Split(reminder, spec.UrlParamStr(paramName))
			urlParams = append(urlParams, fmt.Sprintf("%s:%s", parts[0], paramName))
			reminder = parts[1]
		}
		w.Indent()
		if reminder != `` {
			w.Line(`router.%s("%s", func(w http.ResponseWriter, r *http.Request) {`, nameToPascalCase(operation.Endpoint.Method), reminder)
		} else {
			w.Line(`router.%s("%s", func(w http.ResponseWriter, r *http.Request) {`, nameToPascalCase(operation.Endpoint.Method), strings.Join(urlParams, ""))
		}

		for _, response := range operation.Responses {
			generateOperationMethod(w, response, api.Name.Source, operation)
		}

		w.Unindent()
		w.Line(`  })`)
	}
	w.Line(`}`)
	w.EmptyLine()
}

func nameToPascalCase(name string) string {
	return casee.ToPascalCase(name)
}

func parserMethodName(typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return casee.ToPascalCase(typ.Name)
	case spec.NullableType:
		return parserMethodName(typ.Child) + "Nullable"
	case spec.ArrayType:
		return parserMethodName(typ.Child) + "Array"
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func parserParameterCall(param *spec.NamedParam) string {
	if param.Default != nil {
		defaultValue := DefaultValue(&param.Type.Definition, *param.Default)
		return fmt.Sprintf(`%sDefaulted("%s", %s)`, parserMethodName(&param.Type.Definition), param.Name.Source, defaultValue)
	} else {
		return fmt.Sprintf(`%s("%s")`, parserMethodName(&param.Type.Definition), param.Name.Source)
	}
}

func generateOperationMethod(w *gen.Writer, response spec.NamedResponse, apiName string, operation spec.NamedOperation) {
	if operation.Body != nil {
		w.Line(`  var body %s`, GoType(&response.Type.Definition))
		w.Line(`  json.NewDecoder(r.Body).Decode(&body)`)
		w.Line(`  response := %sService.%s(&body)`, apiName, operation.Name.PascalCase())
		w.Line(`  w.WriteHeader(%s)`, spec.HttpStatusCode(response.Name))
		w.Line(`  json.NewEncoder(w).Encode(response.%s)`, response.Name.PascalCase())
	} else if operation.QueryParams != nil && len(operation.QueryParams) > 0 {
		w.Line(`  query := NewParamsParser(r.URL.Query())`)
		for _, param := range operation.QueryParams {
			w.Line(`  %s := query.%s`, param.Name.CamelCase(), parserParameterCall(&param))
		}
		w.Line(`  if checkErrors(query, w) {`)
		w.Line(`    response := %sService.%s(%s)`, apiName, operation.Name.PascalCase(), strings.Join(addParameters(operation), ", "))
		w.Line(`    w.WriteHeader(%s)`, spec.HttpStatusCode(response.Name))
		w.Line(`    json.NewEncoder(w).Encode(response.%s)`, response.Name.PascalCase())
		w.Line(`  }`)
	}
}

func addParameters(operation spec.NamedOperation) []string{
	urlParams := []string{}

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
