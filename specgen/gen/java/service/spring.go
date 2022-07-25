package service

import (
	"fmt"
	"github.com/pinzolo/casee"
	"github.com/specgen-io/specgen/spec/v2"
	"github.com/specgen-io/specgen/v2/gen/java/imports"
	"github.com/specgen-io/specgen/v2/gen/java/models"
	"github.com/specgen-io/specgen/v2/gen/java/packages"
	"github.com/specgen-io/specgen/v2/gen/java/responses"
	"github.com/specgen-io/specgen/v2/gen/java/types"
	"github.com/specgen-io/specgen/v2/gen/java/writer"
	"github.com/specgen-io/specgen/v2/generator"
	"strings"
)

var Spring = "spring"

type SpringGenerator struct {
	Types  *types.Types
	Models models.Generator
}

func NewSpringGenerator(types *types.Types, models models.Generator) *SpringGenerator {
	return &SpringGenerator{types, models}
}

func (g *SpringGenerator) ServiceImplAnnotation(api *spec.Api) (annotationImport, annotation string) {
	return `org.springframework.stereotype.Service`, fmt.Sprintf(`Service("%s")`, versionServiceName(serviceName(api), api.Http.Version))
}

func (g *SpringGenerator) ServicesControllers(version *spec.Version, mainPackage, thePackage, jsonPackage, modelsVersionPackage, serviceVersionPackage packages.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		serviceVersionSubpackage := serviceVersionPackage.Subpackage(api.Name.SnakeCase())
		files = append(files, g.serviceController(&api, thePackage, jsonPackage, modelsVersionPackage, serviceVersionSubpackage)...)
	}
	files = append(
		files,
		*g.errorsHelpers(thePackage, modelsVersionPackage),
		*g.Models.GenerateJsonParseException(thePackage, modelsVersionPackage),
		*g.contentTypeMismatchException(thePackage),
	)
	return files
}

func (g *SpringGenerator) serviceController(api *spec.Api, thePackage, jsonPackage, modelsVersionPackage, serviceVersionPackage packages.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
	w := writer.NewJavaWriter()
	w.Line(`package %s;`, thePackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(`org.apache.logging.log4j.*`)
	imports.Add(`org.springframework.beans.factory.annotation.Autowired`)
	imports.Add(`org.springframework.format.annotation.DateTimeFormat`)
	imports.Add(`org.springframework.http.*`)
	imports.Add(`org.springframework.web.bind.annotation.*`)
	imports.Add(`javax.servlet.http.HttpServletRequest`)
	imports.AddStatic(thePackage.Subpackage("ErrorsHelpers").PackageStar)
	imports.Add(modelsVersionPackage.PackageStar)
	imports.Add(serviceVersionPackage.PackageStar)
	imports.Add(g.Models.ModelsUsageImports()...)
	imports.Add(g.Types.Imports()...)
	imports.Add(`static org.apache.tomcat.util.http.fileupload.FileUploadBase.CONTENT_TYPE`)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`@RestController("%s")`, versionControllerName(controllerName(api), api.Http.Version))
	className := controllerName(api)
	w.Line(`public class %s {`, className)
	w.Line(`  private static final Logger logger = LogManager.getLogger(%s.class);`, className)
	w.EmptyLine()
	w.Line(`  @Autowired`)
	w.Line(`  private %s %s;`, serviceInterfaceName(api), serviceVarName(api))
	w.EmptyLine()
	g.Models.CreateJsonMapperField(w.Indented(), "@Autowired")
	for _, operation := range api.Operations {
		w.EmptyLine()
		g.controllerMethod(w.Indented(), &operation)
	}
	w.EmptyLine()
	g.checkContentType(w.Indented())
	w.EmptyLine()
	g.errorHandler(w.Indented(), api.Http.Errors)
	w.Line(`}`)

	files = append(files, generator.CodeFile{
		Path:    thePackage.GetPath(fmt.Sprintf("%s.java", className)),
		Content: w.String(),
	})

	return files
}

func (g *SpringGenerator) controllerMethod(w *generator.Writer, operation *spec.NamedOperation) {
	methodName := operation.Endpoint.Method
	url := operation.FullUrl()
	w.Line(`@%sMapping("%s")`, casee.ToPascalCase(methodName), url)
	w.Line(`public ResponseEntity<String> %s(%s) {`, controllerMethodName(operation), strings.Join(springMethodParams(operation, g.Types), ", "))
	w.Indent()
	w.Line(`logger.info("Received request, operationId: %s.%s, method: %s, url: %s");`, operation.Api.Name.Source, operation.Name.Source, methodName, url)
	g.parseBody(w, operation, "bodyStr", "requestBody")
	g.serviceCall(w, operation, "bodyStr", "requestBody", "result")
	g.processResponses(w, operation, "result")
	w.Unindent()
	w.Line(`}`)
}

func (g *SpringGenerator) parseBody(w *generator.Writer, operation *spec.NamedOperation, bodyStringVar, bodyJsonVar string) {
	if operation.BodyIs(spec.BodyString) {
		w.Line(`checkContentType(request, MediaType.TEXT_PLAIN);`)
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`checkContentType(request, MediaType.APPLICATION_JSON);`)
		requestBody, exception := g.Models.ReadJson(bodyStringVar, &operation.Body.Type.Definition)
		w.Line(`%s %s;`, g.Types.Java(&operation.Body.Type.Definition), bodyJsonVar)
		w.Line(`try {`)
		w.Line(`  %s = %s;`, bodyJsonVar, requestBody)
		w.Line(`} catch (%s exception) {`, exception)
		w.Line(`  throw new JsonParseException(exception);`)
		w.Line(`}`)
	}
}

func (g *SpringGenerator) serviceCall(w *generator.Writer, operation *spec.NamedOperation, bodyStringVar, bodyJsonVar, resultVarName string) {
	serviceCall := fmt.Sprintf(`%s.%s(%s)`, serviceVarName(operation.Api), operation.Name.CamelCase(), joinParams(addServiceMethodParams(operation, bodyStringVar, bodyJsonVar)))
	if len(operation.Responses) == 1 && operation.Responses[0].BodyIs(spec.BodyEmpty) {
		w.Line(`%s;`, serviceCall)
	} else {
		w.Line(`var %s = %s;`, resultVarName, serviceCall)
		w.Line(`if (%s == null) {`, resultVarName)
		w.Line(`  throw new RuntimeException("Service implementation didn't return any value'");`)
		w.Line(`}`)
	}
}

func (g *SpringGenerator) processResponses(w *generator.Writer, operation *spec.NamedOperation, resultVarName string) {
	if len(operation.Responses) == 1 {
		g.processResponse(w, &operation.Responses[0].Response, resultVarName)
	}
	if len(operation.Responses) > 1 {
		for _, response := range operation.Responses {
			w.Line(`if (result instanceof %s.%s) {`, responses.InterfaceName(operation), response.Name.PascalCase())
			g.processResponse(w.Indented(), &response.Response, responses.GetBody(&response, resultVarName))
			w.Line(`}`)
		}
		w.EmptyLine()
		w.Line(`throw new RuntimeException("Service implementation didn't return any value'");`)
	}
}

func (g *SpringGenerator) processResponse(w *generator.Writer, response *spec.Response, bodyVar string) {
	if response.BodyIs(spec.BodyEmpty) {
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s);`, response.Name.UpperCase())
		w.Line(`return new ResponseEntity<>(HttpStatus.%s);`, response.Name.UpperCase())
	}
	if response.BodyIs(spec.BodyString) {
		w.Line(`HttpHeaders headers = new HttpHeaders();`)
		w.Line(`headers.add(CONTENT_TYPE, "text/plain");`)
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s);`, response.Name.UpperCase())
		w.Line(`return new ResponseEntity<>(%s, headers, HttpStatus.%s);`, bodyVar, response.Name.UpperCase())
	}
	if response.BodyIs(spec.BodyJson) {
		w.Line(`var bodyJson = %s;`, g.Models.WriteJsonNoCheckedException(bodyVar, &response.Type.Definition))
		w.Line(`HttpHeaders headers = new HttpHeaders();`)
		w.Line(`headers.add(CONTENT_TYPE, "application/json");`)
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s);`, response.Name.UpperCase())
		w.Line(`return new ResponseEntity<>(bodyJson, headers, HttpStatus.%s);`, response.Name.UpperCase())
	}
}

func (g *SpringGenerator) badRequest(w *generator.Writer, operation *spec.NamedOperation, message string) {
	w.Line(`logger.error(%s);`, message)
	w.Line(`logger.info("Completed request with status code: {}", HttpStatus.BAD_REQUEST);`)
	w.Line(`return new ResponseEntity<>(HttpStatus.BAD_REQUEST);`)
}

func (g *SpringGenerator) internalServerError(w *generator.Writer, operation *spec.NamedOperation, message string) {
	w.Line(`logger.error(%s);`, message)
	w.Line(`logger.info("Completed request with status code: {}", HttpStatus.INTERNAL_SERVER_ERROR);`)
	w.Line(`return new ResponseEntity<>(HttpStatus.INTERNAL_SERVER_ERROR);`)
}

func (g *SpringGenerator) checkContentType(w *generator.Writer) {
	w.Lines(`
private void checkContentType(HttpServletRequest request, MediaType expectedContentType) {
	var contentType = request.getHeader("Content-Type");
	if (contentType == null || !contentType.contains(expectedContentType.toString())) {
		throw new ContentTypeMismatchException(expectedContentType.toString(), contentType);
	}
}
`)
}

func (g *SpringGenerator) errorHandler(w *generator.Writer, errors spec.Responses) {
	notFoundError := errors.GetByStatusName(spec.HttpStatusNotFound)
	badRequestError := errors.GetByStatusName(spec.HttpStatusBadRequest)
	internalServerError := errors.GetByStatusName(spec.HttpStatusInternalServerError)
	w.Line(`@ExceptionHandler(Throwable.class)`)
	w.Line(`public ResponseEntity<String> error(HttpServletRequest request, Throwable exception) {`)
	w.Line(`  var notFoundError = getNotFoundError(exception);`)
	w.Line(`  if (notFoundError != null) {`)
	g.processResponse(w.IndentedWith(2), notFoundError, "notFoundError")
	w.Line(`  }`)
	w.Line(`  var badRequestError = getBadRequestError(exception);`)
	w.Line(`  if (badRequestError != null) {`)
	g.processResponse(w.IndentedWith(2), badRequestError, "badRequestError")
	w.Line(`  }`)
	w.Line(`  var internalServerError = new InternalServerError(exception.getMessage());`)
	g.processResponse(w.IndentedWith(1), internalServerError, "internalServerError")
	w.Line(`}`)
}

func (g *SpringGenerator) contentTypeMismatchException(thePackage packages.Module) *generator.CodeFile {
	code := `
package [[.PackageName]];

public class ContentTypeMismatchException extends RuntimeException {
    String expected;
    String actual;
    public ContentTypeMismatchException(String expected, String actual) {
        super(String.format("Expected Content-Type header: '%s' was not provided, found: '%s'", expected, actual));
    }
}
`
	code, _ = generator.ExecuteTemplate(code, struct {
		PackageName string
	}{
		thePackage.PackageName,
	})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("ContentTypeMismatchException.java"),
		Content: strings.TrimSpace(code),
	}
}

func (g *SpringGenerator) errorsHelpers(thePackage, modelsPackage packages.Module) *generator.CodeFile {
	code := `
package [[.PackageName]];

import org.springframework.web.bind.MissingRequestHeaderException;
import org.springframework.web.bind.MissingServletRequestParameterException;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestHeader;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.method.annotation.MethodArgumentTypeMismatchException;

import java.util.List;

import [[.ModelsPackage]].*;

public class ErrorsHelpers {
    private static NotFoundError NOT_FOUND_ERROR = new NotFoundError("Failed to parse url parameters");

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
            var e = (JsonParseException) exception;
            return new BadRequestError("Failed to parse body", ErrorLocation.BODY, e.getErrors());
        }
        if (exception instanceof ContentTypeMismatchException) {
            var error = new ValidationError("Content-Type", "missing", exception.getMessage());
            return new BadRequestError("Failed to parse header", ErrorLocation.HEADER, List.of(error));

        }
        if (exception instanceof MissingServletRequestParameterException) {
            var e = (MissingServletRequestParameterException) exception;
            var message = String.format("Failed to parse query");
            var validation = new ValidationError(e.getParameterName(), "missing", e.getMessage());
            return new BadRequestError(message, ErrorLocation.QUERY, List.of(validation));
        }
        if (exception instanceof MethodArgumentTypeMismatchException) {
            var e = (MethodArgumentTypeMismatchException) exception;
            var message = String.format("Failed to parse query");
            var validation = new ValidationError(e.getName(), "parsing_failed", e.getMessage());
            if (e.getParameter().hasParameterAnnotation(RequestParam.class)) {
                return new BadRequestError(message, ErrorLocation.QUERY, List.of(validation));
            } else if (e.getParameter().hasParameterAnnotation(RequestHeader.class)) {
                return new BadRequestError(message, ErrorLocation.HEADER, List.of(validation));
            }
        }
        if (exception instanceof MissingRequestHeaderException) {
            var e = (MissingRequestHeaderException) exception;
            var message = String.format("Failed to parse header");
            var validation = new ValidationError(e.getHeaderName(), "missing", e.getMessage());
            return new BadRequestError(message, ErrorLocation.HEADER, List.of(validation));
        }
        return null;
    }
}
`

	code, _ = generator.ExecuteTemplate(code, struct {
		PackageName   string
		ModelsPackage string
	}{thePackage.PackageName, modelsPackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("ErrorsHelpers.java"),
		Content: strings.TrimSpace(code),
	}
}

func springMethodParams(operation *spec.NamedOperation, types *types.Types) []string {
	methodParams := []string{"HttpServletRequest request"}

	if operation.Body != nil {
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

	return fmt.Sprintf(`@%s(%s)`, paramAnnotationName, joinParams(annotationParams))
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
