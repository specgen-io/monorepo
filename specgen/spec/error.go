package spec

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/yamlx"
	"gopkg.in/specgen-io/yaml.v3"
)

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

func (response *Response) BodyKind() BodyKind {
	return kindOf(&response.Definition)
}

func (response *Response) BodyIs(kind BodyKind) bool {
	return kindOf(&response.Definition) == kind
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
