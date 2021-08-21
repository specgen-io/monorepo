@file:Suppress("UNUSED_VARIABLE")

plugins {
    `kotlin-dsl`
    `maven-publish`
    alias(libs.plugins.plugin.publish)
}

group = "io.specgen"
version = System.getProperty("project.version") ?: "0.0.0"
description = "A plugin that integrates specgen tool into the Gradle build process."

kotlin {
    explicitApi()

    sourceSets.all {
        languageSettings.useExperimentalAnnotation("kotlin.ExperimentalStdlibApi")
    }
}

repositories.mavenCentral()

gradlePlugin {
    val specgen by plugins.creating {
        id = "io.specgen.gradle"
        displayName = "Gradle Specgen plugin"
        implementationClass = "io.specgen.gradle.SpecgenPlugin"
    }
}

pluginBundle {
    website = "https://github.com/specgen-io/specgen"
    vcsUrl = "https://github.com/specgen-io/specgen.git"
    description = project.description
    tags = listOf("specgen", "codegeneration", "codegen")
}

publishing.repositories.maven("https://specgen.jfrog.io/artifactory/maven") {
    name = "artifactory"

    credentials {
        username = System.getProperty("jfrog.user")
        password = System.getProperty("jfrog.pass")
    }
}
