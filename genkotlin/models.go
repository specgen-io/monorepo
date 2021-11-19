package genkotlin

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func GenerateModels(specification *spec.Spec, packageName string, generatePath string) error {
	if packageName == "" {
		packageName = specification.Name.SnakeCase()
	}
	mainPackage := Package(generatePath, packageName)
	files := generateModels(specification, mainPackage)
	return gen.WriteFiles(files, true)
}

func generateModels(specification *spec.Spec, thePackage Module) []gen.TextFile {
	files := []gen.TextFile{}
	for _, version := range specification.Versions {
		versionPackage := thePackage.Subpackage(version.Version.FlatCase())
		files = append(files, *generateVersionModels(&version, versionPackage))
	}
	files = append(files, *generateJson(thePackage))
	return files
}

func generateJson(thePackage Module) *gen.TextFile {
	code := `
package [[.PackageName]]

import com.fasterxml.jackson.annotation.JsonInclude
import com.fasterxml.jackson.databind.*
import com.fasterxml.jackson.datatype.jsr310.*

fun setupObjectMapper(objectMapper: ObjectMapper): ObjectMapper  {
    objectMapper
        .registerModule(JavaTimeModule())
        .configure(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS, false)
        .setSerializationInclusion(JsonInclude.Include.NON_NULL)
    return objectMapper
}
`

	code, _ = gen.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &gen.TextFile{
		Path:    thePackage.GetPath("Json.kt"),
		Content: strings.TrimSpace(code),
	}
}

func generateVersionModels(version *spec.Version, thePackage Module) *gen.TextFile {
	w := NewKotlinWriter()
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

	return &gen.TextFile{
		Path:    thePackage.GetPath("models.kt"),
		Content: w.String(),
	}
}

func addImports(w *gen.Writer) {
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

func generateObjectModel(w *gen.Writer, model *spec.NamedModel) {
	className := model.Name.PascalCase()
	w.Line(`data class %s(`, className)
	for _, field := range model.Object.Fields {
		w.Line(`  %s`, getJsonPropertyAnnotation(&field))
		w.Line(`  val %s: %s,`, field.Name.CamelCase(), KotlinType(&field.Type.Definition))
	}
	w.Line(`)`)
}

func generateEnumModel(w *gen.Writer, model *spec.NamedModel) {
	enumName := model.Name.PascalCase()
	w.Line(`enum class %s {`, enumName)
	for _, enumItem := range model.Enum.Items {
		w.Line(`  @JsonProperty("%s") %s,`, enumItem.Value, enumItem.Name.UpperCase())
	}
	w.Line(`}`)
}

func generateOneOfModels(w *gen.Writer, model *spec.NamedModel) {
	interfaceName := model.Name.PascalCase()
	if model.OneOf.Discriminator != nil {
		w.Line(`@JsonTypeInfo(use = JsonTypeInfo.Id.NAME, include = JsonTypeInfo.As.PROPERTY, property = "%s")`, *model.OneOf.Discriminator)
	} else {
		w.Line(`@JsonTypeInfo(use = JsonTypeInfo.Id.NAME, include = JsonTypeInfo.As.WRAPPER_OBJECT)`)
	}
	w.Line(`sealed interface %s {`, interfaceName)
	for _, item := range model.OneOf.Items {
		itemType := KotlinType(&item.Type.Definition)
		w.Line(`  @JsonTypeName("%s")`, item.Name.Source)
		w.Line(`  data class %s(@field:JsonIgnore val data: %s): %s {`, item.Name.PascalCase(), itemType, interfaceName)
		w.Line(`    @get:JsonUnwrapped`)
		w.Line(`    private val _data: %s`, itemType)
		w.Line(`      get() = data`)
		w.Line(`  }`)
	}
	w.Line(`}`)
}