package genruby

import (
	"fmt"
	spec "github.com/specgen-io/spec.v2"
	"github.com/vsapronov/gopoetry/ruby"
	"path/filepath"
	"specgen/gen"
	"strings"
)

func GenerateModels(serviceFile string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}

	fileName := specification.Name.SnakeCase() + "_models.rb"
	moduleName := specification.Name.PascalCase()
	modelsPath := filepath.Join(generatePath, fileName)
	models := generateModels(specification, moduleName, modelsPath)

	err = gen.WriteFile(models, true)
	return err
}

func generateModels(specification *spec.Spec, moduleName string, generatePath string) *gen.TextFile {
	unit := ruby.Unit()
	unit.Require("date")
	unit.Require("emery")

	for _, version := range specification.Versions {
		if version.Version.Source == "" {
			module := generateVersionModelsModule(&version, moduleName)
			unit.AddDeclarations(module)
		}
	}

	for _, version := range specification.Versions {
		if version.Version.Source != "" {
			module := generateVersionModelsModule(&version, moduleName)
			unit.AddDeclarations(module)
		}
	}

	return &gen.TextFile{
		Path:    generatePath,
		Content: unit.Code(),
	}
}

func generateVersionModelsModule(version *spec.Version, moduleName string) *ruby.ModuleDeclaration {
	module := ruby.Module(versionedModule(moduleName, version.Version))
	for _, model := range version.ResolvedModels {
		if model.IsObject() {
			module.AddDeclarations(generateObjectModel(model))
		} else if model.IsOneOf() {
			module.AddDeclarations(generateOneOfModel(model))
		} else if model.IsEnum() {
			module.AddDeclarations(generateEnumModel(model))
		}
	}
	return module
}

func generateObjectModel(model *spec.NamedModel) ruby.Writable {
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

func generateEnumModel(model *spec.NamedModel) ruby.Writable {
	class := ruby.Class(model.Name.PascalCase())

	class.AddCode("include Enum")
	for _, enumItem := range model.Enum.Items {
		class.AddCode(fmt.Sprintf("define :%s, '%s'", enumItem.Name.SnakeCase(), enumItem.Value))
	}

	return class
}

func generateOneOfModel(model *spec.NamedModel) ruby.Writable {
	params := []string{}
	for _, item := range model.OneOf.Items {
		params = append(params, fmt.Sprintf("%s: %s", item.Name.SnakeCase(), RubyType(&item.Type.Definition)))
	}
	typedefinition := fmt.Sprintf("%s = T.union(%s)", model.Name.PascalCase(), strings.Join(params, ", "))
	return ruby.Statements().AddLn(typedefinition)
}
