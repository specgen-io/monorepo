package spec

import (
	"fmt"
	"gopkg.in/specgen-io/yaml.v3"
	"yamlx"
)

type HttpErrors struct {
	Responses      Responses
	Models         Models
	ResolvedModels []*NamedModel
}

type Responses []Response

func (responses Responses) GetByStatusName(httpStatus string) *Response {
	for _, response := range responses {
		if response.Name.Source == httpStatus {
			return &response
		}
	}
	return nil
}

func (responses Responses) GetByStatusCode(statusCode string) *Response {
	for _, response := range responses {
		if response.Name.Source == HttpStatusName(statusCode) {
			return &response
		}
	}
	return nil
}

func (responses Responses) HttpStatusCodes() []string {
	codes := []string{}
	for _, response := range responses {
		codes = append(codes, HttpStatusCode(response.Name))
	}
	return codes
}

func (responses *Responses) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "response should be YAML mapping")
	}
	count := len(node.Content) / 2
	array := make([]Response, count)
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		valueNode := node.Content[index*2+1]
		name := Name{}
		err := keyNode.DecodeWith(decodeStrict, &name)
		if err != nil {
			return err
		}
		err = name.Check(SnakeCase)
		if err != nil {
			return err
		}
		if _, ok := httpStatusCode[name.Source]; !ok {
			return yamlError(keyNode, fmt.Sprintf("unknown response name %s", name.Source))
		}
		definition := Definition{}
		err = valueNode.DecodeWith(decodeStrict, &definition)
		if err != nil {
			return err
		}
		array[index] = Response{Name: name, Definition: definition}
	}
	*responses = array
	return nil
}

func (responses Responses) MarshalYAML() (interface{}, error) {
	yamlMap := yamlx.Map()
	for index := 0; index < len(responses); index++ {
		response := responses[index]
		err := yamlMap.AddWithComment(response.Name, response.Definition, response.Description)
		if err != nil {
			return nil, err
		}
	}
	return yamlMap.Node, nil
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

func createErrorResponses() (Responses, error) {
	data := `
bad_request: BadRequestError   # Service will return this if parameters are not provided or couldn't be parsed correctly
not_found: NotFoundError   # Service will return this if the endpoint is not found
internal_server_error: InternalServerError   # Service will return this if unexpected internal error happens
`
	var responses Responses
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &responses)
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
	return &HttpErrors{responses, models, nil}, nil
}
