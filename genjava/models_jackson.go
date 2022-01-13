package genjava

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

var Jackson = "jackson"

type JacksonGenerator struct {
	Type *Types
}

func NewJacksonGenerator(types *Types) *JacksonGenerator {
	return &JacksonGenerator{types}
}

func (g *JacksonGenerator) SetupLibrary(thePackage Module) []sources.CodeFile {
	code := `
package [[.PackageName]];

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.databind.*;
import com.fasterxml.jackson.datatype.jsr310.*;

public class Json {
	public static void setupObjectMapper(ObjectMapper objectMapper) {
		objectMapper
				.registerModule(new JavaTimeModule())
				.configure(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS, false)
				.setSerializationInclusion(JsonInclude.Include.NON_NULL);
	}
}
`

	code, _ = sources.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return []sources.CodeFile{{
		Path:    thePackage.GetPath("Json.java"),
		Content: strings.TrimSpace(code),
	}}
}

func addJacksonImports(w *sources.Writer) {
	w.Line(`import java.time.*;`)
	w.Line(`import java.util.*;`)
	w.Line(`import java.math.BigDecimal;`)
	w.Line(`import com.fasterxml.jackson.databind.JsonNode;`)
	w.Line(`import com.fasterxml.jackson.annotation.*;`)
	w.Line(`import com.fasterxml.jackson.annotation.JsonSubTypes.*;`)
}

func getJacksonJsonPropertyAnnotation(field *spec.NamedDefinition) string {
	required := "false"
	if !field.Type.Definition.IsNullable() {
		required = "true"
	}
	return fmt.Sprintf(`@JsonProperty(value = "%s", required = %s)`, field.Name.Source, required)
}

func (g *JacksonGenerator) ObjectModel(model *spec.NamedModel, thePackage Module) *sources.CodeFile {
	w := NewJavaWriter()
	w.Line(`package %s;`, thePackage.PackageName)
	w.EmptyLine()
	addJacksonImports(w)
	w.EmptyLine()
	className := model.Name.PascalCase()
	w.Line(`public class %s {`, className)
	for _, field := range model.Object.Fields {
		w.EmptyLine()
		w.Line(getJacksonJsonPropertyAnnotation(&field))
		w.Line(`  private %s %s;`, g.Type.JavaType(&field.Type.Definition), field.Name.CamelCase())
	}
	if len(model.Object.Fields) == 0 {
		w.Line(`  public %s() {`, model.Name.PascalCase())
	} else {
		w.Line(`  @JsonCreator`)
		w.Line(`  public %s(`, model.Name.PascalCase())
		for i, field := range model.Object.Fields {
			w.Line(`    %s`, getJacksonJsonPropertyAnnotation(&field))
			ctorParam := fmt.Sprintf(`    %s %s`, g.Type.JavaType(&field.Type.Definition), field.Name.CamelCase())
			if i == len(model.Object.Fields)-1 {
				w.Line(`%s`, ctorParam)
			} else {
				w.Line(`%s,`, ctorParam)
			}
		}
		w.Line(`  ) {`)
	}
	for _, field := range model.Object.Fields {
		if !field.Type.Definition.IsNullable() && g.Type.JavaIsReferenceType(&field.Type.Definition) {
			w.Line(`    if (%s == null) { throw new IllegalArgumentException("null value is not allowed"); }`, field.Name.CamelCase())
		}
		w.Line(`    this.%s = %s;`, field.Name.CamelCase(), field.Name.CamelCase())
	}
	w.Line(`  }`)
	for _, field := range model.Object.Fields {
		w.EmptyLine()
		w.Line(`  public %s %s() {`, g.Type.JavaType(&field.Type.Definition), getterName(&field))
		w.Line(`    return %s;`, field.Name.CamelCase())
		w.Line(`  }`)
		w.EmptyLine()
		w.Line(`  public void %s(%s %s) {`, setterName(&field), g.Type.JavaType(&field.Type.Definition), field.Name.CamelCase())
		w.Line(`    this.%s = %s;`, field.Name.CamelCase(), field.Name.CamelCase())
		w.Line(`  }`)
	}
	w.EmptyLine()
	addObjectModelMethods(w.Indented(), model)
	w.Line(`}`)

	return &sources.CodeFile{
		Path:    thePackage.GetPath(className + ".java"),
		Content: w.String(),
	}
}

func (g *JacksonGenerator) EnumModel(model *spec.NamedModel, thePackage Module) *sources.CodeFile {
	w := NewJavaWriter()
	w.Line(`package %s;`, thePackage.PackageName)
	w.EmptyLine()
	addJacksonImports(w)
	w.EmptyLine()
	enumName := model.Name.PascalCase()
	w.Line(`public enum %s {`, enumName)
	for _, enumItem := range model.Enum.Items {
		w.Line(`  @JsonProperty("%s") %s,`, enumItem.Value, enumItem.Name.UpperCase())
	}
	w.Line(`}`)

	return &sources.CodeFile{
		Path:    thePackage.GetPath(enumName + ".java"),
		Content: w.String(),
	}
}

func (g *JacksonGenerator) OneOfModel(model *spec.NamedModel, thePackage Module) *sources.CodeFile {
	interfaceName := model.Name.PascalCase()
	w := NewJavaWriter()
	w.Line("package %s;", thePackage.PackageName)
	w.EmptyLine()
	addJacksonImports(w)
	w.EmptyLine()
	if model.OneOf.Discriminator != nil {
		w.Line(`@JsonTypeInfo(`)
		w.Line(`  use = JsonTypeInfo.Id.NAME,`)
		w.Line(`  include = JsonTypeInfo.As.PROPERTY,`)
		w.Line(`  property = "%s"`, *model.OneOf.Discriminator)
		w.Line(`)`)
	} else {
		w.Line(`@JsonTypeInfo(`)
		w.Line(`  use = JsonTypeInfo.Id.NAME,`)
		w.Line(`  include = JsonTypeInfo.As.WRAPPER_OBJECT`)
		w.Line(`)`)
	}
	w.Line(`@JsonSubTypes({`)
	for _, item := range model.OneOf.Items {
		w.Line(`  @Type(value = %s.%s.class, name = "%s"),`, interfaceName, oneOfItemClassName(&item), item.Name.Source)
	}
	w.Line(`})`)
	w.Line(`public interface %s {`, interfaceName)
	for index, item := range model.OneOf.Items {
		if index > 0 {
			w.EmptyLine()
		}
		g.generateJacksonOneOfImplementation(w.Indented(), &item, model)
	}
	w.Line(`}`)

	return &sources.CodeFile{
		Path:    thePackage.GetPath(interfaceName + ".java"),
		Content: w.String(),
	}
}

func (g *JacksonGenerator) generateJacksonOneOfImplementation(w *sources.Writer, item *spec.NamedDefinition, model *spec.NamedModel) {
	w.Line(`class %s implements %s {`, oneOfItemClassName(item), model.Name.PascalCase())
	w.Line(`  @JsonUnwrapped`)
	w.Line(`  public %s data;`, g.Type.JavaType(&item.Type.Definition))
	w.EmptyLine()
	w.Line(`  public %s() {`, oneOfItemClassName(item))
	w.Line(`  }`)
	w.EmptyLine()
	w.Line(`  public %s(%s data) {`, oneOfItemClassName(item), g.Type.JavaType(&item.Type.Definition))
	if !item.Type.Definition.IsNullable() && g.Type.JavaIsReferenceType(&item.Type.Definition) {
		w.Line(`    if (data == null) { throw new IllegalArgumentException("null value is not allowed"); }`)
	}
	w.Line(`  	this.data = data;`)
	w.Line(`  }`)
	w.EmptyLine()
	w.Line(`  public %s getData() {`, g.Type.JavaType(&item.Type.Definition))
	w.Line(`    return data;`)
	w.Line(`  }`)
	w.EmptyLine()
	w.Line(`  public void setData(%s data) {`, g.Type.JavaType(&item.Type.Definition))
	if !item.Type.Definition.IsNullable() && g.Type.JavaIsReferenceType(&item.Type.Definition) {
		w.Line(`    if (data == null) { throw new IllegalArgumentException("null value is not allowed"); }`)
	}
	w.Line(`    this.data = data;`)
	w.Line(`  }`)
	w.EmptyLine()
	w.Indent()
	addOneOfModelMethods(w, item)
	w.Unindent()
	w.Line(`}`)
}
