package io.specgen.gradle

import org.gradle.api.provider.Property
import org.gradle.api.tasks.*
import org.gradle.kotlin.dsl.property
import java.io.File


@CacheableTask
public open class SpecgenModelsJavaTask public constructor() : SpecgenBaseTask() {
    @OutputDirectory
    public val outputDirectory: Property<File> =
        project.objects.property<File>().convention(project.buildDir.resolve("generated-src/specgen"))

    @InputFile
    @PathSensitive(value = PathSensitivity.RELATIVE)
    public val specFile: Property<File> =
        project.objects.property<File>().convention(project.projectDir.resolve("spec.yaml"))

    @Input
    @Optional
    public val packageName: Property<String?> = project.objects.property()

    @TaskAction
    public fun execute() {
        val args = mutableListOf(
            "models-java",
            "--spec-file",
            specFile.get().absolutePath,
            "--generate-path",
            outputDirectory.get().absolutePath,
        )

        if (packageName.get() != null) {
            args.add("--package-name")
            args.add(packageName.get())
        }

        runSpecgen(args)
    }
}
