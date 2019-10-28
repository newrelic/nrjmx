package org.newrelic.jmx;

import org.junit.Assert;
import org.junit.Test;
import org.newrelic.nrjmx.JMXFetcher;
import org.newrelic.nrjmx.Logging;
import org.testcontainers.containers.GenericContainer;

import java.io.*;
import java.util.Arrays;
import java.util.logging.Logger;

import static org.junit.Assert.assertEquals;

public class JMXFetcherTest {

    // Runs the JMX-monitored test container without SSL enabled
    private static GenericContainer jmxService() {
        GenericContainer container = new GenericContainer<>("testserver:latest")
            .withExposedPorts(4567, 7199)
            .withEnv("JAVA_OPTS", "-Dcom.sun.management.jmxremote.port=7199 " +
                "-Dcom.sun.management.jmxremote.rmi.port=7199 " +
                "-Djava.rmi.server.hostname=localhost " +
                "-Dcom.sun.management.jmxremote=true " +
                "-Dcom.sun.management.jmxremote.authenticate=false " +
                "-Dcom.sun.management.jmxremote.ssl=false ");
        container.setPortBindings(Arrays.asList("7199:7199", "4567:4567"));
        return container;
    }

    // Runs the JMX-monitored test container with SSL enabled
    private static GenericContainer jmxSSLService() {
        GenericContainer container = new GenericContainer<>("testserver:latest")
            .withEnv("JAVA_OPTS", "-Dcom.sun.management.jmxremote.port=7199 " +
                "-Dcom.sun.management.jmxremote.rmi.port=7199 " +
                "-Djava.rmi.server.hostname=localhost " +
                "-Dcom.sun.management.jmxremote=true " +
                "-Dcom.sun.management.jmxremote.authenticate=false " +
                "-Dcom.sun.management.jmxremote.ssl=true " +
                "-Dcom.sun.management.jmxremote.ssl.need.client.auth=true  " +
                "-Dcom.sun.management.jmxremote.registry.ssl=true  " +
                "-Djavax.net.ssl.keyStore=/serverkeystore  " +
                "-Djavax.net.ssl.keyStorePassword=serverpass  " +
                "-Djavax.net.ssl.trustStore=/servertruststore  " +
                "-Djavax.net.ssl.trustStorePassword=servertrustpass");
        container.setPortBindings(Arrays.asList("7199:7199", "4567:4567"));
        return container;
    }

    @Test(timeout = 20_000)
    public void testJMX() throws Exception {
        GenericContainer container = jmxService();
        try {
            container.start();
            testJMXFetching(new JMXFetcher("localhost", 7199,
                "", "", "", "", "", "", false));
        } finally {
            container.close();
        }
    }

    @Test(timeout = 20_000)
    public void testJMXWithSSL() throws Exception {
        GenericContainer container = jmxSSLService();
        try {
            container.start();
            testJMXFetching(new JMXFetcher("localhost", 7199, "", "",
                getClass().getResource("/clientkeystore").getPath(), "clientpass",
                getClass().getResource("/clienttruststore").getPath(), "clienttrustpass",
                false));
        } finally {
            container.close();
        }
    }


    public void testJMXFetching(JMXFetcher jmxFetcher) throws Exception {
        Logging.setup(Logger.getLogger("nrjmx"), true);
        // Test preparation
        // builds a piped, readable output stream
        PipedOutputStream output = new PipedOutputStream();
        PipedInputStream resultsIs = new PipedInputStream(output);
        BufferedReader results = new BufferedReader(new InputStreamReader(resultsIs));

        // GIVEN a container
        // WITH some monitored objects
        final CatsClient cats = new CatsClient("http://localhost:4567");

        eventually(5_000, new Runnable() {
            @Override
            public void run() {
                Assert.assertEquals("ok!\n", cats.add("Isidoro"));
            }
        });

        Assert.assertEquals("ok!\n", cats.add("Heathcliff"));

        // WHEN queries are submitted
        ByteArrayInputStream queries = new ByteArrayInputStream(
            ("test:*\n" +
                "test:type=Cat,*\n" +
                "this is a wrong query and will be ignored\n" +
                "test:type=Cat,name=Isidoro\n" +
                "test:type=*,name=Heathcliff\n" +
                "test:type=Dog,*\n").getBytes());
        queries.close();

        // AND a JMXFetcher reads them
        jmxFetcher.run(queries, output);

        // THEN the corresponding JMX objects are returned in the same query order,
        // ignoring the invalid queries
        assertEquals("{\"test:type\\u003dCat,name\\u003dIsidoro,attr\\u003dName\":\"Isidoro\"," +
                "\"test:type\\u003dCat,name\\u003dHeathcliff,attr\\u003dName\":\"Heathcliff\"}",
            results.readLine());
        assertEquals("{\"test:type\\u003dCat,name\\u003dIsidoro,attr\\u003dName\":\"Isidoro\"," +
                "\"test:type\\u003dCat,name\\u003dHeathcliff,attr\\u003dName\":\"Heathcliff\"}",
            results.readLine());
        assertEquals("{\"test:type\\u003dCat,name\\u003dIsidoro,attr\\u003dName\":\"Isidoro\"}",
            results.readLine());
        assertEquals("{\"test:type\\u003dCat,name\\u003dHeathcliff,attr\\u003dName\":\"Heathcliff\"}",
            results.readLine());
        assertEquals("{}", results.readLine());

        results.close();
    }

    private static void eventually(long timeoutMs, Runnable r) throws Exception {
        long timeoutNano = timeoutMs * 1_000_000;
        long startTime = System.nanoTime();
        Exception lastException = null;
        while (System.nanoTime() - timeoutNano < startTime) {
            try {
                r.run();
                return;
            } catch (Exception e) {
                lastException = e;
            }
        }
        throw lastException;
    }
}
