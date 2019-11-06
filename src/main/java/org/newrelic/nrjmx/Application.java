package org.newrelic.nrjmx;

import org.apache.commons.cli.HelpFormatter;

import java.util.logging.Level;
import java.util.logging.Logger;

public class Application {
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

        Logger logger = Logger.getLogger("nrjmx");
        Logging.setup(logger, cliArgs.isVerbose());

        // Instantiate a JMXFetcher from the configuration arguments
        JMXFetcher fetcher = cliArgs.getConnectionURL().equals("") ?
                new JMXFetcher(
                        cliArgs.getHostname(), cliArgs.getPort(), cliArgs.getUriPath(),
                        cliArgs.getUsername(), cliArgs.getPassword(),
                        cliArgs.getKeyStore(), cliArgs.getKeyStorePassword(),
                        cliArgs.getTrustStore(), cliArgs.getTrustStorePassword(),
                        cliArgs.getIsRemoteJMX(), cliArgs.getIsRemoteJBossStandalone()
                ) :
                new JMXFetcher(
                        cliArgs.getConnectionURL(),
                        cliArgs.getUsername(), cliArgs.getPassword(),
                        cliArgs.getKeyStore(), cliArgs.getKeyStorePassword(),
                        cliArgs.getTrustStore(), cliArgs.getTrustStorePassword()
                );

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

    private static void logTrace(Arguments cliArgs, Logger logger, Exception e) {
        if (cliArgs.isDebugMode()) {
            logger.info("exception trace for " + e.getClass().getCanonicalName() + ": " + e);
        }
    }
}
