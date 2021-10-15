@file:Suppress("UNUSED_VARIABLE")

package io.specgen.gradle

import org.gradle.api.provider.Property
import org.gradle.api.Plugin
import org.gradle.api.Project
import org.gradle.api.plugins.JavaLibraryPlugin
import org.gradle.api.plugins.JavaPluginExtension
import org.gradle.api.tasks.*
import org.gradle.kotlin.dsl.*
import java.io.File

public abstract class ModelsJavaConfig {
    public abstract val outputDirectory: Property<File>
    public abstract val specFile: Property<File>
    public abstract val packageName: Property<String>
}

public abstract class SpecgenPluginExtension {
    public abstract val modelsJava: Property<ModelsJavaConfig>
}

public class SpecgenPlugin : Plugin<Project> {
    override fun apply(target: Project): Unit = target.run {
        apply<JavaLibraryPlugin>()

        val extension = project.extensions.create<SpecgenPluginExtension>("specgen")

        val generateModelsJava by tasks.registering(SpecgenModelsJavaTask::class)
        val generateServiceSpringJava by tasks.registering(SpecgenServiceSpringJavaTask::class)

        afterEvaluate {
            if (extension.modelsJava.isPresent) {
                project.configure<JavaPluginExtension> {
                    sourceSets.all {
                        java.srcDir(generateModelsJava.get().outputDirectory.get())
                        tasks[compileJavaTaskName]?.dependsOn(generateModelsJava)
                    }
                }
            }
        }
    }

    public companion object {
        public const val SPECGEN_GROUP: String = "specgen"
    }
}
