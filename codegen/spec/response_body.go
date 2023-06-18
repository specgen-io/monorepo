package spec

import "gopkg.in/specgen-io/yaml.v3"

type ResponseBody struct {
	Type        Type
	Description *string
	Location    *yaml.Node
}

func (value *ResponseBody) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.ScalarNode {
		return yamlError(node, "definition has to be scalar value")
	}
	typ, err := parseType(node.Value)
	if err != nil {
		return yamlError(node, err.Error())
	}
	parsed := ResponseBody{
		Type:        Type{*typ, node},
		Description: getDescriptionFromComment(node),
		Location:    node,
	}
	*value = parsed
	return nil
}

func (value ResponseBody) MarshalYAML() (interface{}, error) {
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

type ResponseBodyKind string

const (
	ResponseBodyEmpty  ResponseBodyKind = "empty"
	ResponseBodyString ResponseBodyKind = "string"
	ResponseBodyJson   ResponseBodyKind = "json"
)

func kindOfResponseBody(definition *ResponseBody) ResponseBodyKind {
	if definition != nil {
		if definition.Type.Definition.IsEmpty() {
			return ResponseBodyEmpty
		} else if definition.Type.Definition.Plain == TypeString {
			return ResponseBodyString
		} else {
			return ResponseBodyJson
		}
	}
	return ResponseBodyEmpty
}
