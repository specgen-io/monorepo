package spec

import "gopkg.in/specgen-io/yaml.v3"

type RequestBody struct {
	Type        *Type
	Description *string
	Location    *yaml.Node
}

func (value *RequestBody) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.ScalarNode {
		return yamlError(node, "definition has to be scalar value")
	}
	typ, err := parseType(node.Value)
	if err != nil {
		return yamlError(node, err.Error())
	}
	parsed := RequestBody{
		Type:        &Type{*typ, node},
		Description: getDescriptionFromComment(node),
		Location:    node,
	}
	*value = parsed
	return nil
}

func (value RequestBody) MarshalYAML() (interface{}, error) {
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

//type BodyKind string
//
//const (
//	BodyEmpty  BodyKind = "empty"
//	BodyString BodyKind = "string"
//	BodyJson   BodyKind = "json"
//)

func kindOfRequestBody(definition *RequestBody) BodyKind {
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
