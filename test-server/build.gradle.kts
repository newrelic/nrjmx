plugins {
    java
}

repositories {
    mavenCentral()
}

dependencies {
    implementation ("com.google.code.gson:gson:2.8.5")
    implementation ("com.sparkjava:spark-core:2.9.1")
    implementation ("org.slf4j:slf4j-simple:1.7.26")
}

tasks.register<Copy>("dockerFiles") {
    group = "Docker"
    description = "Groups files together suitable for building an image"
    into ("${buildDir}/dockerFiles")
    from ("src/docker") {
        include ("**" )
    }
    from (tasks.named("jar")) {
        rename("test-server.*\\.jar","test-server.jar")
    }
}

// No need to run `test` as we are only building the JAR to be used elsewhere
tasks.named<Test>("test") {
    enabled = false
}

tasks.assemble {
    dependsOn("dockerFiles")
}