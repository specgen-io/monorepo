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

type RequestBodyKind string

const (
	RequestBodyEmpty  RequestBodyKind = "empty"
	RequestBodyString RequestBodyKind = "string"
	RequestBodyJson   RequestBodyKind = "json"
)

func kindOfRequestBody(body *RequestBody) RequestBodyKind {
	if body != nil {
		if body.Type.Definition.IsEmpty() {
			return RequestBodyEmpty
		} else if body.Type.Definition.Plain == TypeString {
			return RequestBodyString
		} else {
			return RequestBodyJson
		}
	}
	return RequestBodyEmpty
}
