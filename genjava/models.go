package genjava

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"strings"
)

func GenerateModels(specification *spec.Spec, packageName string, generatePath string) error {
	if packageName == "" {
		packageName = specification.Name.SnakeCase()
	}
	mainPackage := Package(generatePath, packageName)
	modelsPackage := mainPackage.Subpackage("models")
	files := generateModels(specification, modelsPackage)
	return gen.WriteFiles(files, true)
}

func generateModels(specification *spec.Spec, thePackage Module) []gen.TextFile {
	files := []gen.TextFile{}
	for _, version := range specification.Versions {
		versionPackage := thePackage.Subpackage(version.Version.FlatCase())
		files = append(files, generateVersionModels(&version, versionPackage)...)
	}
	files = append(files, *generateJsoner(thePackage))
	return files
}

func generateJsoner(thePackage Module) *gen.TextFile {
	code := `
package [[.PackageName]];

import com.fasterxml.jackson.databind.*;
import com.fasterxml.jackson.datatype.jsr310.*;

import java.io.*;

public class Jsoner {

	public static void setupObjectMapper(ObjectMapper objectMapper) {
		objectMapper.registerModule(new JavaTimeModule()).configure(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS, false);
	}

	public static <T> String serialize(ObjectMapper objectMapper, T data) throws IOException {
		return objectMapper.writeValueAsString(data);
	}

	public static <T> T deserialize(ObjectMapper objectMapper, String jsonStr, Class<T> tClass) throws IOException {
		return objectMapper.readValue(jsonStr, tClass);
	}
}
`

	code, _ = gen.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &gen.TextFile{
		Path:    thePackage.GetPath("Jsoner.java"),
		Content: strings.TrimSpace(code),
	}
}

func generateVersionModels(version *spec.Version, thePackage Module) []gen.TextFile {
	files := []gen.TextFile{}
	for _, model := range version.ResolvedModels {
		if model.IsObject() {
			files = append(files, *generateObjectModel(model, thePackage))
		} else if model.IsOneOf() {
			files = append(files, generateOneOfModels(model, thePackage)...)
		} else if model.IsEnum() {
			files = append(files, *generateEnumModel(model, thePackage))
		}
	}
	return files
}

func generateImports(w *gen.Writer) {
	w.Line(`import java.time.*;`)
	w.Line(`import java.util.*;`)
	w.Line(`import java.math.BigDecimal;`)
	w.Line(`import com.fasterxml.jackson.databind.JsonNode;`)
	w.Line(`import com.fasterxml.jackson.annotation.*;`)
	w.Line(`import com.fasterxml.jackson.annotation.JsonSubTypes.*;`)
}

func addType(field spec.NamedDefinition) string {
	if checkType(&field.Type.Definition, spec.TypeJson) {
		return `String`
	}
	return JavaType(&field.Type.Definition)
}

func addGetterBody(field spec.NamedDefinition) string {
	if checkType(&field.Type.Definition, spec.TypeJson) {
		return fmt.Sprintf(` %s == null ? null : %s.toString()`, field.Name.CamelCase(), field.Name.CamelCase())
	}
	return field.Name.CamelCase()
}

func addFieldName(field spec.NamedDefinition) string {
	if checkType(&field.Type.Definition, spec.TypeJson) {
		return `node`
	}
	return field.Name.CamelCase()
}

func addSetterParams(field spec.NamedDefinition) string {
	if checkType(&field.Type.Definition, spec.TypeJson) {
		return fmt.Sprintf(`JsonNode %s`, addFieldName(field))
	}
	return fmt.Sprintf(`%s %s`, addType(field), addFieldName(field))
}

func generateObjectModel(model *spec.NamedModel, thePackage Module) *gen.TextFile {
	w := NewJavaWriter()
	w.Line(`package %s;`, thePackage.PackageName)
	w.EmptyLine()
	generateImports(w)
	w.EmptyLine()
	className := model.Name.PascalCase()
	w.Line(`public class %s {`, className)

	constructParams := []string{}
	for _, field := range model.Object.Fields {
		w.Line(`  @JsonProperty("%s")`, field.Name.SnakeCase())
		if checkType(&field.Type.Definition, spec.TypeJson) {
			w.Line(`  @JsonRawValue`)
		}
		w.Line(`  private %s %s;`, JavaType(&field.Type.Definition), field.Name.CamelCase())
		constructParams = append(constructParams, fmt.Sprintf(`%s %s`, addType(field), field.Name.CamelCase()))
	}

	w.EmptyLine()
	w.Line(`  public %s() {`, model.Name.PascalCase())
	w.Line(`  }`)
	w.EmptyLine()
	w.Line(`  public %s(%s) {`, model.Name.PascalCase(), JoinParams(constructParams))
	for _, field := range model.Object.Fields {
		w.Line(`    this.%s = %s;`, field.Name.CamelCase(), field.Name.CamelCase())
	}
	w.Line(`  }`)

	for _, field := range model.Object.Fields {
		w.EmptyLine()
		w.Line(`  public %s %s() {`, addType(field), getterName(&field))
		w.Line(`    return %s;`, addGetterBody(field))
		w.Line(`  }`)
		w.EmptyLine()
		w.Line(`  public void %s(%s) {`, setterName(&field), addSetterParams(field))
		w.Line(`    this.%s = %s;`, field.Name.CamelCase(), addFieldName(field))
		w.Line(`  }`)
	}
	w.EmptyLine()
	addObjectModelMethods(w.Indented(), model)
	w.Line(`}`)

	return &gen.TextFile{
		Path:    thePackage.GetPath(className + ".java"),
		Content: w.String(),
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
		w.Line(`  int result = Objects.hash(%s);`, JoinParams(hashParams))
		for _, param := range hashCodeParams {
			w.Line(`  result = 31 * result + Arrays.hashCode(%s);`, param)
		}
		w.Line(`  return result;`)
	} else if arrayFieldCount == 0 && nonArrayFieldCount > 0 {
		w.Line(`  return Objects.hash(%s);`, JoinParams(hashParams))
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
	w.Line(`  return String.format("%s{%s}", %s);`, model.Name.PascalCase(), JoinParams(formatParams), JoinParams(formatArgs))
	w.Line(`}`)
}

func generateEnumModel(model *spec.NamedModel, thePackage Module) *gen.TextFile {
	w := NewJavaWriter()
	w.Line(`package %s;`, thePackage.PackageName)
	w.EmptyLine()
	generateImports(w)
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

func generateOneOfModels(model *spec.NamedModel, thePackage Module) []gen.TextFile {
	files := []gen.TextFile{}
	w := NewJavaWriter()
	w.Line("package %s;", thePackage.PackageName)
	w.EmptyLine()
	generateImports(w)
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
		w.Line(`  @Type(value = %s%s.class, name = "%s"),`, model.Name.PascalCase(), item.Name.PascalCase(), item.Name.Source)
	}
	w.Line(`})`)
	interfaceName := model.Name.PascalCase()
	w.Line(`public interface %s {`, interfaceName)
	w.Line(`}`)

	for _, item := range model.OneOf.Items {
		files = append(files, *generateOneOfImplementation(item, model, thePackage))
	}

	files = append(files, gen.TextFile{
		Path:    thePackage.GetPath(interfaceName + ".java"),
		Content: w.String(),
	})
	return files
}

func generateOneOfImplementation(item spec.NamedDefinition, model *spec.NamedModel, thePackage Module) *gen.TextFile {
	className := model.Name.PascalCase() + item.Name.PascalCase()
	w := NewJavaWriter()
	w.Line(`package %s;`, thePackage.PackageName)
	w.EmptyLine()
	generateImports(w)
	w.EmptyLine()
	w.Line(`public class %s implements %s {`, className, model.Name.PascalCase())
	w.Line(`  @JsonUnwrapped`)
	w.Line(`  public %s data;`, JavaType(&item.Type.Definition))
	w.EmptyLine()
	w.Line(`  public %s() {`, className)
	w.Line(`  }`)
	w.EmptyLine()
	w.Line(`  public %s(%s data) {`, className, JavaType(&item.Type.Definition))
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
	addOneOfModelMethods(w, item, model)
	w.Unindent()
	w.Line(`}`)

	return &gen.TextFile{
		Path:    thePackage.GetPath(className + ".java"),
		Content: w.String(),
	}
}

func addOneOfModelMethods(w *gen.Writer, item spec.NamedDefinition, model *spec.NamedModel) {
	modelItemName := model.Name.PascalCase() + item.Name.PascalCase()
	w.Line(`@Override`)
	w.Line(`public boolean equals(Object o) {`)
	w.Line(`  if (this == o) return true;`)
	w.Line(`  if (!(o instanceof %s)) return false;`, modelItemName)
	w.Line(`  %s that = (%s) o;`, modelItemName, modelItemName)
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
