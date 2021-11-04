package spec

import "github.com/specgen-io/specgen/v2/yamlx"

type OneOf struct {
	Items         NamedDefinitions `yaml:"oneOf"`
	Discriminator *string          `yaml:"discriminator"`
}

func (value OneOf) MarshalYAML() (interface{}, error) {
	yamlMap := yamlx.Map()
	yamlMap.AddOmitNil("discriminator", value.Discriminator)
	yamlMap.Add("oneOf", value.Items)
	return yamlMap.Node, nil
}
