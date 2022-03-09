package spec

import sbt._
import sbt.Keys._
import Path.contentOf
import Specgen._

object SpecKeys {
  lazy val specgenSpecFile = settingKey[File]("Path to service specification file")
  lazy val specgenSwaggerFile = settingKey[File]("Path to generated OpenAPI specification file")
  lazy val specgenGeneratePath = settingKey[File]("Path to generate source code into")
  lazy val specgenServicesPath = settingKey[File]("Path to scaffold services code")
  lazy val specgenJsonlib = settingKey[String]("JSON library name, supported: circe")
  lazy val specgenClient = settingKey[String]("HTTP client library name, supported: sttp")
  lazy val specgenServer = settingKey[String]("HTTP server framework name, supported: play")

  lazy val specgenServiceTask = taskKey[Seq[File]]("Run service Play code generation for spec")
  lazy val specgenModelsTask = taskKey[Seq[File]]("Run circe models code generation for spec")
  lazy val specgenClientTask = taskKey[Seq[File]]("Run client Sttp code generation for spec")
}

object CommonKeys extends AutoPlugin {
  val autoImport = SpecKeys
}

object SpecgenService extends AutoPlugin {
  import SpecKeys._

  private def specgenTask = Def.task {
    serviceScala(
      sLog.value,
      specgenServer.value,
      specgenJsonlib.value,
      specgenSpecFile.value,
      specgenSwaggerFile.value,
      specgenGeneratePath.value,
      specgenServicesPath.value,
    )
  }

  override val projectSettings = Seq(
    specgenServer := "play",
    specgenJsonlib := "circe",
    specgenSpecFile := file("spec.yaml"),
    specgenSwaggerFile := baseDirectory.value / "public" / "swagger.yaml",
    specgenGeneratePath := (sourceManaged in Compile).value / "spec",
    specgenServicesPath := (scalaSource in Compile).value,
    specgenServiceTask := specgenTask.value,
    sourceGenerators in Compile += specgenServiceTask,
    Compile / packageSrc / mappings ++= contentOf(specgenGeneratePath.value)
  )
}

object SpecgenModels extends AutoPlugin {
  import SpecKeys._

  private def specgenTask = Def.task {
    modelsScala(
      sLog.value,
      specgenJsonlib.value,
      specgenSpecFile.value,
      specgenGeneratePath.value,
    )
  }

  override val projectSettings = Seq(
    specgenJsonlib := "circe",
    specgenSpecFile := file("spec.yaml"),
    specgenGeneratePath := (sourceManaged in Compile).value / "spec",
    specgenModelsTask := specgenTask.value,
    sourceGenerators in Compile += specgenModelsTask,
    Compile / packageSrc := ((Compile / packageSrc) dependsOn specgenTask).value,
    Compile / packageSrc / mappings ++= contentOf(specgenGeneratePath.value)
  )
}

object SpecgenClient extends AutoPlugin {
  import SpecKeys._

  private def specgenTask = Def.task {
    clientScala(
      sLog.value,
      specgenClient.value,
      specgenJsonlib.value,
      specgenSpecFile.value,
      specgenGeneratePath.value,
    )
  }

  override val projectSettings = Seq(
    specgenClient := "sttp",
    specgenJsonlib := "circe",
    specgenSpecFile := file("spec.yaml"),
    specgenGeneratePath := (sourceManaged in Compile).value / "spec",
    specgenClientTask := specgenTask.value,
    sourceGenerators in Compile += specgenClientTask,
    Compile / packageSrc := ((Compile / packageSrc) dependsOn specgenTask).value,
    Compile / packageSrc / mappings ++= contentOf(specgenGeneratePath.value)
  )
}