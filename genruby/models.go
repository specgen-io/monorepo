package genruby

import (
	"fmt"
	"github.com/ModaOperandi/spec"
	"github.com/vsapronov/gopoetry/ruby"
	"path/filepath"
	"specgen/gen"
)

func GenerateModels(specification *spec.Spec, generatePath string) *gen.TextFile {
	gemName := specification.ServiceName.SnakeCase()+"_client"
	moduleName := specification.ServiceName.PascalCase()
	clientModule := ruby.Module("Client")

	for _, model := range specification.ResolvedModels {
		if model.IsObject() {
			model := generateObjectModel(model)
			clientModule.AddDeclarations(model)
		} else {
			model := generateEnumModel(model)
			clientModule.AddDeclarations(model)
		}
	}

	module := ruby.Module(moduleName).AddDeclarations(clientModule)

	unit := ruby.Unit()
	unit.Require("date")
	unit.Require(gemName+"/tod")
	unit.Require(gemName+"/type")
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

	initialize := class.Initialize()
	for _, field := range model.Object.Fields {
		initialize.KeywordArg(field.Name.SnakeCase())
	}

	initializeBody := initialize.Body()
	initializeBody.AddLn("super method(__method__).parameters.map { |parts| [parts[1], eval(parts[1].to_s)] }.to_h")

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