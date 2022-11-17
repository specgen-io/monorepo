plugins {
    id("org.jetbrains.kotlin.jvm") version "{{versions.kotlin.value}}"
    {{#server.spring}}
    id("org.jetbrains.kotlin.plugin.spring") version "{{versions.kotlin_spring.value}}"
    id("org.springframework.boot") version "{{versions.spring_boot.value}}"
    id("io.spring.dependency-management") version "{{versions.spring_dependency.value}}"
    {{/server.spring}}
    {{#server.micronaut}}
    id("org.jetbrains.kotlin.kapt") version "{{versions.kotlin.value}}"
    id("org.jetbrains.kotlin.plugin.allopen") version "{{versions.kotlin.value}}"
    id("io.micronaut.application") version "{{versions.micronaut_application.value}}"
    {{/server.micronaut}}
    id("io.specgen.kotlin.gradle") version "{{versions.specgen.value}}"
}

group = "{{groupid.value}}"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
}

dependencies {
    implementation("org.jetbrains.kotlin:kotlin-reflect:{{versions.kotlin.value}}")
    implementation("org.jetbrains.kotlin:kotlin-stdlib-jdk8:{{versions.kotlin.value}}")
    {{#server.spring}}
    implementation("org.springframework.boot:spring-boot-starter")
    implementation("org.springframework.boot:spring-boot-starter-web")
    {{#swagger.value}}
    implementation("io.springfox:springfox-boot-starter:{{versions.springfox.value}}")
    {{/swagger.value}}
    {{/server.spring}}
    {{#server.micronaut}}
    kapt("io.micronaut:micronaut-http-validation")
    {{#jsonlib.jackson}}
    kapt("io.micronaut.serde:micronaut-serde-processor")
    implementation("io.micronaut:micronaut-jackson-databind")
    implementation("io.micronaut.serde:micronaut-serde-jackson")
    {{/jsonlib.jackson}}
    implementation("io.micronaut:micronaut-runtime")
    implementation("io.micronaut:micronaut-validation")
    implementation("io.micronaut.kotlin:micronaut-kotlin-runtime")
    implementation("jakarta.annotation:jakarta.annotation-api")
    runtimeOnly("org.slf4j:slf4j-simple")
    {{/server.micronaut}}
    {{#jsonlib.jackson}}
    implementation("com.fasterxml.jackson.module:jackson-module-kotlin")
    {{/jsonlib.jackson}}
    {{#jsonlib.moshi}}
    implementation("com.squareup.moshi:moshi:{{versions.moshi.value}}")
    implementation("com.squareup.moshi:moshi-adapters:{{versions.moshi.value}}")
    implementation("com.squareup.moshi:moshi-kotlin:{{versions.moshi.value}}")
    {{/jsonlib.moshi}}
}

{{#server.micronaut}}
application {
    mainClass.set("{{package.value}}.{{mainclass.value}}Kt")
}
{{/server.micronaut}}

java {
    sourceCompatibility = JavaVersion.toVersion("{{versions.java.value}}")
}

tasks {
    compileKotlin {
        kotlinOptions {
            jvmTarget = "{{versions.java.value}}"
        }
    }
}

specgen {
    serviceKotlin {
        packageName.set("{{package.value}}")
        jsonlib.set("{{jsonlib.value}}")
        server.set("{{server.value}}")
        specFile.set(file("spec.yaml"))
        servicesPath.set(file("src/main/kotlin"))
        {{#swagger.value}}
        swaggerPath.set(file("src/main/resources/static/docs/swagger.yaml"))
        {{/swagger.value}}
    }
}

{{#server.micronaut}}
graalvmNative.toolchainDetection.set(false)
micronaut {
    runtime("netty")
    processing {
        incremental(true)
        annotations("{{package.value}}.*")
    }
}
{{/server.micronaut}}
