package spec

import (
	"github.com/specgen-io/specgen/v2/yamlx"
)

type Object struct {
	Fields      NamedDefinitions `yaml:"object"`
	Description *string          `yaml:"description"`
}

func (value Object) MarshalYAML() (interface{}, error) {
	yamlMap := yamlx.Map()
	yamlMap.Add("object", value.Fields)
	return yamlMap.Node, nil
}
