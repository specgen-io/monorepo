package themodels

import org.junit.jupiter.api.Test
import themodels.models.RawJsonField

class MoshiTest : JsonTest() {

    @Test
    fun jsonType() {
//      val jsonField = """{"the_array":[true,"some string"],"the_object":{"the_bool":true,"the_string":"some value"},"the_scalar":the value}"""

        val theObject = hashMapOf(
            "the_bool" to true,
            "the_string" to "some value"
        )
        val theArray = listOf(true, "some string")
        val map = hashMapOf(
            "the_array" to theArray,
            "the_object" to theObject,
            "the_scalar" to "the value"
        )

        val data = RawJsonField(map)
        val jsonStr =
            """{"json_field":{"the_array":[true,"some string"],"the_scalar":"the value","the_object":{"the_bool":true,"the_string":"some value"}}}"""
        check(data, jsonStr, RawJsonField::class.java)
    }
}