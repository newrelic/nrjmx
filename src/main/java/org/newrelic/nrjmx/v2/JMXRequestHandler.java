package org.newrelic.nrjmx.v2;

import java.util.concurrent.Callable;

/**
 * JMXRequestHandler wraps a JMX request and execute the OnExceptionHandler whenever an
 * Exception is thrown.
 */
public class JMXRequestHandler {

    /**
     * OnExceptionHandler interface provides the behavior for the Exception handler.
     */
    public interface OnExceptionHandler {
        void handle(Exception e) throws Exception;
    }

    public JMXRequestHandler(OnExceptionHandler handler) {
        this.onExceptionHandler = handler;
    }

    /** onExceptionHandler to be called whenever an exception is captured. */
    private OnExceptionHandler onExceptionHandler;

    /**
     * exec performs the task.
     * @param task
     * @return <T>
     * @param <T> the task to be called.
     * @throws Exception
     */
    public <T> T exec(Callable<T> task) throws Exception {
        try {
            return task.call();
        } catch (Exception e) {
            if (onExceptionHandler != null) {
                onExceptionHandler.handle(e);
            }
            throw e;
        }
    }
}
