package spec

type HttpErrors struct {
	Responses      ErrorResponses `yaml:"responses"`
	Models         Models         `yaml:"models"`
	InSpec         *Spec
	ResolvedModels []*NamedModel
}
