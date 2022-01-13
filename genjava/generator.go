package genjava

type Generator struct {
	Jsonlib string
	Types   *Types
	Models  ModelsGenerator
}

func NewGenerator(jsonlib string) *Generator {
	return &Generator{
		jsonlib,
		NewTypes(jsonlib),
		NewModelsGenerator(jsonlib),
	}
}
