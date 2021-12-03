/*
 * Copyright 2020 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package org.newrelic.jmx;

import org.junit.Assert;
import org.junit.Test;
import org.newrelic.nrjmx.JMXFetcher;
import org.newrelic.nrjmx.Logging;
import org.newrelic.nrjmx.v2.nrprotocol.JMXConfig;
import org.slf4j.LoggerFactory;
import org.testcontainers.containers.GenericContainer;
import org.testcontainers.containers.output.Slf4jLogConsumer;

import java.io.*;
import java.util.Arrays;
import java.util.Collections;
import java.util.logging.Logger;

import static org.junit.Assert.assertEquals;

public class JMXFetcherTest {

    // Runs the JMX-monitored test container without SSL enabled
    private static GenericContainer jmxService() {
        GenericContainer container = new GenericContainer<>("testserver:latest")
                .withExposedPorts(4567, 7199)
                .withEnv("JAVA_OPTS", "-Dcom.sun.management.jmxremote.port=7199 " +
                        "-Dcom.sun.management.jmxremote.rmi.port=7199 " +
                        "-Dcom.sun.management.jmxremote.local.only=false" +
                        "-Djava.rmi.server.hostname=0.0.0.0 " +
                        "-Dcom.sun.management.jmxremote=true " +
                        "-Dcom.sun.management.jmxremote.authenticate=false " +
                        "-Dcom.sun.management.jmxremote.ssl=false ");
        container.setPortBindings(Arrays.asList("7199:7199", "4567:4567"));
        Slf4jLogConsumer logConsumer = new Slf4jLogConsumer(LoggerFactory.getLogger("TESTCONT"));
        container.setLogConsumers(Collections.singletonList(logConsumer));
        return container;
    }

    // Runs the JMX-monitored test container with SSL enabled
    private static GenericContainer jmxSSLService() {
        GenericContainer container = new GenericContainer<>("testserver:latest")
                .withExposedPorts(4567, 7199)
                .withEnv("JAVA_OPTS", "-Dcom.sun.management.jmxremote.port=7199 "
                        + "-Dcom.sun.management.jmxremote.local.only=false "
                        + "-Dcom.sun.management.jmxremote.rmi.port=7199 "
                        + "-Dcom.sun.management.jmxremote=true "
                        + "-Dcom.sun.management.jmxremote.authenticate=false "
                        + "-Dcom.sun.management.jmxremote.ssl=true "
                        + "-Dcom.sun.management.jmxremote.ssl.need.client.auth=true  "
                        + "-Dcom.sun.management.jmxremote.registry.ssl=true "
                        + "-Djavax.net.ssl.keyStore=/keystore  "
                        + "-Djavax.net.ssl.keyStorePassword=password "
                        + "-Djavax.net.ssl.trustStore=/truststore "
                        + "-Djavax.net.ssl.trustStorePassword=password");
        container.setPortBindings(Arrays.asList("7199:7199", "4567:4567"));
        Slf4jLogConsumer logConsumer = new Slf4jLogConsumer(LoggerFactory.getLogger("TESTCONT"));
        container.setLogConsumers(Collections.singletonList(logConsumer));
        return container;
    }

    @Test(timeout = 20_000)
    public void testJMX() throws Exception {
        GenericContainer container = jmxService();
        try {
            container.start();
            testJMXFetching(
                    container.getContainerIpAddress(),
                    new JMXFetcher(container.getContainerIpAddress(), 7199,
                            "", "", "", "", "", "", false));
        } finally {
            container.close();
        }
    }


    @Test(timeout = 20_000)
    public void testJMXFromConnectionURL() throws Exception {
        GenericContainer container = jmxService();
        try {
            container.start();
            testJMXFetching(
                    container.getContainerIpAddress(),
                    new JMXFetcher("service:jmx:rmi:///jndi/rmi://" + container.getContainerIpAddress() + ":7199/jmxrmi",
                            "", "", "", "", "", ""));
        } finally {
            container.close();
        }
    }

    @Test(timeout = 60_000)
    public void testJMXWithSSL() throws Exception {
        GenericContainer container = jmxSSLService();
        try {
            container.start();
            testJMXFetching(container.getHost(),
                    new JMXFetcher(container.getHost(), 7199, "", "",
                            getClass().getResource("/keystore").getPath(), "password",
                            getClass().getResource("/truststore").getPath(), "password",
                            false));
        } finally {
            container.close();
        }
    }

    public void testJMXFetching(String host, JMXFetcher jmxFetcher) throws Exception {
        Logging.setup(Logger.getLogger("nrjmx"), true);
        // Test preparation
        // builds a piped, readable output stream
        PipedOutputStream output = new PipedOutputStream();
        PipedInputStream resultsIs = new PipedInputStream(output);
        BufferedReader results = new BufferedReader(new InputStreamReader(resultsIs));

        // GIVEN a container
        // WITH some monitored objects
        final CatsClient cats = new CatsClient("http://" + host + ":4567");

        eventually(10_000, new Runnable() {
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


    @Test(timeout = 20_000)
    public void testGetConnectionURL_Remote_DefaultURIPath() {
        JMXConfig jmxConfig = new JMXConfig()
                .setHostname("localhost")
                .setPort(1234)
                .setIsRemote(true)
                .setIsJBossStandaloneMode(true)
                .setUseSSL(true);
        String actual = org.newrelic.nrjmx.v2.JMXFetcher.buildConnectionString(jmxConfig);

        Assert.assertEquals("Connection URL should match","service:jmx:remote+https://localhost:1234",actual);
    }

    @Test(timeout = 20_000)
    public void testGetConnectionURL_Remote_WithURIPath() {
        JMXConfig jmxConfig = new JMXConfig()
                .setHostname("localhost")
                .setPort(1234)
                .setIsRemote(true)
                .setIsJBossStandaloneMode(true)
                .setUseSSL(true)
                .setUriPath("/something");

        String actual = org.newrelic.nrjmx.v2.JMXFetcher.buildConnectionString(jmxConfig);

        Assert.assertEquals("Connection URL should match","service:jmx:remote+https://localhost:1234/something/",actual);
    }

    @Test(timeout = 20_000)
    public void testGetConnectionURL_Remote_ConnectionURL_Has_Precedence() {
        JMXConfig jmxConfig = new JMXConfig()
                .setConnectionURL("special_connection_URL")
                .setHostname("localhost")
                .setPort(1234)
                .setIsRemote(true)
                .setIsJBossStandaloneMode(true)
                .setUseSSL(true)
                .setUriPath("/something");

        String actual = org.newrelic.nrjmx.v2.JMXFetcher.buildConnectionString(jmxConfig);

        Assert.assertEquals("Connection URL should match","special_connection_URL", actual);
    }

    @Test(timeout = 20_000)
    public void testGetConnectionURL_RMI() {
        JMXConfig jmxConfig = new JMXConfig()
                .setHostname("localhost")
                .setPort(1234);

        String actual = org.newrelic.nrjmx.v2.JMXFetcher.buildConnectionString(jmxConfig);

        Assert.assertEquals("Connection URL should match","service:jmx:rmi:///jndi/rmi://localhost:1234/jmxrmi", actual);
    }

    @Test(timeout = 20_000)
    public void testGetConnectionURL_RMI_With_URIPath() {
        JMXConfig jmxConfig = new JMXConfig()
                .setHostname("localhost")
                .setPort(1234)
                .setUriPath("something");

        String actual = org.newrelic.nrjmx.v2.JMXFetcher.buildConnectionString(jmxConfig);

        Assert.assertEquals("Connection URL should match","service:jmx:rmi:///jndi/rmi://localhost:1234/something", actual);
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
