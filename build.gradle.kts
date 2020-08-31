import org.redline_rpm.header.Architecture.X86_64
import org.redline_rpm.header.Os.LINUX

plugins {
    java
    `maven-publish`
    id("org.beryx.jlink") version ("2.21.2")
    id("org.ysb33r.java.modulehelper") version ("0.9.0")
    id("com.github.sherter.google-java-format") version ("0.8")
    id("nebula.ospackage") version ("8.4.1")
}

allprojects {
    repositories {
        maven {
            url = uri("https://repo.maven.apache.org/maven2")
        }
    }
}

dependencies {
    implementation("commons-cli:commons-cli:1.4")
    implementation("com.google.code.gson:gson:2.8.0")
}

extraJavaModules {
    module("commons-cli-1.4.jar", "commons.cli", "1.4") {
        exports("org.apache.commons.cli")
    }
    module("gson-2.8.0.jar", "com.google.code.gson", "2.8.0") {
        exports("com.google.gson")
    }
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
}

tasks.named<JavaCompile>("compileJava") {
    options.javaModuleVersion.set(provider({ project.version as String }))
}

tasks.test {
    enabled = false
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
    setRelease("1")
    setArch(X86_64)
    setOs(LINUX)
    setVendor("New Relic Infrastructure Team <infrastructure-eng@newrelic.com>")
    setPackageGroup("Application/System")
    setLicense("Apache 2.0")

    dependsOn(tasks.jlink)

    from("src/rpm/usr/bin") {
        into("/usr/bin")
        include("**")
        fileMode = 0x1ED
    }
    from("${buildDir}/image") {
        into("/usr/lib/${project.name}")
    }
    from("LICENSE") {
        into("/usr/share/doc/${project.name}")
        addParentDirs = false
    }
    from("README.md") {
        into("/usr/share/doc/${project.name}")
        addParentDirs = false
    }
}

tasks.register("package") {
    group = "Distribution"
    description = "Builds all packages"
    dependsOn("distTar", "distZip", "buildDeb", "buildRpm")
}