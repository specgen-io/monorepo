val versionValue = System.getProperty("version", "0.0.2")

enablePlugins(SbtPlugin)
sbtPlugin := true

organization := "io.specgen"
name := "sbt-specgen"
version := versionValue

homepage := Some(url("https://specgen.io"))
scmInfo := Some(ScmInfo(url("https://github.com/specgen-io/specgen"), "git@github.com:specgen-io/specgen.git"))
developers := List(Developer("vsapronov", "vsapronov", "vladimir.sapronov@gmail.com", url("https://github.com/vsapronov")))
licenses += ("MIT", url("http://opensource.org/licenses/MIT"))

publishMavenStyle := true
crossPaths := false