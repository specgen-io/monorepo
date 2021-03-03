package genruby

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/vsapronov/gopoetry/ruby"
	"path/filepath"
	"specgen/gen"
	"specgen/static"
	"strings"
)

func GenerateModels(serviceFile string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil { return err }

	folderName := specification.ServiceName.SnakeCase()
	moduleName := specification.ServiceName.PascalCase()
	modelsPath := filepath.Join(generatePath, folderName)
	models := generateModels(specification.ResolvedModels, folderName, moduleName, modelsPath)

	data := static.RubyClient{GemName: folderName, ModuleName: moduleName}
	sources, err := static.RenderTemplate("ruby-models", generatePath, data)
	if err != nil { return err }

	sources = append(sources, *models)
	err = gen.WriteFiles(sources, true)
	return err
}

func generateModels(models spec.Models, folderName string, moduleName string, generatePath string) *gen.TextFile {
	module := ruby.Module(moduleName)

	for _, model := range models {
		if model.IsObject() {
			module.AddDeclarations(generateObjectModel(model))
		} else if model.IsOneOf() {
			module.AddDeclarations(generateOneOfModel(model))
		} else if model.IsEnum() {
			module.AddDeclarations(generateEnumModel(model))
		}
	}

	unit := ruby.Unit()
	unit.Require("date")
	unit.Require(folderName+"/tod")
	unit.Require(folderName+"/type")
	unit.Require(folderName+"/enum")
	unit.Require(folderName+"/dataclass")
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
		class.AddCode(fmt.Sprintf("define :%s, '%s'", enumItem.Name.SnakeCase(), enumItem.Value))
	}

	return class
}

func generateOneOfModel(model spec.NamedModel) ruby.Writable {
	params := []string{}
	for _, item := range model.OneOf.Items {
		params = append(params, fmt.Sprintf("%s: %s", item.Name.SnakeCase(), RubyType(&item.Type.Definition)))
	}
	typedefinition := fmt.Sprintf("%s = T.union(%s)", model.Name.PascalCase(), strings.Join(params, ", "))
	return ruby.Statements().AddLn(typedefinition)
}