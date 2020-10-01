/*
 * Copyright 2020 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package org.newrelic.nrjmx;

import java.util.logging.*;

public class Logging {
  public static void setup(Logger logger, boolean verbose) {
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
}
