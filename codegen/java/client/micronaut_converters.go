package client

import (
	"generator"
	"java/packages"
	"java/writer"
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
import io.micronaut.core.convert.*;
import jakarta.inject.Singleton;

@Singleton
public class [[.ClassName]] implements TypeConverterRegistrar {
	@Override
	public void register(ConversionService<?> conversionService) {
		conversionService.addConverter(
			String[].class,
			CharSequence.class,
			new StringArrayConverter(conversionService)
		);
	}
}
`)
	return w.ToCodeFile()
}

func localDateConverter(thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `LocalDateConverter`)
	w.Lines(`
import io.micronaut.core.convert.*;
import jakarta.inject.Singleton;

import java.time.LocalDate;
import java.time.format.DateTimeFormatter;
import java.util.Optional;

@Singleton
public class [[.ClassName]] implements TypeConverter<LocalDate, String> {

	public [[.ClassName]]() {
		System.out.println("Creating [[.ClassName]]");
	}

	@Override
	public Optional<String> convert(LocalDate value, Class<String> targetType, ConversionContext context) {
		var formatter = DateTimeFormatter.ofPattern("yyyy-MM-dd", context.getLocale());
		return Optional.of(value.format(formatter));
	}
}
`)
	return w.ToCodeFile()
}

func localDateTimeConverter(thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `LocalDateTimeConverter`)
	w.Lines(`
import io.micronaut.core.convert.*;
import jakarta.inject.Singleton;

import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.Optional;

@Singleton
public class [[.ClassName]] implements TypeConverter<LocalDateTime, String> {

	public [[.ClassName]]() {
		System.out.println("Creating [[.ClassName]]");
	}

	@Override
	public Optional<String> convert(LocalDateTime value, Class<String> targetType, ConversionContext context) {
		var formatter = DateTimeFormatter.ofPattern("yyyy-MM-dd'T'HH:mm:ss", context.getLocale());
		return Optional.of(value.format(formatter));
	}
}

`)
	return w.ToCodeFile()
}

func stringArrayConverter(thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `StringArrayConverter`)
	w.Lines(`
import io.micronaut.core.convert.*;
import jakarta.inject.Singleton;

import java.util.*;

@Singleton
public class [[.ClassName]] implements TypeConverter<String[], CharSequence> {
	private final ConversionService<?> conversionService;

	public [[.ClassName]](ConversionService<?> conversionService) {
		this.conversionService = conversionService;
		System.out.println("Creating [[.ClassName]]");
	}

	@Override
	public Optional<CharSequence> convert(String[] value, Class<CharSequence> targetType, ConversionContext context) {
		if (value == null) {
			return Optional.empty();
		}
		var joiner = new StringJoiner(",");
		for (String str : value) {
			joiner.add(str);
		}
		return conversionService.convert(joiner.toString(), targetType, context);
	}
}
`)
	return w.ToCodeFile()
}
