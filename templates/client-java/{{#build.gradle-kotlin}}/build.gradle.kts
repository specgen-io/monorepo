plugins {
    id("java")
    {{#client.micronaut}}
    id("io.micronaut.application") version "{{versions.micronaut_application.value}}"
    {{/client.micronaut}}
    id("io.specgen.java.gradle") version "{{versions.specgen.value}}"
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
    implementation("io.micronaut:micronaut-http-client")
    {{#jsonlib.jackson}}
    implementation("io.micronaut:micronaut-jackson-databind")
    {{/jsonlib.jackson}}
    implementation("io.micronaut:micronaut-runtime")
    implementation("io.micronaut:micronaut-validation")
    implementation("jakarta.annotation:jakarta.annotation-api")
    {{/client.micronaut}}
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
    implementation("org.apache.logging.log4j:log4j-slf4j-impl:{{versions.log4j.value}}")
}

specgen {
    clientJava {
        jsonlib.set("{{jsonlib.value}}")
        client.set("{{client.value}}")
        packageName.set("{{package.value}}")
        specFile.set(file("spec.yaml"))
    }
}
