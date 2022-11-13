package test_service.models

import com.fasterxml.jackson.databind.JsonNode
import com.fasterxml.jackson.module.kotlin.jacksonObjectMapper
import java.math.BigDecimal
import java.time.*
import java.util.*
import kotlin.test.*
import test_service.json.*

internal class ModelsTest {
    private val mapper = setupObjectMapper(jacksonObjectMapper())

    private fun <T> checkDeserialization(data: T, jsonStr: String, tClass: Class<T>) {
        val deserialized: T = mapper.readValue(jsonStr, tClass)
        assertEquals(data, deserialized)
    }

    private fun <T> checkSerialization(data: T, jsonStr: String) {
        val serialized = mapper.writeValueAsString(data)
        assertEquals(jsonStr, serialized)
    }

    private fun <T> check(data: T, jsonStr: String, tClass: Class<T>) {
        checkSerialization(data, jsonStr)
        checkDeserialization(data, jsonStr, tClass)
    }

    @Test
    fun objectModel() {
        val data = Message(42)
        val jsonStr = """{"field":42}"""
        check(data, jsonStr, Message::class.java)
    }

    @Test
    fun objectModelWrongFieldName() {
        assertFailsWith<Throwable> {
            mapper.readValue("""{"wrong_field":42}""", Message::class.java)
        }
    }

    @Test
    fun objectModelMissingField() {
        assertFailsWith<Throwable> {
            mapper.readValue("""{}""", Message::class.java)
        }
    }

    @Test
    fun nestedObject() {
        val data = Parent("the string", Message(123))
        val jsonStr = """{"field":"the string","nested":{"field":123}}"""
        check(data, jsonStr, Parent::class.java)
    }

    @Test
    fun objectFieldNotNull() {
        assertFailsWith<Throwable> {
            val jsonStr = """{"field":"the string","nested":null}"""
            mapper.readValue(jsonStr, Parent::class.java)
        }
    }

    @Test
    fun enumModel() {
        val data = EnumFields(Choice.SECOND_CHOICE)
        val jsonStr = """{"enum_field":"Two"}"""
        check(data, jsonStr, EnumFields::class.java)
    }

    @Test
    fun numericTypes() {
        val data = NumericFields(123, 1234, 1.23f, 1.23, BigDecimal("1.23"))
        val jsonStr = """{"int_field":123,"long_field":1234,"float_field":1.23,"double_field":1.23,"decimal_field":1.23}"""
        check(data, jsonStr, NumericFields::class.java)
    }

    @Test
    fun nonNumericTypes() {
        val data = NonNumericFields(
            true,
            "the string",
            UUID.fromString("123e4567-e89b-12d3-a456-426655440000"),
            LocalDate.parse("2019-11-30"),
            LocalDateTime.parse("2019-11-30T17:45:55")
        )
        val jsonStr =
            """{"boolean_field":true,"string_field":"the string","uuid_field":"123e4567-e89b-12d3-a456-426655440000","date_field":"2019-11-30","datetime_field":"2019-11-30T17:45:55"}"""
        check(data, jsonStr, NonNumericFields::class.java)
    }

    @Test
    fun arrayType() {
        val data = ArrayFields(listOf(1, 2, 3), listOf("one", "two", "three"))
        val jsonStr = """{"int_array_field":[1,2,3],"string_array_field":["one","two","three"]}"""
        check(data, jsonStr, ArrayFields::class.java)
    }

    @Test
    fun mapType() {
        val data = MapFields(
            hashMapOf("one" to 1, "two" to 2),
            hashMapOf("one" to "first", "two" to "second")
        )
        val jsonStr = """{"int_map_field":{"two":2,"one":1},"string_map_field":{"two":"second","one":"first"}}"""
        check(data, jsonStr, MapFields::class.java)
    }

    @Test
    fun optionalTypes() {
        val data = OptionalFields(123, "the string")
        val jsonStr = """{"int_option_field":123,"string_option_field":"the string"}"""
        check(data, jsonStr, OptionalFields::class.java)
    }

    @Test
    fun optionalTypesMissingFields() {
        val data = OptionalFields(null, null)
        val jsonStr = "{}"
        check(data, jsonStr, OptionalFields::class.java)
    }

    @Test
    fun jsonType() {
        val jsonField = """{"the_array":[1,"some string"],"the_object":{"the_bool":true,"the_string":"some value"},"the_scalar":123}"""
        val node: JsonNode = mapper.readTree(jsonField)
        val data = RawJsonField(node)
        val jsonStr = """{"json_field":{"the_array":[1,"some string"],"the_object":{"the_bool":true,"the_string":"some value"},"the_scalar":123}}"""
        check(data, jsonStr, RawJsonField::class.java)
    }

    @Test
    fun oneOfWrapper() {
        val canceled = OrderCanceled(UUID.fromString("123e4567-e89b-12d3-a456-426655440000"))
        val event: OrderEventWrapper = OrderEventWrapper.Canceled(canceled)
        val jsonStr = """{"canceled":{"id":"123e4567-e89b-12d3-a456-426655440000"}}"""
        check(event, jsonStr, OrderEventWrapper::class.java)
    }

    @Test
    fun oneOfWrapperItemNotNull() {
        assertFailsWith<Throwable> {
            val jsonStr = """{"canceled":null}"""
            mapper.readValue(jsonStr, OrderEventWrapper::class.java)
        }
    }

    @Test
    fun oneOfDiscriminatorTest() {
        val canceled = OrderCanceled(UUID.fromString("123e4567-e89b-12d3-a456-426655440000"))
        val event: OrderEventDiscriminator = OrderEventDiscriminator.Canceled(canceled)
        val jsonStr = """{"_type":"canceled","id":"123e4567-e89b-12d3-a456-426655440000"}"""
        check(event, jsonStr, OrderEventDiscriminator::class.java)
    }
}