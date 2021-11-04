package old

import (
	"gopkg.in/specgen-io/yaml.v3"
	"strings"
)

type definitionDefault struct {
	Type        Type    `yaml:"type"`
	Default     *string `yaml:"default"`
	Description *string `yaml:"description"`
	Location    *yaml.Node
}

type DefinitionDefault definitionDefault

func parseDefaultedType(str string) (string, *string) {
	if strings.Contains(str, "=") {
		parts := strings.SplitN(str, "=", 2)
		typeStr := strings.TrimSpace(parts[0])
		defaultValue := strings.TrimSpace(parts[1])
		return typeStr, &defaultValue
	} else {
		return str, nil
	}
}

func (value *DefinitionDefault) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode {
		typeStr, defaultValue := parseDefaultedType(node.Value)
		typ, err := parseType(typeStr)
		if err != nil {
			return yamlError(node, err.Error())
		}
		internal := DefinitionDefault{
			Type:        Type{*typ, node},
			Default:     defaultValue,
			Description: getDescription(node),
		}
		*value = internal
	} else {
		internal := definitionDefault{}
		err := node.DecodeWith(decodeStrict, &internal)
		if err != nil {
			return err
		}
		definition := DefinitionDefault(internal)
		definition.Location = node
		*value = definition
	}
	return nil
}

func (value DefinitionDefault) MarshalYAML() (interface{}, error) {
	yamlValue := value.Type.Definition.String()
	if value.Default != nil {
		yamlValue = yamlValue + " = " + *value.Default
	}
	node := yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: yamlValue,
	}
	if value.Description != nil {
		node.LineComment = "# " + *value.Description
	}
	return node, nil
}

type definition struct {
	Type        Type    `yaml:"type"`
	Description *string `yaml:"description"`
	Location    *yaml.Node
}

type Definition definition

func (value *Definition) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode {
		typ, err := parseType(node.Value)
		if err != nil {
			return yamlError(node, err.Error())
		}
		parsed := Definition{
			Type:        Type{*typ, node},
			Description: getDescription(node),
			Location:    node,
		}
		*value = parsed
	} else {
		internal := definition{}
		err := node.DecodeWith(decodeStrict, &internal)
		if err != nil {
			return err
		}
		definition := Definition(internal)
		definition.Location = node
		*value = definition
	}
	return nil
}

func (value Definition) MarshalYAML() (interface{}, error) {
	yamlValue := value.Type.Definition.String()
	node := yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: yamlValue,
	}
	if value.Description != nil {
		node.LineComment = "# " + *value.Description
	}
	return node, nil
}
