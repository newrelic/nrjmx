package org.newrelic.nrjmx.v2;

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import org.apache.thrift.TException;
import org.apache.thrift.TProcessor;
import org.apache.thrift.protocol.TProtocol;
import org.apache.thrift.server.ServerContext;
import org.apache.thrift.server.TServer;
import org.apache.thrift.transport.TTransport;
import org.apache.thrift.transport.TTransportException;
import org.apache.thrift.transport.layered.TFramedTransport;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.BufferedWriter;
import java.io.FileWriter;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.nio.file.StandardOpenOption;


public class StandardIOServer extends TServer {
    public void write(String message) {

        try {
            message+="\n";
            Files.write(Paths.get("/Users/cciutea/workspace/nr/int/nrjmx/gojmx/cmd/out"), message.getBytes(), StandardOpenOption.APPEND);
        }catch (IOException e) {
            //exception handling left as an exercise for the reader
        }
    }

    /**
     * Simple singlethreaded server implementation.
     */
    private static final Logger LOGGER = LoggerFactory.getLogger(StandardIOServer.class.getName());

    public StandardIOServer(Args args) {
        super(args);
    }

    public void serve() {

        try {
            serverTransport_.listen();
        } catch (TTransportException ttx) {
            LOGGER.error("Error occurred during listening.", ttx);
            return;
        }

        // Run the preServe event
        if (eventHandler_ != null) {
            eventHandler_.preServe();
        }

        setServing(true);

        while (!stopped_) {

            TTransport client = null;
            TProcessor processor = null;
            TTransport inputTransport = null;
            TTransport outputTransport = null;
            TProtocol inputProtocol = null;
            TProtocol outputProtocol = null;
            ServerContext connectionContext = null;
            try {
                client = serverTransport_.accept();
                write("accepted!");
                if (client != null) {
                    processor = processorFactory_.getProcessor(client);

                    inputTransport = new TFramedTransport(inputTransportFactory_.getTransport(client), 8192);
                    outputTransport = new TFramedTransport(outputTransportFactory_.getTransport(client), 8192);
                    inputProtocol = inputProtocolFactory_.getProtocol(inputTransport);
                    outputProtocol = outputProtocolFactory_.getProtocol(outputTransport);
                    if (eventHandler_ != null) {
                        connectionContext = eventHandler_.createContext(inputProtocol, outputProtocol);
                    }
                    while (!stopped_) {
                        write("handle");

                        try {
                            if (eventHandler_ != null) {
                                eventHandler_.processContext(connectionContext, inputTransport, outputTransport);
                            }
                            processor.process(inputProtocol, outputProtocol);
                        } catch (Exception e) {
                            write("exception!");
                            this.stop();

                            return;
                        }
                    }
                }
            } catch (TTransportException ttx) {
                write("TTransportException!");

                // Client died, just move on
                LOGGER.debug("Client Transportation Exception", ttx);
                break;
            } catch (TException tx) {
                write("TException!");

                if (!stopped_) {
                    LOGGER.error("Thrift error occurred during processing of message.", tx);
                }
            } catch (Exception x) {
                write("Exception!");

                if (!stopped_) {
                    LOGGER.error("Error occurred during processing of message.", x);
                }
                break;
            }

            if (eventHandler_ != null) {
                eventHandler_.deleteContext(connectionContext, inputProtocol, outputProtocol);
            }
        }
        setServing(false);
    }

    public void stop() {
        stopped_ = true;
        serverTransport_.interrupt();
    }
}
