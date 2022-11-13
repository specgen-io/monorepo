name := "echo-client"

scalaVersion := "2.13.2"
version := "0.0.1"

enablePlugins(SpecgenClient)
specgenSpecFile := file("./../spec.yaml")
specgenClient := "sttp"
specgenJsonlib := "circe"

libraryDependencies ++= Seq(
  "io.circe" %% "circe-core" % "0.12.3",
  "io.circe" %% "circe-generic-extras" % "0.12.2",
  "io.circe" %% "circe-parser" % "0.12.3",
  "com.beachape" %% "enumeratum" % "1.5.13",
  "com.beachape" %% "enumeratum-circe" % "1.5.22",

  "org.slf4j" % "slf4j-api" % "1.7.28",
  "com.softwaremill.sttp" %% "core" % "1.7.1",

  "org.scalatest" %% "scalatest" % "3.0.8" % Test,
  "com.softwaremill.sttp" %% "akka-http-backend" % "1.7.1" % Test,
  "com.typesafe.akka" %% "akka-stream" % "2.5.23" % Test,
)

val junitxml = System.getProperty("junitxml", "target/junitxml")
(testOptions in Test) += Tests.Argument(TestFrameworks.ScalaTest, "-u", junitxml)
