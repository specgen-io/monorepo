package genruby

import (
	"github.com/ModaOperandi/spec"
	"fmt"
	"gopoetry/ruby"
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
//			generateCirceEnumModel(model, unit)
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

	return class
}

//func generateEnumModel(model spec.NamedModel) ruby.Writable {
//	enumBase := scala.Class(model.Name.PascalCase()).Sealed().Abstract().Extends("StringEnumEntry")
//	enumBaseCtor := enumBase.Contructor()
//	enumBaseCtor.Param("value", "String").Val()
//
//	enumObject := scala.Object(model.Name.PascalCase()).Case().Extends("StringEnum["+model.Name.PascalCase()+"]", "StringCirceEnum["+model.Name.PascalCase()+"]")
//	enumObject_ := enumObject.Define(true)
//	for _, item := range model.Enum.Items {
//		itemObject := scala.Object(item.Name.PascalCase()).Case().Extends(model.Name.PascalCase() + `("` + item.Name.Source + `")`)
//		enumObject_.AddCode(itemObject)
//	}
//	enumObject_.AddLn("val values = findValues")
//
//	unit.AddDeclarations(enumBase)
//	unit.AddDeclarations(enumObject)
//}