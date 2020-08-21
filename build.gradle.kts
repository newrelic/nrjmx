plugins {
    java
    `maven-publish`
    id("org.beryx.jlink") version ("2.21.2")
    id("org.ysb33r.java.modulehelper") version ("0.9.0")
    id("com.github.sherter.google-java-format") version ("0.8")
}

repositories {
    mavenLocal()
    maven {
        url = uri("https://repo.maven.apache.org/maven2")
    }
}

extraJavaModules {
    module("commons-cli-1.4.jar", "commons.cli", "1.4") {
        exports("org.apache.commons.cli")
    }
    module("gson-2.8.0.jar", "com.google.code.gson", "2.8.0") {
        exports("com.google.gson")
    }
//    module("junit-4.12.jar", "junit4", "4.12") {
//        exports("org.junit")
//        exports("org.junit.rules")
//    }
    module("testcontainers-1.14.3.jar", "org.testcontainers", "1.14.3") {
        exports("org.testcontainers.containers")
        exports("org.testcontainers.images.builder")
        exports("org.testcontainers.shaded.okhttp3")
        requires("org.junit.jupiter.api")
    }
//    module("hamcrest-core-1.3.jar", "hamcrest-core", "0")
    module("tcp-unix-socket-proxy-1.0.2.jar", "tcp.unix.socket.proxy", "0")
    module("duct-tape-1.0.8.jar", "duct.tape", "0")
    module("visible-assertions-2.1.2.jar", "visible.assertions", "0")
    module("junixsocket-native-common-2.0.4.jar", "junixsocket.native", "0")
    module("junixsocket-common-2.0.4.jar", "junixsocket", "0")
    module("native-lib-loader-2.0.2.jar", "native.lib.loader", "0")
}

java {
    modularity.inferModulePath.set(true)
    sourceCompatibility = JavaVersion.VERSION_11
    targetCompatibility = JavaVersion.VERSION_11
}

tasks.named<JavaCompile>("compileJava") {
    options.javaModuleVersion.set(provider({ project.version as String }))
}

tasks.named<Test>("test") {
    useJUnitPlatform()
    systemProperty("TEST_SERVER_DOCKER_FILES", File(project(":test-server").buildDir, "dockerFiles"))
}

dependencies {
    implementation("commons-cli:commons-cli:1.4")
    implementation("com.google.code.gson:gson:2.8.0")
    testImplementation("org.junit.jupiter:junit-jupiter-api:5.6.2")
    testRuntimeOnly("org.junit.jupiter:junit-jupiter-engine")
//    testRuntimeOnly("org.junit.jupiter:junit-vintage-engine:5.6.2")
//    testImplementation("junit:junit:4.12")
    testImplementation(project(":fake-junit4-dependencies"))
    testImplementation("org.testcontainers:testcontainers:1.14.3") {
        exclude("junit")
    }
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
    mainModule.set("org.newrelic.nrjmx")
}

jlink {
//    mergedModule {
//        requires 'java.naming'
//        requires 'java.xml'
//    }
//    launcher{
//        name = 'nrjmx'
//        jvmArgs = ['-Dlogback.configurationFile=./logback.xml']
//    }
}