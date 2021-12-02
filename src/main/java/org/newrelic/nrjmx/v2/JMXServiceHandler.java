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

    public void addServer(TServer server) {
        this.server = server;
    }

    @Override
    public void connect(JMXConfig config, long timeoutMs) throws JMXConnectionError, JMXError, TException {
        jmxFetcher.connect(config, timeoutMs);
    }

    @Override
    public void disconnect() throws JMXError, TException {
        if (this.server == null) {
            throw new TException("cannot disconnect, server handler null");
        }
        this.server.stop();
    }

    @Override
    public List<JMXAttribute> queryMbean(String beanName, long timeoutMs) throws
            JMXConnectionError, JMXError, TException {
        //return jmxFetcher.queryMbean(beanName, timeoutMs);
        return null;
    }
}