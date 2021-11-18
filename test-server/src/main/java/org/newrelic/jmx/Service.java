package org.newrelic.jmx;

import com.google.gson.Gson;
import com.google.gson.reflect.TypeToken;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.management.MBeanServer;
import javax.management.ObjectName;
import java.lang.management.ManagementFactory;
import java.util.ArrayList;

import static spark.Spark.post;
import static spark.Spark.put;

public class Service {
    public static void main(String[] args) throws Exception {
        final Logger log = LoggerFactory.getLogger("org.newrelic");

        log.info("JAVA_OPTS={}", System.getenv("JAVA_OPTS"));

        final Gson gson = new Gson();
        final MBeanServer server = ManagementFactory.getPlatformMBeanServer();

        // Registers a cat as an MBean
        post("/cat", (req, res) -> {
            Cat cat = gson.fromJson(req.body(), Cat.class);
            log.info("registering {}", cat);
            server.registerMBean(cat, null);
            return "ok!\n";
        });

        post("/cat_batch", (req, res) -> {

            ArrayList<Cat> cats = gson.fromJson(req.body(),  new TypeToken<ArrayList<Cat>>() {}.getType());
            for(Cat cat : cats) {
                server.registerMBean(cat, null);
            }

            return "ok!\n";
        });

        final ObjectName queryObject = new ObjectName("*:type=Cat,*");

        // Removes all registered MBean cats
        put("/clear", (req, res) -> {
            server.queryNames(queryObject, null).forEach(cat -> {
                log.info("unregistering {}", cat);
                try {
                    server.unregisterMBean(cat);
                } catch (Exception e) {
                    log.error("unregistering", e);
                }
            });
            return "ok!\n";
        });
    }
}
