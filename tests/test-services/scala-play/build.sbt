enablePlugins(SpecgenService)
enablePlugins(PlayScala)

version := "0.0.1"

scalaVersion := "2.13.2"
name := "echo"

specgenSpecFile := file("./../spec.yaml")
specgenServer := "play"
specgenJsonlib := "circe"

libraryDependencies ++= Seq(
  "io.circe" %% "circe-core" % "0.12.3",
  "io.circe" %% "circe-generic-extras" % "0.12.2",
  "io.circe" %% "circe-parser" % "0.12.3",
  "com.beachape" %% "enumeratum" % "1.5.13",
  "com.beachape" %% "enumeratum-circe" % "1.5.22",

  "com.typesafe.play" %% "play" % "2.8.1",
  "org.webjars" % "swagger-ui" % "3.22.2",

  guice
)

PlayKeys.playDefaultPort := 8081