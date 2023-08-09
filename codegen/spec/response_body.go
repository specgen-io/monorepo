package spec

import "gopkg.in/specgen-io/yaml.v3"

type ResponseBody struct {
	Type     Type
	Location *yaml.Node
}

type ResponseBodyKind string

const (
	ResponseBodyEmpty  ResponseBodyKind = "empty"
	ResponseBodyString ResponseBodyKind = "string"
	ResponseBodyJson   ResponseBodyKind = "json"
)

func (body *ResponseBody) Kind() ResponseBodyKind {
	if body != nil {
		if body.Type.Definition.IsEmpty() {
			return ResponseBodyEmpty
		} else if body.Type.Definition.Plain == TypeString {
			return ResponseBodyString
		} else {
			return ResponseBodyJson
		}
	}
	return ResponseBodyEmpty
}

func (body *ResponseBody) Is(kind ResponseBodyKind) bool {
	return body.Kind() == kind
}

func (body *ResponseBody) IsEmpty() bool {
	return body.Kind() == ResponseBodyEmpty
}

func (body *ResponseBody) IsText() bool {
	return body.Kind() == ResponseBodyString
}

func (body *ResponseBody) IsJson() bool {
	return body.Kind() == ResponseBodyJson
}

func (value *ResponseBody) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.ScalarNode {
		return yamlError(node, "definition has to be scalar value")
	}
	typ, err := parseType(node.Value)
	if err != nil {
		return yamlError(node, err.Error())
	}
	parsed := ResponseBody{Type{*typ, node}, node}
	*value = parsed
	return nil
}

func (value ResponseBody) MarshalYAML() (interface{}, error) {
	yamlValue := value.Type.Definition.String()
	node := yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: yamlValue,
	}
	return node, nil
}
