package [[.PackageName]]

import play.api.mvc.QueryStringBindable

object PlayParamsTypesBindings {
  implicit def bindableParser[T](implicit stringBinder: QueryStringBindable[String], codec: Codec[T]): QueryStringBindable[T] = new QueryStringBindable[T] {
    override def bind(key: String, params: Map[String, Seq[String]]): Option[Either[String, T]] =
      for {
        dateStr <- stringBinder.bind(key, params)
      } yield {
        dateStr match {
          case Right(value) =>
            try {
              Right(codec.decode(value))
            } catch {
              case t: Throwable => Left(s"Unable to bind from key: $key, error: ${t.getMessage}")
            }
          case _ => Left(s"Unable to bind from key: $key")
        }
      }

    override def unbind(key: String, value: T): String = stringBinder.unbind(key, codec.encode(value))
  }
}