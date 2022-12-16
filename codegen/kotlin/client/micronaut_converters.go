package client

import (
	"generator"
	"kotlin/packages"
	"kotlin/writer"
)

func converters(convertersPackage packages.Package) []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *convertersRegistrar(convertersPackage))
	files = append(files, *localDateConverter(convertersPackage))
	files = append(files, *localDateTimeConverter(convertersPackage))
	files = append(files, *stringArrayConverter(convertersPackage))
	return files
}

func convertersRegistrar(thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `ConvertersRegistrar`)
	w.Lines(`
import io.micronaut.core.convert.*
import jakarta.inject.Singleton
import java.util.*

@Singleton
class [[.ClassName]]: TypeConverterRegistrar {
    override fun register(conversionService: ConversionService<*>) {
        conversionService.addConverter(Array<String>::class.java, CharSequence::class.java, StringArrayConverter(conversionService))
    }
}
`)
	return w.ToCodeFile()
}

func localDateConverter(thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `LocalDateConverter`)
	w.Lines(`
import io.micronaut.core.convert.*
import jakarta.inject.Singleton

import java.time.format.*
import java.time.LocalDate
import java.util.*

@Singleton
class [[.ClassName]] : TypeConverter<LocalDate, String> {
    init {
        println("Creating [[.ClassName]]")
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
`)
	return w.ToCodeFile()
}

func localDateTimeConverter(thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `LocalDateTimeConverter`)
	w.Lines(`
import io.micronaut.core.convert.*
import jakarta.inject.Singleton

import java.time.format.*
import java.time.LocalDateTime
import java.util.*

@Singleton
class [[.ClassName]] : TypeConverter<LocalDateTime, String> {
    init {
        println("Creating [[.ClassName]]")
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
`)
	return w.ToCodeFile()
}

func stringArrayConverter(thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `StringArrayConverter`)
	w.Lines(`
import io.micronaut.core.convert.*
import java.util.*

class [[.ClassName]](private val conversionService: ConversionService<*>) : TypeConverter<Array<String>, CharSequence> {
    init {
        println("Creating [[.ClassName]]")
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
`)
	return w.ToCodeFile()
}
