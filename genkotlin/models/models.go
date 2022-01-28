package models

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/genkotlin/modules"
	"github.com/specgen-io/specgen/v2/genkotlin/types"
	"github.com/specgen-io/specgen/v2/genkotlin/writer"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func GenerateModels(specification *spec.Spec, packageName string, generatePath string) *sources.Sources {
	sources := sources.NewSources()
	if packageName == "" {
		packageName = specification.Name.SnakeCase()
	}
	mainPackage := modules.Package(generatePath, packageName)

	jsonPackage := mainPackage.Subpackage("json")
	sources.AddGenerated(GenerateJson(jsonPackage))

	for _, version := range specification.Versions {
		versionPackage := mainPackage.Subpackage(version.Version.FlatCase())

		modelsVersionPackage := versionPackage.Subpackage("models")
		sources.AddGenerated(GenerateVersionModels(&version, modelsVersionPackage))
	}

	return sources
}

func GenerateJson(thePackage modules.Module) *sources.CodeFile {
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
	return &sources.CodeFile{
		Path:    thePackage.GetPath("Json.kt"),
		Content: strings.TrimSpace(code),
	}
}

func GenerateVersionModels(version *spec.Version, thePackage modules.Module) *sources.CodeFile {
	w := writer.NewKotlinWriter()
	w.Line(`package %s`, thePackage.PackageName)
	w.EmptyLine()
	addImports(w)

	for _, model := range version.ResolvedModels {
		w.EmptyLine()
		if model.IsObject() {
			generateObjectModel(w, model)
		} else if model.IsOneOf() {
			generateOneOfModels(w, model)
		} else if model.IsEnum() {
			generateEnumModel(w, model)
		}
	}

	return &sources.CodeFile{
		Path:    thePackage.GetPath("models.kt"),
		Content: w.String(),
	}
}

func addImports(w *sources.Writer) {
	w.Line(`import java.time.*`)
	w.Line(`import java.util.*`)
	w.Line(`import java.math.BigDecimal`)
	w.Line(`import com.fasterxml.jackson.databind.JsonNode`)
	w.Line(`import com.fasterxml.jackson.annotation.*`)
}

func getJsonPropertyAnnotation(field *spec.NamedDefinition) string {
	required := "false"
	if !field.Type.Definition.IsNullable() {
		required = "true"
	}
	return fmt.Sprintf(`@JsonProperty(value = "%s", required = %s)`, field.Name.Source, required)
}

func generateObjectModel(w *sources.Writer, model *spec.NamedModel) {
	className := model.Name.PascalCase()
	w.Line(`data class %s(`, className)
	for _, field := range model.Object.Fields {
		w.Line(`  %s`, getJsonPropertyAnnotation(&field))
		w.Line(`  val %s: %s,`, field.Name.CamelCase(), types.KotlinType(&field.Type.Definition))
	}

	if isKotlinArrayType(model) {
		w.Line(`) {`)
		generateObjectModelMethods(w.Indented(), model)
		w.Line(`}`)
	} else {
		w.Line(`)`)
	}
}

func generateObjectModelMethods(w *sources.Writer, model *spec.NamedModel) {
	w.Line(`override fun equals(other: Any?): Boolean {`)
	w.Line(`  if (this === other) return true`)
	w.Line(`  if (other !is %s) return false`, model.Name.PascalCase())
	w.EmptyLine()
	for _, field := range model.Object.Fields {
		if _isKotlinArrayType(&field.Type.Definition) {
			w.Line(`  if (!%s.contentEquals(other.%s)) return false`, field.Name.CamelCase(), field.Name.CamelCase())
		} else {
			w.Line(`  if (%s != other.%s) return false`, field.Name.CamelCase(), field.Name.CamelCase())
		}
	}
	w.EmptyLine()
	w.Line(`  return true`)
	w.Line(`}`)
	w.EmptyLine()
	w.Line(`override fun hashCode(): Int {`)
	hashCodeArrayParams := []string{}
	hashCodeNonArrayParams := []string{}
	var arrayFieldCount, nonArrayFieldCount int
	for _, field := range model.Object.Fields {
		if _isKotlinArrayType(&field.Type.Definition) {
			hashCodeArrayParams = append(hashCodeArrayParams, fmt.Sprintf(`%s.contentHashCode()`, field.Name.CamelCase()))
			arrayFieldCount++
		} else {
			if field.Type.Definition.IsNullable() {
				hashCodeNonArrayParams = append(hashCodeNonArrayParams, fmt.Sprintf(`(%s?.hashCode() ?: 0)`, field.Name.CamelCase()))
			} else {
				hashCodeNonArrayParams = append(hashCodeNonArrayParams, fmt.Sprintf(`%s.hashCode()`, field.Name.CamelCase()))
			}
			nonArrayFieldCount++
		}
	}
	if arrayFieldCount == 1 && nonArrayFieldCount == 0 {
		w.Line(`  return %s`, hashCodeArrayParams[0])
	} else if arrayFieldCount > 1 && nonArrayFieldCount == 0 {
		w.Line(`  var result = %s`, hashCodeArrayParams[0])
		for _, param := range hashCodeArrayParams[1:] {
			w.Line(`  result = 31 * result + %s`, param)
		}
		w.Line(`  return result`)
	} else if arrayFieldCount > 0 && nonArrayFieldCount > 0 {
		w.Line(`  var result = %s`, hashCodeArrayParams[0])
		for _, param := range hashCodeArrayParams[1:] {
			w.Line(`  result = 31 * result + %s`, param)
		}
		for _, param := range hashCodeNonArrayParams {
			w.Line(`  result = 31 * result + %s`, param)
		}
		w.Line(`  return result`)
	}
	w.Line(`}`)
}

func generateEnumModel(w *sources.Writer, model *spec.NamedModel) {
	enumName := model.Name.PascalCase()
	w.Line(`enum class %s {`, enumName)
	for _, enumItem := range model.Enum.Items {
		w.Line(`  @JsonProperty("%s") %s,`, enumItem.Value, enumItem.Name.UpperCase())
	}
	w.Line(`}`)
}

func generateOneOfModels(w *sources.Writer, model *spec.NamedModel) {
	interfaceName := model.Name.PascalCase()
	if model.OneOf.Discriminator != nil {
		w.Line(`@JsonTypeInfo(use = JsonTypeInfo.Id.NAME, include = JsonTypeInfo.As.PROPERTY, property = "%s")`, *model.OneOf.Discriminator)
	} else {
		w.Line(`@JsonTypeInfo(use = JsonTypeInfo.Id.NAME, include = JsonTypeInfo.As.WRAPPER_OBJECT)`)
	}
	w.Line(`sealed interface %s {`, interfaceName)
	for _, item := range model.OneOf.Items {
		itemType := types.KotlinType(&item.Type.Definition)
		w.Line(`  @JsonTypeName("%s")`, item.Name.Source)
		w.Line(`  data class %s(@field:JsonIgnore val data: %s): %s {`, item.Name.PascalCase(), itemType, interfaceName)
		w.Line(`    @get:JsonUnwrapped`)
		w.Line(`    private val _data: %s`, itemType)
		w.Line(`      get() = data`)
		w.Line(`  }`)
	}
	w.Line(`}`)
}
