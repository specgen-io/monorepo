package genjava

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
	"strings"
)

func GenerateModels(serviceFile string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}
	files := generateModels(specification, generatePath)

	return gen.WriteFiles(files, true)
}

func generateModels(specification *spec.Spec, generatePath string) []gen.TextFile {
	files := []gen.TextFile{}
	for _, version := range specification.Versions {
		packageName := fmt.Sprintf("%s.models", specification.Name.SnakeCase())
		versionedPackageName := versionedPackage(version.Version, "spec")
		versionPath := filepath.Join(generatePath, versionedPackageName)
		files = append(files, generateVersionModels(&version, fmt.Sprintf(`%s.%s`, packageName, versionedPackageName), versionPath)...)
		files = append(files, *generateJsoner(packageName, filepath.Join(generatePath, "Jsoner.java")))
	}

	return files
}

func generateJsoner(packageName string, path string) *gen.TextFile {
	code := `
package [[.PackageName]];

import com.fasterxml.jackson.databind.*;
import com.fasterxml.jackson.datatype.jsr310.*;

import java.io.*;

public class Jsoner {

	public static <T> String serialize(ObjectMapper objectMapper, T data) throws IOException {
		ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream();
		objectMapper.registerModule(new JavaTimeModule())
			.configure(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS, false)
			.writeValue(byteArrayOutputStream, data);

		return byteArrayOutputStream.toString();
	}

	public static <T> T deserialize(ObjectMapper objectMapper, String jsonStr, Class<T> tClass) throws IOException {
		InputStream inputStream = new ByteArrayInputStream(jsonStr.getBytes());

		return objectMapper.readValue(inputStream, tClass);
	}
}
`

	code, _ = gen.ExecuteTemplate(code, struct{ PackageName string }{packageName})
	return &gen.TextFile{Path: path, Content: strings.TrimSpace(code)}
}

func generateVersionModels(version *spec.Version, packageName string, generatePath string) []gen.TextFile {
	files := []gen.TextFile{}
	for _, model := range version.ResolvedModels {
		if model.IsObject() {
			files = append(files, *generateObjectModel(model, packageName, generatePath))
		} else if model.IsOneOf() {
			files = append(files, generateOneOfModels(model, packageName, generatePath)...)
		} else if model.IsEnum() {
			files = append(files, *generateEnumModel(model, packageName, generatePath))
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

func generateObjectModel(model *spec.NamedModel, packageName string, generatePath string) *gen.TextFile {
	w := NewJavaWriter()
	w.Line("package %s;", packageName)
	w.EmptyLine()
	generateImports(w)
	w.EmptyLine()
	w.Line(`public class %s {`, model.Name.PascalCase())

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
		w.Line(`  public %s get%s() {`, addType(field), field.Name.PascalCase())
		w.Line(`    return %s;`, addGetterBody(field))
		w.Line(`  }`)
		w.EmptyLine()
		w.Line(`  public void set%s(%s) {`, field.Name.PascalCase(), addSetterParams(field))
		w.Line(`    this.%s = %s;`, field.Name.CamelCase(), addFieldName(field))
		w.Line(`  }`)
	}
	w.Line(`}`)

	return &gen.TextFile{Path: filepath.Join(generatePath, fmt.Sprintf("%s.java", model.Name.Source)), Content: w.String()}
}

func generateEnumModel(model *spec.NamedModel, packageName string, generatePath string) *gen.TextFile {
	w := NewJavaWriter()
	w.Line("package %s;", packageName)
	w.EmptyLine()
	generateImports(w)
	w.EmptyLine()
	w.Line(`public enum %s {`, model.Name.PascalCase())
	for _, enumItem := range model.Enum.Items {
		w.Line(`  %s,`, enumItem.Value)
	}
	w.Line(`}`)

	return &gen.TextFile{Path: filepath.Join(generatePath, fmt.Sprintf("%s.java", model.Name.Source)), Content: w.String()}
}

func generateOneOfModels(model *spec.NamedModel, packageName string, generatePath string) []gen.TextFile {
	files := []gen.TextFile{}
	w := NewJavaWriter()
	w.Line("package %s;", packageName)
	w.EmptyLine()
	generateImports(w)
	w.EmptyLine()
	w.Line(`@JsonTypeInfo(`)
	w.Line(`  use = JsonTypeInfo.Id.NAME,`)
	if model.OneOf.Discriminator != nil {
		w.Line(`  include = JsonTypeInfo.As.PROPERTY`)
	} else {
		w.Line(`  include = JsonTypeInfo.As.WRAPPER_OBJECT`)
	}
	w.Line(`)`)
	w.Line(`@JsonSubTypes({`)
	for _, item := range model.OneOf.Items {
		w.Line(`  @Type(value = %s%s.class, name = "%s"),`, model.Name.PascalCase(), item.Name.PascalCase(), item.Name.Source)
	}
	w.Line(`})`)
	w.Line(`public interface %s {`, model.Name.PascalCase())
	w.Line(`}`)

	for _, item := range model.OneOf.Items {
		files = append(files, *generateOneOfImplementation(item, model, packageName, generatePath))
	}
	files = append(files, gen.TextFile{Path: filepath.Join(generatePath, fmt.Sprintf("%s.java", model.Name.Source)), Content: w.String()})

	return files
}

func generateOneOfImplementation(item spec.NamedDefinition, model *spec.NamedModel, packageName string, generatePath string) *gen.TextFile {
	w := NewJavaWriter()
	w.Line("package %s;", packageName)
	w.EmptyLine()
	generateImports(w)
	w.EmptyLine()
	w.Line(`public class %s%s implements %s {`, model.Name.PascalCase(), item.Name.PascalCase(), model.Name.PascalCase())
	w.Line(`  @JsonUnwrapped`)
	w.Line(`  public %s data;`, item.Type.Definition.Name)
	w.EmptyLine()
	w.Line(`  public %s%s() {`, model.Name.PascalCase(), item.Name.PascalCase())
	w.Line(`  }`)
	w.EmptyLine()
	w.Line(`  public %s%s(%s data) {`, model.Name.PascalCase(), item.Name.PascalCase(), item.Type.Definition.Name)
	w.Line(`  	this.data = data;`)
	w.Line(`  }`)
	w.Line(`}`)

	return &gen.TextFile{Path: filepath.Join(generatePath, fmt.Sprintf("%s.java", model.Name.Source+item.Name.PascalCase())), Content: w.String()}
}
