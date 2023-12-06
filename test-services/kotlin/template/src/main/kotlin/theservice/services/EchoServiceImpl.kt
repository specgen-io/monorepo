package theservice.services

{{#server.micronaut}}
import io.micronaut.context.annotation.Bean
{{/server.micronaut}}
{{#server.spring}}
import org.springframework.stereotype.Service
{{/server.spring}}
import theservice.models.*
import java.math.BigDecimal
import java.util.UUID
import java.time.*
import theservice.services.echo.*

{{#server.micronaut}}
@Bean
{{/server.micronaut}}
{{#server.spring}}
@Service("EchoService")
{{/server.spring}}
class EchoServiceImpl : EchoService {
    override fun echoBodyString(body: String): String {
        return body
    }

    override fun echoBodyModel(body: Message): Message {
        return body
    }

    override fun echoBodyArray(body: List<String>): List<String> {
        return body
    }

    override fun echoBodyMap(body: Map<String, String>): Map<String, String> {
        return body
    }

//    override fun echoFormData(
//        intParam: Int,
//        longParam: Long,
//        floatParam: Float,
//        doubleParam: Double,
//        decimalParam: BigDecimal,
//        boolParam: Boolean,
//        stringParam: String,
//        stringOptParam: String?,
//        stringDefaultedParam: String,
//        stringArrayParam: List<String>,
//        uuidParam: UUID,
//        dateParam: LocalDate,
//        dateArrayParam: List<LocalDate>,
//        datetimeParam: LocalDateTime,
//        enumParam: Choice
//    ): FormParameters {
//        return FormParameters(
//            intParam,
//            longParam,
//            floatParam,
//            doubleParam,
//            decimalParam,
//            boolParam,
//            stringParam,
//            stringOptParam,
//            stringDefaultedParam,
//            stringArrayParam,
//            uuidParam,
//            dateParam,
//            dateArrayParam,
//            datetimeParam,
//            enumParam
//        )
//    }
//
//    override fun echoFormUrlencoded(
//        intParam: Int,
//        longParam: Long,
//        floatParam: Float,
//        doubleParam: Double,
//        decimalParam: BigDecimal,
//        boolParam: Boolean,
//        stringParam: String,
//        stringOptParam: String?,
//        stringDefaultedParam: String,
//        stringArrayParam: List<String>,
//        uuidParam: UUID,
//        dateParam: LocalDate,
//        dateArrayParam: List<LocalDate>,
//        datetimeParam: LocalDateTime,
//        enumParam: Choice
//    ): FormParameters {
//        return FormParameters(
//            intParam,
//            longParam,
//            floatParam,
//            doubleParam,
//            decimalParam,
//            boolParam,
//            stringParam,
//            stringOptParam,
//            stringDefaultedParam,
//            stringArrayParam,
//            uuidParam,
//            dateParam,
//            dateArrayParam,
//            datetimeParam,
//            enumParam
//        )
//    }

    override fun echoQuery(
        intQuery: Int,
        longQuery: Long,
        floatQuery: Float,
        doubleQuery: Double,
        decimalQuery: BigDecimal,
        boolQuery: Boolean,
        stringQuery: String,
        stringOptQuery: String?,
        stringDefaultedQuery: String,
        stringArrayQuery: List<String>,
        uuidQuery: UUID,
        dateQuery: LocalDate,
        dateArrayQuery: List<LocalDate>,
        datetimeQuery: LocalDateTime,
        enumQuery: Choice
    ): Parameters {
        return Parameters(
            intQuery,
            longQuery,
            floatQuery,
            doubleQuery,
            decimalQuery,
            boolQuery,
            stringQuery,
            stringOptQuery,
            stringDefaultedQuery,
            stringArrayQuery,
            uuidQuery,
            dateQuery,
            dateArrayQuery,
            datetimeQuery,
            enumQuery
        )
    }

    override fun echoHeader(
        intHeader: Int,
        longHeader: Long,
        floatHeader: Float,
        doubleHeader: Double,
        decimalHeader: BigDecimal,
        boolHeader: Boolean,
        stringHeader: String,
        stringOptHeader: String?,
        stringDefaultedHeader: String,
        stringArrayHeader: List<String>,
        uuidHeader: UUID,
        dateHeader: LocalDate,
        dateArrayHeader: List<LocalDate>,
        datetimeHeader: LocalDateTime,
        enumHeader: Choice
    ): Parameters {
        return Parameters(
            intHeader,
            longHeader,
            floatHeader,
            doubleHeader,
            decimalHeader,
            boolHeader,
            stringHeader,
            stringOptHeader,
            stringDefaultedHeader,
            stringArrayHeader,
            uuidHeader,
            dateHeader,
            dateArrayHeader,
            datetimeHeader,
            enumHeader
        )
    }

    override fun echoUrlParams(
        intUrl: Int,
        longUrl: Long,
        floatUrl: Float,
        doubleUrl: Double,
        decimalUrl: BigDecimal,
        boolUrl: Boolean,
        stringUrl: String,
        uuidUrl: UUID,
        dateUrl: LocalDate,
        datetimeUrl: LocalDateTime,
        enumUrl: Choice
    ): UrlParameters {
        return UrlParameters(
            intUrl,
            longUrl,
            floatUrl,
            doubleUrl,
            decimalUrl,
            boolUrl,
            stringUrl,
            uuidUrl,
            dateUrl,
            datetimeUrl,
            enumUrl
        )
    }

    override fun echoEverything(
        body: Message,
        floatQuery: Float,
        boolQuery: Boolean,
        uuidHeader: UUID,
        datetimeHeader: LocalDateTime,
        dateUrl: LocalDate,
        decimalUrl: BigDecimal
    ): EchoEverythingResponse {
        return EchoEverythingResponse.Ok(
            Everything(
                body,
                floatQuery,
                boolQuery,
                uuidHeader,
                datetimeHeader,
                dateUrl,
                decimalUrl
            )
        )
    }

    override fun sameOperationName(): SameOperationNameResponse {
        return SameOperationNameResponse.Ok()
    }
}
