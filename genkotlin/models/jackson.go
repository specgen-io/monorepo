package models

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/genkotlin/imports"
	"github.com/specgen-io/specgen/v2/genkotlin/modules"
	"github.com/specgen-io/specgen/v2/genkotlin/types"
	"github.com/specgen-io/specgen/v2/genkotlin/writer"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

var Jackson = "jackson"

type JacksonGenerator struct {
	Types *types.Types
}

func NewJacksonGenerator(types *types.Types) *JacksonGenerator {
	return &JacksonGenerator{types}
}

func (g *JacksonGenerator) JsonImports() []string {
	return []string{
		`com.fasterxml.jackson.databind.*`,
		`com.fasterxml.jackson.annotation.*`,
		`com.fasterxml.jackson.annotation.JsonSubTypes.*`,
		`com.fasterxml.jackson.core.type.*`,
		`com.fasterxml.jackson.core.JsonProcessingException`,
		`com.fasterxml.jackson.module.kotlin.jacksonObjectMapper`,
	}
}

func (g *JacksonGenerator) SetupImport(jsonPackage modules.Module) string {
	return fmt.Sprintf(`%s.setupObjectMapper`, jsonPackage.PackageName)
}

func (g *JacksonGenerator) CreateJsonMapperVar(w *sources.Writer) {
	w.Line(`private var objectMapper: ObjectMapper`)
}

func (g *JacksonGenerator) InitJsonMapper(w *sources.Writer) {
	w.Line(`objectMapper = setupObjectMapper(jacksonObjectMapper()))`)
}

func (g *JacksonGenerator) ReadJson(varJson string, typ *spec.TypeDef) (string, string) {
	return fmt.Sprintf(`objectMapper.readValue(%s, object: TypeReference<%s>(){})`, varJson, g.Types.Kotlin(typ)), `IOException`
}

func (g *JacksonGenerator) WriteJson(varData string, typ *spec.TypeDef) (string, string) {
	return fmt.Sprintf(`objectMapper.writeValueAsString(%s)`, varData), `JsonProcessingException`
}

func (g *JacksonGenerator) SetupLibrary(thePackage modules.Module) []sources.CodeFile {
	code := `
package [[.PackageName]]

import com.fasterxml.jackson.annotation.JsonInclude
import com.fasterxml.jackson.databind.*
import com.fasterxml.jackson.datatype.jsr310.*

fun setupObjectMapper(objectMapper: ObjectMapper): ObjectMapper {
    objectMapper
        .registerModule(JavaTimeModule())
        .configure(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS, false)
        .setSerializationInclusion(JsonInclude.Include.NON_NULL)
    return objectMapper
}
`

	code, _ = sources.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})

	files := []sources.CodeFile{}
	files = append(files, sources.CodeFile{
		Path:    thePackage.GetPath("Json.kt"),
		Content: strings.TrimSpace(code),
	})

	return files
}

func (g *JacksonGenerator) VersionModels(version *spec.Version, thePackage modules.Module, jsonPackage modules.Module) []sources.CodeFile {
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
			g.modelOneOf(w, model)
		} else if model.IsEnum() {
			g.modelEnum(w, model)
		}
	}

	files := []sources.CodeFile{}
	files = append(files, sources.CodeFile{Path: thePackage.GetPath("models.kt"), Content: w.String()})

	return files
}

func (g *JacksonGenerator) jacksonPropertyAnnotation(field *spec.NamedDefinition) string {
	required := "false"
	if !field.Type.Definition.IsNullable() {
		required = "true"
	}
	return fmt.Sprintf(`@JsonProperty(value = "%s", required = %s)`, field.Name.Source, required)
}

func (g *JacksonGenerator) modelObject(w *sources.Writer, model *spec.NamedModel) {
	className := model.Name.PascalCase()
	w.Line(`data class %s(`, className)
	for _, field := range model.Object.Fields {
		w.Line(`  %s`, g.jacksonPropertyAnnotation(&field))
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

func (g *JacksonGenerator) modelEnum(w *sources.Writer, model *spec.NamedModel) {
	enumName := model.Name.PascalCase()
	w.Line(`enum class %s {`, enumName)
	for _, enumItem := range model.Enum.Items {
		w.Line(`  @JsonProperty("%s") %s,`, enumItem.Value, enumItem.Name.UpperCase())
	}
	w.Line(`}`)
}

func (g *JacksonGenerator) modelOneOf(w *sources.Writer, model *spec.NamedModel) {
	interfaceName := model.Name.PascalCase()
	if model.OneOf.Discriminator != nil {
		w.Line(`@JsonTypeInfo(use = JsonTypeInfo.Id.NAME, include = JsonTypeInfo.As.PROPERTY, property = "%s")`, *model.OneOf.Discriminator)
	} else {
		w.Line(`@JsonTypeInfo(use = JsonTypeInfo.Id.NAME, include = JsonTypeInfo.As.WRAPPER_OBJECT)`)
	}
	w.Line(`sealed interface %s {`, interfaceName)
	for _, item := range model.OneOf.Items {
		itemType := g.Types.Kotlin(&item.Type.Definition)
		w.Line(`  @JsonTypeName("%s")`, item.Name.Source)
		w.Line(`  data class %s(@field:JsonIgnore val data: %s): %s {`, item.Name.PascalCase(), itemType, interfaceName)
		w.Line(`    @get:JsonUnwrapped`)
		w.Line(`    private val _data: %s`, itemType)
		w.Line(`      get() = data`)
		w.Line(`  }`)
	}
	w.Line(`}`)
}
