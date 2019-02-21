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
    private String url;
    private String password;
    private String keyStore;
    private String keyStorePassword;
    private String trustStore;
    private String trustStorePassword;
    private boolean verbose;
    private boolean debug;


    Arguments(String[] args) {
        Options options = new Options();
        Option url = Option.builder("U")
            .longOpt("url").desc("JMX url (service:jmx:rmi:///jndi/rmi://%s:%s/jmxrmi)").hasArg().build();
        options.addOption(url);
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


        Option keyStore = Option.builder("keyStore")
                .longOpt("keyStore").desc("SSL keyStore location").hasArg().build();
        options.addOption(keyStore);
        Option keyStorePassword = Option.builder("keyStorePassword")
                .longOpt("keyStorePassword").desc("SSL keyStorePassword").hasArg().build();
        options.addOption(keyStorePassword);
        Option trustStore = Option.builder("trustStore")
                .longOpt("trustStore").desc("SSL trustStore location").hasArg().build();
        options.addOption(trustStore);
        Option trustStorePassword = Option.builder("trustStorePassword")
                .longOpt("trustStorePassword").desc("SSL trustStorePassword").hasArg().build();
        options.addOption(trustStorePassword);

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


        this.url = cmd.getOptionValue("url", "service:jmx:rmi:///jndi/rmi://%s:%s/jmxrmi");
        this.hostname = cmd.getOptionValue("hostname", "localhost");
        this.port     = Integer.parseInt(cmd.getOptionValue("port", "7199"));
        this.username = cmd.getOptionValue("username", "");
        this.password = cmd.getOptionValue("password", "");
        this.keyStore = cmd.getOptionValue("keyStore", "");
        this.keyStorePassword = cmd.getOptionValue("keyStorePassword", "");
        this.trustStore = cmd.getOptionValue("trustStore", "");
        this.trustStorePassword = cmd.getOptionValue("trustStorePassword", "");
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

    public String getKeyStore() {
        return keyStore;
    }

    public String getKeyStorePassword() {
        return keyStorePassword;
    }

    public String getTrustStore() {
        return trustStore;
    }

    public String getTrustStorePassword() {
        return trustStorePassword;
    }

    public boolean isVerbose() {
        return verbose;
    }

    public boolean debugMode() {
        return debug;
    }

	public String getUrl() {
		return url;
	}

	public void setUrl(String url) {
		this.url = url;
	}
}
