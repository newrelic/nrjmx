package org.newrelic.jmx

import com.google.gson.Gson
import org.slf4j.LoggerFactory
import spark.Spark.*
import java.lang.management.ManagementFactory
import javax.management.ObjectName

fun main() {
    val log = LoggerFactory.getLogger("org.newrelic")

    log.info("JAVA_OPTS=${System.getenv("JAVA_OPTS")}")

    val server = ManagementFactory.getPlatformMBeanServer()

    post("/cat") { req, _ ->
        val cat: Cat = Gson().fromJson(req.body(), Cat::class.java)
        log.info("registering $cat")
        server.registerMBean(cat, null)
        "ok!"
    }

    put("/clear") { _, _ ->
        server.queryNames(ObjectName("*:type=Cat,*"), null).forEach {
            log.info("unregistering $it")
        }
        "ok!"
    }
}
