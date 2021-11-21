package org.newrelic.nrjmx.v2;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import org.apache.thrift.transport.TIOStreamTransport;
import org.apache.thrift.transport.TServerTransport;
import org.apache.thrift.transport.TTransport;
import org.apache.thrift.transport.TTransportException;


public class StandardIOTransportServer extends TServerTransport {

    TTransport transport;

    @Override
    public void listen() throws TTransportException {
    }

    @Override
    public void close() {
        transport.close();
    }

    @Override
    protected TTransport acceptImpl() throws TTransportException {
        transport = new TIOStreamTransport(System.in, System.out);
        return transport;
    }

}
