package spec

import java.io.{ByteArrayOutputStream, PrintWriter}
import java.nio.file.{Files, Paths, StandardCopyOption}
import sys.process._

import sbt.{File, Logger}

object Specgen {
  def serviceScala(
    log: Logger,
    server: String,
    jsonlib: String,
    specFile: File,
    swaggerPath: File,
    generatePath: File,
    servicesPath: File,
  ): Seq[File] = {

    log.info(s"Running sbt-specgen service code generation plugin")

    val specToolPath: String = getSpecgenTool(log)

    val specCommand: Seq[String] = Seq(
      specToolPath,
      "service-scala",
      "--spec-file", specFile.getPath,
      "--swagger-path", swaggerPath.getPath,
      "--generate-path", generatePath.getPath,
      "--services-path", servicesPath.getPath
    )

    runSpecgen(log, specCommand)

    val generatedFiles = recursiveListFiles(generatePath)
    generatedFiles.toSeq
  }

  def modelsScala(
    log: Logger,
    jsonlib: String,
    specFile: File,
    generatePath: File,
  ): Seq[File] = {

    log.info(s"Running sbt-specgen models code generation plugin")

    val specToolPath: String = getSpecgenTool(log)

    val specCommand: Seq[String] = Seq(
      specToolPath,
      "models-scala",
      "--spec-file", specFile.getPath,
      "--generate-path", generatePath.getPath,
    )

    runSpecgen(log, specCommand)

    val generatedFiles = recursiveListFiles(generatePath)
    generatedFiles.toSeq
  }

  def clientScala(
    log: Logger,
    client: String,
    jsonlib: String,
    specFile: File,
    generatePath: File,
  ): Seq[File] = {

    log.info(s"Running sbt-specgen client code generation plugin")

    val specToolPath: String = getSpecgenTool(log)

    val specCommand: Seq[String] = Seq(
      specToolPath,
      "client-scala",
      "--spec-file", specFile.getPath,
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

  def getSpecgenTool(log: Logger): String = {
    val osname = getOsName()
    val arch = getArch()

    val specgenRelativePath = s"/dist/${osname}_${arch}/${getExeName("specgen")}"

    val specToolStream = getClass.getResourceAsStream(specgenRelativePath)
    val jarPath = Specgen.getClass.getProtectionDomain.getCodeSource.getLocation.getPath
    val specgenPath = jarPath.substring(0, jarPath.lastIndexOf('.'))+specgenRelativePath

    val specgenFile = new File(specgenPath)
    if (!specgenFile.exists()) {
      log.info(s"Unpacking specgen tool into: ${specgenFile.getPath}")
      val specPathParent = specgenFile.getParentFile
      if (!specPathParent.exists()) {
        specPathParent.mkdirs()
      }
      Files.copy(specToolStream, Paths.get(specgenFile.getPath), StandardCopyOption.REPLACE_EXISTING)
      specgenFile.setExecutable(true)
    }

    specgenFile.getPath
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