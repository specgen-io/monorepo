package old

import (
	"github.com/specgen-io/specgen/v2/yamlx"
	"gopkg.in/specgen-io/yaml.v3"
)

type objectFieldsKeyword struct {
	Fields      NamedDefinitions `yaml:"fields"`
	Description *string          `yaml:"description"`
}

type object struct {
	Fields      NamedDefinitions `yaml:"object"`
	Description *string          `yaml:"description"`
}

type Object object

func (value *Object) UnmarshalYAML(node *yaml.Node) error {
	if getMappingKey(node, "fields") != nil {
		internal := objectFieldsKeyword{}
		err := node.DecodeWith(decodeStrict, &internal)
		if err != nil {
			return err
		}
		*value = Object(internal)
		return nil
	}
	if getMappingKey(node, "object") != nil {
		internal := object{}
		err := node.DecodeWith(decodeStrict, &internal)
		if err != nil {
			return err
		}
		*value = Object(internal)
		return nil
	}
	fields := NamedDefinitions{}
	err := node.DecodeWith(decodeStrict, &fields)
	if err != nil {
		return err
	}
	*value = Object{Fields: fields}
	return nil
}

func (value Object) MarshalYAML() (interface{}, error) {
	yamlMap := yamlx.Map()
	yamlMap.Add("object", value.Fields)
	return yamlMap.Node, nil
}
