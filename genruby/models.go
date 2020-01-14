package genruby

import (
	"fmt"
	"github.com/ModaOperandi/spec"
	"github.com/vsapronov/gopoetry/ruby"
	"path/filepath"
	"specgen/gen"
	"specgen/static"
	"strings"
)

func GenerateClient(serviceFile string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil { return err }

	unit := ruby.Unit().
		Require("json").
		Require("bigdecimal")

	typecheck, err := static.StaticCode("genruby/typecheck.rb")
	if err != nil { return err }
	unit.AddDeclarations(ruby.Code(typecheck))

	generateModels(specification, unit)

	source := &gen.TextFile{
		Path:    filepath.Join(generatePath, "client.rb"),
		Content: unit.Code(),
	}

	err = gen.WriteFile(source, true)
	if err != nil { return err }

	return nil
}

func generateModels(spec *spec.Spec, unit *ruby.UnitDeclaration) {
	for _, model := range spec.Models {
		if model.IsObject() {
			model := generateObjectModel(model)
			unit.AddDeclarations(model)
		} else {
			model := generateEnumModel(model)
			unit.AddDeclarations(model)
		}
	}
}

func generateObjectModel(model spec.NamedModel) ruby.Writable {
	class := ruby.Class(model.Name.PascalCase())

	fields := make([]string, 0)
	for _, field := range model.Object.Fields {
		fields = append(fields, ":" + field.Name.SnakeCase())
	}
	class.AddMembers(ruby.Code(fmt.Sprintf("attr_reader %s", strings.Join(fields, ", "))))

	initialize := class.Initialize()
	for _, field := range model.Object.Fields {
		initialize.KeywordArg(field.Name.SnakeCase())
	}

	initializeBody := initialize.Body()
	for _, field := range model.Object.Fields {
		fieldName := field.Name.SnakeCase()
		typ := RubyType(&field.Type.Definition)
		initializeBody.AddLn(fmt.Sprintf("@%s = Type.check_field('%s', %s, %s)", fieldName, fieldName, typ, fieldName))
	}

	toHash := class.Def("to_hash").NoParams().Body()
	toHash.AddLn("{")
	hash := toHash.Scope()
	for _, field := range model.Object.Fields {
		fieldName := field.Name.SnakeCase()
		hash.AddLn(fmt.Sprintf(":%s => %s,", fieldName, fieldName))
	}
	toHash.AddLn("}")

	toJson := class.Def("to_json").NoParams().Body()
	toJson.AddLn("JSON.dump(to_hash)")

	return class
}

func generateEnumModel(model spec.NamedModel) ruby.Writable {
	class := ruby.Class(model.Name.PascalCase())

	class.AddCode("include Ruby::Enum")
	for _, enumItem := range model.Enum.Items {
		class.AddCode(fmt.Sprintf("define :%s, '%s'", enumItem.Name.SnakeCase(), enumItem.Name.Source))
	}

	return class
}