plugins {
    `java`
    `maven-publish`
}

repositories {
    mavenLocal()
    maven {
        url = uri("http://repo.maven.apache.org/maven2")
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
}

publishing {
    publications {
        create<MavenPublication>("nrjmx") {
            from(components["java"])
        }
    }
}
