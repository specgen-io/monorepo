package spec

import (
	"gopkg.in/specgen-io/yaml.v3"
	"yamlx"
)

type ResponseBody struct {
	Binary   bool
	File     bool
	Type     *Type
	Location *yaml.Node
}

func (body *ResponseBody) Kind() BodyKind {
	if body != nil {
		if body.Binary {
			return BodyBinary
		} else if body.File {
			return BodyFile
		} else if body.Type == nil || body.Type.Definition.IsEmpty() {
			return BodyEmpty
		} else if body.Type.Definition.Plain == TypeString {
			return BodyText
		} else {
			return BodyJson
		}
	}
	return BodyEmpty
}

func (body *ResponseBody) Is(kind BodyKind) bool {
	return body.Kind() == kind
}

func (body *ResponseBody) IsEmpty() bool {
	return body.Kind() == BodyEmpty
}

func (body *ResponseBody) IsText() bool {
	return body.Kind() == BodyText
}

func (body *ResponseBody) IsBinary() bool {
	return body.Kind() == BodyBinary
}

func (body *ResponseBody) IsFile() bool {
	return body.Kind() == BodyFile
}

func (body *ResponseBody) IsJson() bool {
	return body.Kind() == BodyJson
}

func (value *ResponseBody) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.ScalarNode {
		return yamlError(node, "definition has to be scalar value")
	}
	if node.Value == "binary" {
		*value = ResponseBody{Binary: true, Location: node}
	} else if node.Value == "file" {
		*value = ResponseBody{File: true, Location: node}
	} else if node.Value == "empty" {
		*value = ResponseBody{Location: node}
	} else {
		typ, err := parseType(node.Value)
		if err != nil {
			return yamlError(node, err.Error())
		}
		*value = ResponseBody{Type: &Type{*typ, node}, Location: node}
	}
	return nil
}

func (value ResponseBody) MarshalYAML() (interface{}, error) {
	if value.IsBinary() {
		return yamlx.String("binary"), nil
	} else if value.IsFile() {
		return yamlx.String("file"), nil
	} else if value.IsEmpty() {
		return yamlx.String("empty"), nil
	} else {
		return yaml.Node{Kind: yaml.ScalarNode, Value: value.Type.Definition.String()}, nil
	}
}

func (value *ResponseBody) String() string {
	if value.IsEmpty() {
		return "empty"
	} else {
		return value.Type.Definition.String()
	}
}
