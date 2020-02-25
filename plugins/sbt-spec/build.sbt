val versionValue = System.getProperty("version", "0.0.0")

enablePlugins(SbtPlugin)
sbtPlugin := true

organization := "spec"
name := "sbt-spec"
version := versionValue

licenses += ("MIT", url("http://opensource.org/licenses/MIT"))
bintrayOrganization := Some("moda")
bintrayRepository := "sbt-plugins"
bintrayReleaseOnPublish := false