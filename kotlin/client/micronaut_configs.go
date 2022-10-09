package client

import (
	"kotlin/writer"
	"strings"

	"generator"
	"kotlin/packages"
)

func staticConfigFiles(thePackage, jsonPackage packages.Package) []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *objectMapperConfig(thePackage, jsonPackage))
	files = append(files, *clientConfig(thePackage))

	return files
}

func objectMapperConfig(thePackage packages.Package, jsonPackage packages.Package) *generator.CodeFile {
	code := `
package [[.PackageName]]

import com.fasterxml.jackson.databind.ObjectMapper
import com.fasterxml.jackson.module.kotlin.jacksonObjectMapper
import io.micronaut.context.annotation.*
import io.micronaut.jackson.ObjectMapperFactory
import jakarta.inject.Singleton
import [[.JsonPackageName]].*

@Factory
@Replaces(ObjectMapperFactory::class)
class ObjectMapperConfig {
    @Singleton
    @Replaces(ObjectMapper::class)
    fun objectMapper(): ObjectMapper {
        val objectMapper = jacksonObjectMapper()
        setupObjectMapper(objectMapper)
        return objectMapper
    }
}
`

	code, _ = generator.ExecuteTemplate(code, struct {
		PackageName     string
		JsonPackageName string
	}{
		thePackage.PackageName,
		jsonPackage.PackageName,
	})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("ObjectMapperConfig.kt"),
		Content: strings.TrimSpace(code),
	}
}

func clientConfig(thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `ClientConfiguration`)
	w.Lines(`
class ClientConfiguration {
    companion object {
        const val BASE_URL = "http://localhost:8081"
    }
}
`)
	return w.ToCodeFile()
}
