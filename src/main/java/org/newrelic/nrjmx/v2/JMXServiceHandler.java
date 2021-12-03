package org.newrelic.nrjmx.v2;

import java.util.List;

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
        if (server == null) {
            throw new TException("cannot disconnect, server handler null");
        }
        server.stop();
    }

    @Override
    public List<String> getMBeanNames(String mBeanNamePattern, long timeoutMs) throws TException {
        return jmxFetcher.getMBeanNames(mBeanNamePattern, timeoutMs);
    }

    @Override
    public List<String> getMBeanAttrNames(String mBeanName, long timeoutMs) throws TException {
        return jmxFetcher.getMBeanAttrNames(mBeanName, timeoutMs);
    }

    @Override
    public JMXAttribute getMBeanAttr(String mBeanName, String attrName, long timeoutMs) throws TException {
        return jmxFetcher.getMBeanAttr(mBeanName, attrName, timeoutMs);
    }

    public void addServer(TServer server) {
        this.server = server;
    }
}