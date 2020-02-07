package genruby

import (
	"fmt"
	"github.com/ModaOperandi/spec"
	"github.com/vsapronov/gopoetry/ruby"
	"path/filepath"
	"specgen/gen"
	"specgen/static"
)

func GenerateClient(serviceFile string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil { return err }

	gemName := specification.ServiceName.SnakeCase()+"_client"
	moduleName := specification.ServiceName.PascalCase()+"Client"
	models := generateModels(specification, filepath.Join(generatePath, "lib", gemName))

	data := static.RubyGem{GemName: gemName, ModuleName: moduleName}
	sources, err := static.RenderTemplate("genruby", generatePath, data)
	if err != nil { return err }

	sources = append(sources, *models)
	err = gen.WriteFiles(sources, true)
	return err
}

func generateModels(specification *spec.Spec, generatePath string) *gen.TextFile {
	gemName := specification.ServiceName.SnakeCase()+"_client"
	moduleName := specification.ServiceName.PascalCase()+"Client"
	module := ruby.Module(moduleName)
	for _, model := range specification.Models {
		if model.IsObject() {
			model := generateObjectModel(model)
			module.AddDeclarations(model)
		} else {
			model := generateEnumModel(model)
			module.AddDeclarations(model)
		}
	}
	unit := ruby.Unit()
	unit.Require(gemName+"/enum")
	unit.Require(gemName+"/dataclass")
	unit.AddDeclarations(module)
	return &gen.TextFile{
		Path:    filepath.Join(generatePath, "models.rb"),
		Content: unit.Code(),
	}
}

func generateObjectModel(model spec.NamedModel) ruby.Writable {
	class := ruby.Class(model.Name.PascalCase())
	class.AddCode("include DataClass")

	for _, field := range model.Object.Fields {
		typ := RubyType(&field.Type.Definition)
		class.AddCode(fmt.Sprintf("val :%s, %s", field.Name.SnakeCase(), typ))
	}

	//initialize := class.Initialize()
	//for _, field := range model.Object.Fields {
	//	initialize.KeywordArg(field.Name.SnakeCase())
	//}
	//
	//initializeBody := initialize.Body()
	//for _, field := range model.Object.Fields {
	//	fieldName := field.Name.SnakeCase()
	//	typ := RubyType(&field.Type.Definition)
	//	initializeBody.AddLn(fmt.Sprintf("@%s = Type.check_field('%s', %s, %s)", fieldName, fieldName, typ, fieldName))
	//}

	return class
}

func generateEnumModel(model spec.NamedModel) ruby.Writable {
	class := ruby.Class(model.Name.PascalCase())

	class.AddCode("include Enum")
	for _, enumItem := range model.Enum.Items {
		class.AddCode(fmt.Sprintf("define :%s, '%s'", enumItem.Name.SnakeCase(), enumItem.Name.Source))
	}

	return class
}