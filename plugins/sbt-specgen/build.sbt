val versionValue = System.getProperty("version", "0.0.2")

enablePlugins(SbtPlugin)
sbtPlugin := true

organization := "io.specgen"
name := "sbt-specgen"
version := versionValue

homepage := Some(url("https://specgen.io"))
scmInfo := Some(ScmInfo(url("https://github.com/specgen-io"), "git@github.com:specgen-io/monorepo.git"))
developers := List(Developer("specgen", "specgen", "dev@specgen.io", url("https://github.com/specgen-io")))
licenses += ("MIT", url("http://opensource.org/licenses/MIT"))

publishMavenStyle := true
crossPaths := false