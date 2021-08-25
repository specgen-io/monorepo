@file:Suppress("UNUSED_VARIABLE")

package io.specgen.gradle

import org.gradle.api.Plugin
import org.gradle.api.Project
import org.gradle.api.plugins.JavaLibraryPlugin
import org.gradle.api.plugins.JavaPluginExtension
import org.gradle.kotlin.dsl.*

public class SpecgenPlugin : Plugin<Project> {
    override fun apply(target: Project): Unit = target.run {
        apply<JavaLibraryPlugin>()
        val generateModelsJava by tasks.creating(SpecgenModelsJavaTask::class)

        afterEvaluate {
            project.configure<JavaPluginExtension> {
                sourceSets.all {
                    java.srcDir(generateModelsJava.outputDirectory.get())
                    tasks[compileJavaTaskName]?.dependsOn(generateModelsJava)
                }
            }
        }
    }

    public companion object {
        public const val SPECGEN_GROUP: String = "specgen"
    }
}
