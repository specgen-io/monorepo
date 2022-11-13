package test_service.models

import com.squareup.moshi.*

import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Test
import kotlin.test.assertFailsWith

import test_service.json.setupMoshiAdapters

import java.math.BigDecimal
import java.time.*
import java.util.*

internal class ModelsTest {
    private fun createMoshiObject(): Moshi {
        val moshiBuilder = Moshi.Builder()
        setupMoshiAdapters(moshiBuilder)

        return moshiBuilder.build()
    }

    private fun <T> jsonAdapter(tClass: Class<T>): JsonAdapter<T> {
        return createMoshiObject().adapter(tClass)
    }

    private fun <T> check(data: T, jsonStr: String, tClass: Class<T>) {
        val actualJson: String = jsonAdapter(tClass).toJson(data)
        assertEquals(jsonStr, actualJson)

        val actualData: T = jsonAdapter(tClass).fromJson(jsonStr)!!
        assertEquals(actualData, data)
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
            jsonAdapter(Message::class.java).fromJson("""{"wrong_field":42}""")!!
        }
    }

    @Test
    fun objectModelMissingField() {
        assertFailsWith<Throwable> {
            jsonAdapter(Message::class.java).fromJson("""{}""")!!
        }
    }

    @Test
    fun objectModelFieldCases() {
        val data = MessageCases("snake_case value", "camelCase value")
        val jsonStr = """{"snake_case":"snake_case value","camelCase":"camelCase value"}"""
        check(data, jsonStr, MessageCases::class.java)
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
            jsonAdapter(Parent::class.java).fromJson(jsonStr)!!
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
        val jsonStr =
            """{"int_field":123,"long_field":1234,"float_field":1.23,"double_field":1.23,"decimal_field":1.23}"""
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
//      val jsonField = """{"the_array":[true,"some string"],"the_object":{"the_bool":true,"the_string":"some value"},"the_scalar":the value}"""

        val theObject = hashMapOf<String, Any>(
            "the_bool" to true,
            "the_string" to "some value"
        )
        val theArray = listOf<Any>(true, "some string")
        val map = hashMapOf<String, Any>(
            "the_array" to theArray,
            "the_object" to theObject,
            "the_scalar" to "the value"
        )

        val data = RawJsonField(map)
        val jsonStr =
            """{"json_field":{"the_array":[true,"some string"],"the_scalar":"the value","the_object":{"the_bool":true,"the_string":"some value"}}}"""
        check(data, jsonStr, RawJsonField::class.java)
    }

    @Test
    fun oneOfItemType() {
        val data = OrderCreated(UUID.fromString("58d5e212-165b-4ca0-909b-c86b9cee0111"), "SNI/01/136/0500", 3)
        val jsonStr = """{"id":"58d5e212-165b-4ca0-909b-c86b9cee0111","sku":"SNI/01/136/0500","quantity":3}"""
        check(data, jsonStr, OrderCreated::class.java)
    }

    @Test
    fun oneOfWrapper() {
        val data: OrderEventWrapper =
            OrderEventWrapper.Canceled(OrderCanceled(UUID.fromString("123e4567-e89b-12d3-a456-426655440000")))
        val jsonStr = """{"canceled":{"id":"123e4567-e89b-12d3-a456-426655440000"}}"""
        check(data, jsonStr, OrderEventWrapper::class.java)
    }

    @Test
    fun oneOfWrapperItemNotNull() {
        assertFailsWith<Throwable> {
            val jsonStr = """{"canceled":null}"""
            jsonAdapter(OrderEventWrapper::class.java).fromJson(jsonStr)!!
        }
    }

    @Test
    fun oneOfDiscriminator() {
        val data: OrderEventDiscriminator =
            OrderEventDiscriminator.Canceled(OrderCanceled(UUID.fromString("123e4567-e89b-12d3-a456-426655440000")))
        val jsonStr = """{"_type":"canceled","id":"123e4567-e89b-12d3-a456-426655440000"}"""
        check(data, jsonStr, OrderEventDiscriminator::class.java)
    }
}