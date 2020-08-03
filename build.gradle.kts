// TODO: Check generated POM - should be the same as for maven build
// TODO: Adapt maven-assembly-plugin
// TODO: Adapt copy-rename-maven-plugin
// TODO: Build the testserver (dockerfile-maven-plugin)
// TODO: tarball profile equivalent
// TODO: deb
// TODO: rpm

plugins {
    `java`
    `maven-publish`
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
    testImplementation("junit:junit:4.12")
    testImplementation("org.testcontainers:testcontainers:1.11.3")
}

java {
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
