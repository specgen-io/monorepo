plugins {
    id("org.jetbrains.kotlin.jvm") version "{{versions.kotlin.value}}"
    {{#client.micronaut}}
    id("org.jetbrains.kotlin.kapt") version "{{versions.kotlin.value}}"
    id("org.jetbrains.kotlin.plugin.allopen") version "{{versions.kotlin.value}}"
    id("io.micronaut.application") version "{{versions.micronaut_application.value}}"
    {{/client.micronaut}}
    id("io.specgen.kotlin.gradle") version "{{versions.specgen.value}}"
}

group = "{{groupid.value}}"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
}

dependencies {
    {{#client.okhttp}}
    implementation("com.squareup.okhttp3:okhttp:{{versions.okhttp.value}}")
    {{/client.okhttp}}
    {{#client.micronaut}}
    kapt("io.micronaut:micronaut-http-validation")
    implementation("io.micronaut:micronaut-http-client")
    {{#jsonlib.jackson}}
    implementation("io.micronaut:micronaut-jackson-databind")
    {{/jsonlib.jackson}}
    implementation("io.micronaut:micronaut-runtime")
    implementation("io.micronaut:micronaut-validation")
    implementation("io.micronaut.kotlin:micronaut-kotlin-runtime")
    implementation("jakarta.annotation:jakarta.annotation-api")
    {{/client.micronaut}}
    implementation("org.jetbrains.kotlin:kotlin-reflect")
    implementation("org.jetbrains.kotlin:kotlin-stdlib-jdk8")
    {{#jsonlib.jackson}}
    implementation("com.fasterxml.jackson.module:jackson-module-kotlin:{{versions.jackson.value}}")
    implementation("com.fasterxml.jackson.datatype:jackson-datatype-jsr310:{{versions.jackson.value}}")
    {{/jsonlib.jackson}}
    {{#jsonlib.moshi}}
    implementation("com.squareup.moshi:moshi:{{versions.moshi.value}}")
    implementation("com.squareup.moshi:moshi-adapters:{{versions.moshi.value}}")
    implementation("com.squareup.moshi:moshi-kotlin:{{versions.moshi.value}}")
    {{/jsonlib.moshi}}
    implementation("org.apache.logging.log4j:log4j-slf4j-impl:{{versions.log4j.value}}")
}

specgen {
    clientKotlin {
        jsonlib.set("{{jsonlib.value}}")
        client.set("{{client.value}}")
        packageName.set("{{package.value}}")
        specFile.set(file("spec.yaml"))
    }
}
