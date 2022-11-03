package org.newrelic.nrjmx.v2;

import java.util.concurrent.Callable;

public class JMXRequestHandler {

    private long lastOkTimestamp = System.currentTimeMillis();

    public <T> T exec(Callable<T> task, Callable<Void> onError) throws Exception {

        try {
            T result = task.call();
            lastOkTimestamp = System.currentTimeMillis();
            return result;
        } catch (Exception e) {
            if (System.currentTimeMillis() - lastOkTimestamp > 5 * 60 * 1000) {
                onError.call();
            }
            throw e;
        }
    }
}
