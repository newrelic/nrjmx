/*
 * Copyright 2021 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package org.newrelic.nrjmx.v2;

import java.util.List;

import org.apache.thrift.TException;
import org.apache.thrift.server.TServer;
import org.newrelic.nrjmx.v2.nrprotocol.*;

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
    public String getClientVersion() {
        return jmxFetcher.getVersion();
    }

    @Override
    public void connect(JMXConfig config) throws TException {
        this.requestTimeoutMs = config.requestTimeoutMs;
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
    public List<String> queryMBeanNames(String mBeanNamePattern) throws TException {
        return jmxFetcher.queryMBeanNames(mBeanNamePattern, requestTimeoutMs);
    }

    @Override
    public List<String> getMBeanAttributeNames(String mBeanName) throws TException {
        return jmxFetcher.getMBeanAttributeNames(mBeanName, requestTimeoutMs);
    }

    @Override
    public List<AttributeResponse> getMBeanAttributes(String mBeanName, List<String> attributes) throws TException {
        return jmxFetcher.getMBeanAttributes(mBeanName, attributes, requestTimeoutMs);
    }

    @Override
    public List<AttributeResponse> queryMBeanAttributes(String mBeanNamePattern, List<String> attributes) throws TException {
        return jmxFetcher.queryMBeanAttributes(mBeanNamePattern, attributes, requestTimeoutMs);
    }

    public void addServer(TServer server) {
        this.server = server;
    }
}