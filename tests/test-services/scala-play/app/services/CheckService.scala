package services

import javax.inject._
import scala.concurrent._
import services.check._
import models._
import errors._

@Singleton
class CheckService @Inject()()(implicit ec: ExecutionContext) extends ICheckService {
  override def checkEmpty(): Future[Unit] = Future {}

  override def checkEmptyResponse(body: Message): Future[Unit] = Future {}

  override def checkForbidden(): Future[CheckForbiddenResponse] = Future {
    CheckForbiddenResponse.Forbidden()
  }

  override def sameOperationName(): Future[SameOperationNameResponse] = Future {
    SameOperationNameResponse.Ok()
  }

  override def checkBadRequest(): Future[CheckBadRequestResponse] = Future {
    CheckBadRequestResponse.BadRequest(BadRequestError("Error returned from service implementation", ErrorLocation.Unknown, None))
  }
}
