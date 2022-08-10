@file:Suppress("UNUSED_VARIABLE")

plugins {
    `kotlin-dsl`
    `maven-publish`
    alias(libs.plugins.plugin.publish)
}

group = "io.specgen"
version = System.getProperty("project.version") ?: "0.0.0"
description = "A plugin that integrates specgen Java code generation into the Gradle build process."

kotlin {
    explicitApi()

    sourceSets.all {
        languageSettings.useExperimentalAnnotation("kotlin.ExperimentalStdlibApi")
    }
}

repositories.mavenCentral()

dependencies {
    implementation("org.jetbrains.kotlin:kotlin-gradle-plugin:1.6.0")
}

gradlePlugin {
    val specgen by plugins.creating {
        id = "io.specgen.java.gradle"
        displayName = "Gradle Java Specgen plugin"
        implementationClass = "io.specgen.java.gradle.SpecgenPlugin"
    }
}

pluginBundle {
    website = "https://github.com/specgen-io"
    vcsUrl = "https://github.com/specgen-io/monorepo.git"
    description = project.description
    tags = listOf("specgen", "codegeneration", "codegen", "java")
}

publishing.repositories.maven("https://specgen.jfrog.io/artifactory/maven") {
    name = "artifactory"

    credentials {
        username = System.getProperty("jfrog.user")
        password = System.getProperty("jfrog.pass")
    }
}
