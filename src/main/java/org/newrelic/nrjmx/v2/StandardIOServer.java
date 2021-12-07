package org.newrelic.nrjmx.v2;

import org.apache.thrift.TException;
import org.apache.thrift.TProcessor;
import org.apache.thrift.protocol.TProtocol;
import org.apache.thrift.server.TServer;
import org.apache.thrift.server.TSimpleServer;
import org.apache.thrift.transport.TTransport;
import org.apache.thrift.transport.TTransportException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/*
 * Copyright 2021 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

/**
 * Simple singlethreaded server standard io implementation.
 */
public class StandardIOServer extends TServer {
    private static final Logger LOGGER = LoggerFactory.getLogger(StandardIOServer.class.getName());

    public StandardIOServer(Args args) {
        super(args);
    }

    public void listen() throws Exception {
        try {
            serverTransport_.listen();
        } catch (TTransportException ttx) {
            LOGGER.error("Error occurred during listening.", ttx);
            return;
        }

        setServing(true);

        TTransport inputTransport = null;
        TTransport outputTransport = null;
        try {
            TTransport client = serverTransport_.accept();
            if (client != null) {
                TProcessor processor = processorFactory_.getProcessor(client);

                inputTransport = inputTransportFactory_.getTransport(client);
                outputTransport = outputTransportFactory_.getTransport(client);

                TProtocol inputProtocol = inputProtocolFactory_.getProtocol(inputTransport);
                TProtocol outputProtocol = outputProtocolFactory_.getProtocol(outputTransport);

                while (!stopped_) {
                    processor.process(inputProtocol, outputProtocol);
                }
            }
        } finally {
            if (inputTransport != null) {
                inputTransport.close();
            }

            if (outputTransport != null) {
                outputTransport.close();
            }

            setServing(false);
        }
    }

    @Override
    public void serve() {
        try {
            listen();
        } catch (TTransportException ttx) {
            // Client died, just move on
            LOGGER.debug("Client Transportation Exception", ttx);
        } catch (TException tx) {
            if (!stopped_) {
                LOGGER.error("Thrift error occurred during processing of message.", tx);
            }
        } catch (Exception x) {
            if (!stopped_) {
                LOGGER.error("Error occurred during processing of message.", x);
            }
        }
    }

    public void stop() {
        stopped_ = true;
        serverTransport_.interrupt();
    }
}
