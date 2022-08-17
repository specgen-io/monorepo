package service

import (
	"fmt"
	"strings"

	"generator"

	"github.com/pinzolo/casee"
	"java/imports"
	"java/models"
	"java/packages"
	"java/responses"
	"java/types"
	"java/writer"
	"spec"
)

var Micronaut = "micronaut"

type MicronautGenerator struct {
	Types  *types.Types
	Models models.Generator
}

func NewMicronautGenerator(types *types.Types, models models.Generator) *MicronautGenerator {
	return &MicronautGenerator{types, models}
}

func (g *MicronautGenerator) ServiceImplAnnotation(api *spec.Api) (annotationImport, annotation string) {
	return `io.micronaut.context.annotation.Bean`, `Bean`
}

func (g *MicronautGenerator) ServicesControllers(version *spec.Version, mainPackage, thePackage, contentTypePackage, jsonPackage, modelsVersionPackage, serviceVersionPackage packages.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		serviceVersionSubpackage := serviceVersionPackage.Subpackage(api.Name.SnakeCase())
		files = append(files, g.serviceController(&api, thePackage, contentTypePackage, jsonPackage, modelsVersionPackage, serviceVersionSubpackage)...)
	}
	files = append(files, dateConverters(mainPackage)...)
	return files
}

func (g *MicronautGenerator) ServiceImports() []string {
	return []string{
		`org.slf4j.*`,
		`io.micronaut.http.*`,
		`io.micronaut.http.annotation.*`,
	}
}

func (g *MicronautGenerator) ExceptionController(responses *spec.Responses, thePackage, errorsPackage, errorsModelsPackage, jsonPackage packages.Module) *generator.CodeFile {
	w := writer.NewJavaWriter()
	w.Line(`package %s;`, thePackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(g.ServiceImports()...)
	imports.Add(`io.micronaut.http.annotation.Error`)
	imports.Add(`jakarta.inject.Inject`)
	imports.Add(jsonPackage.PackageStar)
	imports.Add(errorsModelsPackage.PackageStar)
	imports.AddStatic(errorsPackage.Subpackage("ErrorsHelpers").PackageStar)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`@Controller`)
	className := `ExceptionController`
	w.Line(`public class %s {`, className)
	w.Line(`  private static final Logger logger = LoggerFactory.getLogger(%s.class);`, className)
	w.EmptyLine()
	w.Line(`  @Inject`)
	w.Line(`  private Json json;`)
	w.EmptyLine()
	g.errorHandler(w.Indented(), *responses)
	w.Line(`}`)

	return &generator.CodeFile{
		Path:    thePackage.GetPath(fmt.Sprintf("%s.java", className)),
		Content: w.String(),
	}
}

func (g *MicronautGenerator) errorHandler(w *generator.Writer, errors spec.Responses) {
	notFoundError := errors.GetByStatusName(spec.HttpStatusNotFound)
	badRequestError := errors.GetByStatusName(spec.HttpStatusBadRequest)
	internalServerError := errors.GetByStatusName(spec.HttpStatusInternalServerError)
	w.Line(`@Error(global = true)`)
	w.Line(`public HttpResponse<?> error(Throwable exception) {`)
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

func (g *MicronautGenerator) serviceController(api *spec.Api, thePackage, contentTypePackage, jsonPackage, modelsVersionPackage, serviceVersionPackage packages.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
	w := writer.NewJavaWriter()
	w.Line(`package %s;`, thePackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(g.ServiceImports()...)
	imports.Add(`io.micronaut.core.annotation.Nullable`)
	imports.Add(`jakarta.inject.Inject`)
	imports.Add(contentTypePackage.PackageStar)
	imports.Add(jsonPackage.PackageStar)
	imports.Add(modelsVersionPackage.PackageStar)
	imports.Add(serviceVersionPackage.PackageStar)
	imports.Add(g.Models.ModelsUsageImports()...)
	imports.Add(g.Types.Imports()...)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`@Controller`)
	className := controllerName(api)
	w.Line(`public class %s {`, className)
	w.Line(`  private static final Logger logger = LoggerFactory.getLogger(%s.class);`, className)
	w.EmptyLine()
	w.Line(`  @Inject`)
	w.Line(`  private %s %s;`, serviceInterfaceName(api), serviceVarName(api))
	w.EmptyLine()
	w.Line(`  @Inject`)
	w.Line(`  private Json json;`)

	for _, operation := range api.Operations {
		w.EmptyLine()
		g.controllerMethod(w.Indented(), &operation)
	}
	w.Line(`}`)

	files = append(files, generator.CodeFile{
		Path:    thePackage.GetPath(fmt.Sprintf("%s.java", className)),
		Content: w.String(),
	})

	return files
}

func (g *MicronautGenerator) controllerMethod(w *generator.Writer, operation *spec.NamedOperation) {
	if operation.BodyIs(spec.BodyString) {
		w.Line(`@Consumes(MediaType.TEXT_PLAIN)`)
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`@Consumes(MediaType.APPLICATION_JSON)`)
	}
	methodName := operation.Endpoint.Method
	url := operation.FullUrl()
	w.Line(`@%s("%s")`, casee.ToPascalCase(methodName), url)
	w.Line(`public HttpResponse<?> %s(%s) {`, controllerMethodName(operation), strings.Join(micronautMethodParams(operation, g.Types), ", "))
	w.Indent()
	w.Line(`logger.info("Received request, operationId: %s.%s, method: %s, url: %s");`, operation.Api.Name.Source, operation.Name.Source, methodName, url)
	g.parseBody(w, operation, "bodyStr", "requestBody")
	g.serviceCall(w, operation, "bodyStr", "requestBody", "result")
	g.processResponses(w, operation, "result")
	w.Unindent()
	w.Line(`}`)
}

func (g *MicronautGenerator) parseBody(w *generator.Writer, operation *spec.NamedOperation, bodyStringVar, bodyJsonVar string) {
	if operation.BodyIs(spec.BodyString) {
		w.Line(`ContentType.check(request, MediaType.TEXT_PLAIN);`)
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`ContentType.check(request, MediaType.APPLICATION_JSON);`)
		w.Line(`%s %s = json.read(%s);`, g.Types.Java(&operation.Body.Type.Definition), bodyJsonVar, g.Models.JsonRead(bodyStringVar, &operation.Body.Type.Definition))
	}
}

func (g *MicronautGenerator) serviceCall(w *generator.Writer, operation *spec.NamedOperation, bodyStringVar, bodyJsonVar, resultVarName string) {
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

func (g *MicronautGenerator) processResponses(w *generator.Writer, operation *spec.NamedOperation, resultVarName string) {
	if len(operation.Responses) == 1 {
		g.processResponse(w, &operation.Responses[0].Response, resultVarName)
	}
	if len(operation.Responses) > 1 {
		for _, response := range operation.Responses {
			w.Line(`if (%s instanceof %s.%s) {`, resultVarName, responses.InterfaceName(operation), response.Name.PascalCase())
			g.processResponse(w.Indented(), &response.Response, responses.GetBody(&response, resultVarName))
			w.Line(`}`)
		}
		w.EmptyLine()
		w.Line(`throw new RuntimeException("Service implementation didn't return any value'");`)
	}
}

func (g *MicronautGenerator) processResponse(w *generator.Writer, response *spec.Response, bodyVar string) {
	if response.BodyIs(spec.BodyEmpty) {
		w.Line(`logger.info("Completed request with status code: HttpStatus.%s");`, response.Name.UpperCase())
		w.Line(`return HttpResponse.status(HttpStatus.%s);`, response.Name.UpperCase())
	}
	if response.BodyIs(spec.BodyString) {
		w.Line(`logger.info("Completed request with status code: HttpStatus.%s");`, response.Name.UpperCase())
		w.Line(`return HttpResponse.status(HttpStatus.%s).body(%s).contentType("text/plain");`, response.Name.UpperCase(), bodyVar)
	}
	if response.BodyIs(spec.BodyJson) {
		w.Line(`var bodyJson = json.write(%s);`, g.Models.JsonWrite(bodyVar, &response.Type.Definition))
		w.Line(`logger.info("Completed request with status code: HttpStatus.%s");`, response.Name.UpperCase())
		w.Line(`return HttpResponse.status(HttpStatus.%s).body(bodyJson).contentType("application/json");`, response.Name.UpperCase())
	}
}

func (g *MicronautGenerator) ContentType(thePackage packages.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *g.contentTypeMismatchException(thePackage))
	files = append(files, *g.checkContentType(thePackage))
	return files
}

func (g *MicronautGenerator) contentTypeMismatchException(thePackage packages.Module) *generator.CodeFile {
	code := `
package [[.PackageName]];

public class ContentTypeMismatchException extends RuntimeException {
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

func (g *MicronautGenerator) checkContentType(thePackage packages.Module) *generator.CodeFile {
	code := `
package [[.PackageName]];

import io.micronaut.http.HttpRequest;

public class ContentType {
	public static void check(HttpRequest<?> request, String expectedContentType) {
		var contentType = request.getHeaders().get("Content-Type");
		if (contentType == null || !contentType.contains(expectedContentType)) {
			throw new ContentTypeMismatchException(expectedContentType, contentType);
		}
	}
}
`
	code, _ = generator.ExecuteTemplate(code, struct {
		PackageName string
	}{
		thePackage.PackageName,
	})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("ContentType.java"),
		Content: strings.TrimSpace(code),
	}
}

func (g *MicronautGenerator) Errors(thePackage, errorsModelsPackage, contentTypePackage, jsonPackage packages.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *g.errorsHelpers(thePackage, errorsModelsPackage, contentTypePackage, jsonPackage))
	files = append(files, *g.Models.ValidationErrorsHelpers(thePackage, errorsModelsPackage, jsonPackage))
	return files
}

func (g *MicronautGenerator) errorsHelpers(thePackage, errorsModelsPackage, contentTypePackage, jsonPackage packages.Module) *generator.CodeFile {
	code := `
package [[.PackageName]];

import io.micronaut.core.annotation.AnnotationValue;
import io.micronaut.core.convert.exceptions.ConversionErrorException;
import io.micronaut.core.type.Argument;
import io.micronaut.web.router.exceptions.*;

import javax.validation.ConstraintViolationException;
import java.util.List;

import [[.ErrorsModelsPackage]].*;
import [[.ContentTypePackage]].*;
import [[.JsonPackage]].*;

import static [[.PackageName]].ValidationErrorsHelpers.extractValidationErrors;

public class ErrorsHelpers {
    private static final NotFoundError NOT_FOUND_ERROR = new NotFoundError("Failed to parse url parameters");

    public static NotFoundError getNotFoundError(Throwable exception) {
		if (exception instanceof UnsatisfiedPathVariableRouteException) {
			return NOT_FOUND_ERROR;
		}
		if (exception instanceof UnsatisfiedPartRouteException) {
			return NOT_FOUND_ERROR;
		}
		if (exception instanceof ConversionErrorException) {
			var e = (ConversionErrorException) exception;
			var annotation = e.getArgument().getAnnotationMetadata().findDeclaredAnnotation("io.micronaut.http.annotation.PathVariable");
			if (annotation.isPresent()) {
				return NOT_FOUND_ERROR;
			}
		}
		return null;
	}

	private static ErrorLocation getLocation(Argument<?> argument) {
		var query = argument.getAnnotationMetadata().findDeclaredAnnotation("io.micronaut.http.annotation.QueryValue");
		if (query.isPresent()) {
			var parameterName = query.get().getValues().get("value");
			return ErrorLocation.QUERY;
		}
		var header = argument.getAnnotationMetadata().findDeclaredAnnotation("io.micronaut.http.annotation.Headers");
		if (header.isPresent()) {
			return ErrorLocation.HEADER;
		}
		return ErrorLocation.BODY;
	}

	private static String getParameterName(Argument<?> argument) {
		var query = argument.getAnnotationMetadata().findDeclaredAnnotation("io.micronaut.http.annotation.QueryValue");
		if (query.isPresent()) {
			return query.get().getValues().get("value").toString();
		}
		var header = argument.getAnnotationMetadata().findDeclaredAnnotation("io.micronaut.http.annotation.Headers");
		if (header.isPresent()) {
			var annotationValues = (AnnotationValue[]) header.get().getValues().get("value");
			return annotationValues[0].getValues().get("value").toString();
		}
		return "unknown";	
	}

	private static BadRequestError argumentBadRequestError(Argument<?> arg, String errorMessage, String code) {
		var location = getLocation(arg);
		var parameterName = getParameterName(arg);	
		var validation = new ValidationError(parameterName, code, errorMessage);
		var message = String.format("Failed to parse %s", location.name().toLowerCase());
		return new BadRequestError(message, location, List.of(validation));
	}

	public static BadRequestError getBadRequestError(Throwable exception) {
		if (exception instanceof JsonParseException) {
			var errors = extractValidationErrors((JsonParseException) exception);
			return new BadRequestError("Failed to parse body", ErrorLocation.BODY, errors);
		}
		if (exception instanceof ContentTypeMismatchException) {
			var error = new ValidationError("Content-Type", "missing", exception.getMessage());
			return new BadRequestError("Failed to parse header", ErrorLocation.HEADER, List.of(error));
		}
		if (exception instanceof UnsatisfiedRouteException) {
			var e = (UnsatisfiedRouteException) exception;
			return argumentBadRequestError(e.getArgument(), e.getMessage(), "missing");
		}
		if (exception instanceof ConversionErrorException) {
			var e = (ConversionErrorException) exception;
			return argumentBadRequestError(e.getArgument(), e.getMessage(), "parsing_failed");
		}
		if (exception instanceof ConstraintViolationException) {
			var message = "Failed to parse body";
			return new BadRequestError(message, ErrorLocation.BODY, null);
		}
		return null;
	}
}
`

	code, _ = generator.ExecuteTemplate(code, struct {
		PackageName         string
		ErrorsModelsPackage string
		ContentTypePackage  string
		JsonPackage         string
	}{
		thePackage.PackageName,
		errorsModelsPackage.PackageName,
		contentTypePackage.PackageName,
		jsonPackage.PackageName,
	})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("ErrorsHelpers.java"),
		Content: strings.TrimSpace(code),
	}
}

func (g *MicronautGenerator) JsonHelpers(thePackage packages.Module) []generator.CodeFile {
	files := []generator.CodeFile{}

	files = append(files, *g.Json(thePackage))
	files = append(files, *g.Models.JsonParseException(thePackage))
	files = append(files, g.Models.SetupLibrary(thePackage)...)

	return files
}

func (g *MicronautGenerator) Json(thePackage packages.Module) *generator.CodeFile {
	w := writer.NewJavaWriter()
	w.Line(`package %s;`, thePackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(g.Models.ModelsUsageImports()...)
	imports.Add(`jakarta.inject.*`)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`@Singleton`)
	className := `Json`
	w.Line(`public class %s {`, className)
	w.EmptyLine()
	g.Models.CreateJsonMapperField(w.Indented(), "@Inject")
	w.Line(g.Models.JsonHelpersMethods())
	w.Line(`}`)

	return &generator.CodeFile{
		Path:    thePackage.GetPath(fmt.Sprintf("%s.java", className)),
		Content: w.String(),
	}
}

func micronautMethodParams(operation *spec.NamedOperation, types *types.Types) []string {
	methodParams := []string{"HttpRequest<?> request"}

	if operation.Body != nil {
		methodParams = append(methodParams, "@Body String bodyStr")
	}

	methodParams = append(methodParams, generateMicronautMethodParam(operation.QueryParams, "QueryValue", types)...)
	methodParams = append(methodParams, generateMicronautMethodParam(operation.HeaderParams, "Header", types)...)
	methodParams = append(methodParams, generateMicronautMethodParam(operation.Endpoint.UrlParams, "PathVariable", types)...)

	return methodParams
}

func generateMicronautMethodParam(namedParams []spec.NamedParam, paramAnnotation string, types *types.Types) []string {
	params := []string{}

	if namedParams != nil && len(namedParams) > 0 {
		for _, param := range namedParams {
			paramType := fmt.Sprintf(`%s %s`, types.Java(&param.Type.Definition), param.Name.CamelCase())
			if param.Type.Definition.IsNullable() {
				paramType = fmt.Sprintf(`@Nullable %s`, paramType)
			}
			params = append(params, fmt.Sprintf(`%s %s`, getMicronautParameterAnnotation(paramAnnotation, &param), paramType))
		}
	}

	return params
}

func getMicronautParameterAnnotation(paramAnnotation string, param *spec.NamedParam) string {
	annotationParams := []string{fmt.Sprintf(`value = "%s"`, param.Name.Source)}

	if param.DefinitionDefault.Default != nil {
		annotationParams = append(annotationParams, fmt.Sprintf(`defaultValue = "%s"`, *param.DefinitionDefault.Default))
	}

	return fmt.Sprintf(`@%s(%s)`, paramAnnotation, joinParams(annotationParams))
}

func dateConverters(thePackage packages.Module) []generator.CodeFile {
	convertersPackage := thePackage.Subpackage("converters")

	files := []generator.CodeFile{}
	files = append(files, *localDateConverter(convertersPackage))
	files = append(files, *localDateTimeConverter(convertersPackage))
	return files
}

func localDateConverter(thePackage packages.Module) *generator.CodeFile {
	code := `
package [[.PackageName]];

import io.micronaut.core.convert.*;
import jakarta.inject.Singleton;

import java.time.LocalDate;
import java.time.format.*;
import java.util.Optional;

@Singleton
public class LocalDateConverter implements TypeConverter<String, LocalDate> {
    @Override
    public Optional<LocalDate> convert(String value, Class<LocalDate> targetType, ConversionContext context) {
        try {
            DateTimeFormatter formatter = DateTimeFormatter.ofPattern("yyyy-MM-dd", context.getLocale());
            LocalDate result = LocalDate.parse(value, formatter);
            return Optional.of(result);
        } catch (DateTimeParseException e) {
            context.reject(value, e);
            return Optional.empty();
        }
    }
}
`

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("LocalDateConverter.java"),
		Content: strings.TrimSpace(code),
	}
}

func localDateTimeConverter(thePackage packages.Module) *generator.CodeFile {
	code := `
package [[.PackageName]];

import io.micronaut.core.convert.*;
import jakarta.inject.Singleton;

import java.time.LocalDateTime;
import java.time.format.*;
import java.util.Optional;

@Singleton
public class LocalDateTimeConverter implements TypeConverter<String, LocalDateTime> {
    @Override
    public Optional<LocalDateTime> convert(String value, Class<LocalDateTime> targetType, ConversionContext context) {
        try {
            DateTimeFormatter formatter = DateTimeFormatter.ofPattern("yyyy-MM-dd'T'HH:mm:ss", context.getLocale());
            LocalDateTime result = LocalDateTime.parse(value, formatter);
            return Optional.of(result);
        } catch (DateTimeParseException e) {
            context.reject(value, e);
            return Optional.empty();
        }
    }
}
`

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("LocalDateTimeConverter.java"),
		Content: strings.TrimSpace(code),
	}
}
