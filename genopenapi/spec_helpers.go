package genopenapi

import "github.com/ModaOperandi/spec"

func NewName(source string) spec.Name {
	return spec.Name{Source: source, Location: nil}
}

func NewField(name string, typ spec.TypeDef, defaultValue *string, description *string) *spec.NamedField {
	return &spec.NamedField{
		Name:              NewName(name),
		DefinitionDefault: spec.DefinitionDefault{spec.Type{Definition: typ}, defaultValue, description, nil},
	}
}

func NewObject(fields spec.Fields, description *string) *spec.Object {
	return &spec.Object{fields, description}
}

func NewParam(name string, typ spec.TypeDef, defaultValue *string, description *string) *spec.NamedParam {
	return &spec.NamedParam{
		Name:              NewName(name),
		DefinitionDefault: spec.DefinitionDefault{spec.Type{Definition: typ}, defaultValue, description, nil},
	}
}

func NewResponse(name string, typ spec.TypeDef, description *string) *spec.NamedResponse {
	return &spec.NamedResponse{
		Name:            NewName(name),
		Definition: spec.Definition{spec.Type{Definition: typ}, description, nil},
	}
}

func NewEnumItem(name string, description *string) *spec.NamedEnumItem {
	return &spec.NamedEnumItem{NewName(name), spec.EnumItem{Description: description}}
}