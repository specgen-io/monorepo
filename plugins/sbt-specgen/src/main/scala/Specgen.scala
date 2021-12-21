package spec

import java.io.{ByteArrayOutputStream, PrintWriter}
import java.nio.file.{Files, Paths, StandardCopyOption}
import sys.process._

import sbt.{File, Logger}

object Specgen {
  def serviceScalaPlay(
    log: Logger,
    specPath: File,
    swaggerPath: File,
    generatePath: File,
    servicesPath: File
  ): Seq[File] = {

    log.info(s"Running sbt-spec code generation plugin")

    val specToolPath: String = getSpecTool(log)

    val specCommand: Seq[String] = Seq(
      specToolPath,
      "service-scala-play",
      "--spec-file", specPath.getPath,
      "--swagger-path", swaggerPath.getPath,
      "--generate-path", generatePath.getPath,
      "--services-path", servicesPath.getPath
    )

    runSpecgen(log, specCommand)

    val generatedFiles = recursiveListFiles(generatePath)
    generatedFiles.toSeq
  }

  def modelsScalaCirce(
    log: Logger,
    specPath: File,
    generatePath: File,
  ): Seq[File] = {

    log.info(s"Running sbt-spec code generation plugin")

    val specToolPath: String = getSpecTool(log)

    val specCommand: Seq[String] = Seq(
      specToolPath,
      "models-scala-circe",
      "--spec-file", specPath.getPath,
      "--generate-path", generatePath.getPath,
    )

    runSpecgen(log, specCommand)

    val generatedFiles = recursiveListFiles(generatePath)
    generatedFiles.toSeq
  }

  def clientSttp(
    log: Logger,
    specPath: File,
    generatePath: File,
  ): Seq[File] = {

    log.info(s"Running sbt-spec code generation plugin")

    val specToolPath: String = getSpecTool(log)

    val specCommand: Seq[String] = Seq(
      specToolPath,
      "client-scala-sttp",
      "--spec-file", specPath.getPath,
      "--generate-path", generatePath.getPath,
    )

    runSpecgen(log, specCommand)

    val generatedFiles = recursiveListFiles(generatePath)
    generatedFiles.toSeq
  }

  def recursiveListFiles(path: File): Array[File] = {
    val current = path.listFiles
    current.filterNot(_.isDirectory) ++ current.filter(_.isDirectory).flatMap(recursiveListFiles)
  }

  def runSpecgen(log: Logger, specgenCommand: Seq[String]) = {
    log.info("Running specgen tool")
    log.info(specgenCommand.mkString(" "))

    val (status: Int, stdout: String, stderr: String) = runCommand(specgenCommand)

    log.info(stdout)
    log.error(stderr)

    if (status != 0) {
      throw new Exception(s"Failed to run specgen tool, exit code: $status")
    }
  }

  def runCommand(cmd: Seq[String]): (Int, String, String) = {
    val stdoutStream = new ByteArrayOutputStream
    val stderrStream = new ByteArrayOutputStream
    val stdoutWriter = new PrintWriter(stdoutStream)
    val stderrWriter = new PrintWriter(stderrStream)
    val exitValue = cmd.!(ProcessLogger(stdoutWriter.println, stderrWriter.println))
    stdoutWriter.close()
    stderrWriter.close()
    (exitValue, stdoutStream.toString, stderrStream.toString)
  }

  def getSpecTool(log: Logger): String = {
    val osname = getOsName()
    val arch = getArch()

    val specToolPath = s"/dist/${osname}_${arch}/${getExeName("specgen")}"

    val specToolStream = getClass.getResourceAsStream(specToolPath)
    val jarPath = SpecPlay.getClass.getProtectionDomain.getCodeSource.getLocation.getPath
    val specPath = jarPath.substring(0, jarPath.lastIndexOf('.'))+specToolPath

    val specPathFile = new File(specPath)
    if (!specPathFile.exists()) {
      log.info(s"Unpacking specgen tool into: ${specPathFile.getPath}")
      val specPathParent = specPathFile.getParentFile
      if (!specPathParent.exists()) {
        specPathParent.mkdirs()
      }
      Files.copy(specToolStream, Paths.get(specPathFile.getPath), StandardCopyOption.REPLACE_EXISTING)
      specPathFile.setExecutable(true)
    }

    specPathFile.getPath
  }

  def getExeName(toolName: String): String =
    if (getOsName() == "windows")
      s"$toolName.exe"
    else
      toolName

  def getOsName(): String = {
    val osName = System.getProperty("os.name").toLowerCase()

    osName match {
      case x if x.contains("win") => "windows"
      case x if x.contains("mac") => "darwin"
      case x if x.contains("nix") || x.contains("nux") => "linux"
    }
  }

  def getArch(): String = {
    val arch = System.getProperty("os.arch")

    arch match {
      case x if x.contains("64") => "amd64"
      case _ => "x86"
    }
  }
}