package io.specgen.gradle

import org.gradle.api.provider.Property
import org.gradle.api.tasks.*
import org.gradle.kotlin.dsl.property
import java.io.File


@CacheableTask
public open class GenerateJavaModelsTask public constructor() : SpecgenRunTask() {
    @OutputDirectory
    public val outputDirectory: Property<File> =
        project.objects.property<File>().convention(project.buildDir.resolve("generated-src/specgen"))

    @InputFile
    @PathSensitive(value = PathSensitivity.RELATIVE)
    public val specFile: Property<File> =
        project.objects.property<File>().convention(project.projectDir.resolve("spec.yaml"))

    @TaskAction
    public fun execute(): Unit = runSpecgen(
        "models-java",
        "--spec-file",
        specFile.get().absolutePath,
        "--generate-path",
        outputDirectory.get().absolutePath,
    )
}
