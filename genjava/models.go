package genjava

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func GenerateModels(specification *spec.Spec, packageName string, generatePath string) *gen.Sources {
	if packageName == "" {
		packageName = specification.Name.SnakeCase()
	}
	mainPackage := Package(generatePath, packageName)
	modelsPackage := mainPackage.Subpackage("models")
	sources := gen.NewSources()
	sources.AddGeneratedAll(generateModels(specification, modelsPackage))
	return sources
}

func generateModels(specification *spec.Spec, thePackage Module) []gen.TextFile {
	files := []gen.TextFile{}
	for _, version := range specification.Versions {
		versionPackage := thePackage.Subpackage(version.Version.FlatCase())
		files = append(files, generateVersionModels(&version, versionPackage)...)
	}
	files = append(files, *generateJson(thePackage))
	return files
}

func generateJson(thePackage Module) *gen.TextFile {
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

	code, _ = gen.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &gen.TextFile{
		Path:    thePackage.GetPath("Json.java"),
		Content: strings.TrimSpace(code),
	}
}

func generateVersionModels(version *spec.Version, thePackage Module) []gen.TextFile {
	files := []gen.TextFile{}
	for _, model := range version.ResolvedModels {
		if model.IsObject() {
			files = append(files, *generateObjectModel(model, thePackage))
		} else if model.IsOneOf() {
			files = append(files, *generateOneOfModels(model, thePackage))
		} else if model.IsEnum() {
			files = append(files, *generateEnumModel(model, thePackage))
		}
	}
	return files
}

func addImports(w *gen.Writer) {
	w.Line(`import java.time.*;`)
	w.Line(`import java.util.*;`)
	w.Line(`import java.math.BigDecimal;`)
	w.Line(`import com.fasterxml.jackson.databind.JsonNode;`)
	w.Line(`import com.fasterxml.jackson.annotation.*;`)
	w.Line(`import com.fasterxml.jackson.annotation.JsonSubTypes.*;`)
}

func addObjectModelCtors(w *gen.Writer, model *spec.NamedModel) {
	if len(model.Object.Fields) == 0 {
		w.Line(`public %s() {`, model.Name.PascalCase())
	} else {
		w.Line(`@JsonCreator`)
		w.Line(`public %s(`, model.Name.PascalCase())
		for i, field := range model.Object.Fields {
			w.Line(`  %s`, getJsonPropertyAnnotation(&field))
			if i == len(model.Object.Fields)-1 {
				w.Line(`  %s %s`, JavaType(&field.Type.Definition), field.Name.CamelCase())
			} else {
				w.Line(`  %s %s,`, JavaType(&field.Type.Definition), field.Name.CamelCase())
			}
		}
		w.Line(`) {`)
	}
	for _, field := range model.Object.Fields {
		if !field.Type.Definition.IsNullable() && JavaIsReferenceType(&field.Type.Definition) {
			w.Line(`  if (%s == null) { throw new IllegalArgumentException("null value is not allowed"); }`, field.Name.CamelCase())
		}
		w.Line(`  this.%s = %s;`, field.Name.CamelCase(), field.Name.CamelCase())
	}
	w.Line(`}`)
}

func getJsonPropertyAnnotation(field *spec.NamedDefinition) string {
	required := "false"
	if !field.Type.Definition.IsNullable() {
		required = "true"
	}
	return fmt.Sprintf(`@JsonProperty(value = "%s", required = %s)`, field.Name.Source, required)
}

func addObjectModelFields(w *gen.Writer, model *spec.NamedModel) {
	for _, field := range model.Object.Fields {
		w.EmptyLine()
		w.Line(getJsonPropertyAnnotation(&field))
		w.Line(`private %s %s;`, JavaType(&field.Type.Definition), field.Name.CamelCase())
	}
	for _, field := range model.Object.Fields {
		w.EmptyLine()
		w.Line(`public %s %s() {`, JavaType(&field.Type.Definition), getterName(&field))
		w.Line(`  return %s;`, field.Name.CamelCase())
		w.Line(`}`)
		w.EmptyLine()
		w.Line(`public void %s(%s %s) {`, setterName(&field), JavaType(&field.Type.Definition), field.Name.CamelCase())
		if !field.Type.Definition.IsNullable() && JavaIsReferenceType(&field.Type.Definition) {
			w.Line(`  if (%s == null) { throw new IllegalArgumentException("null value is not allowed"); }`, field.Name.CamelCase())
		}
		w.Line(`  this.%s = %s;`, field.Name.CamelCase(), field.Name.CamelCase())
		w.Line(`}`)
	}
}

func addObjectModelMethods(w *gen.Writer, model *spec.NamedModel) {
	w.Line(`@Override`)
	w.Line(`public boolean equals(Object o) {`)
	w.Line(`  if (this == o) return true;`)
	w.Line(`  if (!(o instanceof %s)) return false;`, model.Name.PascalCase())
	w.Line(`  %s that = (%s) o;`, model.Name.PascalCase(), model.Name.PascalCase())
	equalsParams := []string{}
	for _, field := range model.Object.Fields {
		equalsParam := equalsExpression(&field.Type.Definition, fmt.Sprintf(`%s()`, getterName(&field)), fmt.Sprintf(`%s()`, getterName(&field)))
		equalsParams = append(equalsParams, equalsParam)
	}
	w.Line(`  return %s;`, strings.Join(equalsParams, " && "))
	w.Line(`}`)
	w.EmptyLine()
	w.Line(`@Override`)
	w.Line(`public int hashCode() {`)
	hashCodeParams := []string{}
	hashParams := []string{}
	var arrayFieldCount, nonArrayFieldCount int
	for _, field := range model.Object.Fields {
		if isJavaArrayType(&field.Type.Definition) {
			hashCodeParams = append(hashCodeParams, fmt.Sprintf(`%s()`, getterName(&field)))
			arrayFieldCount++
		} else {
			hashParams = append(hashParams, fmt.Sprintf(`%s()`, getterName(&field)))
			nonArrayFieldCount++
		}
	}
	if arrayFieldCount > 0 && nonArrayFieldCount == 0 {
		w.Line(`  int result = Arrays.hashCode(%s);`, hashCodeParams[0])
		for _, param := range hashCodeParams[1:] {
			w.Line(`  result = 31 * result + Arrays.hashCode(%s);`, param)
		}
		w.Line(`  return result;`)
	} else if arrayFieldCount > 0 && nonArrayFieldCount > 0 {
		w.Line(`  int result = Objects.hash(%s);`, JoinDelimParams(hashParams))
		for _, param := range hashCodeParams {
			w.Line(`  result = 31 * result + Arrays.hashCode(%s);`, param)
		}
		w.Line(`  return result;`)
	} else if arrayFieldCount == 0 && nonArrayFieldCount > 0 {
		w.Line(`  return Objects.hash(%s);`, JoinDelimParams(hashParams))
	}
	w.Line(`}`)
	w.EmptyLine()
	w.Line(`@Override`)
	w.Line(`public String toString() {`)
	formatParams := []string{}
	formatArgs := []string{}
	for _, field := range model.Object.Fields {
		formatParams = append(formatParams, fmt.Sprintf(`%s=%s`, field.Name.CamelCase(), "%s"))
		formatArgs = append(formatArgs, toStringExpression(&field.Type.Definition, field.Name.CamelCase()))
	}
	w.Line(`  return String.format("%s{%s}", %s);`, model.Name.PascalCase(), JoinDelimParams(formatParams), JoinDelimParams(formatArgs))
	w.Line(`}`)
}

func generateObjectModel(model *spec.NamedModel, thePackage Module) *gen.TextFile {
	w := NewJavaWriter()
	w.Line(`package %s;`, thePackage.PackageName)
	w.EmptyLine()
	addImports(w)
	w.EmptyLine()
	className := model.Name.PascalCase()
	w.Line(`public class %s {`, className)
	addObjectModelCtors(w.Indented(), model)
	addObjectModelFields(w.Indented(), model)
	w.EmptyLine()
	addObjectModelMethods(w.Indented(), model)
	w.Line(`}`)

	return &gen.TextFile{
		Path:    thePackage.GetPath(className + ".java"),
		Content: w.String(),
	}
}

func generateEnumModel(model *spec.NamedModel, thePackage Module) *gen.TextFile {
	w := NewJavaWriter()
	w.Line(`package %s;`, thePackage.PackageName)
	w.EmptyLine()
	addImports(w)
	w.EmptyLine()
	enumName := model.Name.PascalCase()
	w.Line(`public enum %s {`, enumName)
	for _, enumItem := range model.Enum.Items {
		w.Line(`  @JsonProperty("%s") %s,`, enumItem.Value, enumItem.Name.UpperCase())
	}
	w.Line(`}`)

	return &gen.TextFile{
		Path:    thePackage.GetPath(enumName + ".java"),
		Content: w.String(),
	}
}

func generateOneOfModels(model *spec.NamedModel, thePackage Module) *gen.TextFile {
	interfaceName := model.Name.PascalCase()
	w := NewJavaWriter()
	w.Line("package %s;", thePackage.PackageName)
	w.EmptyLine()
	addImports(w)
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
		generateOneOfImplementation(w.Indented(), &item, model)
	}
	w.Line(`}`)

	return &gen.TextFile{
		Path:    thePackage.GetPath(interfaceName + ".java"),
		Content: w.String(),
	}
}

func oneOfItemClassName(item *spec.NamedDefinition) string {
	return item.Name.PascalCase()
}

func generateOneOfImplementation(w *gen.Writer, item *spec.NamedDefinition, model *spec.NamedModel) {
	w.Line(`class %s implements %s {`, oneOfItemClassName(item), model.Name.PascalCase())
	w.Line(`  @JsonUnwrapped`)
	w.Line(`  public %s data;`, JavaType(&item.Type.Definition))
	w.EmptyLine()
	w.Line(`  public %s() {`, oneOfItemClassName(item))
	w.Line(`  }`)
	w.EmptyLine()
	w.Line(`  public %s(%s data) {`, oneOfItemClassName(item), JavaType(&item.Type.Definition))
	if !item.Type.Definition.IsNullable() && JavaIsReferenceType(&item.Type.Definition) {
		w.Line(`    if (data == null) { throw new IllegalArgumentException("null value is not allowed"); }`)
	}
	w.Line(`  	this.data = data;`)
	w.Line(`  }`)
	w.EmptyLine()
	w.Line(`  public %s getData() {`, JavaType(&item.Type.Definition))
	w.Line(`    return data;`)
	w.Line(`  }`)
	w.EmptyLine()
	w.Line(`  public void setData(%s data) {`, JavaType(&item.Type.Definition))
	if !item.Type.Definition.IsNullable() && JavaIsReferenceType(&item.Type.Definition) {
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

func addOneOfModelMethods(w *gen.Writer, item *spec.NamedDefinition) {
	w.Line(`@Override`)
	w.Line(`public boolean equals(Object o) {`)
	w.Line(`  if (this == o) return true;`)
	w.Line(`  if (!(o instanceof %s)) return false;`, oneOfItemClassName(item))
	w.Line(`  %s that = (%s) o;`, oneOfItemClassName(item), oneOfItemClassName(item))
	w.Line(`  return Objects.equals(getData(), that.getData());`)
	w.Line(`}`)
	w.EmptyLine()
	w.Line(`@Override`)
	w.Line(`public int hashCode() {`)
	w.Line(`  return Objects.hash(getData());`)
	w.Line(`}`)
	w.EmptyLine()
	w.Line(`@Override`)
	w.Line(`public String toString() {`)
	w.Line(`  return String.format("%s{data=%s}", data);`, JavaType(&item.Type.Definition), "%s")
	w.Line(`}`)
}
