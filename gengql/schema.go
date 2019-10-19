package gengql

import (
	"github.com/ModaOperandi/spec"
	"github.com/vsapronov/gopoetry/graphql"
	"path/filepath"
	"specgen/gen"
)

func GenerateSchema(serviceFile string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}

	schemaFile := generateSchema(specification, generatePath)
	gen.WriteFile(schemaFile, true)

	return nil
}

func generateSchema(spec *spec.Spec, outPath string) *gen.TextFile {
	unit := graphql.Unit()

	for _, model := range spec.Models {
		generateModelSchema(model, unit)
	}

	return &gen.TextFile{
		Path:    filepath.Join(outPath, "schema.graphql"),
		Content: unit.Code(),
	}
}

func generateModelSchema(model spec.NamedModel, unit *graphql.UnitDeclaration) {
	typ := graphql.Type(model.Name.PascalCase())
	for _, field := range model.Object.Fields {
		typ.Field(field.Name.CamelCase(), GraphqlType(&field.Type))
	}
	unit.AddDeclarations(typ)
}
