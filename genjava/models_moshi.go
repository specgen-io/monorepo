package genjava

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

var Moshi = "moshi"

func addMoshiImports(w *sources.Writer) {
	w.Line(`import com.squareup.moshi.Json;`)
	w.Line(`import java.math.BigDecimal;`)
	w.Line(`import java.time.*;`)
	w.Line(`import java.util.*;`)
}

type MoshiGenerator struct {
	Type *Types
}

func NewMoshiGenerator(types *Types) *MoshiGenerator {
	return &MoshiGenerator{types}
}

func (g *MoshiGenerator) ReadJson(varJson string, typeJava string) string {
	panic("This is not implemented yet!!!")
}

func (g *MoshiGenerator) WriteJson(varData string) string {
	panic("This is not implemented yet!!!")
}

func (g *MoshiGenerator) VersionModels(version *spec.Version, thePackage Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, model := range version.ResolvedModels {
		if model.IsObject() {
			files = append(files, *g.ObjectModel(model, thePackage))
		} else if model.IsOneOf() {
			files = append(files, *g.OneOfModel(model, thePackage))
		} else if model.IsEnum() {
			files = append(files, *g.EnumModel(model, thePackage))
		}
	}
	return files
}

func (g *MoshiGenerator) ObjectModel(model *spec.NamedModel, thePackage Module) *sources.CodeFile {
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
		w.Line(`  private %s %s;`, g.Type.JavaType(&field.Type.Definition), field.Name.CamelCase())
	}
	w.EmptyLine()
	ctorParams := []string{}
	for _, field := range model.Object.Fields {
		ctorParams = append(ctorParams, fmt.Sprintf(`%s %s`, g.Type.JavaType(&field.Type.Definition), field.Name.CamelCase()))
	}
	w.Line(`  public %s(%s) {`, model.Name.PascalCase(), JoinDelimParams(ctorParams))
	for _, field := range model.Object.Fields {
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

func (g *MoshiGenerator) EnumModel(model *spec.NamedModel, thePackage Module) *sources.CodeFile {
	w := NewJavaWriter()
	w.Line(`package %s;`, thePackage.PackageName)
	w.EmptyLine()
	addMoshiImports(w)
	w.EmptyLine()
	enumName := model.Name.PascalCase()
	w.Line(`public enum %s {`, enumName)
	for _, enumItem := range model.Enum.Items {
		w.Line(`  @Json(name = "%s") %s,`, enumItem.Value, enumItem.Name.UpperCase())
	}
	w.Line(`}`)

	return &sources.CodeFile{
		Path:    thePackage.GetPath(enumName + ".java"),
		Content: w.String(),
	}
}

func (g *MoshiGenerator) OneOfModel(model *spec.NamedModel, thePackage Module) *sources.CodeFile {
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
		g.oneOfImplementation(w.Indented(), &item, model)
	}
	w.Line(`}`)

	return &sources.CodeFile{
		Path:    thePackage.GetPath(interfaceName + ".java"),
		Content: w.String(),
	}
}

func (g *MoshiGenerator) oneOfImplementation(w *sources.Writer, item *spec.NamedDefinition, model *spec.NamedModel) {
	w.Line(`class %s implements %s {`, oneOfItemClassName(item), model.Name.PascalCase())
	w.Line(`  public %s data;`, g.Type.JavaType(&item.Type.Definition))
	w.EmptyLine()
	w.Line(`  public %s() {`, oneOfItemClassName(item))
	w.Line(`  }`)
	w.EmptyLine()
	w.Line(`  public %s(%s data) {`, oneOfItemClassName(item), g.Type.JavaType(&item.Type.Definition))
	w.Line(`  	this.data = data;`)
	w.Line(`  }`)
	w.EmptyLine()
	w.Line(`  public %s getData() {`, g.Type.JavaType(&item.Type.Definition))
	w.Line(`    return data;`)
	w.Line(`  }`)
	w.EmptyLine()
	w.Line(`  public void setData(%s data) {`, g.Type.JavaType(&item.Type.Definition))
	w.Line(`    this.data = data;`)
	w.Line(`  }`)
	w.EmptyLine()
	w.Indent()
	addOneOfModelMethods(w, item)
	w.Unindent()
	w.Line(`}`)
}

func (g *MoshiGenerator) SetupLibrary(thePackage Module) []sources.CodeFile {
	adaptersPackage := thePackage.Subpackage("adapters")
	files := []sources.CodeFile{}
	files = append(files, *generateBigDecimalAdapter(adaptersPackage))
	files = append(files, *generateLocalDateAdapter(adaptersPackage))
	files = append(files, *generateLocalDateTimeAdapter(adaptersPackage))
	files = append(files, *generateUuidAdapter(adaptersPackage))
	return files
}

func generateBigDecimalAdapter(thePackage Module) *sources.CodeFile {
	code := `
package [[.PackageName]];

import com.squareup.moshi.*;
import okio.Okio;

import java.io.*;
import java.math.BigDecimal;
import java.nio.charset.StandardCharsets;

public class BigDecimalAdapter {
	@FromJson
	public BigDecimal fromJson(JsonReader reader) throws IOException {
		var token = reader.peek();
		if (token != JsonReader.Token.NUMBER) {
			throw new JsonDataException("BigDecimal should be represented as number in JSON, found: "+token.name());
		}
		var source = reader.nextSource();
		return new BigDecimal(new String(source.readByteArray(), StandardCharsets.UTF_8));
	}

	@ToJson
	public void toJson(JsonWriter writer, BigDecimal value) throws IOException {
		var source = Okio.source(new ByteArrayInputStream(value.toString().getBytes()));
		var buffer = Okio.buffer(source);
		writer.value(buffer);
	}
}
`

	code, _ = sources.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &sources.CodeFile{
		Path:    thePackage.GetPath("BigDecimalAdapter.java"),
		Content: strings.TrimSpace(code),
	}
}

func generateLocalDateAdapter(thePackage Module) *sources.CodeFile {
	code := `
package [[.PackageName]];

import com.squareup.moshi.*;

import java.time.LocalDate;

public class LocalDateAdapter {
	@FromJson
	private LocalDate fromJson(String string) {
		return LocalDate.parse(string);
	}

	@ToJson
	private String toJson(LocalDate value) {
		return value.toString();
	}
}
`

	code, _ = sources.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &sources.CodeFile{
		Path:    thePackage.GetPath("LocalDateAdapter.java"),
		Content: strings.TrimSpace(code),
	}
}

func generateLocalDateTimeAdapter(thePackage Module) *sources.CodeFile {
	code := `
package [[.PackageName]];

import com.squareup.moshi.*;

import java.time.LocalDateTime;

public class LocalDateTimeAdapter {
	@FromJson
	private LocalDateTime fromJson(String string) {
		return LocalDateTime.parse(string);
	}

	@ToJson
	private String toJson(LocalDateTime value) {
		return value.toString();
	}
}
`

	code, _ = sources.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &sources.CodeFile{
		Path:    thePackage.GetPath("LocalDateTimeAdapter.java"),
		Content: strings.TrimSpace(code),
	}
}

func generateUuidAdapter(thePackage Module) *sources.CodeFile {
	code := `
package [[.PackageName]];

import com.squareup.moshi.*;

import java.util.UUID;

public class UuidAdapter {
	@FromJson
	private UUID fromJson(String string) {
		return UUID.fromString(string);
	}

	@ToJson
	private String toJson(UUID value) {
		return value.toString();
	}
}
`

	code, _ = sources.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &sources.CodeFile{
		Path:    thePackage.GetPath("UuidAdapter.java"),
		Content: strings.TrimSpace(code),
	}
}
