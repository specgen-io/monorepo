package genjava

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

var Moshi = "moshi"

func addMoshiImports(w *sources.Writer) {
	w.Line(`import com.squareup.moshi.Json;`)
	w.Line(`import java.math.BigDecimal;`)
	w.Line(`import java.time.*;`)
	w.Line(`import java.util.*;`)
}

func generateMoshiObjectModel(model *spec.NamedModel, thePackage Module) *sources.CodeFile {
	w := NewJavaWriter()
	w.Line(`package %s;`, thePackage.PackageName)
	w.EmptyLine()
	addMoshiImports(w)
	w.EmptyLine()
	className := model.Name.PascalCase()
	w.Line(`public class %s {`, className)
	for _, field := range model.Object.Fields {
		w.EmptyLine()
		w.Line(`  @Json(name = "%s")`, field.Name.Source)
		w.Line(`  private %s %s;`, JavaType(&field.Type.Definition), field.Name.CamelCase())
	}
	w.EmptyLine()
	ctorParams := []string{}
	for _, field := range model.Object.Fields {
		ctorParams = append(ctorParams, fmt.Sprintf(`%s %s`, JavaType(&field.Type.Definition), field.Name.CamelCase()))
	}
	w.Line(`  public %s(%s) {`, model.Name.PascalCase(), JoinDelimParams(ctorParams))
	for _, field := range model.Object.Fields {
		w.Line(`    this.%s = %s;`, field.Name.CamelCase(), field.Name.CamelCase())
	}
	w.Line(`  }`)
	addObjectModelProperties(w.Indented(), model)
	w.EmptyLine()
	addObjectModelMethods(w.Indented(), model)
	w.Line(`}`)

	return &sources.CodeFile{
		Path:    thePackage.GetPath(className + ".java"),
		Content: w.String(),
	}
}

func generateMoshiOneOfModels(model *spec.NamedModel, thePackage Module) *sources.CodeFile {
	interfaceName := model.Name.PascalCase()
	w := NewJavaWriter()
	w.Line("package %s;", thePackage.PackageName)
	w.EmptyLine()
	addMoshiImports(w)
	w.EmptyLine()
	w.Line(`public interface %s {`, interfaceName)
	for index, item := range model.OneOf.Items {
		if index > 0 {
			w.EmptyLine()
		}
		generateMoshiOneOfImplementation(w.Indented(), &item, model)
	}
	w.Line(`}`)

	return &sources.CodeFile{
		Path:    thePackage.GetPath(interfaceName + ".java"),
		Content: w.String(),
	}
}

func generateMoshiOneOfImplementation(w *sources.Writer, item *spec.NamedDefinition, model *spec.NamedModel) {
	w.Line(`class %s implements %s {`, oneOfItemClassName(item), model.Name.PascalCase())
	w.Line(`  public %s data;`, JavaType(&item.Type.Definition))
	w.EmptyLine()
	w.Line(`  public %s() {`, oneOfItemClassName(item))
	w.Line(`  }`)
	w.EmptyLine()
	w.Line(`  public %s(%s data) {`, oneOfItemClassName(item), JavaType(&item.Type.Definition))
	w.Line(`  	this.data = data;`)
	w.Line(`  }`)
	w.EmptyLine()
	w.Line(`  public %s getData() {`, JavaType(&item.Type.Definition))
	w.Line(`    return data;`)
	w.Line(`  }`)
	w.EmptyLine()
	w.Line(`  public void setData(%s data) {`, JavaType(&item.Type.Definition))
	w.Line(`    this.data = data;`)
	w.Line(`  }`)
	w.EmptyLine()
	w.Indent()
	addOneOfModelMethods(w, item)
	w.Unindent()
	w.Line(`}`)
}
