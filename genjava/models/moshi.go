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
	generatedSetupMoshiMethods []string
	Types                      *types.Types
}

func NewMoshiGenerator(types *types.Types) *MoshiGenerator {
	return &MoshiGenerator{[]string{}, types}
}

func (g *MoshiGenerator) JsonImports() []string {
	return []string{
		`com.squareup.moshi.Json`,
		`com.squareup.moshi.Moshi`,
		`com.squareup.moshi.Types`,
	}
}

func (g *MoshiGenerator) CreateJsonMapperField(w *sources.Writer) {
	w.Line(`private Moshi moshi;`)
}

func (g *MoshiGenerator) InitJsonMapper(w *sources.Writer) {
	w.Line(`Moshi.Builder moshiBuilder = new Moshi.Builder();`)
	w.Line(`setupMoshiAdapters(moshiBuilder);`)
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

func (g *MoshiGenerator) VersionModels(version *spec.Version, thePackage packages.Module, jsonPackage packages.Module) []sources.CodeFile {
	g.generatedSetupMoshiMethods = append(g.generatedSetupMoshiMethods, fmt.Sprintf(`%s.Json.setupMoshiOneOfAdapters`, thePackage.PackageName))

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

	adaptersPackage := jsonPackage.Subpackage("adapters")
	for range g.generatedSetupMoshiMethods {
		files = append(files, *g.setupOneOfAdapters(version, thePackage, adaptersPackage))
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

func (g *MoshiGenerator) SetupLibrary(thePackage packages.Module) []sources.CodeFile {
	adaptersPackage := thePackage.Subpackage("adapters")

	files := []sources.CodeFile{}
	files = append(files, *g.setupAdapters(thePackage, adaptersPackage))
	files = append(files, *bigDecimalAdapter(adaptersPackage))
	files = append(files, *localDateAdapter(adaptersPackage))
	files = append(files, *localDateTimeAdapter(adaptersPackage))
	files = append(files, *uuidAdapter(adaptersPackage))
	files = append(files, *unionAdapterFactory(adaptersPackage))
	files = append(files, *unwrapFieldAdapterFactory(adaptersPackage))
	return files
}

func (g *MoshiGenerator) setupAdapters(thePackage packages.Module, adaptersPackage packages.Module) *sources.CodeFile {
	w := writer.NewJavaWriter()
	w.Line("package %s;", thePackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(`com.squareup.moshi.Moshi`)
	imports.Add(adaptersPackage.PackageStar)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`public class Json {`)
	w.Line(`  public static void setupMoshiAdapters(Moshi.Builder moshiBuilder) {`)
	w.Line(`    moshiBuilder`)
	w.Line(`      .add(new BigDecimalAdapter())`)
	w.Line(`      .add(new UuidAdapter())`)
	w.Line(`      .add(new LocalDateAdapter())`)
	w.Line(`      .add(new LocalDateTimeAdapter());`)
	w.EmptyLine()
	for _, setupMoshiMethod := range g.generatedSetupMoshiMethods {
		w.Line(`    %s(moshiBuilder);`, setupMoshiMethod)
	}
	w.Line(`  }`)
	w.Line(`}`)

	return &sources.CodeFile{
		Path:    thePackage.GetPath("Json.java"),
		Content: w.String(),
	}
}

func (g *MoshiGenerator) setupOneOfAdapters(version *spec.Version, thePackage packages.Module, adaptersPackage packages.Module) *sources.CodeFile {
	w := writer.NewJavaWriter()
	w.Line("package %s;", thePackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(`com.squareup.moshi.Moshi`)
	imports.Add(adaptersPackage.PackageStar)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`public class Json {`)
	w.Line(`  public static void setupMoshiOneOfAdapters(Moshi.Builder moshiBuilder) {`)
	for _, model := range version.ResolvedModels {
		if model.IsOneOf() {
			w.IndentWith(2)
			w.Line(`moshiBuilder`)
			modelName := model.Name.PascalCase()
			for _, item := range model.OneOf.Items {
				w.Line(`  .add(new UnwrapFieldAdapterFactory(%s.%s.class))`, modelName, oneOfItemClassName(&item))
			}
			addUnionAdapterFactory := fmt.Sprintf(`  .add(UnionAdapterFactory.of(%s.class)`, modelName)
			if model.OneOf.Discriminator != nil {
				w.Line(`%s.withDiscriminator("%s")`, addUnionAdapterFactory, *model.OneOf.Discriminator)
			} else {
				w.Line(addUnionAdapterFactory)
			}
			for _, item := range model.OneOf.Items {
				w.Line(`    .withSubtype(%s.%s.class, "%s")`, modelName, oneOfItemClassName(&item), item.Name.Source)
			}
			w.Line(`  );`)
			w.UnindentWith(2)
		}
	}
	w.Line(`  }`)
	w.Line(`}`)

	return &sources.CodeFile{
		Path:    thePackage.GetPath("Json.java"),
		Content: w.String(),
	}
}

func bigDecimalAdapter(thePackage packages.Module) *sources.CodeFile {
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

func localDateAdapter(thePackage packages.Module) *sources.CodeFile {
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

func localDateTimeAdapter(thePackage packages.Module) *sources.CodeFile {
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

func uuidAdapter(thePackage packages.Module) *sources.CodeFile {
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

func unionAdapterFactory(thePackage packages.Module) *sources.CodeFile {
	code := `
package [[.PackageName]];

import com.squareup.moshi.*;
import java.io.IOException;
import java.lang.annotation.Annotation;
import java.lang.reflect.Type;
import java.util.*;
import javax.annotation.*;

public final class UnionAdapterFactory<T> implements JsonAdapter.Factory {
    final Class<T> baseType;
    final String discriminator;
    final List<String> tags;
    final List<Type> subtypes;
    @Nullable final JsonAdapter<Object> fallbackAdapter;

    UnionAdapterFactory(
            Class<T> baseType,
            String discriminator,
            List<String> tags,
            List<Type> subtypes,
            @Nullable JsonAdapter<Object> fallbackAdapter) {
        this.baseType = baseType;
        this.discriminator = discriminator;
        this.tags = tags;
        this.subtypes = subtypes;
        this.fallbackAdapter = fallbackAdapter;
    }

    /**
     * @param baseType The base type for which this factory will create adapters. Cannot be Object.
     *     JSON object.
     */
    @CheckReturnValue
    public static <T> UnionAdapterFactory<T> of(Class<T> baseType) {
        if (baseType == null) throw new NullPointerException("baseType == null");
        return new UnionAdapterFactory<>(baseType, null, Collections.<String>emptyList(), Collections.<Type>emptyList(), null);
    }

    /** Returns a new factory that decodes instances of {@code subtype}. */
    public UnionAdapterFactory<T> withDiscriminator(String discriminator) {
        if (discriminator == null) throw new NullPointerException("discriminator == null");
        return new UnionAdapterFactory<>(baseType, discriminator, tags, subtypes, fallbackAdapter);
    }

    /** Returns a new factory that decodes instances of {@code subtype}. */
    public UnionAdapterFactory<T> withSubtype(Class<? extends T> subtype, String tag) {
        if (subtype == null) throw new NullPointerException("subtype == null");
        if (tag == null) throw new NullPointerException("tag == null");
        if (tags.contains(tag)) {
            throw new IllegalArgumentException("Tags must be unique.");
        }
        List<String> newTags = new ArrayList<>(tags);
        newTags.add(tag);
        List<Type> newSubtypes = new ArrayList<>(subtypes);
        newSubtypes.add(subtype);
        return new UnionAdapterFactory<>(baseType, discriminator, newTags, newSubtypes, fallbackAdapter);
    }

    /**
     * Returns a new factory that with default to {@code fallbackJsonAdapter.fromJson(reader)} upon
     * decoding of unrecognized tags.
     *
     * <p>The {@link JsonReader} instance will not be automatically consumed, so make sure to consume
     * it within your implementation of {@link JsonAdapter#fromJson(JsonReader)}
     */
    public UnionAdapterFactory<T> withFallbackAdapter(@Nullable JsonAdapter<Object> fallbackJsonAdapter) {
        return new UnionAdapterFactory<>(baseType, discriminator, tags, subtypes, fallbackJsonAdapter);
    }

    /**
     * Returns a new factory that will default to {@code defaultValue} upon decoding of unrecognized
     * tags. The default value should be immutable.
     */
    public UnionAdapterFactory<T> withDefaultValue(@Nullable T defaultValue) {
        return withFallbackAdapter(buildFallbackAdapter(defaultValue));
    }

    private JsonAdapter<Object> buildFallbackAdapter(final T defaultValue) {
        return new JsonAdapter<Object>() {
            @Override
            public @Nullable Object fromJson(JsonReader reader) throws IOException {
                reader.skipValue();
                return defaultValue;
            }

            @Override
            public void toJson(JsonWriter writer, Object value) throws IOException {
                throw new IllegalArgumentException("Expected one of " + subtypes + " but found " + value + ", a " + value.getClass() + ". Register this subtype.");
            }
        };
    }

    @Override
    public JsonAdapter<?> create(Type type, Set<? extends Annotation> annotations, Moshi moshi) {
        if (Types.getRawType(type) != baseType || !annotations.isEmpty()) {
            return null;
        }

        List<JsonAdapter<Object>> jsonAdapters = new ArrayList<>(subtypes.size());
        for (int i = 0, size = subtypes.size(); i < size; i++) {
            jsonAdapters.add(moshi.adapter(subtypes.get(i)));
        }

        if (discriminator != null) {
            return new UnionDiscriminatorAdapter(discriminator, tags, subtypes, jsonAdapters, fallbackAdapter).nullSafe();
        } else {
            return new UnionWrapperAdapter(tags, subtypes, jsonAdapters, fallbackAdapter).nullSafe();
        }
    }

    static final class UnionDiscriminatorAdapter extends JsonAdapter<Object> {
        final String discriminator;
        final List<String> tags;
        final List<Type> subtypes;
        final List<JsonAdapter<Object>> adapters;
        @Nullable final JsonAdapter<Object> fallbackAdapter;
        final JsonReader.Options discriminatorOptions;
        final JsonReader.Options tagsOptions;

        UnionDiscriminatorAdapter(
                String discriminator,
                List<String> tags,
                List<Type> subtypes,
                List<JsonAdapter<Object>> adapters,
                @Nullable JsonAdapter<Object> fallbackAdapter) {
            this.discriminator = discriminator;
            this.tags = tags;
            this.subtypes = subtypes;
            this.adapters = adapters;
            this.fallbackAdapter = fallbackAdapter;

            this.discriminatorOptions = JsonReader.Options.of(discriminator);
            this.tagsOptions = JsonReader.Options.of(tags.toArray(new String[0]));
        }

        @Override
        public Object fromJson(JsonReader reader) throws IOException {
            final int tagIndex = getTagIndex(reader);
            JsonAdapter<Object> adapter = fallbackAdapter;
            if (tagIndex != -1) {
                adapter = adapters.get(tagIndex);
            }
            return adapter.fromJson(reader);
        }

        @Override
        public void toJson(JsonWriter writer, Object value) throws IOException {
            int tagIndex = getTagIndex(value);
            if (tagIndex == -1) {
                fallbackAdapter.toJson(writer, value);
            } else {
                final JsonAdapter<Object> adapter = adapters.get(tagIndex);
                writer.beginObject();
                writer.name(discriminator).value(tags.get(tagIndex));
                int flattenToken = writer.beginFlatten();
                adapter.toJson(writer, value);
                writer.endFlatten(flattenToken);
                writer.endObject();
            }
        }

        private int getTagIndex(JsonReader reader) throws IOException {
            JsonReader peeked = reader.peekJson();
            peeked.setFailOnUnknown(false);
            try {
                peeked.beginObject();
                while (peeked.hasNext()) {
                    if (peeked.selectName(discriminatorOptions) == -1) {
                        peeked.skipName();
                        peeked.skipValue();
                        continue;
                    }

                    int tagIndex = peeked.selectString(tagsOptions);
                    if (tagIndex == -1 && this.fallbackAdapter == null) {
                        throw new JsonDataException("Expected one of " + tags + " for key '" + discriminator + "' but found '" + peeked.nextString() + "'. Register a subtype for this tag.");
                    }
                    return tagIndex;
                }

                throw new JsonDataException("Missing discriminator field " + discriminator);
            } finally {
                peeked.close();
            }
        }

        private int getTagIndex(Object value) {
            Class<?> type = value.getClass();
            int tagIndex = subtypes.indexOf(type);
            if (tagIndex == -1 && fallbackAdapter == null) {
                throw new IllegalArgumentException("Expected one of " + subtypes + " but found " + value + ", a " + value.getClass() + ". Register this subtype.");
            } else {
                return tagIndex;
            }
        }

        @Override
        public String toString() {
            return "UnionDiscriminatorAdapter(" + discriminator + ")";
        }
    }

    static final class UnionWrapperAdapter extends JsonAdapter<Object> {
        final List<String> tags;
        final List<Type> subtypes;
        final List<JsonAdapter<Object>> adapters;
        @Nullable final JsonAdapter<Object> fallbackAdapter;
        final JsonReader.Options tagsOptions;

        UnionWrapperAdapter(
                List<String> tags,
                List<Type> subtypes,
                List<JsonAdapter<Object>> adapters,
                @Nullable JsonAdapter<Object> fallbackAdapter) {
            this.tags = tags;
            this.subtypes = subtypes;
            this.adapters = adapters;
            this.fallbackAdapter = fallbackAdapter;

            this.tagsOptions = JsonReader.Options.of(tags.toArray(new String[0]));
        }

        @Override
        public Object fromJson(JsonReader reader) throws IOException {
            int tagIndex = getTagIndex(reader);
            if (tagIndex == -1) {
                return this.fallbackAdapter.fromJson(reader);
            } else {
                reader.beginObject();
                reader.skipName();
                final Object value = adapters.get(tagIndex).fromJson(reader);
                reader.endObject();
                return value;
            }
        }

        @Override
        public void toJson(JsonWriter writer, Object value) throws IOException {
            int tagIndex = getTagIndex(value);
            if (tagIndex == -1) {
                fallbackAdapter.toJson(writer, value);
            } else {
                final JsonAdapter<Object> adapter = adapters.get(tagIndex);
                writer.beginObject();
                writer.name(tags.get(tagIndex));
                adapter.toJson(writer, value);
                writer.endObject();
            }
        }

        private int getTagIndex(JsonReader reader) throws IOException {
            JsonReader peeked = reader.peekJson();
            peeked.setFailOnUnknown(false);
            try {
                peeked.beginObject();
                int tagIndex = peeked.selectName(tagsOptions);
                if (tagIndex == -1 && this.fallbackAdapter == null) {
                    throw new JsonDataException("Expected one of keys:" + tags + "' but found '" + peeked.nextString() + "'. Register a subtype for this tag.");
                }
                return tagIndex;
            } finally {
                peeked.close();
            }
        }

        private int getTagIndex(Object value) {
            Class<?> type = value.getClass();
            int tagIndex = subtypes.indexOf(type);
            if (tagIndex == -1 && fallbackAdapter == null) {
                throw new IllegalArgumentException("Expected one of " + subtypes + " but found " + value + ", a " + value.getClass() + ". Register this subtype.");
            }
            return tagIndex;
        }

        @Override
        public String toString() {
            return "UnionWrapperAdapter";
        }
    }
}
`

	code, _ = sources.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &sources.CodeFile{
		Path:    thePackage.GetPath("UnionAdapterFactory.java"),
		Content: strings.TrimSpace(code),
	}
}

func unwrapFieldAdapterFactory(thePackage packages.Module) *sources.CodeFile {
	code := `
package [[.PackageName]];

import com.squareup.moshi.*;
import org.jetbrains.annotations.Nullable;
import java.io.IOException;
import java.lang.annotation.Annotation;
import java.lang.reflect.*;
import java.util.Set;

public final class UnwrapFieldAdapterFactory<T> implements JsonAdapter.Factory {
    final Class<T> type;

    public UnwrapFieldAdapterFactory(Class<T> type) {
        this.type = type;
    }

    @Nullable
    @Override
    public JsonAdapter<?> create(Type type, Set<? extends Annotation> annotations, Moshi moshi) {
        if (Types.getRawType(type) != this.type || !annotations.isEmpty()) {
            return null;
        }

        final Field[] fields = this.type.getDeclaredFields();
        if (fields.length != 1) {
            throw new RuntimeException("Type "+type.getTypeName()+" has "+fields.length+" fields, unwrap adapter can be used only with single-field types");
        }
        var field = fields[0];

        Constructor<T> constructor;
        try {
            constructor = this.type.getDeclaredConstructor(field.getType());
        } catch (NoSuchMethodException e) {
            throw new RuntimeException("Type "+type.getTypeName()+" does not have constructor with single parameter of type "+ field.getType().getName()+", it's required for unwrap adapter");
        }

        JsonAdapter<?> fieldAdapter = moshi.adapter(field.getType());
        return new UnwrapFieldAdapter(constructor, field, fieldAdapter);
    }

    public class UnwrapFieldAdapter<O, I> extends JsonAdapter<Object> {
        private Constructor<O> constructor;
        private Field field;
        private JsonAdapter<I> fieldAdapter;

        public UnwrapFieldAdapter(Constructor<O> constructor, Field field, JsonAdapter<I> fieldAdapter) {
            this.constructor = constructor;
            this.field = field;
            this.fieldAdapter = fieldAdapter;
        }

        @Override
        public Object fromJson(JsonReader reader) throws IOException {
            I fieldValue = fieldAdapter.fromJson(reader);
            try {
                return constructor.newInstance(fieldValue);
            } catch (Throwable e) {
                throw new IOException("Failed to create object with constructor "+constructor.getName(), e);
            }
        }

        @Override
        public void toJson(JsonWriter writer, Object value) throws IOException {
            I fieldValue;
            try {
                fieldValue = (I)field.get(value);
            } catch (IllegalAccessException e) {
                throw new IOException("Failed to get value of field "+field.getName(), e);
            }
            fieldAdapter.toJson(writer, fieldValue);
        }

        @Override
        public String toString() {
            return "UnwrapFieldAdapter(" + field.getName() + ")";
        }
    }
}
`

	code, _ = sources.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &sources.CodeFile{
		Path:    thePackage.GetPath("UnwrapFieldAdapterFactory.java"),
		Content: strings.TrimSpace(code),
	}
}
