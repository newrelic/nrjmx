package org.newrelic.jmx;

import static spark.Spark.post;
import static spark.Spark.put;

import com.google.gson.Gson;
import java.lang.management.ManagementFactory;
import javax.management.MBeanServer;
import javax.management.ObjectName;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class Service {

  public static void main(String[] args) throws Exception {
    final Logger log = LoggerFactory.getLogger("org.newrelic");

    log.info("JAVA_OPTS={}", System.getenv("JAVA_OPTS"));

    final Gson gson = new Gson();
    final MBeanServer server = ManagementFactory.getPlatformMBeanServer();

    // Registers a cat as an MBean
    post(
        "/cat",
        (req, res) -> {
          Cat cat = gson.fromJson(req.body(), Cat.class);
          log.info("registering {}", cat);
          server.registerMBean(cat, null);
          return "ok!\n";
        });

    final ObjectName queryObject = new ObjectName("*:type=Cat,*");

    // Removes all registered MBean cats
    put(
        "/clear",
        (req, res) -> {
          server
              .queryNames(queryObject, null)
              .forEach(
                  cat -> {
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
