package org.newrelic.nrjmx;

import com.google.gson.Gson;
import org.apache.commons.cli.HelpFormatter;
import org.newrelic.nrjmx.JMXFetcher.ConnectionError;
import org.newrelic.nrjmx.JMXFetcher.QueryError;

import javax.management.ObjectInstance;
import java.util.Scanner;
import java.util.Set;
import java.util.logging.*;

public class Application {
    private static final Logger logger = Logger.getLogger("nrjmx");

    private static void setupLogging(boolean verbose) {
        logger.setUseParentHandlers(false);
        Handler consoleHandler = new ConsoleHandler();
        logger.addHandler(consoleHandler);

        consoleHandler.setFormatter(new SimpleFormatter());

        if (verbose) {
            logger.setLevel(Level.FINE);
            consoleHandler.setLevel(Level.FINE);
        } else {
            logger.setLevel(Level.INFO);
            consoleHandler.setLevel(Level.INFO);
        }
    }

    public static void printHelp() {
        new HelpFormatter().printHelp("nrjmx", Arguments.options());
    }

    public static void main(String[] args) {
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

        setupLogging(cliArgs.isVerbose());

        // TODO: move all the code below to a testable class
        JMXFetcher fetcher = null;
        try {
            fetcher = new JMXFetcher(
                cliArgs.getHostname(), cliArgs.getPort(),
                cliArgs.getUsername(), cliArgs.getPassword(),
                cliArgs.getKeyStore(), cliArgs.getKeyStorePassword(),
                cliArgs.getTrustStore(), cliArgs.getTrustStorePassword(),
                cliArgs.getIsRemoteJMX()
            );
        } catch (ConnectionError e) {
            logger.severe(e.getMessage());
            logger.log(Level.FINE, e.getMessage(), e);
            System.exit(1);
        } catch (Exception e) {
            if (cliArgs.isDebugMode()) {
                e.printStackTrace();
            } else {
                System.out.println(e.getClass().getCanonicalName());
                logger.severe(e.getClass().getCanonicalName() + ": " + e.getMessage());
                logger.log(Level.FINE, e.getMessage(), e);
            }
            System.exit(1);
        }

        Gson gson = new Gson();

        try (Scanner input = new Scanner(System.in)) {
            while (input.hasNextLine()) {
                String beanName = input.nextLine();

                Set<ObjectInstance> beanInstances;
                try {
                    beanInstances = fetcher.query(beanName);
                } catch (QueryError e) {
                    logger.warning(e.getMessage());
                    logger.log(Level.FINE, e.getMessage(), e);
                    continue;
                }

                for (ObjectInstance instance : beanInstances) {
                    try {
                        fetcher.queryAttributes(instance);
                    } catch (QueryError e) {
                        logger.warning(e.getMessage());
                        logger.log(Level.FINE, e.getMessage(), e);
                    }
                }

                System.out.println(gson.toJson(fetcher.popResults()));
            }
            logger.info("Stopped receiving data, leaving...\n");
        }

    }
}
