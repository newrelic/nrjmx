package org.newrelic.nrjmx.v2;

import java.util.List;

import org.apache.thrift.TException;
import org.apache.thrift.server.TServer;
import org.newrelic.nrjmx.v2.nrprotocol.*;

/*
 * Copyright 2021 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

/**
 * JMXServiceHandler is the implementation for nrjmx thrift service.
 */
public class JMXServiceHandler implements JMXService.Iface {

    private final JMXFetcher jmxFetcher;
    private TServer server;
    private long requestTimeoutMs = 0;

    public JMXServiceHandler(JMXFetcher jmxFetcher) {
        this.jmxFetcher = jmxFetcher;
    }

    @Override
    public void ping() throws TException {
    }

    @Override
    public void connect(JMXConfig config) throws TException {
        this.requestTimeoutMs = config.requestTimoutMs;
        jmxFetcher.connect(config, requestTimeoutMs);
    }

    @Override
    public void disconnect() throws TException {
        if (server == null) {
            throw new TException("cannot disconnect, nrjmx handler null");
        }
        server.stop();
    }

    @Override
    public List<String> getMBeanNames(String mBeanNamePattern) throws TException {
        return jmxFetcher.getMBeanNames(mBeanNamePattern, requestTimeoutMs);
    }

    @Override
    public List<String> getMBeanAttrNames(String mBeanName) throws TException {
        return jmxFetcher.getMBeanAttrNames(mBeanName, requestTimeoutMs);
    }

    @Override
    public JMXAttribute getMBeanAttr(String mBeanName, String attrName) throws TException {
        return jmxFetcher.getMBeanAttr(mBeanName, attrName, requestTimeoutMs);
    }

    public void addServer(TServer server) {
        this.server = server;
    }
}