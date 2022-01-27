package models

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/genjava/imports"
	"github.com/specgen-io/specgen/v2/genjava/packages"
	"github.com/specgen-io/specgen/v2/genjava/types"
	"github.com/specgen-io/specgen/v2/genjava/writer"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

var Moshi = "moshi"

type MoshiGenerator struct {
	Types *types.Types
}

func NewMoshiGenerator(types *types.Types) *MoshiGenerator {
	return &MoshiGenerator{types}
}

func (g *MoshiGenerator) JsonImports() []string {
	return []string{`com.squareup.moshi.*`}
}

func (g *MoshiGenerator) CreateJsonMapperField(w *sources.Writer) {
	w.Line(`private Moshi moshi;`)
}

func (g *MoshiGenerator) InitJsonMapper(w *sources.Writer) {
	w.Line(`Moshi.Builder moshiBuilder = new Moshi.Builder();`)
	w.Line(`Json.setupMoshiAdapters(moshiBuilder);`)
	w.Line(`this.moshi = moshiBuilder.build();`)
}

func (g *MoshiGenerator) ReadJson(varJson string, typ *spec.TypeDef) (string, string) {
	adapter := fmt.Sprintf(`adapter(%s.class)`, g.Types.Java(typ))
	if typ.Node == spec.MapType {
		typeJava := g.Types.Java(typ.Child)
		adapter = fmt.Sprintf(`<Map<String, %s>>adapter(Types.newParameterizedType(Map.class, String.class, %s.class))`, typeJava, typeJava)
	}

	return fmt.Sprintf(`moshi.%s.fromJson(%s)`, adapter, varJson), `IOException`
}

func (g *MoshiGenerator) WriteJson(varData string, typ *spec.TypeDef) (string, string) {
	adapterParam := fmt.Sprintf(`%s.class`, g.Types.Java(typ))
	if typ.Node == spec.MapType {
		typeJava := g.Types.Java(typ.Child)
		adapterParam = fmt.Sprintf(`Types.newParameterizedType(Map.class, String.class, %s.class)`, typeJava)
	}

	return fmt.Sprintf(`moshi.adapter(%s).toJson(%s)`, adapterParam, varData), `AssertionError`
}

func (g *MoshiGenerator) VersionModels(version *spec.Version, thePackage packages.Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, model := range version.ResolvedModels {
		if model.IsObject() {
			files = append(files, *g.modelObject(model, thePackage))
		} else if model.IsOneOf() {
			files = append(files, *g.modelOneOf(model, thePackage))
		} else if model.IsEnum() {
			files = append(files, *g.modelEnum(model, thePackage))
		}
	}
	return files
}

func (g *MoshiGenerator) modelObject(model *spec.NamedModel, thePackage packages.Module) *sources.CodeFile {
	w := writer.NewJavaWriter()
	w.Line(`package %s;`, thePackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(g.JsonImports()...)
	imports.Add(g.Types.Imports()...)
	imports.Write(w)
	w.EmptyLine()
	className := model.Name.PascalCase()
	w.Line(`public class %s {`, className)
	for _, field := range model.Object.Fields {
		w.EmptyLine()
		w.Line(`  @Json(name = "%s")`, field.Name.Source)
		w.Line(`  private %s %s;`, g.Types.Java(&field.Type.Definition), field.Name.CamelCase())
	}
	w.EmptyLine()
	ctorParams := []string{}
	for _, field := range model.Object.Fields {
		ctorParams = append(ctorParams, fmt.Sprintf(`%s %s`, g.Types.Java(&field.Type.Definition), field.Name.CamelCase()))
	}
	w.Line(`  public %s(%s) {`, model.Name.PascalCase(), strings.Join(ctorParams, ", "))
	for _, field := range model.Object.Fields {
		w.Line(`    this.%s = %s;`, field.Name.CamelCase(), field.Name.CamelCase())
	}
	w.Line(`  }`)
	for _, field := range model.Object.Fields {
		w.EmptyLine()
		w.Line(`  public %s %s() {`, g.Types.Java(&field.Type.Definition), getterName(&field))
		w.Line(`    return %s;`, field.Name.CamelCase())
		w.Line(`  }`)
		w.EmptyLine()
		w.Line(`  public void %s(%s %s) {`, setterName(&field), g.Types.Java(&field.Type.Definition), field.Name.CamelCase())
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

func (g *MoshiGenerator) modelEnum(model *spec.NamedModel, thePackage packages.Module) *sources.CodeFile {
	w := writer.NewJavaWriter()
	w.Line(`package %s;`, thePackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(g.JsonImports()...)
	imports.Add(g.Types.Imports()...)
	imports.Write(w)
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

func (g *MoshiGenerator) modelOneOf(model *spec.NamedModel, thePackage packages.Module) *sources.CodeFile {
	interfaceName := model.Name.PascalCase()
	w := writer.NewJavaWriter()
	w.Line("package %s;", thePackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(g.JsonImports()...)
	imports.Add(g.Types.Imports()...)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`public interface %s {`, interfaceName)
	for index, item := range model.OneOf.Items {
		if index > 0 {
			w.EmptyLine()
		}
		g.modelOneOfImplementation(w.Indented(), &item, model)
	}
	w.Line(`}`)

	return &sources.CodeFile{
		Path:    thePackage.GetPath(interfaceName + ".java"),
		Content: w.String(),
	}
}

func (g *MoshiGenerator) modelOneOfImplementation(w *sources.Writer, item *spec.NamedDefinition, model *spec.NamedModel) {
	w.Line(`class %s implements %s {`, oneOfItemClassName(item), model.Name.PascalCase())
	w.Line(`  public %s data;`, g.Types.Java(&item.Type.Definition))
	w.EmptyLine()
	w.Line(`  public %s() {`, oneOfItemClassName(item))
	w.Line(`  }`)
	w.EmptyLine()
	w.Line(`  public %s(%s data) {`, oneOfItemClassName(item), g.Types.Java(&item.Type.Definition))
	w.Line(`  	this.data = data;`)
	w.Line(`  }`)
	w.EmptyLine()
	w.Line(`  public %s getData() {`, g.Types.Java(&item.Type.Definition))
	w.Line(`    return data;`)
	w.Line(`  }`)
	w.EmptyLine()
	w.Line(`  public void setData(%s data) {`, g.Types.Java(&item.Type.Definition))
	w.Line(`    this.data = data;`)
	w.Line(`  }`)
	w.EmptyLine()
	w.Indent()
	addOneOfModelMethods(w, item)
	w.Unindent()
	w.Line(`}`)
}

func (g *MoshiGenerator) SetupLibrary(thePackage packages.Module) []sources.CodeFile {
	adaptersPackage := thePackage.Subpackage("adapters")

	code := `
package [[.PackageName]];

import com.squareup.moshi.Moshi;
import [[.PackageName]].adapters.*;

public class Json {
	public static void setupMoshiAdapters(Moshi.Builder moshiBuilder) {
		moshiBuilder
			.add(new BigDecimalAdapter())
			.add(new UuidAdapter())
			.add(new LocalDateAdapter())
			.add(new LocalDateTimeAdapter());
	}
}
`

	code, _ = sources.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})

	files := []sources.CodeFile{}
	files = append(files, sources.CodeFile{
		Path:    thePackage.GetPath("Json.java"),
		Content: strings.TrimSpace(code),
	})
	files = append(files, *moshiBigDecimalAdapter(adaptersPackage))
	files = append(files, *moshiLocalDateAdapter(adaptersPackage))
	files = append(files, *moshiLocalDateTimeAdapter(adaptersPackage))
	files = append(files, *moshiUuidAdapter(adaptersPackage))
	return files
}

func moshiBigDecimalAdapter(thePackage packages.Module) *sources.CodeFile {
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

func moshiLocalDateAdapter(thePackage packages.Module) *sources.CodeFile {
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

func moshiLocalDateTimeAdapter(thePackage packages.Module) *sources.CodeFile {
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

func moshiUuidAdapter(thePackage packages.Module) *sources.CodeFile {
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
