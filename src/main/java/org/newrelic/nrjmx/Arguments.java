/*
 * Copyright 2020 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package org.newrelic.nrjmx;

import org.apache.commons.cli.CommandLine;
import org.apache.commons.cli.DefaultParser;
import org.apache.commons.cli.Option;
import org.apache.commons.cli.Options;
import org.apache.commons.cli.ParseException;

class Arguments {

  private static Options options = null;
  private String hostname;
  private String connectionURL;
  private int port;
  private String uriPath;
  private String username;
  private String password;
  private String keyStore;
  private String keyStorePassword;
  private String trustStore;
  private String trustStorePassword;
  private boolean verbose;
  private boolean debug;
  private boolean isRemoteJMX;
  private boolean isRemoteJBossStandalone;
  private boolean help;

  private Arguments() {}

  static Options options() {
    if (options == null) {
      options = new Options();
      Option connectionURL =
          Option.builder("C")
              .longOpt("connURL")
              .desc("full connection URL. Default none.")
              .hasArg()
              .build();
      options.addOption(connectionURL);
      Option hostname =
          Option.builder("H").longOpt("hostname").desc("JMX hostname (localhost)").hasArg().build();
      options.addOption(hostname);
      Option port = Option.builder("P").longOpt("port").desc("JMX port (7199)").hasArg().build();
      options.addOption(port);
      Option uriPath =
          Option.builder("U")
              .longOpt("uriPath")
              .desc("path for the JMX service URI. Defaults to jmxrmi")
              .hasArg()
              .build();
      options.addOption(uriPath);
      Option username =
          Option.builder("u").longOpt("username").desc("JMX username").hasArg().build();
      options.addOption(username);
      Option password =
          Option.builder("p").longOpt("password").desc("JMX password").hasArg().build();
      options.addOption(password);

      Option keyStore =
          Option.builder("keyStore")
              .longOpt("keyStore")
              .desc("SSL keyStore location")
              .hasArg()
              .build();
      options.addOption(keyStore);
      Option keyStorePassword =
          Option.builder("keyStorePassword")
              .longOpt("keyStorePassword")
              .desc("SSL keyStorePassword")
              .hasArg()
              .build();
      options.addOption(keyStorePassword);
      Option trustStore =
          Option.builder("trustStore")
              .longOpt("trustStore")
              .desc("SSL trustStore location")
              .hasArg()
              .build();
      options.addOption(trustStore);
      Option trustStorePassword =
          Option.builder("trustStorePassword")
              .longOpt("trustStorePassword")
              .desc("SSL trustStorePassword")
              .hasArg()
              .build();
      options.addOption(trustStorePassword);

      Option verbose =
          Option.builder("v").longOpt("verbose").desc("Verbose output").hasArg(false).build();
      options.addOption(verbose);
      Option debug = Option.builder("d").longOpt("debug").desc("Debug mode").hasArg(false).build();
      options.addOption(debug);
      Option help = Option.builder("h").longOpt("help").desc("Show help").hasArg(false).build();
      options.addOption(help);

      Option remote =
          Option.builder("r").longOpt("remote").desc("Remote JMX mode").hasArg(false).build();
      options.addOption(remote);
      Option remoteJBossStandalone =
          Option.builder("s")
              .longOpt("remoteJBossStandalone")
              .desc("Remote JBoss Standalone mode")
              .hasArg(false)
              .build();
      options.addOption(remoteJBossStandalone);
    }
    return options;
  }

  static Arguments from(String[] args) throws ParseException {
    CommandLine cmd = new DefaultParser().parse(options(), args);

    Arguments argsObj = new Arguments();
    argsObj.connectionURL = cmd.getOptionValue("connURL", "");
    argsObj.hostname = cmd.getOptionValue("hostname", "localhost");
    argsObj.port = Integer.parseInt(cmd.getOptionValue("port", "7199"));
    argsObj.uriPath = cmd.getOptionValue("uriPath", "jmxrmi");
    argsObj.username = cmd.getOptionValue("username", "");
    argsObj.password = cmd.getOptionValue("password", "");
    argsObj.keyStore = cmd.getOptionValue("keyStore", "");
    argsObj.keyStorePassword = cmd.getOptionValue("keyStorePassword", "");
    argsObj.trustStore = cmd.getOptionValue("trustStore", "");
    argsObj.trustStorePassword = cmd.getOptionValue("trustStorePassword", "");
    argsObj.verbose = cmd.hasOption("verbose");
    argsObj.help = cmd.hasOption("help");
    argsObj.debug = cmd.hasOption("debug");
    argsObj.isRemoteJMX = cmd.hasOption("remote");
    argsObj.isRemoteJBossStandalone = cmd.hasOption("remoteJBossStandalone");
    return argsObj;
  }

  String getConnectionURL() {
    return connectionURL;
  }

  String getHostname() {
    return hostname;
  }

  int getPort() {
    return port;
  }

  String getUriPath() {
    return uriPath;
  }

  String getUsername() {
    return username;
  }

  boolean getIsRemoteJMX() {
    return isRemoteJMX;
  }

  boolean getIsRemoteJBossStandalone() {
    return isRemoteJBossStandalone;
  }

  String getPassword() {
    return password;
  }

  String getKeyStore() {
    return keyStore;
  }

  String getKeyStorePassword() {
    return keyStorePassword;
  }

  String getTrustStore() {
    return trustStore;
  }

  String getTrustStorePassword() {
    return trustStorePassword;
  }

  boolean isVerbose() {
    return verbose;
  }

  boolean isDebugMode() {
    return debug;
  }

  boolean isHelp() {
    return help;
  }
}
