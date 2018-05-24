package org.newrelic.nrjmx;

import org.apache.commons.cli.CommandLine;
import org.apache.commons.cli.CommandLineParser;
import org.apache.commons.cli.DefaultParser;
import org.apache.commons.cli.HelpFormatter;
import org.apache.commons.cli.Option;
import org.apache.commons.cli.Options;
import org.apache.commons.cli.ParseException;

public class Arguments {
    private String hostname;
    private int port;
    private String username;
    private String password;
    private boolean verbose;
    private boolean debug;


    Arguments(String[] args) {
        Options options = new Options();
        Option hostname = Option.builder("H")
            .longOpt("hostname").desc("JMX hostname (localhost)").hasArg().build();
        options.addOption(hostname);
        Option port = Option.builder("P")
            .longOpt("port").desc("JMX port (7199)").hasArg().build();
        options.addOption(port);
        Option username = Option.builder("u")
            .longOpt("username").desc("JMX username").hasArg().build();
        options.addOption(username);
        Option password = Option.builder("p")
            .longOpt("password").desc("JMX password").hasArg().build();
        options.addOption(password);

        Option verbose = Option.builder("v")
            .longOpt("verbose").desc("Verbose output").hasArg(false).build();
        options.addOption(verbose);
        Option debug = Option.builder("d")
                .longOpt("debug").desc("Debug mode").hasArg(false).build();
        options.addOption(debug);
        Option help = Option.builder("h")
            .longOpt("help").desc("Show help").hasArg(false).build();
        options.addOption(help);

        HelpFormatter formatter = new HelpFormatter();
        CommandLineParser parser = new DefaultParser();
        CommandLine cmd = null;

        try {
            cmd = parser.parse(options, args);
        } catch (ParseException e) {
            System.err.println(e.getMessage());
            formatter.printHelp("nrjmx", options);
            System.exit(1);
        }
        if (cmd.hasOption("help")) {
            formatter.printHelp("nrjmx", options);
            System.exit(0);
        }

        this.hostname = cmd.getOptionValue("hostname", "localhost");
        this.port     = Integer.parseInt(cmd.getOptionValue("port", "7199"));
        this.username = cmd.getOptionValue("username", "");
        this.password = cmd.getOptionValue("password", "");
        this.verbose  = cmd.hasOption("verbose");
        this.debug    = cmd.hasOption("debug");
    }

    public String getHostname() {
        return hostname;
    }

    public int getPort() {
        return port;
    }

    public String getUsername() {
        return username;
    }

    public String getPassword() {
        return password;
    }

    public boolean isVerbose() {
        return verbose;
    }

    public boolean debugMode() {
        return debug;
    }
}
