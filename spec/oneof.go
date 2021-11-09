package spec

type OneOf struct {
	Discriminator *string          `yaml:"discriminator,omitempty"`
	Items         NamedDefinitions `yaml:"oneOf"`
}
