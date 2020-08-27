plugins {
    java
    application
}

repositories {
    mavenCentral()
}

dependencies {
    implementation ("com.google.code.gson:gson:2.8.5")
    implementation ("com.sparkjava:spark-core:2.9.1")
    implementation ("org.slf4j:slf4j-simple:1.7.26")
}

tasks.installDist {
    from ("src/docker") {
        include ("**" )
    }
}

// No need to run `test` as we are only building the JAR to be used elsewhere
tasks.test {
    enabled = false
}

// No need to run `distTar`
tasks.distTar {
    enabled = false
}

application {
    mainClass.set("org.newrelic.jmx.Service")
}

tasks.assemble {
    dependsOn("dockerFiles")
}