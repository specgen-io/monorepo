package themodels

import org.junit.jupiter.api.Test
import kotlin.test.assertFailsWith
import themodels.models.*
import java.math.BigDecimal
import java.time.*
import java.util.*

class JsonTest {
    val json = createJson()

    @Test
    fun objectModel() {
        val data = Message(42)
        val jsonStr = """{"field":42}"""
        check(data, jsonStr, Message::class.java)
    }

    @Test
    fun objectModelWrongFieldName() {
        assertFailsWith<Throwable> {
            json.read("""{"wrong_field":42}""", Message::class.java)
        }
    }

    @Test
    fun nestedObject() {
        val data = Parent("the string", Message(123))
        val jsonStr = """{"field":"the string","nested":{"field":123}}"""
        check(data, jsonStr, Parent::class.java)
    }

    @Test
    fun objectModelFieldCases() {
        val data = MessageCases("snake_case value", "camelCase value")
        val jsonStr = """{"snake_case":"snake_case value","camelCase":"camelCase value"}"""
        check(data, jsonStr, MessageCases::class.java)
    }

    @Test
    fun objectModelMissingField() {
        assertFailsWith<Throwable> {
            json.read("""{}""", Message::class.java)
        }
    }

    @Test
    fun objectFieldIsNull() {
        assertFailsWith<Throwable> {
            json.read("""{"field":"the string","nested":null}""", Parent::class.java)
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
    fun oneOfWrapper() {
        val data: OrderEventWrapper =
            OrderEventWrapper.Canceled(OrderCanceled(UUID.fromString("123e4567-e89b-12d3-a456-426655440000")))
        val jsonStr = """{"canceled":{"id":"123e4567-e89b-12d3-a456-426655440000"}}"""
        check(data, jsonStr, OrderEventWrapper::class.java)
    }

    @Test
    fun oneOfWrapperItemNotNull() {
        assertFailsWith<Throwable> {
            json.read("""{"canceled":null}""", OrderEventWrapper::class.java)
        }
    }

    @Test
    fun oneOfDiscriminatorTest() {
        val data: OrderEventDiscriminator =
            OrderEventDiscriminator.Canceled(OrderCanceled(UUID.fromString("123e4567-e89b-12d3-a456-426655440000")))
        val jsonStr = """{"_type":"canceled","id":"123e4567-e89b-12d3-a456-426655440000"}"""
        check(data, jsonStr, OrderEventDiscriminator::class.java)
    }
}