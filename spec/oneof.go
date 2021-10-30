package spec

type OneOf struct {
	Items         NamedDefinitions `yaml:"oneOf"`
	Discriminator *string          `yaml:"discriminator"`
	Description   *string          `yaml:"description"`
}

func (value OneOf) MarshalYAML() (interface{}, error) {
	yamlMap := NewYamlMap()
	yamlMap.AddOmitNil("discriminator", value.Discriminator)
	yamlMap.Add("oneOf", value.Items)
	return yamlMap.Node, nil
}
