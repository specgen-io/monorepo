package themodels

import com.fasterxml.jackson.databind.JsonNode
import com.fasterxml.jackson.databind.ObjectMapper
import org.junit.jupiter.api.Test
import themodels.models.RawJsonField

class JacksonTest {

    @Test
    fun jsonType() {
        val jsonField =
            """{"the_array":[1,"some string"],"the_object":{"the_bool":true,"the_string":"some value"},"the_scalar":123}"""
        val node: JsonNode = ObjectMapper().readTree(jsonField)
        val data = RawJsonField(node)
        val jsonStr =
            """{"json_field":{"the_array":[1,"some string"],"the_object":{"the_bool":true,"the_string":"some value"},"the_scalar":123}}"""
        check(data, jsonStr, RawJsonField::class.java)
    }
}