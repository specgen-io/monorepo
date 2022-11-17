package services

import java.time.LocalDate
import java.util.UUID
import javax.inject._
import scala.concurrent._
import services.echo._
import models._

@Singleton
class EchoService @Inject()()(implicit ec: ExecutionContext) extends IEchoService {
  override def echoBodyString(body: String): Future[String] = Future { body }

  override def echoBodyModel(body: Message): Future[Message] = Future { body }

  override def echoBodyArray(body: List[String]): Future[List[String]] = Future { body }

  override def echoBodyMap(body: Map[String, String]): Future[Map[String, String]] = Future { body }

  override def echoQuery(
                          intQuery: Int,
                          longQuery: Long,
                          floatQuery: Float,
                          doubleQuery: Double,
                          decimalQuery: BigDecimal,
                          boolQuery: Boolean,
                          stringQuery: String,
                          stringOptQuery: Option[String],
                          stringDefaultedQuery: String,
                          stringArrayQuery: List[String],
                          uuidQuery: java.util.UUID,
                          dateQuery: java.time.LocalDate,
                          dateArrayQuery: List[java.time.LocalDate],
                          datetimeQuery: java.time.LocalDateTime,
                          enumQuery: Choice): Future[Parameters] = Future {
    Parameters(
      intField = intQuery,
      longField = longQuery,
      floatField = floatQuery,
      doubleField = doubleQuery,
      decimalField = decimalQuery,
      boolField = boolQuery,
      stringField = stringQuery,
      stringOptField = stringOptQuery,
      stringDefaultedField = stringDefaultedQuery,
      stringArrayField = stringArrayQuery,
      uuidField = uuidQuery,
      dateField = dateQuery,
      dateArrayField = dateArrayQuery,
      datetimeField = datetimeQuery,
      enumField = enumQuery
    )
  }

  override def echoHeader(
                           intHeader: Int,
                           longHeader: Long,
                           floatHeader: Float,
                           doubleHeader: Double,
                           decimalHeader: BigDecimal,
                           boolHeader: Boolean,
                           stringHeader: String,
                           stringOptHeader: Option[String],
                           stringDefaultedHeader: String,
                           stringArrayHeader: List[String],
                           uuidHeader: java.util.UUID,
                           dateHeader: java.time.LocalDate,
                           dateArrayHeader: List[java.time.LocalDate],
                           datetimeHeader: java.time.LocalDateTime,
                           enumHeader: Choice): Future[Parameters] = Future {
    Parameters(
      intField = intHeader,
      longField = longHeader,
      floatField = floatHeader,
      doubleField = doubleHeader,
      decimalField = decimalHeader,
      boolField = boolHeader,
      stringField = stringHeader,
      stringOptField = stringOptHeader,
      stringDefaultedField = stringDefaultedHeader,
      stringArrayField = stringArrayHeader,
      uuidField = uuidHeader,
      dateField = dateHeader,
      dateArrayField = dateArrayHeader,
      datetimeField = datetimeHeader,
      enumField = enumHeader
    )
  }

  override def echoUrlParams(
                              intUrl: Int,
                              longUrl: Long,
                              floatUrl: Float,
                              doubleUrl: Double,
                              decimalUrl: BigDecimal,
                              boolUrl: Boolean,
                              stringUrl: String,
                              uuidUrl: java.util.UUID,
                              dateUrl: java.time.LocalDate,
                              datetimeUrl: java.time.LocalDateTime,
                              enumUrl: Choice): Future[UrlParameters] = Future {
    UrlParameters(
      intField = intUrl,
      longField = longUrl,
      floatField = floatUrl,
      doubleField = doubleUrl,
      decimalField = decimalUrl,
      boolField = boolUrl,
      stringField = stringUrl,
      uuidField = uuidUrl,
      dateField = dateUrl,
      datetimeField = datetimeUrl,
      enumField = enumUrl
    )
  }

  override def echoEverything(
                               uuidHeader: java.util.UUID,
                               datetimeHeader: java.time.LocalDateTime,
                               body: Message,
                               dateUrl: java.time.LocalDate,
                               decimalUrl: BigDecimal,
                               floatQuery: Float,
                               boolQuery: Boolean): Future[EchoEverythingResponse] = Future {
    EchoEverythingResponse.Ok(
      Everything(
        bodyField = body,
        floatQuery = floatQuery,
        boolQuery = boolQuery,
        uuidHeader = uuidHeader,
        datetimeHeader = datetimeHeader,
        dateUrl = dateUrl,
        decimalUrl = decimalUrl
      )
    )
  }

  override def sameOperationName(): Future[SameOperationNameResponse] = Future {
    SameOperationNameResponse.Ok()
  }
}
