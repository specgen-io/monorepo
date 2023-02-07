package theclient

import org.assertj.core.api.Assertions.*
import org.junit.jupiter.api.Assertions.*
import org.junit.jupiter.api.Test
import theclient.clients.echo.*
import theclient.models.*
import java.math.BigDecimal
import java.time.*
import java.util.*

class EchoClientTest {
    private val client = EchoClient("http://localhost:8081")

    private val intValue: Int = 123
    private val longValue: Long = 12345
    private val floatValue: Float = 1.23f
    private val doubleValue: Double = 12.345
    private val decimalValue: BigDecimal = BigDecimal("12345")
    private val boolValue: Boolean = true
    private val stringValue: String = "the value"
    private val stringOptValue: String = "the value"
    private val stringDefaultedValue: String = "value"
    private val stringArrayValue: List<String> = listOf("the str1", "the str2")
    private val uuidValue: UUID = UUID.fromString("123e4567-e89b-12d3-a456-426655440000")
    private val dateValue: LocalDate = LocalDate.parse("2020-01-01")
    private val dateArrayValue: List<LocalDate> = listOf(LocalDate.parse("2020-01-01"), LocalDate.parse("2020-01-02"))
    private val datetimeValue: LocalDateTime = LocalDateTime.parse("2019-11-30T17:45:55")
    private val enumValue = Choice.SECOND_CHOICE

    private val bodyString = "TWFueSBoYW5kcyBtYWtlIGxpZ2h0IHdvcmsu"
    private val bodyMessage = Message(intValue, stringValue)
    private val stringList = listOf("the str1", "the str2")
    private val stringHashMap = hashMapOf(
        "string_field" to "the value",
        "string_field_2" to "the value_2"
    )

    @Test
    fun echoBodyString_responseIsEqualToRequest() {
        val response = client.echoBodyString(bodyString)
        assertEquals(bodyString, response)
    }

    @Test
    fun echoBodyString_doesntThrowException() {
        assertDoesNotThrow<String> { client.echoBodyString(bodyString) }
    }

    @Test
    fun echoBodyModel_responseIsEqualToRequest() {
        val response = client.echoBodyModel(bodyMessage)
        assertEquals(bodyMessage, response)
    }

    @Test
    fun echoBodyModel_doesntThrowException() {
        assertDoesNotThrow<Message> { client.echoBodyModel(bodyMessage) }
    }

    @Test
    fun echoBodyArray_responseIsEqualToRequest() {
        val response = client.echoBodyArray(stringList)
        assertEquals(stringList, response)
    }

    @Test
    fun echoBodyArray_doesntThrowException() {
        assertDoesNotThrow<List<String>> { client.echoBodyArray(stringList) }
    }

    @Test
    fun echoBodyMap_responseIsEqualToRequest() {
        val response = client.echoBodyMap(stringHashMap)
        assertEquals(stringHashMap, response)
    }

    @Test
    fun echoBodyMap_doesntThrowException() {
        assertDoesNotThrow<Map<String, String>> { client.echoBodyMap(stringHashMap) }
    }

    @Test
    fun echoQuery_responseIsEqualToRequest() {
        val request = Parameters(
            intValue,
            longValue,
            floatValue,
            doubleValue,
            decimalValue,
            boolValue,
            stringValue,
            stringOptValue,
            stringDefaultedValue,
            stringArrayValue,
            uuidValue,
            dateValue,
            dateArrayValue,
            datetimeValue,
            enumValue
        )
        val response = client.echoQuery(
            intValue,
            longValue,
            floatValue,
            doubleValue,
            decimalValue,
            boolValue,
            stringValue,
            stringOptValue,
            stringDefaultedValue,
            stringArrayValue,
            uuidValue,
            dateValue,
            dateArrayValue,
            datetimeValue,
            enumValue
        )
        assertEquals(request, response)
    }

    @Test
    fun echoQuery_doesntThrowException() {
        assertDoesNotThrow<Parameters> {
            client.echoQuery(
                intValue,
                longValue,
                floatValue,
                doubleValue,
                decimalValue,
                boolValue,
                stringValue,
                stringOptValue,
                stringDefaultedValue,
                stringArrayValue,
                uuidValue,
                dateValue,
                dateArrayValue,
                datetimeValue,
                enumValue
            )
        }
    }

    @Test
    fun echoHeader_responseIsEqualToRequest() {
        val request = Parameters(
            intValue,
            longValue,
            floatValue,
            doubleValue,
            decimalValue,
            boolValue,
            stringValue,
            stringOptValue,
            stringDefaultedValue,
            stringArrayValue,
            uuidValue,
            dateValue,
            dateArrayValue,
            datetimeValue,
            enumValue
        )
        val response = client.echoHeader(
            intValue,
            longValue,
            floatValue,
            doubleValue,
            decimalValue,
            boolValue,
            stringValue,
            stringOptValue,
            stringDefaultedValue,
            stringArrayValue,
            uuidValue,
            dateValue,
            dateArrayValue,
            datetimeValue,
            enumValue
        )
        assertEquals(request, response)
    }

    @Test
    fun echoHeader_doesntThrowException() {
        assertDoesNotThrow<Parameters> {
            client.echoHeader(
                intValue,
                longValue,
                floatValue,
                doubleValue,
                decimalValue,
                boolValue,
                stringValue,
                stringOptValue,
                stringDefaultedValue,
                stringArrayValue,
                uuidValue,
                dateValue,
                dateArrayValue,
                datetimeValue,
                enumValue
            )
        }
    }

    @Test
    fun echoUrlParams_responseIsEqualToRequest() {
        val request = UrlParameters(
            intValue,
            longValue,
            floatValue,
            doubleValue,
            decimalValue,
            boolValue,
            stringValue,
            uuidValue,
            dateValue,
            datetimeValue,
            enumValue
        )
        val response = client.echoUrlParams(
            intValue,
            longValue,
            floatValue,
            doubleValue,
            decimalValue,
            boolValue,
            stringValue,
            uuidValue,
            dateValue,
            datetimeValue,
            enumValue
        )
        assertEquals(request, response)
    }

    @Test
    fun echoUrlParams_doesntThrowException() {
        assertDoesNotThrow<UrlParameters> {
            client.echoUrlParams(
                intValue,
                longValue,
                floatValue,
                doubleValue,
                decimalValue,
                boolValue,
                stringValue,
                uuidValue,
                dateValue,
                datetimeValue,
                enumValue
            )
        }
    }

    @Test
    fun echoEverything_responseIsEqualToRequest() {
        val request = Everything(
            bodyMessage,
            floatValue,
            boolValue,
            uuidValue,
            datetimeValue,
            dateValue,
            decimalValue
        )
        val response = client.echoEverything(
            bodyMessage,
            floatValue,
            boolValue,
            uuidValue,
            datetimeValue,
            dateValue,
            decimalValue
        )
        assertThat(request).usingRecursiveComparison().isEqualTo(response)
    }

    @Test
    fun echoEverything_doesntThrowException() {
        assertDoesNotThrow {
            client.echoEverything(
                bodyMessage,
                floatValue,
                boolValue,
                uuidValue,
                datetimeValue,
                dateValue,
                decimalValue
            )
        }
    }

    @Test
    fun echoSuccess_ok() {
        val expected = EchoSuccessResponse.Ok(OkResult("ok"))
        val response = client.echoSuccess("ok")

        assertThat(expected).usingRecursiveComparison().isEqualTo(response)
    }

    @Test
    fun echoSuccess_created() {
        val expected = EchoSuccessResponse.Created(CreatedResult("created"))
        val response = client.echoSuccess("created")

        assertThat(expected).usingRecursiveComparison().isEqualTo(response)
    }

    @Test
    fun echoSuccess_accepted() {
        val expected = EchoSuccessResponse.Accepted(AcceptedResult("accepted"))
        val response = client.echoSuccess("accepted")

        assertThat(expected).usingRecursiveComparison().isEqualTo(response)
    }
}
