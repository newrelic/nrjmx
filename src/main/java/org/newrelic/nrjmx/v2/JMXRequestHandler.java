package org.newrelic.nrjmx.v2;

import java.rmi.ConnectException;
import java.util.Arrays;
import java.util.HashSet;
import java.util.Set;
import java.util.concurrent.Callable;

public class JMXRequestHandler {

    private static Set<String> knownConnectionExceptions = new HashSet<>(Arrays.asList(
            "org.jboss.remoting3.NotOpenException"
    ));

    public <T> T exec(Callable<T> task) throws Exception {

        try {
            return task.call();
        } catch (Exception e) {
            if (knownConnectionExceptions.contains(e.getClass().getName())) {
                throw new ConnectException(e.getMessage());
            }
            throw e;
        }
    }
}
