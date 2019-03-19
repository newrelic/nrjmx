package org.newrelic.nrjmx;

import java.util.NoSuchElementException;
import java.util.Scanner;
import java.util.Set;
import java.util.logging.ConsoleHandler;
import java.util.logging.Handler;
import java.util.logging.Level;
import java.util.logging.Logger;
import java.util.logging.SimpleFormatter;

import javax.management.ObjectInstance;

import org.newrelic.nrjmx.JMXFetcher.ConnectionError;
import org.newrelic.nrjmx.JMXFetcher.QueryError;

import com.google.gson.Gson;

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

    public static void main(String[] args) {
        Arguments cliArgs = new Arguments(args);
        setupLogging(cliArgs.isVerbose());

        JMXFetcher fetcher;
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
            System.exit(1);
            return;
        } catch (Exception e) {
            if(cliArgs.debugMode()) {
                e.printStackTrace();
                System.exit(1);
                return;
            }else{
                System.out.println( e.getClass().getCanonicalName());
                logger.severe(e.getClass().getCanonicalName() + ": " + e.getMessage());
                System.exit(1);
                return;
            }
        }

        Gson gson = new Gson();

        Scanner input = new Scanner(System.in);
        while (true) {
            String beanName;

            try {
                beanName = input.nextLine();
            } catch (NoSuchElementException e) {
                logger.info("Stopped receiving data, leaving...\n");
                input.close();
                break;
            }

            Set<ObjectInstance> beanInstances;
            try {
                beanInstances = fetcher.query(beanName);
            } catch (QueryError e) {
                logger.warning(e.getMessage());
                continue;
            }

            for (ObjectInstance instance : beanInstances) {
                try {
                    fetcher.queryAttributes(instance);
                } catch (QueryError e) {
                    logger.warning(e.getMessage());
                }
            }
            
            System.out.println(gson.toJson(fetcher.popResults()));
        }
    }
}
