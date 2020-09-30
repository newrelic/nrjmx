plugins {
    java
}

dependencies {
    testImplementation(rootProject)
    testImplementation("org.junit.jupiter:junit-jupiter-api:5.6.2")
    testImplementation("org.junit.jupiter:junit-jupiter-params:5.6.2")
    testImplementation("org.testcontainers:testcontainers:1.14.3")
    testRuntimeOnly("org.junit.jupiter:junit-jupiter-engine")
    testRuntimeOnly("ch.qos.logback:logback-classic:1.2.3")
}

tasks.test {
    useJUnitPlatform()
    systemProperty("TEST_SERVER_DOCKER_FILES", File(project(":test-server").buildDir, "install/test-server"))
    dependsOn(":test-server:installDist")
}
