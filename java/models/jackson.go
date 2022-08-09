package models

import (
	"fmt"
	"strings"

	"generator"
	"java/imports"
	"java/packages"
	"java/types"
	"java/writer"
	"spec"
)

var Jackson = "jackson"

type JacksonGenerator struct {
	Types *types.Types
}

func NewJacksonGenerator(types *types.Types) *JacksonGenerator {
	return &JacksonGenerator{types}
}

func (g *JacksonGenerator) ModelsDefinitionsImports() []string {
	return []string{
		`com.fasterxml.jackson.databind.*`,
		`com.fasterxml.jackson.annotation.*`,
		`com.fasterxml.jackson.annotation.JsonSubTypes.*`,
		`com.fasterxml.jackson.core.type.TypeReference`,
	}
}

func (g *JacksonGenerator) ModelsUsageImports() []string {
	return []string{
		`com.fasterxml.jackson.databind.ObjectMapper`,
		`com.fasterxml.jackson.core.type.TypeReference`,
	}
}

func (g *JacksonGenerator) SetupImport(jsonPackage packages.Module) string {
	return fmt.Sprintf(`static %s.Json.setupObjectMapper`, jsonPackage.PackageName)
}

func (g *JacksonGenerator) CreateJsonMapperField(w *generator.Writer, annotation string) {
	if annotation != "" {
		w.Line(annotation)
	}
	w.Line(`private ObjectMapper objectMapper;`)
	w.EmptyLine()
	w.Line(`private String writeJson(Object result) {`)
	w.Line(`  try {`)
	w.Line(`    return objectMapper.writeValueAsString(result);`)
	w.Line(`  } catch (Exception exception) {`)
	w.Line(`    throw new RuntimeException(exception);`)
	w.Line(`  }`)
	w.Line(`}`)
}

func (g *JacksonGenerator) InitJsonMapper(w *generator.Writer) {
	w.Line(`this.objectMapper = new ObjectMapper();`)
	w.Line(`setupObjectMapper(objectMapper);`)
}

func (g *JacksonGenerator) ReadJson(varJson string, typ *spec.TypeDef) (string, string) {
	return fmt.Sprintf(`objectMapper.readValue(%s, new TypeReference<%s>() {})`, varJson, g.Types.Java(typ)), `IOException`
}

func (g *JacksonGenerator) WriteJson(varData string, typ *spec.TypeDef) (string, string) {
	return fmt.Sprintf(`objectMapper.writeValueAsString(%s)`, varData), `Exception`
}

func (g *JacksonGenerator) WriteJsonNoCheckedException(varData string, typ *spec.TypeDef) string {
	return fmt.Sprintf(`writeJson(%s)`, varData)
}

func (g *JacksonGenerator) SetupLibrary(thePackage packages.Module) []generator.CodeFile {
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

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return []generator.CodeFile{{
		Path:    thePackage.GetPath("Json.java"),
		Content: strings.TrimSpace(code),
	}}
}

func (g *JacksonGenerator) GenerateJsonParseException(thePackage, modelsPackage packages.Module) *generator.CodeFile {
	code := `
package [[.PackageName]];

import com.fasterxml.jackson.databind.exc.InvalidFormatException;
import java.util.*;
import [[.ModelsPackage]];

public class JsonParseException extends RuntimeException {
    private List<ValidationError> errors;
    public List<ValidationError> getErrors() {
        return errors;
    }

    public JsonParseException(Throwable exception) {
        super ("Failed to parse body: "+exception.getMessage(), exception);
        this.errors = extractErrors(exception);
    }

	private static List<ValidationError> extractErrors(Throwable exception) {
		if (exception instanceof InvalidFormatException) {
			var jsonPath = getJsonPath((InvalidFormatException)exception);
			var validation = new ValidationError(jsonPath, "parsing_failed", exception.getMessage());
			return List.of(validation);
		}
		return null;
	}

	private static String getJsonPath(InvalidFormatException exception) {
		var path = new StringBuilder("");
		for (int i = 0; i < exception.getPath().size(); i++) {
			var reference = exception.getPath().get(i);
			if (reference.getIndex() != -1) {
				path.append("[").append(reference.getIndex()).append("]");
			} else {
				if (i != 0) {
					path.append(".");
				}
				path.append(reference.getFieldName());
			}
		}
		return path.toString();
	}
}
`
	code, _ = generator.ExecuteTemplate(code, struct {
		PackageName   string
		ModelsPackage string
	}{
		thePackage.PackageName,
		modelsPackage.PackageStar,
	})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("JsonParseException.java"),
		Content: strings.TrimSpace(code),
	}
}

func (g *JacksonGenerator) VersionModels(version *spec.Version, thePackage packages.Module, jsonPackage packages.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
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

func jacksonJsonPropertyAnnotation(field *spec.NamedDefinition) string {
	required := "false"
	if !field.Type.Definition.IsNullable() {
		required = "true"
	}
	return fmt.Sprintf(`@JsonProperty(value = "%s", required = %s)`, field.Name.Source, required)
}

func (g *JacksonGenerator) modelObject(model *spec.NamedModel, thePackage packages.Module) *generator.CodeFile {
	w := writer.NewJavaWriter()
	w.Line(`package %s;`, thePackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(g.ModelsDefinitionsImports()...)
	imports.Add(g.Types.Imports()...)
	imports.Write(w)
	w.EmptyLine()
	className := model.Name.PascalCase()
	w.Line(`public class %s {`, className)
	for _, field := range model.Object.Fields {
		w.EmptyLine()
		w.Line(`  %s`, jacksonJsonPropertyAnnotation(&field))
		w.Line(`  private %s %s;`, g.Types.Java(&field.Type.Definition), field.Name.CamelCase())
	}
	w.EmptyLine()
	if len(model.Object.Fields) == 0 {
		w.Line(`  public %s() {`, model.Name.PascalCase())
	} else {
		w.Line(`  @JsonCreator`)
		w.Line(`  public %s(`, model.Name.PascalCase())
		for i, field := range model.Object.Fields {
			w.Line(`    %s`, jacksonJsonPropertyAnnotation(&field))
			ctorParam := fmt.Sprintf(`    %s %s`, g.Types.Java(&field.Type.Definition), field.Name.CamelCase())
			if i == len(model.Object.Fields)-1 {
				w.Line(`%s`, ctorParam)
			} else {
				w.Line(`%s,`, ctorParam)
			}
		}
		w.Line(`  ) {`)
	}
	for _, field := range model.Object.Fields {
		if !field.Type.Definition.IsNullable() && g.Types.IsReference(&field.Type.Definition) {
			w.Line(`    if (%s == null) { throw new IllegalArgumentException("null value is not allowed"); }`, field.Name.CamelCase())
		}
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

	return &generator.CodeFile{
		Path:    thePackage.GetPath(className + ".java"),
		Content: w.String(),
	}
}

func (g *JacksonGenerator) modelEnum(model *spec.NamedModel, thePackage packages.Module) *generator.CodeFile {
	w := writer.NewJavaWriter()
	w.Line(`package %s;`, thePackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(g.ModelsDefinitionsImports()...)
	imports.Add(g.Types.Imports()...)
	imports.Write(w)
	w.EmptyLine()
	enumName := model.Name.PascalCase()
	w.Line(`public enum %s {`, enumName)
	for _, enumItem := range model.Enum.Items {
		w.Line(`  @JsonProperty("%s") %s,`, enumItem.Value, enumItem.Name.UpperCase())
	}
	w.Line(`}`)

	return &generator.CodeFile{
		Path:    thePackage.GetPath(enumName + ".java"),
		Content: w.String(),
	}
}

func (g *JacksonGenerator) modelOneOf(model *spec.NamedModel, thePackage packages.Module) *generator.CodeFile {
	interfaceName := model.Name.PascalCase()
	w := writer.NewJavaWriter()
	w.Line("package %s;", thePackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(g.ModelsDefinitionsImports()...)
	imports.Add(g.Types.Imports()...)
	imports.Write(w)
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
		g.modelOneOfImplementation(w.Indented(), &item, model)
	}
	w.Line(`}`)

	return &generator.CodeFile{
		Path:    thePackage.GetPath(interfaceName + ".java"),
		Content: w.String(),
	}
}

func (g *JacksonGenerator) modelOneOfImplementation(w *generator.Writer, item *spec.NamedDefinition, model *spec.NamedModel) {
	w.Line(`class %s implements %s {`, oneOfItemClassName(item), model.Name.PascalCase())
	w.Line(`  @JsonUnwrapped`)
	w.Line(`  public %s data;`, g.Types.Java(&item.Type.Definition))
	w.EmptyLine()
	w.Line(`  public %s() {`, oneOfItemClassName(item))
	w.Line(`  }`)
	w.EmptyLine()
	w.Line(`  public %s(%s data) {`, oneOfItemClassName(item), g.Types.Java(&item.Type.Definition))
	if !item.Type.Definition.IsNullable() && g.Types.IsReference(&item.Type.Definition) {
		w.Line(`    if (data == null) { throw new IllegalArgumentException("null value is not allowed"); }`)
	}
	w.Line(`  	this.data = data;`)
	w.Line(`  }`)
	w.EmptyLine()
	w.Line(`  public %s getData() {`, g.Types.Java(&item.Type.Definition))
	w.Line(`    return data;`)
	w.Line(`  }`)
	w.EmptyLine()
	w.Line(`  public void setData(%s data) {`, g.Types.Java(&item.Type.Definition))
	if !item.Type.Definition.IsNullable() && g.Types.IsReference(&item.Type.Definition) {
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
