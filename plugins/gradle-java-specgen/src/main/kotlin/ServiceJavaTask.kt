package io.specgen.gradle.java

import org.gradle.api.Project
import org.gradle.api.provider.Property
import org.gradle.api.tasks.*
import org.gradle.kotlin.dsl.*
import java.io.File
import javax.inject.Inject

public open class ServiceJavaConfig @Inject constructor(project: Project) {
    @OutputDirectory
    public val outputDirectory: Property<File> =
        project.objects.property<File>().convention(project.buildDir.resolve("generated-src/specgen"))

    @Input
    @Optional
    public val jsonlib: Property<String> =
        project.objects.property<String>().convention("jackson")

    @Input
    @Optional
    public val server: Property<String> =
        project.objects.property<String>().convention("spring")

    @InputFile
    @PathSensitive(value = PathSensitivity.RELATIVE)
    public val specFile: Property<File> =
        project.objects.property<File>().convention(project.projectDir.resolve("spec.yaml"))

    @Input
    @Optional
    public val packageName: Property<String> = project.objects.property()

    @OutputDirectory
    @Optional
    public val servicesPath: Property<File> = project.objects.property()

    @OutputFile
    @Optional
    public val swaggerPath: Property<File> = project.objects.property()
}

public open class SpecgenServiceJavaTask public constructor() : SpecgenBaseTask() {
    @TaskAction
    public fun execute() {
        val extension = project.extensions.findByType<SpecgenPluginExtension>()
        // TODO: Check if there are nulls below
        val config = extension!!.configServiceJava!!

        val commandlineArgs = mutableListOf(
            "service-java",
            "--jsonlib", config.jsonlib.get(),
            "--server", config.server.get(),
            "--spec-file", config.specFile.get().absolutePath,
            "--generate-path", config.outputDirectory.get().absolutePath
        )

        if (config.packageName.isPresent) {
            commandlineArgs.add("--package-name")
            commandlineArgs.add(config.packageName.get())
        }
        if (config.servicesPath.isPresent) {
            commandlineArgs.add("--services-path")
            commandlineArgs.add(config.servicesPath.get().absolutePath)
        }
        if (config.swaggerPath.isPresent) {
            commandlineArgs.add("--swagger-path")
            commandlineArgs.add(config.swaggerPath.get().absolutePath)
        }

        runSpecgen(commandlineArgs)
    }
}
