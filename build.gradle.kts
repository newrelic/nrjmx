plugins {
    java
    `maven-publish`
    id("org.beryx.jlink") version ("2.21.2")
    id("org.ysb33r.java.modulehelper") version ("0.9.0")
    id("com.github.sherter.google-java-format") version ("0.8")
    id("nebula.ospackage") version ("8.4.1")
}

repositories {
    mavenLocal()
    maven {
        url = uri("https://repo.maven.apache.org/maven2")
    }
}

dependencies {
    implementation("commons-cli:commons-cli:1.4")
    implementation("com.google.code.gson:gson:2.8.0")
    testImplementation("org.junit.jupiter:junit-jupiter-api:5.6.2")
    testRuntimeOnly("org.junit.jupiter:junit-jupiter-engine")
    testImplementation(project(":fake-junit4-dependencies"))
    testImplementation("org.testcontainers:testcontainers:1.14.3") {
        exclude("junit")
    }
}

extraJavaModules {
    module("commons-cli-1.4.jar", "commons.cli", "1.4") {
        exports("org.apache.commons.cli")
    }
    module("gson-2.8.0.jar", "com.google.code.gson", "2.8.0") {
        exports("com.google.gson")
    }
    module("testcontainers-1.14.3.jar", "org.testcontainers", "1.14.3") {
        exports("org.testcontainers.containers")
        exports("org.testcontainers.images.builder")
        exports("org.testcontainers.shaded.okhttp3")
        requires("org.junit.jupiter.api")
    }
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

tasks.named<JavaCompile>("compileJava") {
    options.javaModuleVersion.set(provider({ project.version as String }))
}

tasks.named<Test>("test") {
    useJUnitPlatform()
    systemProperty("TEST_SERVER_DOCKER_FILES", File(project(":test-server").buildDir, "dockerFiles"))
    dependsOn(":test-server:dockerFiles")
}

tasks.buildDeb {
    dependsOn(tasks.jlink)

    from("src/deb/usr/bin") {
        into("/usr/bin")
        include("**")
        fileMode = 0x1ED
    }
    from("${buildDir}/image") {
        into("/usr/lib/${project.name}")
    }
    from("LICENSE") {
        into("/usr/share/doc/${project.name}")
    }
    from("README.md") {
        into("/usr/share/doc/${project.name}")
    }
}

tasks.buildRpm {
    dependsOn(tasks.jlink)

    from("src/deb/usr/bin") {
        into("/usr/bin")
        include("**")
        fileMode = 0x1ED
    }
    from("${buildDir}/image") {
        into("/usr/lib/${project.name}")
    }
    from("LICENSE") {
        into("/usr/share/doc/${project.name}")
    }
    from("README.md") {
        into("/usr/share/doc/${project.name}")
    }
}

tasks.register("package") {
    group = "Distribution"
    description = "Builds all packages"
    dependsOn("distTar", "distZip", "buildDeb", "buildRpm")
}