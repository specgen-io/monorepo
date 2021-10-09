package io.specgen.gradle

import org.gradle.api.DefaultTask
import org.gradle.work.DisableCachingByDefault
import java.io.File
import java.io.IOException
import java.nio.file.Files
import java.nio.file.StandardCopyOption

private class JarHandle private constructor()

@DisableCachingByDefault(because = "Abstract super-class, not to be instantiated directly")
public abstract class SpecgenBaseTask : DefaultTask() {
    init {
        group = SpecgenPlugin.SPECGEN_GROUP
    }

    protected fun runSpecgen(args: List<String>) {
        try {
            project.exec {
                val jarPath = JarHandle::class.java.protectionDomain.codeSource.location.path
                val osName: String = getOsName()
                val archName: String = getArchName()
                val exeName: String = getExeName()
                val specgenRelativePath = "/dist/${osName}_$archName/$exeName"
                val specgenPath = jarPath.substring(0, jarPath.lastIndexOf('.')) + specgenRelativePath
                val executable = File(specgenPath)

                if (!executable.exists()) {
                    if (logger.isDebugEnabled) logger.debug("Unpacking specgen tool into: ${executable.path}.")

                    executable.parentFile.let {
                        if (!it.exists()) it.mkdirs()
                    }

                    try {
                        checkNotNull(this@SpecgenBaseTask.javaClass.getResourceAsStream(specgenRelativePath)) { "specgen executable isn't present in JAR." }
                            .use { specgenToolStream ->
                                Files.copy(
                                    specgenToolStream,
                                    executable.toPath(),
                                    StandardCopyOption.REPLACE_EXISTING,
                                )
                            }
                    } catch (e: IOException) {
                        throw IllegalStateException("Failed to copy specgen file.", e)
                    }

                    executable.setExecutable(true)
                }
                executable(executable)
                args(args)
            }.assertNormalExitValue()
        } catch (error: Exception) {
            throw SpecgenException("Source generation failed.", error)
        }
    }

    private fun getOsName(): String {
        val osName = System.getProperty("os.name").toLowerCase()

        return when {
            "win" in osName -> "windows"
            "mac" in osName -> "darwin"
            "nix" in osName || "nux" in osName -> "linux"
            else -> ""
        }
    }

    private fun getArchName(): String {
        val archName = System.getProperty("os.arch")

        return if ("64" in archName) {
            "amd64"
        } else {
            "x86"
        }
    }

    private fun getExeName(): String = if (getOsName() == "windows") {
        "specgen.exe"
    } else {
        "specgen"
    }
}
