package models

import (
	"fmt"
	"generator"
	"java/imports"
	"java/packages"
	"java/types"
	"java/writer"
	"spec"
)

var Jackson = "jackson"
var jacksonCustomObjectMapper = `CustomObjectMapper`

type JacksonGenerator struct {
	Types    *types.Types
	Packages *Packages
}

func NewJacksonGenerator(types *types.Types, packages *Packages) *JacksonGenerator {
	return &JacksonGenerator{types, packages}
}

func (g *JacksonGenerator) Models(version *spec.Version) []generator.CodeFile {
	return g.models(version.ResolvedModels, g.Packages.Models(version))
}

func (g *JacksonGenerator) ErrorModels(httperrors *spec.HttpErrors) []generator.CodeFile {
	return g.models(httperrors.ResolvedModels, g.Packages.ErrorsModels)
}

func (g *JacksonGenerator) models(models []*spec.NamedModel, thePackage packages.Package) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, model := range models {
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

func (g *JacksonGenerator) modelObject(model *spec.NamedModel, thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, model.Name.PascalCase())
	imports := imports.New()
	imports.Add(g.modelsDefinitionsImports()...)
	imports.Add(g.Types.Imports()...)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`public class [[.ClassName]] {`)
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
	return w.ToCodeFile()
}

func (g *JacksonGenerator) modelEnum(model *spec.NamedModel, thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, model.Name.PascalCase())
	imports := imports.New()
	imports.Add(g.modelsDefinitionsImports()...)
	imports.Add(g.Types.Imports()...)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`public enum [[.ClassName]] {`)
	for _, enumItem := range model.Enum.Items {
		w.Line(`  @JsonProperty("%s") %s,`, enumItem.Value, enumItem.Name.UpperCase())
	}
	w.Line(`}`)
	return w.ToCodeFile()
}

func (g *JacksonGenerator) modelOneOf(model *spec.NamedModel, thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, model.Name.PascalCase())
	imports := imports.New()
	imports.Add(g.modelsDefinitionsImports()...)
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
		w.Line(`  @Type(value = [[.ClassName]].%s.class, name = "%s"),`, oneOfItemClassName(&item), item.Name.Source)
	}
	w.Line(`})`)
	w.Line(`public interface [[.ClassName]] {`)
	for index, item := range model.OneOf.Items {
		if index > 0 {
			w.EmptyLine()
		}
		g.modelOneOfImplementation(w.Indented(), &item, model)
	}
	w.Line(`}`)
	return w.ToCodeFile()
}

func (g *JacksonGenerator) modelOneOfImplementation(w generator.Writer, item *spec.NamedDefinition, model *spec.NamedModel) {
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

func (g *JacksonGenerator) JsonRead(varJson string, typ *spec.TypeDef) string {
	return fmt.Sprintf(`read(%s, new TypeReference<%s>() {})`, varJson, g.Types.Java(typ))
}

func (g *JacksonGenerator) JsonWrite(varData string, typ *spec.TypeDef) string {
	return fmt.Sprintf(`write(%s)`, varData)
}

func (g *JacksonGenerator) modelsDefinitionsImports() []string {
	return []string{
		`com.fasterxml.jackson.databind.*`,
		`com.fasterxml.jackson.annotation.*`,
		`com.fasterxml.jackson.annotation.JsonSubTypes.*`,
		`com.fasterxml.jackson.core.type.TypeReference`,
	}
}

func (g *JacksonGenerator) ModelsUsageImports() []string {
	return []string{
		`com.fasterxml.jackson.core.type.TypeReference`,
		`com.fasterxml.jackson.databind.ObjectMapper`,
	}
}

func (g *JacksonGenerator) ValidationErrorsHelpers() *generator.CodeFile {
	w := writer.New(g.Packages.Errors, `ValidationErrorsHelpers`)
	w.Template(
		map[string]string{
			`JsonPackage`:         g.Packages.Json.PackageName,
			`ErrorsModelsPackage`: g.Packages.ErrorsModels.PackageName,
		}, `
import com.fasterxml.jackson.databind.exc.InvalidFormatException;
import [[.JsonPackage]].*;
import [[.ErrorsModelsPackage]].*;

import java.util.List;

public class [[.ClassName]] {
	public static List<ValidationError> extractValidationErrors(JsonParseException exception) {
		var causeException = exception.getCause();
		if (causeException instanceof InvalidFormatException) {
			var jsonPath = getJsonPath((InvalidFormatException) causeException);
			var validation = new ValidationError(jsonPath, "parsing_failed", exception.getMessage());
			return List.of(validation);
		}
		return null;
	}

	private static String getJsonPath(InvalidFormatException exception) {
		var path = new StringBuilder();
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
`)
	return w.ToCodeFile()
}

func (g *JacksonGenerator) JsonHelpers() []generator.CodeFile {
	files := []generator.CodeFile{}

	files = append(files, *g.json())
	files = append(files, *g.jsonParseException())
	files = append(files, g.setupLibrary()...)

	return files
}

func (g *JacksonGenerator) json() *generator.CodeFile {
	w := writer.New(g.Packages.Json, `Json`)
	w.Lines(`
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;

import java.io.IOException;

public class Json {
	private final ObjectMapper objectMapper;

	public Json(ObjectMapper objectMapper) {
		this.objectMapper = objectMapper;
	}

	public String write(Object data) {
		try {
			return objectMapper.writeValueAsString(data);
		} catch (Exception exception) {
			throw new RuntimeException(exception);
		}
	}

	public <T> T read(String jsonStr, TypeReference<T> typeReference) {
		try {
			return objectMapper.readValue(jsonStr, typeReference);
		} catch (IOException exception) {
			throw new JsonParseException(exception);
		}
	}
}
`)
	return w.ToCodeFile()
}

func (g *JacksonGenerator) jsonParseException() *generator.CodeFile {
	w := writer.New(g.Packages.Json, `JsonParseException`)
	w.Lines(`
public class [[.ClassName]] extends RuntimeException {
	public JsonParseException(Throwable exception) {
		super("Failed to parse body: " + exception.getMessage(), exception);
	}
}
`)
	return w.ToCodeFile()
}

func (g *JacksonGenerator) setupLibrary() []generator.CodeFile {
	w := writer.New(g.Packages.Json, jacksonCustomObjectMapper)
	w.Lines(`
import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.databind.*;
import com.fasterxml.jackson.datatype.jsr310.*;

public class [[.ClassName]] {
	public static void setup(ObjectMapper objectMapper) {
		objectMapper
			.registerModule(new JavaTimeModule())
			.configure(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS, false)
			.setSerializationInclusion(JsonInclude.Include.NON_NULL);
	}
}
`)
	return []generator.CodeFile{*w.ToCodeFile()}
}

func (g *JacksonGenerator) CreateJsonHelper(name string) string {
	return fmt.Sprintf(`
ObjectMapper objectMapper = new ObjectMapper();
%s.setup(objectMapper);
%s = new Json(objectMapper);
`, jacksonCustomObjectMapper, name)
}
