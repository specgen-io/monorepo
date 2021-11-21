@file:Suppress("UNUSED_VARIABLE")

package io.specgen.gradle

import org.gradle.api.*
import org.gradle.api.model.*
import org.gradle.api.plugins.*
import org.gradle.api.tasks.*
import org.gradle.kotlin.dsl.*

import org.jetbrains.kotlin.gradle.dsl.*

import javax.inject.Inject

public open class SpecgenPluginExtension @Inject constructor(private val objectFactory: ObjectFactory) {
    public var configModelsJava: ModelsJavaConfig? = null
    public var configServiceJavaSpring: ServiceJavaSpringConfig? = null
    public var configModelsKotlin: ModelsKotlinConfig? = null

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

    @Nested
    public fun modelsKotlin(action: Action<in ModelsKotlinConfig>) {
        val config = objectFactory.newInstance(ModelsKotlinConfig::class.java)
        action.execute(config)
        configModelsKotlin = config
    }

}

public class SpecgenPlugin : Plugin<Project> {
    override fun apply(project: Project): Unit {
        project.apply<JavaLibraryPlugin>()

        val extension = project.extensions.create<SpecgenPluginExtension>(SPECGEN_EXTENSION)

        val specgenModelsJava by project.tasks.registering(SpecgenModelsJavaTask::class)
        val specgenServiceJavaSpring by project.tasks.registering(SpecgenServiceJavaSpringTask::class)
        val specgenModelsKotlin by project.tasks.registering(SpecgenModelsKotlinTask::class)

        project.afterEvaluate {
            extension.configModelsJava?.let { config ->
                project.configure<JavaPluginExtension> {
                    sourceSets.all {
                        java.srcDir(config.outputDirectory.get())
                        tasks[compileJavaTaskName]?.dependsOn(specgenModelsJava)
                    }
                }
            }
            extension.configServiceJavaSpring?.let { config ->
                project.configure<JavaPluginExtension> {
                    sourceSets.all {
                        java.srcDir(config.outputDirectory.get())
                        tasks[compileJavaTaskName]?.dependsOn(specgenServiceJavaSpring)
                    }
                }
            }

            extension.configModelsKotlin?.let { config ->
                project.configure<KotlinProjectExtension> {
                    sourceSets.all {
                        kotlin.srcDir(config.outputDirectory.get())
                        tasks["compileKotlin"]?.dependsOn(specgenModelsKotlin)
                    }
                }
            }

        }
    }

    public companion object {
        public const val SPECGEN_GROUP: String = "specgen"
        public const val SPECGEN_EXTENSION: String = "specgen"
    }
}
