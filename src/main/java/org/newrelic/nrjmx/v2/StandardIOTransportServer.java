package org.newrelic.nrjmx.v2;

import org.apache.thrift.transport.TIOStreamTransport;
import org.apache.thrift.transport.TServerTransport;
import org.apache.thrift.transport.TTransport;
import org.apache.thrift.transport.TTransportException;

/*
 * Copyright 2021 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

/**
 * StandardIOTransportServer is a TServerTransport implementation for stdin/stdout communication.
 */
public class StandardIOTransportServer extends TServerTransport {

    TTransport transport;

    @Override
    public void listen() throws TTransportException {
    }

    @Override
    public void close() {
        if (transport != null) {
            transport.close();
        }
    }

    @Override
    public TTransport accept() throws TTransportException {
        transport = new TIOStreamTransport(System.in, System.out);
        return transport;
    }
}
