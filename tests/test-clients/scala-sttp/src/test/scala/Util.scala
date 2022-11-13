package util

object Util {
  def service_url(): String = {
    val service_url = sys.env.get("SERVICE_URL")
    service_url match {
      case Some(service_url) => service_url
      case None => throw new Exception("There's no SERVICE_URL environment variable provided, it's required for tests to connect to service")
    }
  }
}
