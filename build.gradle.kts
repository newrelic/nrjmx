// TODO: Check generated POM - should be the same as for maven build
// TODO: Adapt maven-assembly-plugin
// TODO: Adapt copy-rename-maven-plugin
// TODO: Build the testserver (dockerfile-maven-plugin)
// TODO: tarball profile equivalent
// TODO: deb
// TODO: rpm

plugins {
    java
    `maven-publish`
    id("org.beryx.jlink") version("2.21.2")
    id( "org.ysb33r.java.modulehelper") version("1.0.0-SNAPSHOT")
}

repositories {
    mavenLocal()
    maven {
        url = uri("https://repo.maven.apache.org/maven2")
    }
}

extraJavaModules {
    module("commons-cli-1.4.jar","commons.cli","1.4") {
        exports("org.apache.commons.cli")
    }
    module("gson-2.8.0.jar","com.google.code.gson","2.8.0") {
        exports("com.google.gson")
    }
}

java {
    modularity.inferModulePath.set(true)
    sourceCompatibility = JavaVersion.VERSION_11
    targetCompatibility = JavaVersion.VERSION_11
}

tasks.named<JavaCompile>("compileJava") {
    options.javaModuleVersion.set(provider({ project.version as String }))
}

dependencies {
    implementation("commons-cli:commons-cli:1.4")
    implementation("com.google.code.gson:gson:2.8.0")
    testImplementation("junit:junit:4.12")
    testImplementation("org.testcontainers:testcontainers:1.11.3")
}

publishing {
    publications {
        create<MavenPublication>("nrjmx") {
            from(components["java"])
        }
    }
}

application {
    mainClass.set("org.newrelic.nrjmx.Application")
//    mainModule.set("org.newrelic.njrmx")
}

