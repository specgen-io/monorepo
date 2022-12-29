plugins {
    id("java")
    id("io.specgen.java.gradle") version "{{versions.specgen.value}}"
}

group = "{{groupid.value}}"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
}

dependencies {
    {{#jsonlib.jackson}}
    implementation("com.fasterxml.jackson.core:jackson-core:{{versions.jackson.value}}")
    implementation("com.fasterxml.jackson.core:jackson-annotations:{{versions.jackson.value}}")
    implementation("com.fasterxml.jackson.core:jackson-databind:{{versions.jackson.value}}")
    implementation("com.fasterxml.jackson.datatype:jackson-datatype-jsr310:{{versions.jackson.value}}")
    {{/jsonlib.jackson}}
    {{#jsonlib.moshi}}
    implementation("com.squareup.moshi:moshi:{{versions.moshi.value}}")
    implementation("com.squareup.moshi:moshi-adapters:{{versions.moshi.value}}")
    implementation("com.google.code.findbugs:jsr305:{{versions.jsr305.value}}")
    {{/jsonlib.moshi}}
}

specgen {
    modelsJava {
        jsonlib.set("{{jsonlib.value}}")
        packageName.set("{{package.value}}")
        specFile.set(file("spec.yaml"))
    }
}
