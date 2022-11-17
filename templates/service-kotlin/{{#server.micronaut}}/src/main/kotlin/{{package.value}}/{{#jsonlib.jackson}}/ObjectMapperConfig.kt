package {{package.value}}

import com.fasterxml.jackson.databind.ObjectMapper
import com.fasterxml.jackson.module.kotlin.jacksonObjectMapper
import io.micronaut.context.annotation.*
import io.micronaut.jackson.ObjectMapperFactory
import {{package.value}}.json.*

@Factory
@Replaces(ObjectMapperFactory::class)
class ObjectMapperConfig {
    @Bean
    @Replaces(ObjectMapper::class)
    fun objectMapper(): ObjectMapper {
        val objectMapper = jacksonObjectMapper()
        setupObjectMapper(objectMapper)
        return objectMapper
    }

    @Bean
    fun json(): Json {
        return Json(objectMapper())
    }
}