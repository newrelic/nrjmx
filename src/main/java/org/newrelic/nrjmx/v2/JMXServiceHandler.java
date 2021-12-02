package org.newrelic.nrjmx.v2;

import java.util.List;
import java.util.concurrent.ExecutorService;

import org.apache.thrift.TException;
import org.apache.thrift.server.TServer;
import org.newrelic.nrjmx.v2.nrprotocol.*;

public class JMXServiceHandler implements JMXService.Iface {

    private JMXFetcher jmxFetcher;
    private TServer server;

    public JMXServiceHandler(JMXFetcher jmxFetcher) {
        this.jmxFetcher = jmxFetcher;
    }

    @Override
    public void ping() throws TException {
        return;
    }

    @Override
    public void connect(JMXConfig config, long timeoutMs) throws TException {
        jmxFetcher.connect(config, timeoutMs);
    }

    @Override
    public void disconnect() throws TException {
        if (this.server == null) {
            throw new TException("cannot disconnect, server handler null");
        }
        this.server.stop();
    }

    @Override
    public List<JMXAttribute> queryMbean(String beanName, long timeoutMs) throws TException {
        return jmxFetcher.queryMbean(beanName, timeoutMs);
    }

    public void addServer(TServer server) {
        this.server = server;
    }
}