import org.jetbrains.kotlin.gradle.tasks.KotlinCompile

val kotlinVersion = "1.8.21"

plugins {
    kotlin("jvm") version("1.8.21")
    kotlin("plugin.serialization") version("1.8.21")
    application
}

repositories {
    mavenCentral()
}
   
dependencies {
    val coroutinesVersion = "1.7.3"
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:$coroutinesVersion")

    val serializationVersion = "1.6.3"
    implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:$serializationVersion")

    val okhttpVersion = "4.12.0"
    implementation("com.squareup.okhttp3:okhttp:$okhttpVersion")
}

application {
   mainClass.set("io.webrpc.client.MainKt")
}

tasks.withType<JavaCompile> {
    options.encoding = "UTF-8"
    sourceCompatibility = "11"
    targetCompatibility = "11"
}

tasks.withType<KotlinCompile> {
    kotlinOptions.jvmTarget = "11"
}
