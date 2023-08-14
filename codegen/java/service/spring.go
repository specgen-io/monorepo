package service

import (
	"fmt"
	"generator"
	"github.com/pinzolo/casee"
	"java/models"
	"java/types"
	"java/writer"
	"spec"
	"strings"
)

var Spring = "spring"

type SpringGenerator struct {
	Types    *types.Types
	Models   models.Generator
	Packages *Packages
}

func NewSpringGenerator(types *types.Types, models models.Generator, servicePackages *Packages) *SpringGenerator {
	return &SpringGenerator{types, models, servicePackages}
}

func (g *SpringGenerator) ServiceImplAnnotation(api *spec.Api) (annotationImport, annotation string) {
	return `org.springframework.stereotype.Service`, fmt.Sprintf(`Service("%s")`, versionServiceName(serviceName(api), api.InHttp.InVersion))
}

func (g *SpringGenerator) ServicesControllers(version *spec.Version) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		files = append(files, *g.serviceController(&api))
	}
	return files
}

func (g *SpringGenerator) ServiceImports() []string {
	return []string{
		`org.apache.logging.log4j.*`,
		`org.springframework.beans.factory.annotation.Autowired`,
		`org.springframework.format.annotation.*`,
		`org.springframework.http.*`,
		`org.springframework.web.bind.annotation.*`,
	}
}

func (g *SpringGenerator) ExceptionController(responses *spec.ErrorResponses) *generator.CodeFile {
	w := writer.New(g.Packages.RootControllers, `ExceptionController`)
	w.Imports.Add(g.ServiceImports()...)
	w.Imports.Star(g.Packages.Json)
	w.Imports.Star(g.Packages.ErrorsModels)
	w.Imports.StaticStar(g.Packages.Errors.Subpackage(ErrorsHelpersClassName))
	w.Imports.AddStatic(`org.apache.tomcat.util.http.fileupload.FileUploadBase.CONTENT_TYPE`)
	w.EmptyLine()
	w.Line(`@ControllerAdvice`)
	w.Line(`public class [[.ClassName]] {`)
	w.Line(`  private static final Logger logger = LogManager.getLogger([[.ClassName]].class);`)
	w.EmptyLine()
	w.Line(`  @Autowired`)
	w.Line(`  private Json json;`)
	w.EmptyLine()
	g.errorHandler(w.Indented(), *responses)
	w.Line(`}`)
	return w.ToCodeFile()
}

func (g *SpringGenerator) errorHandler(w *writer.Writer, errors spec.ErrorResponses) {
	notFoundError := errors.GetByStatusName(spec.HttpStatusNotFound)
	badRequestError := errors.GetByStatusName(spec.HttpStatusBadRequest)
	internalServerError := errors.GetByStatusName(spec.HttpStatusInternalServerError)
	w.Line(`@ExceptionHandler(Throwable.class)`)
	w.Line(`public ResponseEntity<String> error(Throwable exception) {`)
	w.Line(`  var notFoundError = getNotFoundError(exception);`)
	w.Line(`  if (notFoundError != null) {`)
	g.processResponse(w.IndentedWith(2), &notFoundError.Response, "notFoundError")
	w.Line(`  }`)
	w.Line(`  var badRequestError = getBadRequestError(exception);`)
	w.Line(`  if (badRequestError != null) {`)
	g.processResponse(w.IndentedWith(2), &badRequestError.Response, "badRequestError")
	w.Line(`  }`)
	w.Line(`  var internalServerError = new InternalServerError(exception.getMessage());`)
	g.processResponse(w.IndentedWith(1), &internalServerError.Response, "internalServerError")
	w.Line(`}`)
}

func (g *SpringGenerator) serviceController(api *spec.Api) *generator.CodeFile {
	w := writer.New(g.Packages.Controllers(api.InHttp.InVersion), controllerName(api))
	w.Imports.Add(g.ServiceImports()...)
	w.Imports.Add(`javax.servlet.http.HttpServletRequest`)
	w.Imports.Star(g.Packages.ContentType)
	w.Imports.Star(g.Packages.Json)
	w.Imports.Star(g.Packages.ErrorsModels)
	w.Imports.Star(g.Packages.Models(api.InHttp.InVersion))
	w.Imports.Star(g.Packages.ServicesApi(api))
	w.Imports.Add(g.Models.ModelsUsageImports()...)
	w.Imports.Add(g.Types.Imports()...)
	w.Imports.AddStatic(`org.apache.tomcat.util.http.fileupload.FileUploadBase.CONTENT_TYPE`)
	w.EmptyLine()
	w.Line(`@RestController("%s")`, versionControllerName(controllerName(api), api.InHttp.InVersion))
	w.Line(`public class [[.ClassName]] {`)
	w.Line(`  private static final Logger logger = LogManager.getLogger([[.ClassName]].class);`)
	w.EmptyLine()
	w.Line(`  @Autowired`)
	w.Line(`  private %s %s;`, serviceInterfaceName(api), serviceVarName(api))
	w.EmptyLine()
	w.Line(`  @Autowired`)
	w.Line(`  private Json json;`)

	for _, operation := range api.Operations {
		w.EmptyLine()
		g.controllerMethod(w.Indented(), &operation)
	}
	w.Line(`}`)
	return w.ToCodeFile()
}

func (g *SpringGenerator) controllerMethod(w *writer.Writer, operation *spec.NamedOperation) {
	methodName := operation.Endpoint.Method
	url := operation.FullUrl()
	w.Line(`@%sMapping("%s")`, casee.ToPascalCase(methodName), url)
	w.Line(`public ResponseEntity<String> %s(%s) {`, controllerMethodName(operation), strings.Join(springMethodParams(operation, g.Types), ", "))
	w.Indent()
	w.Line(`logger.info("Received request, operationId: %s.%s, method: %s, url: %s");`, operation.InApi.Name.Source, operation.Name.Source, methodName, url)
	g.parseBody(w, operation, "bodyStr", "requestBody")
	g.serviceCall(w, operation, "bodyStr", "requestBody", "result")
	g.processResponses(w, operation, "result")
	w.Unindent()
	w.Line(`}`)
}

func (g *SpringGenerator) parseBody(w *writer.Writer, operation *spec.NamedOperation, bodyStringVar, bodyJsonVar string) {
	if operation.BodyIs(spec.RequestBodyString) {
		w.Line(`ContentType.check(request, MediaType.TEXT_PLAIN);`)
	}
	if operation.BodyIs(spec.RequestBodyJson) {
		w.Line(`ContentType.check(request, MediaType.APPLICATION_JSON);`)
		typ := g.Types.Java(&operation.Body.Type.Definition)
		w.Line(`%s %s = json.%s;`, typ, bodyJsonVar, g.Models.JsonRead(bodyStringVar, &operation.Body.Type.Definition))
	}
}

func (g *SpringGenerator) serviceCall(w *writer.Writer, operation *spec.NamedOperation, bodyStringVar, bodyJsonVar, resultVarName string) {
	serviceCall := fmt.Sprintf(`%s.%s(%s)`, serviceVarName(operation.InApi), operation.Name.CamelCase(), strings.Join(addServiceMethodParams(operation, bodyStringVar, bodyJsonVar), ", "))
	if len(operation.Responses) == 1 && operation.Responses[0].Body.Is(spec.ResponseBodyEmpty) {
		w.Line(`%s;`, serviceCall)
	} else {
		w.Line(`var %s = %s;`, resultVarName, serviceCall)
		w.Line(`if (%s == null) {`, resultVarName)
		w.Line(`  throw new RuntimeException("Service responseImpl didn't return any value'");`)
		w.Line(`}`)
	}
}

func (g *SpringGenerator) processResponses(w *writer.Writer, operation *spec.NamedOperation, resultVarName string) {
	if len(operation.Responses) == 1 {
		g.processResponse(w, &operation.Responses[0].Response, resultVarName)
	}
	if len(operation.Responses) > 1 {
		for _, response := range operation.Responses {
			w.Line(`if (result instanceof %s.%s) {`, responseInterfaceName(operation), response.Name.PascalCase())
			g.processResponse(w.Indented(), &response.Response, getResponseBody(&response, resultVarName))
			w.Line(`}`)
		}
		w.EmptyLine()
		w.Line(`throw new RuntimeException("Service responseImpl didn't return any value'");`)
	}
}

func (g *SpringGenerator) processResponse(w *writer.Writer, response *spec.Response, bodyVar string) {
	if response.Body.Is(spec.ResponseBodyEmpty) {
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s);`, response.Name.UpperCase())
		w.Line(`return new ResponseEntity<>(HttpStatus.%s);`, response.Name.UpperCase())
	}
	if response.Body.Is(spec.ResponseBodyString) {
		w.Line(`HttpHeaders headers = new HttpHeaders();`)
		w.Line(`headers.add(CONTENT_TYPE, "text/plain");`)
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s);`, response.Name.UpperCase())
		w.Line(`return new ResponseEntity<>(%s, headers, HttpStatus.%s);`, bodyVar, response.Name.UpperCase())
	}
	if response.Body.Is(spec.ResponseBodyJson) {
		w.Line(`var bodyJson = json.%s;`, g.Models.JsonWrite(bodyVar, &response.Body.Type.Definition))
		w.Line(`HttpHeaders headers = new HttpHeaders();`)
		w.Line(`headers.add(CONTENT_TYPE, "application/json");`)
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s);`, response.Name.UpperCase())
		w.Line(`return new ResponseEntity<>(bodyJson, headers, HttpStatus.%s);`, response.Name.UpperCase())
	}
}

func (g *SpringGenerator) ContentType() []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *g.contentTypeMismatchException())
	files = append(files, *g.checkContentType())
	return files
}

func (g *SpringGenerator) contentTypeMismatchException() *generator.CodeFile {
	w := writer.New(g.Packages.ContentType, `ContentTypeMismatchException`)
	w.Lines(`
public class ContentTypeMismatchException extends RuntimeException {
    public ContentTypeMismatchException(String expected, String actual) {
        super(String.format("Expected Content-Type header: '%s' was not provided, found: '%s'", expected, actual));
    }
}
`)
	return w.ToCodeFile()
}

func (g *SpringGenerator) checkContentType() *generator.CodeFile {
	w := writer.New(g.Packages.ContentType, `ContentType`)
	w.Lines(`
import org.springframework.http.MediaType;
import javax.servlet.http.HttpServletRequest;

public class ContentType {
	public static void check(HttpServletRequest request, MediaType expectedContentType) {
		var contentType = request.getHeader("Content-Type");
		if (contentType == null || !contentType.contains(expectedContentType.toString())) {
			throw new ContentTypeMismatchException(expectedContentType.toString(), contentType);
		}
	}
}
`)
	return w.ToCodeFile()
}

func (g *SpringGenerator) ErrorsHelpers() *generator.CodeFile {
	w := writer.New(g.Packages.Errors, ErrorsHelpersClassName)
	w.Template(
		map[string]string{
			`ErrorsModelsPackage`: g.Packages.ErrorsModels.PackageName,
			`ContentTypePackage`:  g.Packages.ContentType.PackageName,
			`JsonPackage`:         g.Packages.Json.PackageName,
			`ErrorsPackage`:       g.Packages.Errors.PackageName,
		}, `
import org.springframework.web.bind.*;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.method.annotation.MethodArgumentTypeMismatchException;

import java.util.List;

import [[.ErrorsModelsPackage]].*;
import [[.ContentTypePackage]].*;
import [[.JsonPackage]].*;

import static [[.ErrorsPackage]].ValidationErrorsHelpers.extractValidationErrors;

public class [[.ClassName]] {
    private static final NotFoundError NOT_FOUND_ERROR = new NotFoundError("Failed to parse url parameters");

    public static NotFoundError getNotFoundError(Throwable exception) {
        if (exception instanceof MethodArgumentTypeMismatchException) {
            var e = (MethodArgumentTypeMismatchException) exception;
            if (e.getParameter().hasParameterAnnotation(PathVariable.class)) {
                return NOT_FOUND_ERROR;
            }
        }
        return null;
    }

    public static BadRequestError getBadRequestError(Throwable exception) {
        if (exception instanceof JsonParseException) {
            var errors = extractValidationErrors((JsonParseException)exception);
            return new BadRequestError("Failed to parse body", ErrorLocation.BODY, errors);
        }
        if (exception instanceof ContentTypeMismatchException) {
            var error = new ValidationError("Content-Type", "missing", exception.getMessage());
            return new BadRequestError("Failed to parse header", ErrorLocation.HEADER, List.of(error));
        }
        if (exception instanceof MissingServletRequestParameterException) {
            var e = (MissingServletRequestParameterException) exception;
            var message = "Failed to parse query";
            var validation = new ValidationError(e.getParameterName(), "missing", e.getMessage());
            return new BadRequestError(message, ErrorLocation.QUERY, List.of(validation));
        }
        if (exception instanceof MethodArgumentTypeMismatchException) {
            var e = (MethodArgumentTypeMismatchException) exception;
            var message = "Failed to parse query";
            var validation = new ValidationError(e.getName(), "parsing_failed", e.getMessage());
            if (e.getParameter().hasParameterAnnotation(RequestParam.class)) {
                return new BadRequestError(message, ErrorLocation.QUERY, List.of(validation));
            } else if (e.getParameter().hasParameterAnnotation(RequestHeader.class)) {
                return new BadRequestError(message, ErrorLocation.HEADER, List.of(validation));
            }
        }
        if (exception instanceof MissingRequestHeaderException) {
            var e = (MissingRequestHeaderException) exception;
            var message = "Failed to parse header";
            var validation = new ValidationError(e.getHeaderName(), "missing", e.getMessage());
            return new BadRequestError(message, ErrorLocation.HEADER, List.of(validation));
        }
        return null;
    }
}
`)
	return w.ToCodeFile()
}

func springMethodParams(operation *spec.NamedOperation, types *types.Types) []string {
	methodParams := []string{"HttpServletRequest request"}

	if operation.BodyIs(spec.RequestBodyString) || operation.BodyIs(spec.RequestBodyJson) {
		methodParams = append(methodParams, "@RequestBody String bodyStr")
	}
	methodParams = append(methodParams, generateSpringMethodParam(operation.QueryParams, "RequestParam", types)...)
	methodParams = append(methodParams, generateSpringMethodParam(operation.HeaderParams, "RequestHeader", types)...)
	methodParams = append(methodParams, generateSpringMethodParam(operation.Endpoint.UrlParams, "PathVariable", types)...)

	return methodParams
}

func generateSpringMethodParam(namedParams []spec.NamedParam, paramAnnotationName string, types *types.Types) []string {
	params := []string{}

	if namedParams != nil && len(namedParams) > 0 {
		for _, param := range namedParams {
			paramAnnotation := getSpringParameterAnnotation(paramAnnotationName, &param)
			paramType := fmt.Sprintf(`%s %s`, types.Java(&param.Type.Definition), param.Name.CamelCase())
			dateFormatAnnotation := dateFormatSpringAnnotation(&param.Type.Definition)
			if dateFormatAnnotation != "" {
				params = append(params, fmt.Sprintf(`%s %s %s`, paramAnnotation, dateFormatAnnotation, paramType))
			} else {
				params = append(params, fmt.Sprintf(`%s %s`, paramAnnotation, paramType))
			}
		}
	}

	return params
}

func getSpringParameterAnnotation(paramAnnotationName string, param *spec.NamedParam) string {
	annotationParams := []string{fmt.Sprintf(`name = "%s"`, param.Name.Source)}

	if param.Type.Definition.IsNullable() {
		annotationParams = append(annotationParams, `required = false`)
	}
	if param.DefinitionDefault.Default != nil {
		annotationParams = append(annotationParams, fmt.Sprintf(`defaultValue = "%s"`, *param.DefinitionDefault.Default))
	}

	return fmt.Sprintf(`@%s(%s)`, paramAnnotationName, strings.Join(annotationParams, ", "))
}

func dateFormatSpringAnnotation(typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return dateFormatSpringAnnotationPlain(typ.Plain)
	case spec.NullableType:
		return dateFormatSpringAnnotation(typ.Child)
	case spec.ArrayType:
		return dateFormatSpringAnnotation(typ.Child)
	case spec.MapType:
		return dateFormatSpringAnnotation(typ.Child)
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func dateFormatSpringAnnotationPlain(typ string) string {
	switch typ {
	case spec.TypeDate:
		return `@DateTimeFormat(iso = DateTimeFormat.ISO.DATE)`
	case spec.TypeDateTime:
		return `@DateTimeFormat(iso = DateTimeFormat.ISO.DATE_TIME)`
	default:
		return ``
	}
}
