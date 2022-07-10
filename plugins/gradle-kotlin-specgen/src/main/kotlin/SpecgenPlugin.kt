@file:Suppress("UNUSED_VARIABLE")

package io.specgen.kotlin.gradle

import org.gradle.api.*
import org.gradle.api.model.*
import org.gradle.api.plugins.*
import org.gradle.api.tasks.*
import org.gradle.kotlin.dsl.*

import org.jetbrains.kotlin.gradle.dsl.*

import javax.inject.Inject

public open class SpecgenPluginExtension @Inject constructor(private val objectFactory: ObjectFactory) {
    public var configModelsKotlin: ModelsKotlinConfig? = null
    public var configClientKotlin: ClientKotlinConfig? = null
    public var configServiceKotlin: ServiceKotlinConfig? = null

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

    @Nested
    public fun serviceKotlin(action: Action<in ServiceKotlinConfig>) {
        val config = objectFactory.newInstance(ServiceKotlinConfig::class.java)
        action.execute(config)
        configServiceKotlin = config
    }
}

public class SpecgenPlugin : Plugin<Project> {
    override fun apply(project: Project) {
        project.apply<JavaLibraryPlugin>()

        val extension = project.extensions.create<SpecgenPluginExtension>(SPECGEN_EXTENSION)

        val specgenModelsKotlin by project.tasks.registering(SpecgenModelsKotlinTask::class)
        val specgenClientKotlin by project.tasks.registering(SpecgenClientKotlinTask::class)
        val specgenServiceKotlin by project.tasks.registering(SpecgenServiceKotlinTask::class)

        project.afterEvaluate {
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
            extension.configServiceKotlin?.let { config ->
                project.configure<KotlinProjectExtension> {
                    sourceSets.all {
                        kotlin.srcDir(config.outputDirectory.get())
                        tasks["compileKotlin"]?.dependsOn(specgenServiceKotlin)
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
