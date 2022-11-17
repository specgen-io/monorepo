plugins {
    id("org.jetbrains.kotlin.jvm") version "{{versions.kotlin.value}}"
    id("io.specgen.kotlin.gradle") version "{{versions.specgen.value}}"
}

group = "{{groupid.value}}"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
}

dependencies {
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
    implementation("com.squareup.okhttp3:okhttp:{{versions.okhttp.value}}")
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
