package v2.services

import javax.inject._
import scala.concurrent._
import v2.services.echo._
import v2.models._

@Singleton
class EchoService @Inject()()(implicit ec: ExecutionContext) extends IEchoService {
  override def echoBodyModel(body: Message): Future[Message] = Future {
    Message(body.boolField, body.stringField)
  }
}
