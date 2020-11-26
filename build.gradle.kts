import fi.linuxbox.gradle.download.Download
import org.redline_rpm.header.Architecture.X86_64
import org.redline_rpm.header.Os.LINUX

plugins {
    java
    `maven-publish`
    id("org.beryx.jlink") version ("2.21.2")
    id("org.ysb33r.java.modulehelper") version ("0.9.0")
    id("com.github.sherter.google-java-format") version ("0.9")
    id("nebula.ospackage") version ("8.4.1")
    id("fi.linuxbox.download") version ("0.6")
}

val jmxTermVersion: String by project

allprojects {
    repositories {
        mavenCentral()
        maven {
            url = uri("https://repo.maven.apache.org/maven2")
        }
    }
}

configurations.create("jmxterm")

dependencies {
    implementation("commons-cli:commons-cli:1.4")
    implementation("com.google.code.gson:gson:2.8.0")
}

extraJavaModules {
    module("commons-cli-1.4.jar", "commons.cli", "1.4") {
        exports("org.apache.commons.cli")
    }
    module("gson-2.8.0.jar", "com.google.code.gson", "2.8.6") {
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

tasks.register<Download>("downloadJmxTerm") {
    group = "JmxTerm"
    description = "Downloads jmxterm"
    from("https://github.com/jiaqi/jmxterm/releases/download/v${jmxTermVersion}/jmxterm-${jmxTermVersion}-uber.jar")
    to("${buildDir}/jmxterm/lib/jmxterm-uber.jar")
}

tasks.register<CreateStartScripts>("jmxtermScripts") {
    group = "JmxTerm"
    description = "Scripts to wrap jmxterm"
    applicationName = "jmxterm"
    outputDir = file("${buildDir}/jmxterm/bin")
    mainClassName = ""
    (unixStartScriptGenerator as TemplateBasedScriptGenerator).template = project.resources.text.fromFile(file("src/jmxterm/jmxterm.template.sh"))
    (windowsStartScriptGenerator as TemplateBasedScriptGenerator).template = project.resources.text.fromFile(file("src/jmxterm/jmxterm.template.bat"))
}

tasks.register<Jar>("noarchJar") {
    dependsOn(configurations.runtimeClasspath)
    destinationDirectory.set(file("${buildDir}/distributions"))
    archiveClassifier.set("noarch")

    manifest.attributes("Main-Class" to "org.newrelic.nrjmx.Application")

    from(sourceSets.main.get().output)

    from({
        configurations.runtimeClasspath.get().filter { it.name.endsWith("jar") }.map { zipTree(it) }
    })
}

tasks.register<Zip>("jlinkDistZip") {
    dependsOn(tasks.jlink, "downloadJmxTerm", "jmxtermScripts")
    destinationDirectory.set(file("${buildDir}/distributions"))
    // nrjmx_windows_x.y.z_amd64.zip
    archiveFileName.set("${project.name}_windows_${project.version}_amd64.zip")

    into("${project.name}-${project.version}")

    from("LICENSE")
    from("README.md")
    from("${buildDir}/image")

    from("${buildDir}/jmxterm/lib") {
        include("jmxterm-uber.jar")
    }
    from("${buildDir}/jmxterm/bin") {
        into("bin")
        fileMode = 0x1ED
    }
    from("bin/nrjmx.bat") {
        fileMode = 0x1ED
    }
    from("bin/jmxterm.bat") {
        fileMode = 0x1ED
    }
}

tasks.register<Tar>("jlinkDistTar") {
    dependsOn(tasks.jlink, "downloadJmxTerm", "jmxtermScripts")
    destinationDirectory.set(file("${buildDir}/distributions"))
    // nrjmx_linux_x.y.z_amd64.tar.gz
    archiveFileName.set("${project.name}_linux_${project.version}_amd64.tar.gz")
    compression = Compression.GZIP

    from("LICENSE") {
        into("usr/share/doc/nrjmx")
    }

    from("README.md") {
        into("usr/share/doc/nrjmx")
    }

    from("${buildDir}/image") {
        into("usr/lib/nrjmx")
    }

    from("${buildDir}/jmxterm/lib") {
        include("jmxterm-uber.jar")
        into("usr/lib/nrjmx/lib")
    }

    from("${buildDir}/jmxterm/bin") {
        into("usr/lib/nrjmx/bin")
        fileMode = 0x1ED
    }

    from("src/rpm/usr/bin") {
        into("usr/bin")
        include("**")
        fileMode = 0x1ED
    }
}

tasks.buildDeb {
    dependsOn(tasks.jlink, "downloadJmxTerm", "jmxtermScripts")

    setRelease("1")
    setArch("amd64")
    setOs(LINUX)
    setVendor("New Relic Infrastructure Team <infrastructure-eng@newrelic.com>")
    setPackageGroup("Application/System")
    setLicense("Apache 2.0")

    from("${buildDir}/jmxterm") {
        include("lib/jmxterm-uber.jar")
        into("/usr/lib/${project.name}")
    }

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
    dependsOn(tasks.jlink, "downloadJmxTerm", "jmxtermScripts")

    setRelease("1")
    setArch(X86_64)
    setOs(LINUX)
    setVendor("New Relic Infrastructure Team <infrastructure-eng@newrelic.com>")
    setPackageGroup("Application/System")
    setLicense("Apache 2.0")

    from("${buildDir}/jmxterm") {
        include("lib/jmxterm-uber.jar")
        into("/usr/lib/${project.name}")
    }

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

tasks.distZip {
    dependsOn("downloadJmxTerm", "jmxtermScripts")
    from("${buildDir}/jmxterm") {
        include("**")
        into("${project.name}-${project.version}")
    }
}

tasks.register("package") {
    group = "Distribution"
    description = "Builds all packages"
    dependsOn(
            "package-linux",
            "package-windows")
}

tasks.register("package-linux") {
    group = "Distribution"
    description = "Builds all packages for Linux"
    dependsOn(
            "noarchJar",
            "buildDeb",
            "buildRpm",
            "jlinkDistTar")
}

tasks.register("package-windows") {
    group = "Distribution"
    description = "Builds all packages for Windows"
    dependsOn(
            "jlinkDistZip")
}
