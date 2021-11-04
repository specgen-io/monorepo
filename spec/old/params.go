package old

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/yamlx"
	"gopkg.in/specgen-io/yaml.v3"
)

type UrlParams Params
type QueryParams Params
type HeaderParams Params

type Params []NamedParam

func (params *Params) paramsUnmarshalYAML(node *yaml.Node, paramsName string) error {
	if node.Kind != yaml.MappingNode {
		return yamlError(node, fmt.Sprintf("%s parameters should be YAML mapping", paramsName))
	}
	count := len(node.Content) / 2
	array := make([]NamedParam, count)
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		valueNode := node.Content[index*2+1]
		name := Name{}
		err := keyNode.DecodeWith(decodeStrict, &name)
		if err != nil {
			return err
		}
		err = name.Check(HttpParams)
		if err != nil {
			return err
		}
		definition := DefinitionDefault{}
		err = valueNode.DecodeWith(decodeStrict, &definition)
		if err != nil {
			return err
		}
		if definition.Description == nil {
			definition.Description = getDescription(keyNode)
		}
		array[index] = NamedParam{Name: name, DefinitionDefault: definition}
	}
	*params = array
	return nil
}

func (value *QueryParams) UnmarshalYAML(node *yaml.Node) error {
	params := &Params{}
	err := params.paramsUnmarshalYAML(node, "query")
	if err != nil {
		return err
	}
	*value = []NamedParam(*params)
	return nil
}

func (value *HeaderParams) UnmarshalYAML(node *yaml.Node) error {
	params := &Params{}
	err := params.paramsUnmarshalYAML(node, "header")
	if err != nil {
		return err
	}
	*value = []NamedParam(*params)
	return nil

}

func paramsMarshalYAML(params []NamedParam) (interface{}, error) {
	yamlMap := yamlx.Map()
	for index := 0; index < len(params); index++ {
		param := params[index]
		err := yamlMap.AddWithComment(param.Name, param.DefinitionDefault, param.Description)
		if err != nil {
			return nil, err
		}
	}
	return yamlMap.Node, nil
}

func (params QueryParams) MarshalYAML() (interface{}, error) {
	return paramsMarshalYAML(params)
}

func (params HeaderParams) MarshalYAML() (interface{}, error) {
	return paramsMarshalYAML(params)
}
