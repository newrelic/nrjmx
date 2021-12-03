/*
 * Copyright 2020 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package org.newrelic.nrjmx;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.nio.file.StandardOpenOption;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.logging.Logger;

import org.apache.commons.cli.HelpFormatter;
import org.apache.thrift.TProcessor;
import org.apache.thrift.protocol.TCompactProtocol;
import org.apache.thrift.server.TServer;
import org.apache.thrift.server.TServer.Args;
import org.apache.thrift.server.TSimpleServer;
import org.apache.thrift.transport.*;
import org.apache.thrift.transport.layered.TFramedTransport;
import org.newrelic.nrjmx.v2.JMXServiceHandler;
import org.newrelic.nrjmx.v2.StandardIOServer;
import org.newrelic.nrjmx.v2.StandardIOTransportServer;
import org.newrelic.nrjmx.v2.nrprotocol.JMXService;

public class Application {

    public static void printHelp() {
        new HelpFormatter().printHelp("nrjmx", Arguments.options());
    }

    public static void main(String[] args) throws TTransportException {
        Arguments cliArgs = null;
        try {
            cliArgs = Arguments.from(args);
        } catch (Exception e) {
            System.err.println(e.getMessage());
            printHelp();
            System.exit(1);
        }

        if (cliArgs.isHelp()) {
            printHelp();
            System.exit(0);
        }

        if (!cliArgs.isProtocolV2()) {
            runV1(cliArgs);
        } else {
            runV2(cliArgs);
        }
    }

    private static void runV1(Arguments cliArgs) {
        Logger logger = Logger.getLogger("nrjmx");
        Logging.setup(logger, cliArgs.isVerbose());

        // Instantiate a JMXFetcher from the configuration arguments
        JMXFetcher fetcher =
                cliArgs.getConnectionURL().equals("")
                        ? new JMXFetcher(
                        cliArgs.getHostname(),
                        cliArgs.getPort(),
                        cliArgs.getUriPath(),
                        cliArgs.getUsername(),
                        cliArgs.getPassword(),
                        cliArgs.getKeyStore(),
                        cliArgs.getKeyStorePassword(),
                        cliArgs.getTrustStore(),
                        cliArgs.getTrustStorePassword(),
                        cliArgs.getIsRemoteJMX(),
                        cliArgs.getIsRemoteJBossStandalone())
                        : new JMXFetcher(
                        cliArgs.getConnectionURL(),
                        cliArgs.getUsername(),
                        cliArgs.getPassword(),
                        cliArgs.getKeyStore(),
                        cliArgs.getKeyStorePassword(),
                        cliArgs.getTrustStore(),
                        cliArgs.getTrustStorePassword());

        try {
            fetcher.run(System.in, System.out);
        } catch (JMXFetcher.ConnectionError e) {
            logger.severe("jmx connection error: " + e.getMessage());
            logTrace(cliArgs, logger, e);
            System.exit(1);
        } catch (Exception e) {
            logger.severe("error running nrjmx: " + e.getMessage());
            logTrace(cliArgs, logger, e);
            System.exit(1);
        }
    }

    public static void write(String message) {

        try {
            message += "\n";
            Files.write(Paths.get("/Users/cciutea/workspace/nr/int/nrjmx/gojmx/cmd/out2"), message.getBytes(), StandardOpenOption.APPEND);
        } catch (IOException e) {
            //exception handling left as an exercise for the reader
        }

    }

    private static void runV2(Arguments cliArgs) {
        Runtime.getRuntime().addShutdownHook(new Thread() {
            public void run() {
                write("Shutdown Hook is running !");
            }
        });

        ExecutorService executor = Executors.newSingleThreadExecutor();
        org.newrelic.nrjmx.v2.JMXFetcher jmxFetcher = new org.newrelic.nrjmx.v2.JMXFetcher(executor);

        JMXServiceHandler handler = new JMXServiceHandler(jmxFetcher);
        TProcessor processor = new JMXService.Processor<JMXServiceHandler>(handler);

        TServerTransport serverTransport = new StandardIOTransportServer();
        TServer server = new StandardIOServer(
                new Args(serverTransport).processor(processor).protocolFactory(new TCompactProtocol.Factory()));

        handler.addServer(server);
        server.serve();

        serverTransport.close();
        executor.shutdownNow();
    }

    private static void runV3(Arguments cliArgs) throws TTransportException {
        ExecutorService executor = Executors.newSingleThreadExecutor();
        org.newrelic.nrjmx.v2.JMXFetcher jmxFetcher = new org.newrelic.nrjmx.v2.JMXFetcher(executor);

        JMXServiceHandler handler = new JMXServiceHandler(jmxFetcher);
        TProcessor processor = new JMXService.Processor<JMXServiceHandler>(handler);

        TServerTransport serverTransport = new TServerSocket(9090);
        TServer server = new TSimpleServer(
                new Args(serverTransport).processor(processor)
                        .inputTransportFactory(new TFramedTransport.Factory(8192))
                        .outputTransportFactory(new TFramedTransport.Factory(8192))
                        .protocolFactory(new TCompactProtocol.Factory()));

        handler.addServer(server);
        server.serve();

        serverTransport.close();
    }

    private static void logTrace(Arguments cliArgs, Logger logger, Exception e) {
        if (cliArgs.isDebugMode()) {
            logger.info("exception trace for " + e.getClass().getCanonicalName() + ": " + e);
        }
    }
}
