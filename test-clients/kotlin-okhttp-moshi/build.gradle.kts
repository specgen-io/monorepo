plugins {
    kotlin("jvm") version "1.7.20"
    id("com.adarshr.test-logger") version "3.2.0"
    id("io.specgen.kotlin.gradle")
}

group = "io.specgen"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
}

dependencies {
    implementation(kotlin("stdlib"))
    implementation("com.squareup.moshi:moshi:1.14.0")
    implementation("com.squareup.moshi:moshi-adapters:1.14.0")
    implementation("com.squareup.moshi:moshi-kotlin:1.14.0")
    implementation("com.squareup.okhttp3:okhttp:4.10.0")
    implementation("org.apache.logging.log4j:log4j-slf4j-impl:2.19.0")

    testImplementation(kotlin("test"))
    testImplementation("org.assertj:assertj-core:3.23.1")
}

specgen {
    clientKotlin {
        jsonlib.set("moshi")
        client.set("okhttp")
        packageName.set("test_client")
        specFile.set(file("../spec.yaml"))
    }
}

tasks.test {
    useJUnitPlatform()
}