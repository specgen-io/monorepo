package spec

import (
	"gopkg.in/specgen-io/yaml.v3"
	"yamlx"
)

type NamedDefinition struct {
	Name Name
	Definition
}

type NamedDefinitions []NamedDefinition

func (value *NamedDefinitions) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "named definitions should be YAML mapping")
	}
	count := len(node.Content) / 2
	array := make([]NamedDefinition, count)
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		valueNode := node.Content[index*2+1]
		name := Name{}
		err := keyNode.DecodeWith(decodeStrict, &name)
		if err != nil {
			return err
		}
		err = name.Check(JsonField)
		if err != nil {
			return err
		}
		definition := Definition{}
		err = valueNode.DecodeWith(decodeStrict, &definition)
		if err != nil {
			return err
		}
		array[index] = NamedDefinition{Name: name, Definition: definition}
	}
	*value = array
	return nil
}

func (value NamedDefinitions) MarshalYAML() (interface{}, error) {
	yamlMap := yamlx.Map()
	for index := 0; index < len(value); index++ {
		definition := value[index]
		err := yamlMap.AddWithComment(definition.Name, definition.Definition, definition.Description)
		if err != nil {
			return nil, err
		}
	}
	return yamlMap.Node, nil
}
