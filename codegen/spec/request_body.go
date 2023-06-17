package spec

import (
	"gopkg.in/specgen-io/yaml.v3"
	"yamlx"
)

type FormDataParams Params
type FormUrlEncodedParams Params

func (value *FormDataParams) UnmarshalYAML(node *yaml.Node) error {
	params := &Params{}
	err := params.paramsUnmarshalYAML(node, "form-data")
	if err != nil {
		return err
	}
	*value = []NamedParam(*params)
	return nil
}

func (params FormDataParams) MarshalYAML() (interface{}, error) {
	return paramsMarshalYAML(params)
}

func (value *FormUrlEncodedParams) UnmarshalYAML(node *yaml.Node) error {
	params := &Params{}
	err := params.paramsUnmarshalYAML(node, "form-urlencoded")
	if err != nil {
		return err
	}
	*value = []NamedParam(*params)
	return nil
}

func (params FormUrlEncodedParams) MarshalYAML() (interface{}, error) {
	return paramsMarshalYAML(params)
}

type RequestBody struct {
	Type           *Type
	FormData       FormDataParams
	FormUrlEncoded FormUrlEncodedParams
	Description    *string
	Location       *yaml.Node
}

func (value *RequestBody) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode {
		typ, err := parseType(node.Value)
		if err != nil {
			return yamlError(node, err.Error())
		}
		parsed := RequestBody{Type: &Type{*typ, node}, Description: getDescriptionFromComment(node), Location: node}
		*value = parsed
		return nil
	} else if node.Kind == yaml.MappingNode {
		if len(node.Content) != 2 {
			return yamlError(node, `body has to be either type or an object with single field: form-data or form-urlencoded`)
		}
		if paramsNode := getMappingValue(node, "form-data"); paramsNode != nil {
			params := FormDataParams{}
			err := paramsNode.DecodeWith(decodeLooze, &params)
			if err != nil {
				return yamlError(node, err.Error())
			}
			parsed := RequestBody{FormData: params, Description: getDescriptionFromComment(node), Location: node}
			*value = parsed
			return nil
		} else if paramsNode := getMappingValue(node, "form-urlencoded"); paramsNode != nil {
			params := FormUrlEncodedParams{}
			err := paramsNode.DecodeWith(decodeLooze, &params)
			if err != nil {
				return yamlError(node, err.Error())
			}
			parsed := RequestBody{FormUrlEncoded: params, Description: getDescriptionFromComment(node), Location: node}
			*value = parsed
			return nil
		}
	}

	return yamlError(node, "body has to be either type or params")
}

func (value RequestBody) MarshalYAML() (interface{}, error) {
	var node yaml.Node
	if value.Type != nil {
		yamlValue := value.Type.Definition.String()
		node = yaml.Node{Kind: yaml.ScalarNode, Value: yamlValue}
	} else if value.FormData != nil {
		yamlMap := yamlx.Map()
		yamlMap.Add("form-data", value.FormData)
		node = yamlMap.Node
	} else if value.FormUrlEncoded != nil {
		yamlMap := yamlx.Map()
		yamlMap.Add("form-urlencoded", value.FormUrlEncoded)
		node = yamlMap.Node
	}
	if value.Description != nil {
		node.LineComment = "# " + *value.Description
	}
	return node, nil
}

type RequestBodyKind string

const (
	RequestBodyEmpty          RequestBodyKind = "empty"
	RequestBodyString         RequestBodyKind = "string"
	RequestBodyJson           RequestBodyKind = "json"
	RequestBodyFormData       RequestBodyKind = "form-data"
	RequestBodyFormUrlEncoded RequestBodyKind = "form-urlencoded"
)

func kindOfRequestBody(body *RequestBody) RequestBodyKind {
	if body != nil {
		if body.Type != nil {
			if body.Type.Definition.IsEmpty() {
				return RequestBodyEmpty
			} else if body.Type.Definition.Plain == TypeString {
				return RequestBodyString
			} else {
				return RequestBodyJson
			}
		}
		if body.FormData != nil {
			return RequestBodyFormData
		}
		if body.FormUrlEncoded != nil {
			return RequestBodyFormUrlEncoded
		}
	}
	return RequestBodyEmpty
}
