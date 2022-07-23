package models

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen/kotlin/imports"
	"github.com/specgen-io/specgen/v2/gen/kotlin/modules"
	"github.com/specgen-io/specgen/v2/gen/kotlin/types"
	"github.com/specgen-io/specgen/v2/gen/kotlin/writer"
	"github.com/specgen-io/specgen/v2/generator"
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

func (g *MoshiGenerator) ModelsDefinitionsImports() []string {
	return []string{
		`com.squareup.moshi.*`,
	}
}

func (g *MoshiGenerator) ModelsUsageImports() []string {
	return []string{
		`com.squareup.moshi.*`,
	}
}

func (g *MoshiGenerator) SetupImport(jsonPackage modules.Module) string {
	return fmt.Sprintf(`%s.setupMoshiAdapters`, jsonPackage.PackageName)
}

func (g *MoshiGenerator) CreateJsonMapperField() string {
	return `private val moshi: Moshi`
}

func (g *MoshiGenerator) InitJsonMapper(w *generator.Writer) {
	w.Line(`val moshiBuilder = Moshi.Builder()`)
	w.Line(`setupMoshiAdapters(moshiBuilder)`)
	w.Line(`moshi = moshiBuilder.build()`)
}

func (g *MoshiGenerator) ReadJson(varJson string, typ *spec.TypeDef) (string, string) {
	adapter := fmt.Sprintf(`adapter(%s::class.java)`, g.Types.Kotlin(typ))
	if typ.Node == spec.MapType {
		typeJava := g.Types.Kotlin(typ.Child)
		adapter = fmt.Sprintf(`adapter<Map<String, %s>>(Types.newParameterizedType(MutableMap::class.java, String::class.java, %s::class.java))`, typeJava, typeJava)
	}
	if typ.Node == spec.ArrayType {
		typeJava := g.Types.Kotlin(typ.Child)
		adapter = fmt.Sprintf(`adapter<List<%s>>(Types.newParameterizedType(List::class.java, %s::class.java))`, typeJava, typeJava)
	}

	return fmt.Sprintf(`moshi.%s.fromJson(%s)!!`, adapter, varJson), `JsonDataException`

}

func (g *MoshiGenerator) WriteJson(varData string, typ *spec.TypeDef) (string, string) {
	adapterParam := fmt.Sprintf(`adapter(%s::class.java)`, g.Types.Kotlin(typ))
	if typ.Node == spec.MapType {
		typeJava := g.Types.Kotlin(typ.Child)
		adapterParam = fmt.Sprintf(`adapter<Map<String, %s>>(Types.newParameterizedType(MutableMap::class.java, String::class.java, %s::class.java))`, typeJava, typeJava)
	}
	if typ.Node == spec.ArrayType {
		typeJava := g.Types.Kotlin(typ.Child)
		adapterParam = fmt.Sprintf(`adapter<List<%s>>(Types.newParameterizedType(List::class.java, %s::class.java))`, typeJava, typeJava)
	}

	return fmt.Sprintf(`moshi.%s.toJson(%s)`, adapterParam, varData), `IOException`

}

func (g *MoshiGenerator) GenerateJsonParseException(thePackage, modelsPackage modules.Module) *generator.CodeFile {
	code := `
package [[.PackageName]]

import com.squareup.moshi.JsonDataException
import java.util.regex.Pattern
import [[.ModelsPackage]]

class JsonParseException(exception: Throwable) :
    RuntimeException("Failed to parse body: " + exception.message, exception) {
    val errors: List<ValidationError>?

    init {
        errors = extractErrors(exception)
    }

    companion object {
        private val pathPattern = Pattern.compile("\\$\\.([^ ]+)")
        private fun extractErrors(exception: Throwable): List<ValidationError>? {
            if (exception is JsonDataException) {
                val matcher = pathPattern.matcher(exception.message)
                if (matcher.find()) {
                    val jsonPath = matcher.group(1)
                    return listOf(ValidationError(jsonPath, "parsing_failed", exception.message))
                }
            }
            return null
        }
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
		Path:    thePackage.GetPath("JsonParseException.kt"),
		Content: strings.TrimSpace(code),
	}
}

func (g *MoshiGenerator) VersionModels(version *spec.Version, thePackage modules.Module, jsonPackage modules.Module) []generator.CodeFile {
	w := writer.NewKotlinWriter()
	w.Line(`package %s`, thePackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(g.ModelsDefinitionsImports()...)
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

	files := []generator.CodeFile{}
	files = append(files, generator.CodeFile{Path: thePackage.GetPath("models.kt"), Content: w.String()})

	g.generatedSetupMoshiMethods = append(g.generatedSetupMoshiMethods, fmt.Sprintf(`%s.setupModelsMoshiAdapters`, thePackage.PackageName))
	adaptersPackage := jsonPackage.Subpackage("adapters")
	for range g.generatedSetupMoshiMethods {
		files = append(files, *g.setupOneOfAdapters(version, thePackage, adaptersPackage))
	}

	return files
}

func (g *MoshiGenerator) modelObject(w *generator.Writer, model *spec.NamedModel) {
	className := model.Name.PascalCase()
	w.Line(`data class %s(`, className)
	for _, field := range model.Object.Fields {
		w.Line(`  @Json(name = "%s")`, field.Name.Source)
		w.Line(`  val %s: %s,`, field.Name.CamelCase(), g.Types.Kotlin(&field.Type.Definition))
	}
	w.Line(`)`)
}

func (g *MoshiGenerator) modelEnum(w *generator.Writer, model *spec.NamedModel) {
	enumName := model.Name.PascalCase()
	w.Line(`enum class %s {`, enumName)
	for _, enumItem := range model.Enum.Items {
		w.Line(`  @Json(name = "%s")`, enumItem.Value)
		w.Line(`  %s,`, enumItem.Name.UpperCase())
	}
	w.Line(`}`)
}

func (g *MoshiGenerator) modelOneOf(w *generator.Writer, model *spec.NamedModel) {
	sealedClassName := model.Name.PascalCase()
	w.Line(`sealed class %s {`, sealedClassName)
	for _, item := range model.OneOf.Items {
		w.Line(`  data class %s(val data: %s): %s()`, oneOfItemClassName(&item), g.Types.Kotlin(&item.Type.Definition), sealedClassName)
	}
	w.Line(`}`)
}

func (g *MoshiGenerator) SetupLibrary(thePackage modules.Module) []generator.CodeFile {
	adaptersPackage := thePackage.Subpackage("adapters")

	files := []generator.CodeFile{}
	files = append(files, *g.setupAdapters(thePackage, adaptersPackage))
	files = append(files, *bigDecimalAdapter(adaptersPackage))
	files = append(files, *localDateAdapter(adaptersPackage))
	files = append(files, *localDateTimeAdapter(adaptersPackage))
	files = append(files, *uuidAdapter(adaptersPackage))
	files = append(files, *unionAdapterFactory(adaptersPackage))
	files = append(files, *unwrapFieldAdapterFactory(adaptersPackage))
	return files
}

func (g *MoshiGenerator) setupAdapters(thePackage modules.Module, adaptersPackage modules.Module) *generator.CodeFile {
	w := writer.NewKotlinWriter()
	w.Line("package %s", thePackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(`com.squareup.moshi.Moshi`)
	imports.Add(`com.squareup.moshi.kotlin.reflect.KotlinJsonAdapterFactory`)
	imports.Add(adaptersPackage.PackageStar)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`fun setupMoshiAdapters(moshiBuilder: Moshi.Builder) {`)
	w.Line(`  moshiBuilder`)
	w.Line(`    .add(BigDecimalAdapter())`)
	w.Line(`    .add(UuidAdapter())`)
	w.Line(`    .add(LocalDateAdapter())`)
	w.Line(`    .add(LocalDateTimeAdapter())`)
	w.EmptyLine()
	for _, setupMoshiMethod := range g.generatedSetupMoshiMethods {
		w.Line(`    %s(moshiBuilder);`, setupMoshiMethod)
	}
	w.EmptyLine()
	w.Line(`  moshiBuilder`)
	w.Line(`    .add(KotlinJsonAdapterFactory())`)
	w.Line(`}`)

	return &generator.CodeFile{
		Path:    thePackage.GetPath("Json.kt"),
		Content: w.String(),
	}
}

func (g *MoshiGenerator) setupOneOfAdapters(version *spec.Version, thePackage modules.Module, adaptersPackage modules.Module) *generator.CodeFile {
	w := writer.NewKotlinWriter()
	w.Line(`package %s`, thePackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(`com.squareup.moshi.Moshi`)
	imports.Add(adaptersPackage.PackageStar)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`fun setupModelsMoshiAdapters(moshiBuilder: Moshi.Builder) {`)
	for _, model := range version.ResolvedModels {
		if model.IsOneOf() {
			w.Indent()
			w.Line(`moshiBuilder`)
			modelName := model.Name.PascalCase()
			for _, item := range model.OneOf.Items {
				w.Line(`  .add(UnwrapFieldAdapterFactory(%s.%s::class.java))`, modelName, oneOfItemClassName(&item))
			}
			addUnionAdapterFactory := fmt.Sprintf(`  .add(UnionAdapterFactory.of(%s::class.java)`, modelName)
			if model.OneOf.Discriminator != nil {
				w.Line(`%s.withDiscriminator("%s")`, addUnionAdapterFactory, *model.OneOf.Discriminator)
			} else {
				w.Line(addUnionAdapterFactory)
			}
			for _, item := range model.OneOf.Items {
				w.Line(`    .withSubtype(%s.%s::class.java, "%s")`, modelName, oneOfItemClassName(&item), item.Name.Source)
			}
			w.Line(`)`)
			w.Unindent()
		}
	}
	w.Line(`}`)

	return &generator.CodeFile{
		Path:    thePackage.GetPath("Json.kt"),
		Content: w.String(),
	}
}

func bigDecimalAdapter(thePackage modules.Module) *generator.CodeFile {
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

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("BigDecimalAdapter.kt"),
		Content: strings.TrimSpace(code),
	}
}

func localDateAdapter(thePackage modules.Module) *generator.CodeFile {
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

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("LocalDateAdapter.kt"),
		Content: strings.TrimSpace(code),
	}
}

func localDateTimeAdapter(thePackage modules.Module) *generator.CodeFile {
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

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("LocalDateTimeAdapter.kt"),
		Content: strings.TrimSpace(code),
	}
}

func uuidAdapter(thePackage modules.Module) *generator.CodeFile {
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

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("UuidAdapter.kt"),
		Content: strings.TrimSpace(code),
	}
}

func unionAdapterFactory(thePackage modules.Module) *generator.CodeFile {
	code := `
package [[.PackageName]]

import com.squareup.moshi.*
import java.lang.reflect.Type

class UnionAdapterFactory<T> internal constructor(
    private val baseType: Class<T>,
    private val discriminator: String?,
    private val tags: List<String>,
    private val subtypes: List<Type>,
    private val fallbackAdapter: JsonAdapter<Any>?
) : JsonAdapter.Factory {

    companion object {
        fun <T> of(baseType: Class<T>): UnionAdapterFactory<T> {
            return UnionAdapterFactory(baseType, null, emptyList(), emptyList(), null)
        }
    }

    fun withDiscriminator(discriminator: String?): UnionAdapterFactory<T> {
        if (discriminator == null) throw NullPointerException("discriminator == null")
        return UnionAdapterFactory(baseType, discriminator, tags, subtypes, fallbackAdapter)
    }

    fun withSubtype(subtype: Class<out T>, tag: String): UnionAdapterFactory<T> {
        require(!tags.contains(tag)) { "Tags must be unique." }
        val newTags: MutableList<String> = ArrayList(tags).also {
            it.add(tag)
        }
        val newSubtypes: MutableList<Type> = ArrayList(subtypes).also {
            it.add(subtype)
        }
        return UnionAdapterFactory(baseType, discriminator, newTags, newSubtypes, fallbackAdapter)
    }

    private fun withFallbackAdapter(fallbackJsonAdapter: JsonAdapter<Any>?): UnionAdapterFactory<T> {
        return UnionAdapterFactory(baseType, discriminator, tags, subtypes, fallbackJsonAdapter)
    }

    fun withDefaultValue(defaultValue: T): UnionAdapterFactory<T> {
        return withFallbackAdapter(buildFallbackAdapter(defaultValue))
    }

    private fun buildFallbackAdapter(defaultValue: T): JsonAdapter<Any> {
        return object : JsonAdapter<Any>() {
            override fun fromJson(reader: JsonReader): Any? {
                reader.skipValue()
                return defaultValue
            }

            override fun toJson(writer: JsonWriter, value: Any?) {
                throw IllegalArgumentException("Expected one of " + subtypes + " but found " + value + ", a " + value!!.javaClass + ". Register this subtype.")
            }
        }
    }

    override fun create(type: Type, annotations: Set<Annotation>, moshi: Moshi): JsonAdapter<*>? {
        if (Types.getRawType(type) != baseType || annotations.isNotEmpty()) {
            return null
        }

        val jsonAdapters: MutableList<JsonAdapter<Any>> = java.util.ArrayList(subtypes.size)

        for (element in subtypes) {
            jsonAdapters.add(moshi.adapter(element))
        }

        return if (discriminator != null) {
            UnionDiscriminatorAdapter(discriminator, tags, subtypes, jsonAdapters, fallbackAdapter).nullSafe()
        } else {
            UnionWrapperAdapter(tags, subtypes, jsonAdapters, fallbackAdapter).nullSafe()
        }
    }

    internal class UnionDiscriminatorAdapter(
        private val discriminator: String,
        private val tags: List<String>,
        private val subtypes: List<Type>,
        private val adapters: List<JsonAdapter<Any>>,
        private val fallbackAdapter: JsonAdapter<Any>?
    ) : JsonAdapter<Any>() {
        private val discriminatorOptions: JsonReader.Options = JsonReader.Options.of(discriminator)
        private val tagsOptions: JsonReader.Options = JsonReader.Options.of(*tags.toTypedArray())

        override fun fromJson(reader: JsonReader): Any? {
            val tagIndex = getTagIndex(reader)
            var adapter = fallbackAdapter
            if (tagIndex != -1) {
                adapter = adapters[tagIndex]
            }
            return adapter!!.fromJson(reader)
        }

        override fun toJson(writer: JsonWriter, value: Any?) {
            val tagIndex = getTagIndex(value!!)
            if (tagIndex == -1) {
                fallbackAdapter!!.toJson(writer, value)
            } else {
                val adapter = adapters[tagIndex]
                writer.beginObject()
                writer.name(discriminator).value(tags[tagIndex])
                val flattenToken = writer.beginFlatten()
                adapter.toJson(writer, value)
                writer.endFlatten(flattenToken)
                writer.endObject()
            }
        }

        private fun getTagIndex(reader: JsonReader): Int {
            val peeked = reader.peekJson()
            peeked.setFailOnUnknown(false)
            peeked.use {
                it.beginObject()
                while (it.hasNext()) {
                    if (it.selectName(discriminatorOptions) == -1) {
                        it.skipName()
                        it.skipValue()
                        continue
                    }
                    val tagIndex = it.selectString(tagsOptions)
                    if (tagIndex == -1 && fallbackAdapter == null) {
                        throw JsonDataException("Expected one of " + tags + " for key '" + discriminator + "' but found '" + it.nextString() + "'. Register a subtype for this tag.")
                    }
                    return tagIndex
                }
                throw JsonDataException("Missing discriminator field $discriminator")
            }
        }

        private fun getTagIndex(value: Any): Int {
            val type: Class<*> = value.javaClass
            val tagIndex = subtypes.indexOf(type)
            return if (tagIndex == -1 && fallbackAdapter == null) {
                throw IllegalArgumentException("Expected one of " + subtypes + " but found " + value + ", a " + value.javaClass + ". Register this subtype.")
            } else {
                tagIndex
            }
        }

        override fun toString(): String {
            return "UnionDiscriminatorAdapter($discriminator)"
        }
    }

    internal class UnionWrapperAdapter(
        private val tags: List<String>,
        private val subtypes: List<Type>,
        private val adapters: List<JsonAdapter<Any>>,
        private val fallbackAdapter: JsonAdapter<Any>?
    ) : JsonAdapter<Any>() {
        private val tagsOptions: JsonReader.Options = JsonReader.Options.of(*tags.toTypedArray())

        override fun fromJson(reader: JsonReader): Any? {
            val tagIndex: Int = getTagIndex(reader)
            return if (tagIndex == -1) {
                fallbackAdapter!!.fromJson(reader)
            } else {
                reader.beginObject()
                reader.skipName()
                val value = adapters[tagIndex].fromJson(reader)
                reader.endObject()
                value
            }
        }

        override fun toJson(writer: JsonWriter, value: Any?) {
            val tagIndex: Int = getTagIndex(value!!)
            if (tagIndex == -1) {
                fallbackAdapter!!.toJson(writer, value)
            } else {
                val adapter = adapters[tagIndex]
                writer.beginObject()
                writer.name(tags[tagIndex])
                adapter.toJson(writer, value)
                writer.endObject()
            }
        }

        private fun getTagIndex(reader: JsonReader): Int {
            val peeked = reader.peekJson()
            peeked.setFailOnUnknown(false)
            return peeked.use {
                it.beginObject()
                val tagIndex = it.selectName(tagsOptions)
                if (tagIndex == -1 && fallbackAdapter == null) {
                    throw JsonDataException("Expected one of keys:" + tags + "' but found '" + it.nextString() + "'. Register a subtype for this tag.")
                }
                tagIndex
            }
        }

        private fun getTagIndex(value: Any): Int {
            val type: Class<*> = value.javaClass
            val tagIndex = subtypes.indexOf(type)
            require(!(tagIndex == -1 && fallbackAdapter == null)) { "Expected one of " + subtypes + " but found " + value + ", a " + value.javaClass + ". Register this subtype." }
            return tagIndex
        }

        override fun toString(): String {
            return "UnionWrapperAdapter"
        }
    }
}
`

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("UnionAdapterFactory.kt"),
		Content: strings.TrimSpace(code),
	}
}

func unwrapFieldAdapterFactory(thePackage modules.Module) *generator.CodeFile {
	code := `
package [[.PackageName]]

import com.squareup.moshi.*
import java.io.IOException
import java.lang.reflect.*

class UnwrapFieldAdapterFactory<T>(private val type: Class<T>) : JsonAdapter.Factory {
    override fun create(type: Type, annotations: Set<Annotation?>, moshi: Moshi): JsonAdapter<*>? {
        if (Types.getRawType(type) != this.type || annotations.isNotEmpty()) {
            return null
        }

        val fields = this.type.declaredFields
        if (fields.size != 1) {
            throw RuntimeException("Type " + type.typeName + " has " + fields.size + " fields, unwrap adapter can be used only with single-field types")
        }
        val field = fields[0]
        val getterName = "get" + field.name.replaceFirstChar { it.uppercase() }

        val getter = try {
            this.type.getDeclaredMethod(getterName)
        } catch (e: NoSuchMethodException) {
            throw RuntimeException("Type " + type.typeName + " field " + field.name + " does not have getter method " + field.type.name + ", it's required for unwrap adapter", e)
        }

        val constructor: Constructor<T> = try {
            this.type.getDeclaredConstructor(field.type)
        } catch (e: NoSuchMethodException) {
            throw RuntimeException("Type " + type.typeName + " does not have constructor with single parameter of type " + field.type.name + ", it's required for unwrap adapter")
        }

        val fieldAdapter: JsonAdapter<*>  = moshi.adapter(field.type)
        return UnwrapFieldAdapter(constructor, getter, fieldAdapter)
    }

    inner class UnwrapFieldAdapter<O, I>(
        private val constructor: Constructor<O>,
        private val getter: Method,
        private val fieldAdapter: JsonAdapter<I>
    ) : JsonAdapter<Any>() {

        override fun fromJson(reader: JsonReader): Any? {
            val fieldValue = fieldAdapter.fromJson(reader)
            return try {
                constructor.newInstance(fieldValue)
            } catch (e: Throwable) {
                throw IOException("Failed to create object with constructor " + constructor.name, e)
            }
        }

        @Suppress("UNCHECKED_CAST")
        override fun toJson(writer: JsonWriter, value: Any?) {
            if (value != null) {
                val fieldValue: I = try {
                    getter.invoke(value) as I
                } catch (e: IllegalAccessException) {
                    throw IOException("Failed to get value of field " + getter.name, e)
                }
                fieldAdapter.toJson(writer, fieldValue)
            } else {
                fieldAdapter.toJson(writer, null)
            }
        }

        override fun toString(): String {
            return "UnwrapFieldAdapter(" + getter.name + ")"
        }
    }
}
`

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("UnwrapFieldAdapterFactory.kt"),
		Content: strings.TrimSpace(code),
	}
}
