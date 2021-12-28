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
    public var configClientJava: ClientJavaConfig? = null
    public var configServiceJava: ServiceJavaConfig? = null
    public var configModelsKotlin: ModelsKotlinConfig? = null
    public var configClientKotlin: ClientKotlinConfig? = null

    @Nested
    public fun modelsJava(action: Action<in ModelsJavaConfig>) {
        val config = objectFactory.newInstance(ModelsJavaConfig::class.java)
        action.execute(config)
        configModelsJava = config
    }

    @Nested
    public fun clientJava(action: Action<in ClientJavaConfig>) {
        val config = objectFactory.newInstance(ClientJavaConfig::class.java)
        action.execute(config)
        configClientJava = config
    }

    @Nested
    public fun serviceJava(action: Action<in ServiceJavaConfig>) {
        val config = objectFactory.newInstance(ServiceJavaConfig::class.java)
        action.execute(config)
        configServiceJava = config
    }

    @Nested
    public fun modelsKotlin(action: Action<in ModelsKotlinConfig>) {
        val config = objectFactory.newInstance(ModelsKotlinConfig::class.java)
        action.execute(config)
        configModelsKotlin = config
    }

    @Nested
    public fun clientKotlin(action: Action<in ClientKotlinConfig>) {
        val config = objectFactory.newInstance(ClientKotlinConfig::class.java)
        action.execute(config)
        configClientKotlin = config
    }
}

public class SpecgenPlugin : Plugin<Project> {
    override fun apply(project: Project) {
        project.apply<JavaLibraryPlugin>()

        val extension = project.extensions.create<SpecgenPluginExtension>(SPECGEN_EXTENSION)

        val specgenModelsJava by project.tasks.registering(SpecgenModelsJavaTask::class)
        val specgenClientJava by project.tasks.registering(SpecgenClientJavaTask::class)
        val specgenServiceJava by project.tasks.registering(SpecgenServiceJavaTask::class)
        val specgenModelsKotlin by project.tasks.registering(SpecgenModelsKotlinTask::class)
        val specgenClientKotlin by project.tasks.registering(SpecgenClientKotlinTask::class)

        project.afterEvaluate {
            extension.configModelsJava?.let { config ->
                project.configure<JavaPluginExtension> {
                    sourceSets.all {
                        java.srcDir(config.outputDirectory.get())
                        tasks[compileJavaTaskName]?.dependsOn(specgenModelsJava)
                    }
                }
            }
            extension.configClientJava?.let { config ->
                project.configure<JavaPluginExtension> {
                    sourceSets.all {
                        java.srcDir(config.outputDirectory.get())
                        tasks[compileJavaTaskName]?.dependsOn(specgenClientJava)
                    }
                }
            }
            extension.configServiceJava?.let { config ->
                project.configure<JavaPluginExtension> {
                    sourceSets.all {
                        java.srcDir(config.outputDirectory.get())
                        tasks[compileJavaTaskName]?.dependsOn(specgenServiceJava)
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
            extension.configClientKotlin?.let { config ->
                project.configure<KotlinProjectExtension> {
                    sourceSets.all {
                        kotlin.srcDir(config.outputDirectory.get())
                        tasks["compileKotlin"]?.dependsOn(specgenClientKotlin)
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
