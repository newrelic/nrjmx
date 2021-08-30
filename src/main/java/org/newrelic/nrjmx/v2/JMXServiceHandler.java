package org.newrelic.nrjmx.v2;

import java.util.List;

import org.apache.thrift.TException;
import org.apache.thrift.server.TServer;
import org.newrelic.nrjmx.v2.jmx.JMXAttribute;
import org.newrelic.nrjmx.v2.jmx.JMXConfig;
import org.newrelic.nrjmx.v2.jmx.JMXLoggerMessage;
import org.newrelic.nrjmx.v2.jmx.JMXService;

public class JMXServiceHandler implements JMXService.Iface {

    private JMXFetcher jmxFetcher;
    private TServer server;

    @Override
    public boolean connect(JMXConfig config) throws TException {
        this.jmxFetcher = new JMXFetcher(config);
        return jmxFetcher.connect();
    }

    @Override
    public List<JMXAttribute> queryMbean(String beanName) throws TException {
        return jmxFetcher.queryMbean(beanName);
    }

    @Override
    public void disconnect() throws TException {
        if (this.server == null) {
            throw new TException("cannot disconnect, server handler null");
        }
        this.server.stop();
    }

    public void addServer(TServer server) {
        this.server = server;
    }

    @Override
    public List<JMXLoggerMessage> getLogs() throws TException {
        return Logger.getDefaultLogger().getLogs();
    }

}