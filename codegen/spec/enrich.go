package spec

import (
	"errors"
	"fmt"
	"gopkg.in/specgen-io/yaml.v3"
)

func enrich(specification *Spec) (*Messages, error) {
	messages := NewMessages()
	addRequiredHttpErrors(specification, messages)
	if specification.HttpErrors != nil {
		httpErrors := specification.HttpErrors
		httpErrors.InSpec = specification
		httpErrors.ResolvedModels = enrichModels(httpErrors.Models, messages)
		errorModels := buildModelsMap(httpErrors.Models)
		enricher := &httpEnricher{errorModels, messages}
		enricher.httpErrors(httpErrors)
	}
	for index := range specification.Versions {
		version := &specification.Versions[index]
		version.InSpec = specification
		version.ResolvedModels = enrichModels(version.Models, messages)
		models := buildModelsMap(version.Models)
		if specification.HttpErrors != nil {
			errorModels := buildModelsMap(specification.HttpErrors.Models)
			for name, model := range errorModels {
				models[name] = model
			}
		}
		enricher := &httpEnricher{models, messages}
		enricher.version(version)
	}
	if messages.ContainsLevel(LevelError) {
		return messages, errors.New("failed to parse specification")
	}
	return messages, nil
}

func addRequiredHttpErrors(specification *Spec, messages *Messages) {
	if specification.HttpErrors != nil {
		hasRequiredErrorModels :=
			errorModelShouldNotBeDeclared(specification.HttpErrors, messages, BadRequestError) ||
				errorModelShouldNotBeDeclared(specification.HttpErrors, messages, ValidationError) ||
				errorModelShouldNotBeDeclared(specification.HttpErrors, messages, ErrorLocation) ||
				errorModelShouldNotBeDeclared(specification.HttpErrors, messages, InternalServerError) ||
				errorModelShouldNotBeDeclared(specification.HttpErrors, messages, NotFoundError)

		hasRequiredErrorResponses :=
			errorResponseShouldNotBeDeclared(specification.HttpErrors, messages, HttpStatusBadRequest) ||
				errorResponseShouldNotBeDeclared(specification.HttpErrors, messages, HttpStatusNotFound) ||
				errorResponseShouldNotBeDeclared(specification.HttpErrors, messages, HttpStatusInternalServerError)

		if hasRequiredErrorModels || hasRequiredErrorResponses {
			return
		}
	}

	if hasApis(specification) {
		if specification.HttpErrors == nil {
			specification.HttpErrors = &HttpErrors{ErrorResponses{}, Models{}, nil, nil}
		}
		errorResponses, err := createErrorResponses()
		if err != nil {
			messages.Add(Error(`failed to add required error responses`))
			return
		}
		errorModels, err := createErrorModels()
		if err != nil {
			messages.Add(Error(`failed to add required error responses models`))
			return
		}
		specification.HttpErrors.Responses = append(specification.HttpErrors.Responses, errorResponses...)
		specification.HttpErrors.Models = append(specification.HttpErrors.Models, errorModels...)
	}
}

func errorResponseShouldNotBeDeclared(httpErrors *HttpErrors, messages *Messages, httpStatusName string) bool {
	if httpErrors.Responses != nil {
		errorResponse := httpErrors.Responses.GetByStatusName(httpStatusName)
		if errorResponse != nil {
			messages.Add(Error(`error response '%s' is declared but should not`, httpStatusName).At(locationFromNode(errorResponse.Name.Location)))
			return true
		}
	}
	return false
}

func errorModelShouldNotBeDeclared(httpErrors *HttpErrors, messages *Messages, name string) bool {
	if httpErrors.Models != nil {
		for _, model := range httpErrors.Models {
			if model.Name.Source == name {
				messages.Add(Error(`error model '%s' is declared but should not`, name).At(locationFromNode(model.Location)))
				return true
			}
		}
	}
	return false
}

func hasApis(specification *Spec) bool {
	hasHttp := false
	for _, version := range specification.Versions {
		if len(version.Http.Apis) > 0 {
			hasHttp = true
			break
		}
	}
	return hasHttp
}

type httpEnricher struct {
	models   ModelsMap
	Messages *Messages
}

func (enricher *httpEnricher) version(version *Version) {
	for modIndex := range version.Models {
		version.Models[modIndex].InVersion = version
	}

	http := &version.Http
	http.InVersion = version

	for apiIndex := range http.Apis {
		api := &version.Http.Apis[apiIndex]
		api.InHttp = http
		for opIndex := range api.Operations {
			operation := &api.Operations[opIndex]
			operation.InApi = api
			enricher.operation(operation)
		}
	}
}

func (enricher *httpEnricher) httpErrors(httpErrors *HttpErrors) {
	for index := range httpErrors.Models {
		httpErrors.Models[index].InHttpErrors = httpErrors
	}

	for index := range httpErrors.Responses {
		enricher.responseBody(&httpErrors.Responses[index].Body)
	}
}

func (enricher *httpEnricher) operation(operation *NamedOperation) {
	enricher.params(operation.Endpoint.UrlParams)
	enricher.params(operation.QueryParams)
	enricher.params(operation.HeaderParams)

	if operation.Body != nil {
		enricher.requestBody(operation.Body)
	}

	for index := range operation.Responses {
		operation.Responses[index].Operation = operation
	}

	httpErrors := operation.InApi.InHttp.InVersion.InSpec.HttpErrors
	if httpErrors != nil {
		for index := range operation.Responses {
			responseName := operation.Responses[index].Name.Source
			errorResponse := httpErrors.Responses.GetByStatusName(responseName)
			operation.Responses[index].ErrorResponse = errorResponse
		}
	}

	for index := range operation.Responses {
		enricher.responseBody(&operation.Responses[index].Body)
	}
}

func (enricher *httpEnricher) params(params []NamedParam) {
	for index := range params {
		enricher.typ(&params[index].DefinitionDefault.Type)
	}
}

func (enricher *httpEnricher) definition(definition *Definition) {
	if definition != nil {
		enricher.typ(&definition.Type)
	}
}

func (enricher *httpEnricher) requestBody(body *RequestBody) {
	if body != nil {
		if body.Type != nil {
			enricher.typ(body.Type)
		}
		if body.FormData != nil {
			enricher.params(body.FormData)
		}
		if body.FormUrlEncoded != nil {
			enricher.params(body.FormUrlEncoded)
		}
	}
}

func (enricher *httpEnricher) responseBody(body *ResponseBody) {
	if body != nil && body.Type != nil {
		enricher.typ(body.Type)
	}
}

func (enricher *httpEnricher) typ(typ *Type) {
	if typ != nil {
		enricher.typeDef(typ, &typ.Definition)
	}
}

func (enricher *httpEnricher) typeDef(starter *Type, typ *TypeDef) {
	if typ != nil {
		switch typ.Node {
		case PlainType:
			if model, found := enricher.models[typ.Plain]; found {
				typ.Info = ModelTypeInfo(model)
			} else {
				if info, found := Types[typ.Plain]; found {
					typ.Info = &info
				} else {
					e := Error("unknown type: %s", typ.Plain).At(locationFromNode(starter.Location))
					enricher.Messages.Add(e)
				}
			}
		case NullableType:
			enricher.typeDef(starter, typ.Child)
			typ.Info = NullableTypeInfo(typ.Child.Info)
		case ArrayType:
			enricher.typeDef(starter, typ.Child)
			typ.Info = ArrayTypeInfo()
		case MapType:
			enricher.typeDef(starter, typ.Child)
			typ.Info = MapTypeInfo()
		default:
			panic(fmt.Sprintf("unknown kind of type: %v", typ))
		}
	}
}

func createErrorModels() (Models, error) {
	data := `
BadRequestError:
  object:
    message: string
    location: ErrorLocation
    errors: ValidationError[]?

ValidationError:
  object:
    path: string
    code: string
    message: string?

ErrorLocation:
  enum:
    - unknown
    - query
    - header
    - parameters
    - body

NotFoundError:
  object:
    message: string

InternalServerError:
  object:
    message: string
`
	var models Models
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &models)
	if err != nil {
		return nil, err
	}
	return models, nil
}

const InternalServerError string = "InternalServerError"
const BadRequestError string = "BadRequestError"
const NotFoundError string = "NotFound"
const ValidationError string = "ValidationError"
const ErrorLocation string = "ErrorLocation"

func createErrorResponses() (ErrorResponses, error) {
	data := `
bad_request: BadRequestError   # Service will return this if parameters are not provided or couldn't be parsed correctly
not_found: NotFoundError   # Service will return this if the endpoint is not found
internal_server_error: InternalServerError   # Service will return this if unexpected internal error happens
`
	var responses ErrorResponses
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &responses)
	for index := range responses {
		responses[index].Required = true
	}
	if err != nil {
		return nil, err
	}
	return responses, nil
}
