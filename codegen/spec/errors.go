package spec

import (
	"gopkg.in/specgen-io/yaml.v3"
)

type HttpErrors struct {
	Responses      ErrorResponses
	Models         Models
	InSpec         *Spec
	ResolvedModels []*NamedModel
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

func createHttpErrors() (*HttpErrors, error) {
	models, err := createErrorModels()
	if err != nil {
		return nil, err
	}
	responses, err := createErrorResponses()
	if err != nil {
		return nil, err
	}
	return &HttpErrors{responses, models, nil, nil}, nil
}
