package genscala

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func GenerateSttpClient(specification *spec.Spec, packageName string, generatePath string) *sources.Sources {
	if packageName == "" {
		packageName = specification.Name.FlatCase()
	}
	mainPackage := NewPackage(generatePath, packageName, "")
	jsonPackage := mainPackage
	paramsPackage := mainPackage

	sources := sources.NewSources()
	jsonHelpers := generateJson(jsonPackage)
	taggedUnion := generateTaggedUnion(jsonPackage)
	sources.AddGenerated(taggedUnion, jsonHelpers)

	scalaHttpStaticFile := generateStringParams(paramsPackage)
	sources.AddGenerated(scalaHttpStaticFile)

	for _, version := range specification.Versions {
		versionClientPackage := mainPackage.Subpackage(version.Version.FlatCase())
		versionModelsPackage := versionClientPackage.Subpackage("models")
		clientImplementations := generateClientImplementations(&version, versionClientPackage, versionModelsPackage, jsonPackage, paramsPackage)
		sources.AddGeneratedAll(clientImplementations)
		models := generateCirceModels(&version, versionModelsPackage, jsonPackage)
		sources.AddGenerated(models)
	}

	return sources
}

func generateClientImplementations(version *spec.Version, thepackage, modelsPackage, jsonPackage, paramsPackage Package) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, api := range version.Http.Apis {
		apiPackage := thepackage.Subpackage(api.Name.FlatCase())
		apiClient := generateApiClientApi(&api, apiPackage, modelsPackage, jsonPackage, paramsPackage)
		files = append(files, *apiClient)
	}
	return files
}

func generateApiClientApi(api *spec.Api, thepackage, modelsPackage, jsonPackage, stringParamsPackage Package) *sources.CodeFile {
	w := NewScalaWriter()

	w.Line(`package %s`, thepackage.PackageName)
	w.EmptyLine()
	w.Line(`import scala.concurrent._`)
	w.Line(`import org.slf4j._`)
	w.Line(`import com.softwaremill.sttp._`)
	w.Line(`import %s.ParamsTypesBindings._`, stringParamsPackage.PackageName)
	if jsonPackage.PackageName != thepackage.PackageName {
		w.Line(`import %s.Jsoner`, jsonPackage.PackageName)
	}
	w.Line(`import %s`, modelsPackage.PackageStar)

	w.EmptyLine()
	generateClientApiTrait(w, api)

	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 {
			w.EmptyLine()
			generateResponse(w, &operation)
		}
	}

	w.EmptyLine()
	generateClientApiClass(w, api)

	return &sources.CodeFile{
		Path:    thepackage.GetPath("Client.scala"),
		Content: w.String(),
	}
}

func createParams(params []spec.NamedParam, defaulted bool) []string {
	methodParams := []string{}
	for _, param := range params {
		if !defaulted && param.Default == nil {
			methodParams = append(methodParams, fmt.Sprintf(`%s: %s`, param.Name.CamelCase(), ScalaType(&param.Type.Definition)))
		}
		if defaulted && param.Default != nil {
			defaultValue := DefaultValue(&param.Type.Definition, *param.Default)
			methodParams = append(methodParams, fmt.Sprintf(`%s: %s = %s`, param.Name.CamelCase(), ScalaType(&param.Type.Definition), defaultValue))
		}
	}
	return methodParams
}

func createUrlParams(urlParams []spec.NamedParam) []string {
	methodParams := []string{}
	for _, param := range urlParams {
		methodParams = append(methodParams, fmt.Sprintf(`%s: %s`, param.Name.CamelCase(), ScalaType(&param.Type.Definition)))
	}
	return methodParams
}

func generateClientOperationSignature(operation *spec.NamedOperation) string {
	methodParams := []string{}
	methodParams = append(methodParams, createParams(operation.HeaderParams, false)...)
	if operation.Body != nil {
		methodParams = append(methodParams, fmt.Sprintf(`body: %s`, ScalaType(&operation.Body.Type.Definition)))
	}
	methodParams = append(methodParams, createUrlParams(operation.Endpoint.UrlParams)...)
	methodParams = append(methodParams, createParams(operation.QueryParams, false)...)
	methodParams = append(methodParams, createParams(operation.HeaderParams, true)...)
	methodParams = append(methodParams, createParams(operation.QueryParams, true)...)

	return fmt.Sprintf(`def %s(%s): Future[%s]`, operation.Name.CamelCase(), JoinParams(methodParams), responseType(operation))
}

func generateClientApiTrait(w *sources.Writer, api *spec.Api) {
	apiTraitName := clientTraitName(api.Name)
	w.Line(`trait %s {`, apiTraitName)
	for _, operation := range api.Operations {
		w.Line(`  %s`, generateClientOperationSignature(&operation))
	}
	w.Line(`}`)
}

func clientTraitName(apiName spec.Name) string {
	return "I" + apiName.PascalCase() + "Client"
}

func clientClassName(apiName spec.Name) string {
	return apiName.PascalCase() + "Client"
}

func addParamsWriting(w *sources.Writer, params []spec.NamedParam, paramsName string) {
	if params != nil && len(params) > 0 {
		w.Line(`val %s = new StringParamsWriter()`, paramsName)
		for _, p := range params {
			w.Line(`%s.write("%s", %s)`, paramsName, p.Name.Source, p.Name.CamelCase())
		}
	}
}

func generateClientOperationImplementation(w *sources.Writer, operation *spec.NamedOperation) {
	httpMethod := strings.ToLower(operation.Endpoint.Method)
	url := operation.FullUrl()
	for _, param := range operation.Endpoint.UrlParams {
		url = strings.Replace(url, spec.UrlParamStr(&param), "${stringify("+param.Name.CamelCase()+")}", -1)
	}

	addParamsWriting(w, operation.QueryParams, "query")
	if operation.QueryParams != nil && len(operation.QueryParams) > 0 {
		w.Line(`val url = Uri.parse(baseUrl+s"%s").get.params(query.params:_*)`, url)
	} else {
		w.Line(`val url = Uri.parse(baseUrl+s"%s").get`, url)
	}

	addParamsWriting(w, operation.HeaderParams, "headers")
	if operation.Body != nil {
		if operation.Body.Type.Definition.Plain == spec.TypeString {
			w.Line(`logger.debug(s"Request to url: ${url}, body: ${body}")`)
		} else {
			w.Line(`val bodyJson = Jsoner.write(body)`)
			w.Line(`logger.debug(s"Request to url: ${url}, body: ${bodyJson}")`)
		}
	} else {
		w.Line(`logger.debug(s"Request to url: ${url}")`)
	}
	w.Line(`val response: Future[Response[String]] =`)
	w.Line(`  sttp`)
	w.Line(`    .%s(url)`, httpMethod)

	if operation.HeaderParams != nil && len(operation.HeaderParams) > 0 {
		w.Line(`    .headers(headers.params:_*)`)
	}
	if operation.Body != nil {
		if operation.Body.Type.Definition.Plain == spec.TypeString {
			w.Line(`    .header("Content-Type", "text/plain")`)
			w.Line(`    .body(body)`)
		} else {
			w.Line(`    .header("Content-Type", "application/json")`)
			w.Line(`    .body(bodyJson)`)
		}
	}
	w.Line(`    .parseResponseIf { status => status < 500 }`)
	w.Line(`    .send()`)

	w.Line(`response.map {`)
	w.Line(`  response: Response[String] =>`)
	w.Line(`    response.body match {`)
	w.Line(`      case Right(body) =>`)
	w.Line(`        logger.debug(s"Response status: ${response.code}, body: ${body}")`)
	w.Line(`        response.code match {`)
	generateClientResponses(w.IndentedWith(5), operation)
	w.Line(`          case _ => `)
	w.Line(`            val errorMessage = s"Request returned unexpected status code: ${response.code}, body: ${new String(body)}"`)
	w.Line(`            logger.error(errorMessage)`)
	w.Line(`            throw new RuntimeException(errorMessage)`)
	w.Line(`        }`)
	w.Line(`      case Left(errorData) =>`)
	w.Line(`        val errorMessage = s"Request failed, status code: ${response.code}, body: ${new String(errorData)}"`)
	w.Line(`        logger.error(errorMessage)`)
	w.Line(`        throw new RuntimeException(errorMessage)`)
	w.Line(`    }`)
	w.Line(`}`)
}

func generateClientResponses(w *sources.Writer, operation *spec.NamedOperation) {
	if len(operation.Responses) == 1 {
		response := operation.Responses[0]
		responseParam := `()`
		if !response.Type.Definition.IsEmpty() {
			if response.Type.Definition.Plain == spec.TypeString {
				responseParam = `body`
			} else {
				responseParam = fmt.Sprintf(`Jsoner.readThrowing[%s](body)`, ScalaType(&response.Type.Definition))
			}
		}
		w.Line(`case %s => %s`, spec.HttpStatusCode(response.Name), responseParam)
	} else {
		for _, response := range operation.Responses {
			responseParam := ``
			if !response.Type.Definition.IsEmpty() {
				if response.Type.Definition.Plain == spec.TypeString {
					responseParam = `body`
				} else {
					responseParam = fmt.Sprintf(`Jsoner.readThrowing[%s](body)`, ScalaType(&response.Type.Definition))
				}
			}
			w.Line(`case %s => %s.%s(%s)`, spec.HttpStatusCode(response.Name), responseTypeName(operation), response.Name.PascalCase(), responseParam)
		}
	}
}

func generateClientApiClass(w *sources.Writer, api *spec.Api) {
	apiClassName := clientClassName(api.Name)
	apiTraitName := clientTraitName(api.Name)

	w.Line(`class %s(baseUrl: String)(implicit backend: SttpBackend[Future, Nothing]) extends %s {`, apiClassName, apiTraitName)
	w.Line(`  import ExecutionContext.Implicits.global`)
	w.Line(`  private val logger: Logger = LoggerFactory.getLogger(this.getClass)`)
	for _, operation := range api.Operations {
		w.Line(`  %s = {`, generateClientOperationSignature(&operation))
		generateClientOperationImplementation(w.IndentedWith(2), &operation)
		w.Line(`  }`)
	}
	w.Line(`}`)
}
