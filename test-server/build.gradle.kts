import org.gradle.jvm.tasks.Jar

plugins {
    id("org.jetbrains.kotlin.jvm").version("1.3.21")

    application
}


repositories {
    jcenter()
}

dependencies {
    implementation("com.sparkjava:spark-core:2.9.1")
    implementation("com.google.code.gson:gson:2.8.5")
    implementation("org.slf4j:slf4j-simple:1.7.26")

    // Use the Kotlin JDK 8 standard library.
    implementation("org.jetbrains.kotlin:kotlin-stdlib-jdk8")

    // Use the Kotlin test library.
    testImplementation("org.jetbrains.kotlin:kotlin-test")

    // Use the Kotlin JUnit integration.
    testImplementation("org.jetbrains.kotlin:kotlin-test-junit")
}

application {
    // Define the main class for the application.
    mainClassName = "org.newrelic.jmx.AppKt"
}

val fatJar = task("fatJar", type = Jar::class) {
    baseName = "${project.name}-fat"
    manifest {
        attributes["Implementation-Title"] = "Test JMX service"
        attributes["Implementation-Version"] = "1.0.0"
        attributes["Main-Class"] = application.mainClassName
    }
    from(configurations.runtimeClasspath.get().map { if (it.isDirectory) it else zipTree(it) })
    with(tasks.jar.get() as CopySpec)
}

tasks {
    "build" {
        dependsOn(fatJar)
    }
}
