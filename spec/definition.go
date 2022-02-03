package spec

import (
	"gopkg.in/specgen-io/yaml.v3"
	"strings"
)

type DefinitionDefault struct {
	Type        Type
	Default     *string
	Description *string
	Location    *yaml.Node
}

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
	if node.Kind != yaml.ScalarNode {
		return yamlError(node, "definition with default has to be scalar value")
	}
	typeStr, defaultValue := parseDefaultedType(node.Value)
	typ, err := parseType(typeStr)
	if err != nil {
		return yamlError(node, err.Error())
	}
	internal := DefinitionDefault{
		Type:        Type{*typ, node},
		Default:     defaultValue,
		Description: getDescriptionFromComment(node),
	}
	*value = internal
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

type Definition struct {
	Type        Type
	Description *string
	Location    *yaml.Node
}

func (value *Definition) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.ScalarNode {
		return yamlError(node, "definition has to be scalar value")
	}
	typ, err := parseType(node.Value)
	if err != nil {
		return yamlError(node, err.Error())
	}
	parsed := Definition{
		Type:        Type{*typ, node},
		Description: getDescriptionFromComment(node),
		Location:    node,
	}
	*value = parsed
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

type BodyKind string

const (
	BodyEmpty  BodyKind = "empty"
	BodyString BodyKind = "string"
	BodyJson   BodyKind = "json"
)

func kindOf(definition *Definition) BodyKind {
	if definition != nil {
		if definition.Type.Definition.IsEmpty() {
			return BodyEmpty
		} else if definition.Type.Definition.Plain == TypeString {
			return BodyString
		} else {
			return BodyJson
		}
	}
	return BodyEmpty
}
