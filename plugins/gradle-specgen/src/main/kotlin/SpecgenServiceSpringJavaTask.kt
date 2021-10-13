package io.specgen.gradle

import org.gradle.api.provider.Property
import org.gradle.api.tasks.*
import org.gradle.kotlin.dsl.property
import java.io.File

@CacheableTask
public open class SpecgenServiceSpringJavaTask public constructor() : SpecgenBaseTask() {
    @OutputDirectory
    public val outputDirectory: Property<File> =
        project.objects.property<File>().convention(project.buildDir.resolve("generated-src/specgen"))

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

    @TaskAction
    public fun execute() {
        val commandlineArgs = mutableListOf(
            "service-java-spring",
            "--spec-file",
            specFile.get().absolutePath,
            "--generate-path",
            outputDirectory.get().absolutePath
        )

        if (packageName.isPresent) {
            commandlineArgs.add("--package-name")
            commandlineArgs.add(packageName.get())
        }
        if (servicesPath.isPresent) {
            commandlineArgs.add("--services-path")
            commandlineArgs.add(servicesPath.get().absolutePath)
        }
        if (swaggerPath.isPresent) {
            commandlineArgs.add("--swagger-path")
            commandlineArgs.add(swaggerPath.get().absolutePath)
        }

        runSpecgen(commandlineArgs)
    }
}
