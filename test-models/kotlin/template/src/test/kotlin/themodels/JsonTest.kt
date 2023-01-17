package themodels

import com.fasterxml.jackson.module.kotlin.jacksonObjectMapper
import org.junit.jupiter.api.Assertions.assertEquals
import themodels.json.Json
import themodels.json.setupObjectMapper

open class JsonTest {
    protected val json = createJson()

    private fun createJson(): Json {
        val objectMapper = jacksonObjectMapper()
        setupObjectMapper(objectMapper)
        return Json(objectMapper)
    }

    fun <T> check(data: T, jsonStr: String, tClass: Class<T>) {
        val actualJson = json.read(jsonStr, tClass)
        assertEquals(data, actualJson)

        val actualData = json.write(data!!)
        assertEquals(jsonStr, actualData)
    }
}