package models

import (
	"github.com/specgen-io/specgen/v2/genkotlin/imports"
	"github.com/specgen-io/specgen/v2/genkotlin/modules"
	"github.com/specgen-io/specgen/v2/genkotlin/types"
	"github.com/specgen-io/specgen/v2/genkotlin/writer"
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
	return []string{
		`com.squareup.moshi.*`,
	}
}

func (g *MoshiGenerator) InitJsonMapper(w *sources.Writer) {
	panic("This is not implemented yet!!!")
}

func (g *MoshiGenerator) ReadJson(varJson string, typ *spec.TypeDef) (string, string) {
	panic("This is not implemented yet!!!")
}

func (g *MoshiGenerator) WriteJson(varData string, typ *spec.TypeDef) (string, string) {
	panic("This is not implemented yet!!!")
}

func (g *MoshiGenerator) VersionModels(version *spec.Version, thePackage modules.Module) *sources.CodeFile {
	w := writer.NewKotlinWriter()
	w.Line(`package %s`, thePackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(g.JsonImports()...)
	imports.Add(g.Types.Imports()...)
	imports.Write(w)

	for _, model := range version.ResolvedModels {
		w.EmptyLine()
		if model.IsObject() {
			g.modelObject(w, model)
		} else if model.IsOneOf() {
			//g.modelOneOf(w, model)
		} else if model.IsEnum() {
			g.modelEnum(w, model)
		}
	}

	return &sources.CodeFile{
		Path:    thePackage.GetPath("models.kt"),
		Content: w.String(),
	}
}

func (g *MoshiGenerator) modelObject(w *sources.Writer, model *spec.NamedModel) {
	className := model.Name.PascalCase()
	w.Line(`data class %s(`, className)
	for _, field := range model.Object.Fields {
		w.Line(`  @Json(name = "%s")`, field.Name.Source)
		w.Line(`  val %s: %s,`, field.Name.CamelCase(), g.Types.Kotlin(&field.Type.Definition))
	}

	if isKotlinArrayType(model) {
		w.Line(`) {`)
		addObjectModelMethods(w.Indented(), model)
		w.Line(`}`)
	} else {
		w.Line(`)`)
	}
}

func (g *MoshiGenerator) modelEnum(w *sources.Writer, model *spec.NamedModel) {
	enumName := model.Name.PascalCase()
	w.Line(`enum class %s {`, enumName)
	for _, enumItem := range model.Enum.Items {
		w.Line(`  @Json(name = "%s")`, enumItem.Value)
		w.Line(`  %s,`, enumItem.Name.UpperCase())
	}
	w.Line(`}`)
}

func (g *MoshiGenerator) SetupLibrary(thePackage modules.Module) []sources.CodeFile {
	adaptersPackage := thePackage.Subpackage("adapters")

	files := []sources.CodeFile{}
	files = append(files, *g.setupAdapters(thePackage, adaptersPackage))
	files = append(files, *bigDecimalAdapter(adaptersPackage))
	files = append(files, *localDateAdapter(adaptersPackage))
	files = append(files, *localDateTimeAdapter(adaptersPackage))
	files = append(files, *uuidAdapter(adaptersPackage))
	return files
}

func (g *MoshiGenerator) setupAdapters(thePackage modules.Module, adaptersPackage modules.Module) *sources.CodeFile {
	w := writer.NewKotlinWriter()
	w.Line("package %s", thePackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(`com.squareup.moshi.Moshi`)
	imports.Add(`com.squareup.moshi.kotlin.reflect.KotlinJsonAdapterFactory`)
	imports.Add(adaptersPackage.PackageStar)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`object Json {`)
	w.Line(`  fun setupMoshiAdapters(moshiBuilder: Moshi.Builder) {`)
	w.Line(`    moshiBuilder`)
	w.Line(`      .add(BigDecimalAdapter())`)
	w.Line(`      .add(UuidAdapter())`)
	w.Line(`      .add(LocalDateAdapter())`)
	w.Line(`      .add(LocalDateTimeAdapter())`)
	w.Line(`      .add(KotlinJsonAdapterFactory())`)
	w.EmptyLine()
	w.Line(`  }`)
	w.Line(`}`)

	return &sources.CodeFile{
		Path:    thePackage.GetPath("Json.kt"),
		Content: w.String(),
	}
}

func bigDecimalAdapter(thePackage modules.Module) *sources.CodeFile {
	code := `
package [[.PackageName]];

import com.squareup.moshi.*
import okio.*

import java.io.ByteArrayInputStream
import java.math.BigDecimal
import java.nio.charset.StandardCharsets

class BigDecimalAdapter {
    @FromJson
    fun fromJson(reader: JsonReader): BigDecimal {
        val token = reader.peek()
        if (token != JsonReader.Token.NUMBER) {
            throw JsonDataException("BigDecimal should be represented as number in JSON, found: " + token.name)
        }
        val source = reader.nextSource()
        return BigDecimal(String(source.readByteArray(), StandardCharsets.UTF_8))
    }

    @ToJson
    fun toJson(writer: JsonWriter, value: BigDecimal) {
        val source = ByteArrayInputStream(value.toString().toByteArray()).source()
        val buffer = source.buffer()
        writer.value(buffer)
    }
}
`

	code, _ = sources.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &sources.CodeFile{
		Path:    thePackage.GetPath("BigDecimalAdapter.kt"),
		Content: strings.TrimSpace(code),
	}
}

func localDateAdapter(thePackage modules.Module) *sources.CodeFile {
	code := `
package [[.PackageName]];

import com.squareup.moshi.*

import java.time.*

class LocalDateAdapter {
    @FromJson
    private fun fromJson(string: String): LocalDate {
        return LocalDate.parse(string)
    }

    @ToJson
    private fun toJson(value: LocalDate): String {
        return value.toString()
    }
}
`

	code, _ = sources.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &sources.CodeFile{
		Path:    thePackage.GetPath("LocalDateAdapter.kt"),
		Content: strings.TrimSpace(code),
	}
}

func localDateTimeAdapter(thePackage modules.Module) *sources.CodeFile {
	code := `
package [[.PackageName]];

import com.squareup.moshi.*

import java.time.*

class LocalDateTimeAdapter {
    @FromJson
    private fun fromJson(string: String): LocalDateTime {
        return LocalDateTime.parse(string)
    }

    @ToJson
    private fun toJson(value: LocalDateTime): String {
        return value.toString()
    }
}
`

	code, _ = sources.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &sources.CodeFile{
		Path:    thePackage.GetPath("LocalDateTimeAdapter.kt"),
		Content: strings.TrimSpace(code),
	}
}

func uuidAdapter(thePackage modules.Module) *sources.CodeFile {
	code := `
package [[.PackageName]];

import com.squareup.moshi.*

import java.util.UUID

class UuidAdapter {
    @FromJson
    private fun fromJson(string: String): UUID {
        return UUID.fromString(string)
    }

    @ToJson
    private fun toJson(value: UUID): String {
        return value.toString()
    }
}
`

	code, _ = sources.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &sources.CodeFile{
		Path:    thePackage.GetPath("UuidAdapter.kt"),
		Content: strings.TrimSpace(code),
	}
}
