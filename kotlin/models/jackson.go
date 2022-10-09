package models

import (
	"fmt"
	"kotlin/imports"
	"kotlin/writer"
	"strings"

	"generator"
	"kotlin/packages"
	"kotlin/types"
	"spec"
)

var Jackson = "jackson"

type JacksonGenerator struct {
	Types    *types.Types
	Packages *Packages
}

func NewJacksonGenerator(types *types.Types, packages *Packages) *JacksonGenerator {
	return &JacksonGenerator{types, packages}
}

func (g *JacksonGenerator) Models(version *spec.Version) []generator.CodeFile {
	return []generator.CodeFile{*g.models(version.ResolvedModels, g.Packages.Models(version))}
}

func (g *JacksonGenerator) ErrorModels(httperrors *spec.HttpErrors) []generator.CodeFile {
	return []generator.CodeFile{*g.models(httperrors.ResolvedModels, g.Packages.ErrorsModels)}
}

func (g *JacksonGenerator) models(models []*spec.NamedModel, modelsPackage packages.Package) *generator.CodeFile {
	w := writer.New(modelsPackage, `models`)
	imports := imports.New()
	imports.Add(g.ModelsDefinitionsImports()...)
	imports.Add(g.Types.Imports()...)
	imports.Write(w)

	for _, model := range models {
		w.EmptyLine()
		if model.IsObject() {
			g.modelObject(w, model)
		} else if model.IsOneOf() {
			g.modelOneOf(w, model)
		} else if model.IsEnum() {
			g.modelEnum(w, model)
		}
	}
	return w.ToCodeFile()
}

func jacksonPropertyAnnotation(field *spec.NamedDefinition) string {
	required := "false"
	if !field.Type.Definition.IsNullable() {
		required = "true"
	}
	return fmt.Sprintf(`@JsonProperty(value = "%s", required = %s)`, field.Name.Source, required)
}

func (g *JacksonGenerator) modelObject(w generator.Writer, model *spec.NamedModel) {
	className := model.Name.PascalCase()
	w.Line(`data class %s(`, className)
	for _, field := range model.Object.Fields {
		w.Line(`  %s`, jacksonPropertyAnnotation(&field))
		w.Line(`  val %s: %s,`, field.Name.CamelCase(), g.Types.Kotlin(&field.Type.Definition))
	}
	w.Line(`)`)
}

func (g *JacksonGenerator) modelEnum(w generator.Writer, model *spec.NamedModel) {
	enumName := model.Name.PascalCase()
	w.Line(`enum class %s {`, enumName)
	for _, enumItem := range model.Enum.Items {
		w.Line(`  @JsonProperty("%s") %s,`, enumItem.Value, enumItem.Name.UpperCase())
	}
	w.Line(`}`)
}

func (g *JacksonGenerator) modelOneOf(w generator.Writer, model *spec.NamedModel) {
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

func (g *JacksonGenerator) JsonRead(varJson string, typ *spec.TypeDef) string {
	return fmt.Sprintf(`%s, object : TypeReference<%s>() {}`, varJson, g.Types.Kotlin(typ))
}

func (g *JacksonGenerator) JsonWrite(varData string, typ *spec.TypeDef) string {
	return varData
}

func (g *JacksonGenerator) ReadJson(varJson string, typ *spec.TypeDef) (string, string) {
	return fmt.Sprintf(`objectMapper.readValue(%s, object: TypeReference<%s>(){})`, varJson, g.Types.Kotlin(typ)), `IOException`
}

func (g *JacksonGenerator) WriteJson(varData string, typ *spec.TypeDef) (string, string) {
	return fmt.Sprintf(`objectMapper.writeValueAsString(%s)`, varData), `Exception`
}

func (g *JacksonGenerator) ModelsDefinitionsImports() []string {
	return []string{
		`com.fasterxml.jackson.annotation.*`,
		`com.fasterxml.jackson.annotation.JsonSubTypes.*`,
		`com.fasterxml.jackson.core.type.*`,
		`com.fasterxml.jackson.core.JsonProcessingException`,
		`com.fasterxml.jackson.databind.*`,
		`com.fasterxml.jackson.module.kotlin.jacksonObjectMapper`,
	}
}

func (g *JacksonGenerator) ModelsUsageImports() []string {
	return []string{
		`com.fasterxml.jackson.databind.*`,
		`com.fasterxml.jackson.core.type.*`,
		`com.fasterxml.jackson.module.kotlin.jacksonObjectMapper`,
	}
}

func (g *JacksonGenerator) ValidationErrorsHelpers() *generator.CodeFile {
	code := `
package [[.PackageName]];

import com.fasterxml.jackson.databind.exc.InvalidFormatException
import [[.ErrorsModelsPackage]].*
import [[.JsonPackage]].*

object ValidationErrorsHelpers {
    fun extractValidationErrors(exception: JsonParseException): List<ValidationError>? {
        val causeException = exception.cause
        if (causeException is InvalidFormatException) {
            val jsonPath = getJsonPath(causeException)
            val validation = ValidationError(jsonPath, "parsing_failed", exception.message)
            return listOf(validation)
        }
        return null
    }

    private fun getJsonPath(exception: InvalidFormatException): String {
        val path = StringBuilder()
        for (i in exception.path.indices) {
            val reference = exception.path[i]
            if (reference.index != -1) {
                path.append("[").append(reference.index).append("]")
            } else {
                if (i != 0) {
                    path.append(".")
                }
                path.append(reference.fieldName)
            }
        }
        return path.toString()
    }
}
`

	code, _ = generator.ExecuteTemplate(code, struct {
		PackageName         string
		ErrorsModelsPackage string
		JsonPackage         string
	}{
		g.Packages.Errors.PackageName,
		g.Packages.ErrorsModels.PackageName,
		g.Packages.Json.PackageName,
	})
	return &generator.CodeFile{
		Path:    g.Packages.Errors.GetPath("ValidationErrorsHelpers.kt"),
		Content: strings.TrimSpace(code),
	}
}

func (g *JacksonGenerator) CreateJsonMapperField(annotation string) string {
	objectMapperVar := `private val objectMapper: ObjectMapper`
	if annotation != "" {
		return fmt.Sprintf(`@%s %s`, annotation, objectMapperVar)
	}
	return objectMapperVar
}

func (g *JacksonGenerator) InitJsonMapper(w generator.Writer) {
	w.Line(`objectMapper = setupObjectMapper(jacksonObjectMapper())`)
}

func (g *JacksonGenerator) JsonHelpersMethods() string {
	return `
	fun write(data: Any): String {
        return try {
            objectMapper.writeValueAsString(data)
        } catch (exception: IOException) {
            throw RuntimeException(exception)
        }
    }

    fun <T> read(jsonStr: String, typeReference: TypeReference<T>): T {
        return try {
            objectMapper.readValue(jsonStr, typeReference)
        } catch (exception: IOException) {
            throw JsonParseException(exception)
        }
    }
`
}

func (g *JacksonGenerator) SetupLibrary() []generator.CodeFile {
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

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{g.Packages.Json.PackageName})
	return []generator.CodeFile{{
		Path:    g.Packages.Json.GetPath("CustomObjectMapper.kt"),
		Content: strings.TrimSpace(code),
	}}
}
