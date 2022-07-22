package spec

import (
	"fmt"

	"github.com/specgen-io/specgen/yamlx/v2"
	"gopkg.in/specgen-io/yaml.v3"
)

type Response struct {
	Name Name
	Definition
}

type OperationResponse struct {
	Response
	Operation *NamedOperation
}

type OperationResponses []OperationResponse

func (responses OperationResponses) Get(httpStatus string) *OperationResponse {
	for _, response := range responses {
		if response.Name.Source == httpStatus {
			return &response
		}
	}
	return nil
}

func (responses OperationResponses) GetByStatusCode(statusCode string) *OperationResponse {
	for _, response := range responses {
		if response.Name.Source == HttpStatusName(statusCode) {
			return &response
		}
	}
	return nil
}

func (responses OperationResponses) HttpStatusCodes() []string {
	codes := []string{}
	for _, response := range responses {
		codes = append(codes, HttpStatusCode(response.Name))
	}
	return codes
}

func (value *OperationResponses) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "response should be YAML mapping")
	}
	count := len(node.Content) / 2
	array := make([]OperationResponse, count)
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
		array[index] = OperationResponse{Response{Name: name, Definition: definition}, nil}
	}
	*value = array
	return nil
}

func (value OperationResponses) MarshalYAML() (interface{}, error) {
	yamlMap := yamlx.Map()
	for index := 0; index < len(value); index++ {
		response := value[index]
		err := yamlMap.AddWithComment(response.Name, response.Definition, response.Description)
		if err != nil {
			return nil, err
		}
	}
	return yamlMap.Node, nil
}
