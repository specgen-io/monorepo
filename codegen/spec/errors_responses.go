package spec

import (
	"fmt"
	"gopkg.in/specgen-io/yaml.v3"
	"yamlx"
)

type ErrorResponse struct {
	Response
	Required bool
}

type ErrorResponses []ErrorResponse

func (responses ErrorResponses) GetByStatusName(httpStatus string) *ErrorResponse {
	for _, response := range responses {
		if response.Name.Source == httpStatus {
			return &response
		}
	}
	return nil
}

func (responses ErrorResponses) GetByStatusCode(statusCode string) *ErrorResponse {
	for _, response := range responses {
		if response.Name.Source == HttpStatusName(statusCode) {
			return &response
		}
	}
	return nil
}

func (responses ErrorResponses) HttpStatusCodes() []string {
	codes := []string{}
	for _, response := range responses {
		codes = append(codes, HttpStatusCode(response.Name))
	}
	return codes
}

func (responses ErrorResponses) Required() []*ErrorResponse {
	result := []*ErrorResponse{}
	for index := range responses {
		if responses[index].Required {
			result = append(result, &responses[index])
		}
	}
	return result
}

func (responses ErrorResponses) NonRequired() []*ErrorResponse {
	result := []*ErrorResponse{}
	for index := range responses {
		if !responses[index].Required {
			result = append(result, &responses[index])
		}
	}
	return result
}

func (value *ErrorResponses) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "response should be YAML mapping")
	}
	count := len(node.Content) / 2
	array := make([]ErrorResponse, count)
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
		body := ResponseBody{}
		err = valueNode.DecodeWith(decodeStrict, &body)
		if err != nil {
			return err
		}
		response := Response{name, body, getDescriptionFromComment(valueNode)}
		array[index] = ErrorResponse{response, false}
	}
	*value = array
	return nil
}

func (value ErrorResponses) MarshalYAML() (interface{}, error) {
	yamlMap := yamlx.Map()
	for index := 0; index < len(value); index++ {
		response := value[index]
		err := yamlMap.AddWithComment(response.Name, response.Body, response.Description)
		if err != nil {
			return nil, err
		}
	}
	return yamlMap.Node, nil
}
