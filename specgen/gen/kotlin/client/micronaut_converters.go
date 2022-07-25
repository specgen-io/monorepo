package client

import (
	"strings"

	"github.com/specgen-io/specgen/generator/v2"
	"github.com/specgen-io/specgen/v2/gen/kotlin/modules"
)

func converters(thePackage modules.Module) []generator.CodeFile {
	convertersPackage := thePackage.Subpackage("converters")

	files := []generator.CodeFile{}
	files = append(files, *convertersRegistrar(convertersPackage))
	files = append(files, *localDateConverter(convertersPackage))
	files = append(files, *localDateTimeConverter(convertersPackage))
	files = append(files, *stringArrayConverter(convertersPackage))

	return files
}

func convertersRegistrar(thePackage modules.Module) *generator.CodeFile {
	code := `
package [[.PackageName]]

import io.micronaut.core.convert.*
import jakarta.inject.Singleton
import java.util.*

@Singleton
class ConvertersRegistrar: TypeConverterRegistrar {
    override fun register(conversionService: ConversionService<*>) {
        conversionService.addConverter(Array<String>::class.java, CharSequence::class.java, StringArrayConverter(conversionService))
    }
}
`

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("ConvertersRegistrar.kt"),
		Content: strings.TrimSpace(code),
	}
}

func localDateConverter(thePackage modules.Module) *generator.CodeFile {
	code := `
package [[.PackageName]]

import io.micronaut.core.convert.*
import jakarta.inject.Singleton

import java.time.format.*
import java.time.LocalDate
import java.util.*

@Singleton
class LocalDateConverter : TypeConverter<LocalDate, String> {
    init {
        println("Creating LocalDateConverter")
    }

    override fun convert(
        value: LocalDate,
        targetType: Class<String>,
        context: ConversionContext
    ): Optional<String> {
        val formatter = DateTimeFormatter.ofPattern("yyyy-MM-dd", context.locale)
        return Optional.of(value.format(formatter))
    }
}
`

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("LocalDateConverter.kt"),
		Content: strings.TrimSpace(code),
	}
}

func localDateTimeConverter(thePackage modules.Module) *generator.CodeFile {
	code := `
package [[.PackageName]]

import io.micronaut.core.convert.*
import jakarta.inject.Singleton

import java.time.format.*
import java.time.LocalDateTime
import java.util.*

@Singleton
class LocalDateTimeConverter : TypeConverter<LocalDateTime, String> {
    init {
        println("Creating LocalDateTimeConverter")
    }

    override fun convert(
        value: LocalDateTime,
        targetType: Class<String>,
        context: ConversionContext
    ): Optional<String> {
        val formatter = DateTimeFormatter.ofPattern("yyyy-MM-dd'T'HH:mm:ss", context.locale)
        return Optional.of(value.format(formatter))
    }
}
`

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("LocalDateTimeConverter.kt"),
		Content: strings.TrimSpace(code),
	}
}

func stringArrayConverter(thePackage modules.Module) *generator.CodeFile {
	code := `
package [[.PackageName]]

import io.micronaut.core.convert.*
import java.util.*

class StringArrayConverter(private val conversionService: ConversionService<*>) : TypeConverter<Array<String>, CharSequence> {
    init {
        println("Creating StringArrayConverter")
    }

    override fun convert(
        value: Array<String>,
        targetType: Class<CharSequence>,
        context: ConversionContext
    ): Optional<CharSequence> {
        if (value.isEmpty()) {
            return Optional.empty<CharSequence>()
        }
        val joiner = StringJoiner(",")
        for (string in value) {
            joiner.add(string)
        }
        return conversionService.convert(joiner.toString(), targetType, context)
    }
}
`

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("StringArrayConverter.kt"),
		Content: strings.TrimSpace(code),
	}
}
