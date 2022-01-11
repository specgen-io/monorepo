package io.specgen.gradle

import org.gradle.api.Project
import org.gradle.api.provider.Property
import org.gradle.api.tasks.*
import org.gradle.kotlin.dsl.*
import java.io.File
import javax.inject.Inject

public open class ClientJavaConfig @Inject constructor(project: Project) {
    @OutputDirectory
    public val outputDirectory: Property<File> =
        project.objects.property<File>().convention(project.buildDir.resolve("generated-src/specgen"))

    @Input
    @Optional
    public val jsonlib: Property<String> = project.objects.property()

    @InputFile
    @PathSensitive(value = PathSensitivity.RELATIVE)
    public val specFile: Property<File> =
        project.objects.property<File>().convention(project.projectDir.resolve("spec.yaml"))

    @Input
    @Optional
    public val packageName: Property<String> = project.objects.property()
}

public open class SpecgenClientJavaTask public constructor() : SpecgenBaseTask() {
    @TaskAction
    public fun execute() {
        val extension = project.extensions.findByType<SpecgenPluginExtension>()
        // TODO: Check if there are nulls below
        val config = extension!!.configClientJava!!

        val commandlineArgs = mutableListOf(
            "client-java",
            "--jsonlib", config.jsonlib.get(),
            "--spec-file", config.specFile.get().absolutePath,
            "--generate-path", config.outputDirectory.get().absolutePath,
        )

        if (config.packageName.isPresent) {
            commandlineArgs.add("--package-name")
            commandlineArgs.add(config.packageName.get())
        }

        runSpecgen(commandlineArgs)
    }
}
