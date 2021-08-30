package org.newrelic.nrjmx.v2;

import java.util.ArrayList;
import java.util.List;

import org.newrelic.nrjmx.v2.jmx.JMXLoggerMessage;
import org.newrelic.nrjmx.v2.jmx.JMXLoggerMessageLevel;

public class Logger {

    private static Logger defaultLogger;

    private List<JMXLoggerMessage> messages;

    public Logger() {
        messages = new ArrayList<>();
    }

    public static Logger getDefaultLogger() {
        if (defaultLogger == null) {
            defaultLogger = new Logger();
        }
        return defaultLogger;
    }

    public List<JMXLoggerMessage> getLogs() {
        List<JMXLoggerMessage> messages = this.messages;
        this.messages = new ArrayList<>();
        return messages;
    }
    
    public void Error(String format, Object... args) {
        JMXLoggerMessage message = new JMXLoggerMessage();
        message.level = JMXLoggerMessageLevel.ERROR;
        message.message = String.format(format, args);
        messages.add(message);
    }

    public void Info(String format, Object... args) {
        JMXLoggerMessage message = new JMXLoggerMessage();
        message.level = JMXLoggerMessageLevel.INFO;
        message.message = String.format(format, args);
        messages.add(message);
    }

    public void Warn(String format, Object... args) {
        JMXLoggerMessage message = new JMXLoggerMessage();
        message.level = JMXLoggerMessageLevel.WARNING;
        message.message = String.format(format, args);
        messages.add(message);
    }
    public void Debug(String format, Object... args) {
        JMXLoggerMessage message = new JMXLoggerMessage();
        message.level = JMXLoggerMessageLevel.DEBUG;
        message.message = String.format(format, args);
        messages.add(message);
    }
}
