name := "echo-client"

scalaVersion := "2.13.2"
version := "0.0.1"

enablePlugins(SpecgenModels)
specgenSpecFile := file("./../spec.yaml")
specgenJsonlib := "circe"

libraryDependencies ++= Seq(
  "io.circe" %% "circe-core" % "0.12.3",
  "io.circe" %% "circe-generic-extras" % "0.12.2",
  "io.circe" %% "circe-parser" % "0.12.3",
  "com.beachape" %% "enumeratum" % "1.5.13",
  "com.beachape" %% "enumeratum-circe" % "1.5.22",

  "org.scalatest" %% "scalatest" % "3.0.8" % Test
)

val junitxml = System.getProperty("junitxml", "target/junitxml")
(testOptions in Test) += Tests.Argument(TestFrameworks.ScalaTest, "-u", junitxml)
