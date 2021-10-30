package spec

import (
	"fmt"
	"gopkg.in/specgen-io/yaml.v3"
)

type NamedResponse struct {
	Name Name
	Definition
	Operation *NamedOperation
}

type Responses []NamedResponse

func (value *Responses) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "response should be YAML mapping")
	}
	count := len(node.Content) / 2
	array := make([]NamedResponse, count)
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
		array[index] = NamedResponse{Name: name, Definition: definition}
	}
	*value = array
	return nil
}

func (value Responses) MarshalYAML() (interface{}, error) {
	yamlMap := NewYamlMap()
	for index := 0; index < len(value); index++ {
		response := value[index]
		err := yamlMap.AddWithComment(response.Name, response.Definition, response.Description)
		if err != nil {
			return nil, err
		}
	}
	return yamlMap.Node, nil
}
