/*
 * Copyright 2020 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package org.newrelic.nrjmx;

import org.apache.commons.cli.HelpFormatter;
import org.apache.thrift.TProcessor;
import org.apache.thrift.protocol.TCompactProtocol;
import org.apache.thrift.server.TServer.Args;
import org.apache.thrift.transport.TServerTransport;
import org.apache.thrift.transport.layered.TFramedTransport;
import org.newrelic.nrjmx.v2.JMXServiceHandler;
import org.newrelic.nrjmx.v2.StandardIOServer;
import org.newrelic.nrjmx.v2.StandardIOTransportServer;
import org.newrelic.nrjmx.v2.nrprotocol.JMXService;

import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.logging.Logger;

public class Application {

    public static void printHelp() {
        new HelpFormatter().printHelp("nrjmx", Arguments.options());
    }

    public static void main(String[] args) throws Exception {
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
            runV2();
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

    private static void runV2() {
        ExecutorService executor = Executors.newSingleThreadExecutor();
        org.newrelic.nrjmx.v2.JMXFetcher jmxFetcher = new org.newrelic.nrjmx.v2.JMXFetcher(executor);

        JMXServiceHandler handler = new JMXServiceHandler(jmxFetcher);
        TProcessor processor = new JMXService.Processor<>(handler);

        TServerTransport serverTransport = new StandardIOTransportServer();
        StandardIOServer server = new StandardIOServer(
                new Args(serverTransport)
                        .processor(processor)
                        .inputTransportFactory(new TFramedTransport.Factory(8192))
                        .outputTransportFactory(new TFramedTransport.Factory(8192))
                        .protocolFactory(new TCompactProtocol.Factory()));

        handler.addServer(server);

        try {
            server.listen();
        } catch (Exception e) {
            e.printStackTrace();
            System.exit(1);
        } finally {
            serverTransport.close();
            executor.shutdownNow();
        }

        // Add ShutdownHook to disconnect the fetcher.
        Runtime.getRuntime().addShutdownHook(
                new Thread() {
                    @Override
                    public void run() {
                        if (jmxFetcher != null) {
                            try {
                                jmxFetcher.disconnect();
                            } catch (Exception e) {
                            }
                        }
                    }
                }
        );
    }

    private static void logTrace(Arguments cliArgs, Logger logger, Exception e) {
        if (cliArgs.isDebugMode()) {
            logger.info("exception trace for " + e.getClass().getCanonicalName() + ": " + e);
        }
    }
}
