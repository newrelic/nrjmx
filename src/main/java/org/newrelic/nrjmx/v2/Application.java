package org.newrelic.nrjmx.v2;

import org.apache.thrift.TProcessor;
import org.apache.thrift.protocol.TJSONProtocol;
import org.apache.thrift.server.TServer;
import org.apache.thrift.server.TServer.Args;
import org.apache.thrift.transport.TServerTransport;
import org.newrelic.nrjmx.v2.jmx.JMXService;

public class Application {

    public static void main(String[] args) {
        JMXServiceHandler handler = new JMXServiceHandler();
        TProcessor processor = new JMXService.Processor<JMXServiceHandler>(handler);

        TServerTransport serverTransport = new StandardIOTransportServer();
        TServer server = new StandardIOServer(
                new Args(serverTransport).processor(processor).protocolFactory(new TJSONProtocol.Factory()));

        handler.addServer(server);
        server.serve();

        serverTransport.close();
    }
}