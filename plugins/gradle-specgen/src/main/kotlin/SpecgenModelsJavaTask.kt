package io.specgen.gradle

import org.gradle.api.Project
import org.gradle.api.provider.Property
import org.gradle.api.tasks.*
import org.gradle.kotlin.dsl.property
import java.io.File
import javax.inject.Inject

public open class ModelsJavaConfig @Inject constructor(project: Project) {
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
}

@CacheableTask
public open class SpecgenModelsJavaTask public constructor() : SpecgenBaseTask() {
    @Internal
    public var config: ModelsJavaConfig? = null

    @TaskAction
    public fun execute() {
        val config = this.config!!

        val commandlineArgs = mutableListOf(
            "models-java",
            "--spec-file",
            config.specFile.get().absolutePath,
            "--generate-path",
            config.outputDirectory.get().absolutePath,
        )

        if (config.packageName.isPresent) {
            commandlineArgs.add("--package-name")
            commandlineArgs.add(config.packageName.get())
        }

        runSpecgen(commandlineArgs)
    }
}
