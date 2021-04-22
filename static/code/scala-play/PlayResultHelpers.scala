package [[.PackageName]]

import play.api.mvc.Result
import play.api.mvc.Results._
import spec.services.OperationResult

object PlayResultHelpers {
  implicit class ResponsePlay(response: OperationResult) {
    def toPlay(): Result = {
      val status = new Status(response.status)
      response.body match {
        case Some(body) => status(body)
        case None => status
      }
    }
  }
}