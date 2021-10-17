@file:Suppress("UNUSED_VARIABLE")

package io.specgen.gradle

import org.gradle.api.*
import org.gradle.api.model.*
import org.gradle.api.plugins.*
import org.gradle.api.tasks.*
import org.gradle.kotlin.dsl.*

import javax.inject.Inject

public open class SpecgenPluginExtension @Inject constructor(private val objectFactory: ObjectFactory) {
    public var configModelsJava: ModelsJavaConfig? = null
    public var configServiceJavaSpring: ServiceJavaSpringConfig? = null

    @Nested
    public fun modelsJava(action: Action<in ModelsJavaConfig>) {
        val config = objectFactory.newInstance(ModelsJavaConfig::class.java)
        action.execute(config)
        configModelsJava = config
    }

    @Nested
    public fun serviceJavaSpring(action: Action<in ServiceJavaSpringConfig>) {
        val config = objectFactory.newInstance(ServiceJavaSpringConfig::class.java)
        action.execute(config)
        configServiceJavaSpring = config
    }
}

public class SpecgenPlugin : Plugin<Project> {
    override fun apply(project: Project): Unit {
        project.apply<JavaLibraryPlugin>()

        val extension = project.extensions.create<SpecgenPluginExtension>("specgen")

        project.afterEvaluate {
            if (extension.configModelsJava != null) {
                val task by project.tasks.creating(SpecgenModelsJavaTask::class) {
                    config = extension.configModelsJava
                }
                project.configure<JavaPluginExtension> {
                    sourceSets.all {
                        java.srcDir(task.config!!.outputDirectory.get())
                        tasks[compileJavaTaskName]?.dependsOn(task)
                    }
                }
            }
            if (extension.configServiceJavaSpring != null) {
                val task by project.tasks.creating(SpecgenServiceJavaSpringTask::class) {
                    config = extension.configServiceJavaSpring
                }
                project.configure<JavaPluginExtension> {
                    sourceSets.all {
                        java.srcDir(task.config!!.outputDirectory.get())
                        tasks[compileJavaTaskName]?.dependsOn(task)
                    }
                }
            }
        }
    }

    public companion object {
        public const val SPECGEN_GROUP: String = "specgen"
    }
}
