package service

import (
	"fmt"
	"github.com/pinzolo/casee"
	"github.com/specgen-io/specgen/v2/gen/java/imports"
	"github.com/specgen-io/specgen/v2/gen/java/models"
	"github.com/specgen-io/specgen/v2/gen/java/packages"
	"github.com/specgen-io/specgen/v2/gen/java/responses"
	"github.com/specgen-io/specgen/v2/gen/java/types"
	"github.com/specgen-io/specgen/v2/gen/java/writer"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

var Micronaut = "micronaut"

type MicronautGenerator struct {
	Types  *types.Types
	Models *models.JacksonGenerator
}

func NewMicronautGenerator(types *types.Types, models *models.JacksonGenerator) *MicronautGenerator {
	return &MicronautGenerator{types, models}
}

func (g *MicronautGenerator) ServiceImplAnnotation(api *spec.Api) (annotationImport, annotation string) {
	return `io.micronaut.context.annotation.Bean`, `Bean`
}

func (g *MicronautGenerator) ServicesControllers(version *spec.Version, mainPackage, thePackage, jsonPackage, modelsVersionPackage, serviceVersionPackage packages.Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, api := range version.Http.Apis {
		serviceVersionSubpackage := serviceVersionPackage.Subpackage(api.Name.SnakeCase())
		files = append(files, g.serviceController(&api, thePackage, jsonPackage, modelsVersionPackage, serviceVersionSubpackage)...)
	}
	files = append(files, *g.errorsHelpers(thePackage, modelsVersionPackage), *g.badRequestException(thePackage, modelsVersionPackage))
	files = append(files, dateConverters(mainPackage)...)
	return files
}

func (g *MicronautGenerator) serviceController(api *spec.Api, apiPackage, jsonPackage, modelsVersionPackage, serviceVersionPackage packages.Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	w := writer.NewJavaWriter()
	w.Line(`package %s;`, apiPackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(`org.slf4j.*`)
	imports.Add(`io.micronaut.core.convert.format.Format`)
	imports.Add(`io.micronaut.core.annotation.Nullable`)
	imports.Add(`io.micronaut.http.*`)
	imports.Add(`io.micronaut.http.annotation.*`)
	imports.Add(`io.micronaut.http.annotation.Error`)
	imports.Add(`jakarta.inject.Inject`)

	imports.AddStatic(apiPackage.Subpackage("ErrorsHelpers").PackageStar)
	imports.Add(modelsVersionPackage.PackageStar)
	imports.Add(serviceVersionPackage.PackageStar)
	imports.Add(g.Models.JsonImports()...)
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
	g.Models.CreateJsonMapperField(w.Indented())

	for _, operation := range api.Operations {
		w.EmptyLine()
		g.controllerMethod(w.Indented(), &operation)
	}
	w.EmptyLine()
	g.errorHandler(w.Indented(), api.Http.Errors)
	w.EmptyLine()
	g.responseHelpers(w.Indented())
	w.Line(`}`)

	files = append(files, sources.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.java", className)),
		Content: w.String(),
	})

	return files
}

func (g *MicronautGenerator) controllerMethod(w *sources.Writer, operation *spec.NamedOperation) {
	if operation.BodyIs(spec.BodyString) {
		w.Line(`@Consumes(MediaType.TEXT_PLAIN)`)
	}
	methodName := operation.Endpoint.Method
	url := operation.FullUrl()
	w.Line(`@%s("%s")`, casee.ToPascalCase(methodName), url)
	w.Line(`public HttpResponse<?> %s(%s) {`, controllerMethodName(operation), joinParams(addMicronautMethodParams(operation, g.Types)))
	w.Indent()
	w.Line(`logger.info("Received request, operationId: %s.%s, method: %s, url: %s");`, operation.Api.Name.Source, operation.Name.Source, methodName, url)
	g.parseBody(w, operation, "bodyStr", "requestBody")
	g.serviceCall(w, operation, "bodyStr", "requestBody", "result")
	g.processResponses(w, operation, "result")
	w.Unindent()
	w.Line(`}`)
}

func (g *MicronautGenerator) parseBody(w *sources.Writer, operation *spec.NamedOperation, bodyStringVar, bodyJsonVar string) {
	if operation.BodyIs(spec.BodyJson) {
		requestBody, exception := g.Models.ReadJson(bodyStringVar, &operation.Body.Type.Definition)
		w.Line(`%s %s;`, g.Types.Java(&operation.Body.Type.Definition), bodyJsonVar)
		w.Line(`try {`)
		w.Line(`  %s = %s;`, bodyJsonVar, requestBody)
		w.Line(`} catch (%s e) {`, exception)
		w.Line(`  throw new BadRequestException(jacksonBodyBadRequestError(e));`)
		w.Line(`}`)
	}
}

func (g *MicronautGenerator) serviceCall(w *sources.Writer, operation *spec.NamedOperation, bodyStringVar, bodyJsonVar, responseVarName string) {
	serviceCall := fmt.Sprintf(`%s.%s(%s)`, serviceVarName(operation.Api), operation.Name.CamelCase(), joinParams(addServiceMethodParams(operation, bodyStringVar, bodyJsonVar)))
	if len(operation.Responses) == 1 && operation.Responses[0].BodyIs(spec.BodyEmpty) {
		w.Line(`%s;`, serviceCall)
	} else {
		w.Line(`var %s = %s;`, responseVarName, serviceCall)
		w.Line(`if (%s == null) {`, responseVarName)
		w.Line(`  throw new RuntimeException("Service implementation didn't return any value'");`)
		w.Line(`}`)
	}
}

func (g *MicronautGenerator) processResponses(w *sources.Writer, operation *spec.NamedOperation, responseVarName string) {
	if len(operation.Responses) == 1 {
		g.processResponse(w, &operation.Responses[0].Response, "result")
	}
	if len(operation.Responses) > 1 {
		for _, response := range operation.Responses {
			w.Line(`if (%s instanceof %s.%s) {`, responseVarName, responses.InterfaceName(operation), response.Name.PascalCase())
			g.processResponse(w.Indented(), &response.Response, responses.GetBody(&response, responseVarName))
			w.Line(`}`)
		}
		w.EmptyLine()
		w.Line(`throw new RuntimeException("Service implementation didn't return any value'");`)
	}
}

func (g *MicronautGenerator) processResponse(w *sources.Writer, response *spec.Response, bodyVar string) {
	if response.BodyIs(spec.BodyEmpty) {
		w.Line(`return responseEmpty(HttpStatus.%s);`, response.Name.UpperCase())
	}
	if response.BodyIs(spec.BodyString) {
		w.Line(`return responseText(HttpStatus.%s, %s);`, response.Name.UpperCase(), bodyVar)
	}
	if response.BodyIs(spec.BodyJson) {
		w.Line(`return responseJson(HttpStatus.%s, %s);`, response.Name.UpperCase(), bodyVar)
	}
}

func (g *MicronautGenerator) errorHandler(w *sources.Writer, errors spec.Responses) {
	notFoundError := errors.GetByStatusName(spec.HttpStatusNotFound)
	badRequestError := errors.GetByStatusName(spec.HttpStatusBadRequest)
	internalServerError := errors.GetByStatusName(spec.HttpStatusInternalServerError)
	w.Line(`@Error(global = false)`)
	w.Line(`public HttpResponse<?> error(HttpRequest request, Throwable exception) {`)
	w.Line(`  var notFoundError = getNotFoundError(exception);`)
	w.Line(`  if (notFoundError != null) {`)
	w.Line(`    return responseJson(HttpStatus.%s, %s);`, notFoundError.Name.UpperCase(), `notFoundError`)
	w.Line(`  }`)
	w.Line(`  var badRequestError = getBadRequestError(exception);`)
	w.Line(`  if (badRequestError != null) {`)
	w.Line(`    return responseJson(HttpStatus.%s, %s);`, badRequestError.Name.UpperCase(), `badRequestError`)
	w.Line(`  }`)
	w.Line(`  var internalServerError = new InternalServerError(exception.getMessage());`)
	w.Line(`  return responseJson(HttpStatus.%s, %s);`, internalServerError.Name.UpperCase(), `internalServerError`)
	w.Line(`}`)
}

func (g *MicronautGenerator) responseHelpers(w *sources.Writer) {
	w.Lines(`
private HttpResponse<?> responseEmpty(HttpStatus status) {
  logger.info("Completed request with status code: {}", status);
  return HttpResponse.status(status);
}

private HttpResponse<?> responseText(HttpStatus status, String body) {
  logger.info("Completed request with status code: {}", HttpStatus.OK);
  return HttpResponse.status(HttpStatus.OK).body(body).contentType("text/plain");
}

private HttpResponse<?> responseJson(HttpStatus status, Object body) {
  try {
    var bodyJson = objectMapper.writeValueAsString(body);
    logger.info("Completed request with status code: {}", status);
    return HttpResponse.status(status).body(bodyJson).contentType("application/json");
  } catch (Exception e) {
    throw new RuntimeException("Failed to serialize response body: "+e.getMessage(), e);
  }
}
`)
}

func (g *MicronautGenerator) badRequestException(thePackage, modelsPackage packages.Module) *sources.CodeFile {
	code := `
package [[.PackageName]];

import [[.ModelsPackage]].BadRequestError;

public class BadRequestException extends RuntimeException {
    private BadRequestError error;
    public BadRequestError getError() {
        return error;
    }

    public BadRequestException(BadRequestError error) {
        super (error.getMessage());
        this.error = error;
    }
}
`

	code, _ = sources.ExecuteTemplate(code, struct {
		PackageName   string
		ModelsPackage string
	}{thePackage.PackageName, modelsPackage.PackageName})
	return &sources.CodeFile{
		Path:    thePackage.GetPath("BadRequestException.java"),
		Content: strings.TrimSpace(code),
	}
}

func (g *MicronautGenerator) errorsHelpers(thePackage, modelsPackage packages.Module) *sources.CodeFile {
	code := `
package [[.PackageName]];

import io.micronaut.core.convert.exceptions.ConversionErrorException;
import io.micronaut.core.type.Argument;
import io.micronaut.web.router.exceptions.*;

import javax.validation.ConstraintViolationException;
import com.fasterxml.jackson.databind.exc.InvalidFormatException;
import java.io.IOException;

import [[.ModelsPackage]].*;

public class ErrorsHelpers {
    private static NotFoundError NOT_FOUND_ERROR = new NotFoundError("Failed to parse url parameters");

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
        var header = argument.getAnnotationMetadata().findDeclaredAnnotation("io.micronaut.http.annotation.Header");
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
        var header = argument.getAnnotationMetadata().findDeclaredAnnotation("io.micronaut.http.annotation.Header");
        if (header.isPresent()) {
            return header.get().getValues().get("value").toString();
        }
        return "unknown";
    }

    private static BadRequestError argumentBadRequestError(Argument<?> arg, String errorMessage, String code) {
        var location = getLocation(arg);
        var parameterName = getParameterName(arg);
        var validation = new ValidationError(parameterName, code, errorMessage);
        var message = String.format("Failed to parse %s", location.name().toLowerCase());
        return new BadRequestError(message, location, new ValidationError[]{validation});
    }

    private static String getJsonPath(InvalidFormatException exception) {
        var path = new StringBuilder("$");
        for (var reference: exception.getPath()) {
            if (reference.getIndex() != -1) {
                path.append("[").append(reference.getIndex()).append("]");
            } else {
                path.append(".").append(reference.getFieldName());
            }
        }
        return path.toString();
    }

    public static BadRequestError jacksonBodyBadRequestError(IOException exception) {
        var location = ErrorLocation.BODY;
        var message = "Failed to parse body";
        ValidationError[] errors = null;
        if (exception instanceof InvalidFormatException) {
            var jsonPath = getJsonPath((InvalidFormatException)exception);
            var validation = new ValidationError(jsonPath, "missing", exception.getMessage());
            errors = new ValidationError[]{validation};
        }
        return new BadRequestError(message, location, errors);
    }

    public static BadRequestError getBadRequestError(Throwable exception) {
        if (exception instanceof BadRequestException) {
            return ((BadRequestException)exception).getError();
        }
        if (exception instanceof UnsatisfiedRouteException) {
            var e = (UnsatisfiedRouteException) exception;
            return argumentBadRequestError(e.getArgument(), e.getMessage(), "missing");
        }
        if (exception instanceof ConversionErrorException) {
            var e = (ConversionErrorException) exception;
            return argumentBadRequestError(e.getArgument(), e.getMessage(), "wrong_format");
        }
        if (exception instanceof ConstraintViolationException) {
            var message = "Failed to parse params";
            return new BadRequestError(message, ErrorLocation.BODY, null);
        }
        return null;
    }
}
`

	code, _ = sources.ExecuteTemplate(code, struct {
		PackageName   string
		ModelsPackage string
	}{thePackage.PackageName, modelsPackage.PackageName})
	return &sources.CodeFile{
		Path:    thePackage.GetPath("ErrorsHelpers.java"),
		Content: strings.TrimSpace(code),
	}
}

func addMicronautMethodParams(operation *spec.NamedOperation, types *types.Types) []string {
	methodParams := []string{}

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

func dateConverters(thePackage packages.Module) []sources.CodeFile {
	convertersPackage := thePackage.Subpackage("converters")

	files := []sources.CodeFile{}
	files = append(files, *localDateConverter(convertersPackage))
	files = append(files, *localDateTimeConverter(convertersPackage))
	return files
}

func localDateConverter(thePackage packages.Module) *sources.CodeFile {
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

	code, _ = sources.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &sources.CodeFile{
		Path:    thePackage.GetPath("LocalDateConverter.java"),
		Content: strings.TrimSpace(code),
	}
}

func localDateTimeConverter(thePackage packages.Module) *sources.CodeFile {
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

	code, _ = sources.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &sources.CodeFile{
		Path:    thePackage.GetPath("LocalDateTimeConverter.java"),
		Content: strings.TrimSpace(code),
	}
}
