plugins {
    kotlin("jvm") version "1.6.10"
    id("com.adarshr.test-logger") version "3.1.0"
    id("io.specgen.kotlin.gradle")
}

group = "io.specgen"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
}

dependencies {
    implementation(kotlin("stdlib"))
    implementation("com.squareup.moshi:moshi:1.13.0")
    implementation("com.squareup.moshi:moshi-adapters:1.13.0")
    implementation("com.squareup.moshi:moshi-kotlin:1.13.0")
    testImplementation(kotlin("test"))
}

specgen {
    modelsKotlin {
        jsonlib.set("moshi")
        specFile.set(file("../spec.yaml"))
    }
}

tasks.test {
    useJUnitPlatform()
}