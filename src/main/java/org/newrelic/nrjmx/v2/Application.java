package org.newrelic.nrjmx.v2;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

import com.fasterxml.jackson.annotation.JsonFormat.Feature;
import com.fasterxml.jackson.core.JsonParseException;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.DeserializationContext;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.JsonMappingException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.deser.std.StdDeserializer;
import com.fasterxml.jackson.databind.module.SimpleModule;

import org.apache.thrift.TProcessor;
import org.apache.thrift.protocol.TJSONProtocol;
import org.apache.thrift.server.TServer;
import org.apache.thrift.server.TSimpleServer;
import org.apache.thrift.server.TServer.Args;
import org.apache.thrift.transport.TServerTransport;
import org.newrelic.nrjmx.v2.jmx.JMXService;

public class Application {

    public static String json = "[{\"__typename\": \"A\",\"val\": 2},{\"__typename\": \"B\",\"val\": 3}]";

    public static interface I {
        // int getVal();
    }

    public static class A implements I {
        // public String type;
        public int val;
    }

    public static class B implements I {
        public int val;
    }

    public static class IDeserialize extends StdDeserializer<I> {

        public IDeserialize() {
            super(I.class);
        }

        // protected IDeserialize(Class<?> vc) {
        // super(vc);
        // }

        @Override
        public I deserialize(com.fasterxml.jackson.core.JsonParser p, DeserializationContext ctxt)
                throws IOException, JsonProcessingException {
            ObjectMapper om = new ObjectMapper().configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false);
            // TODO Auto-generated method stub
            JsonNode node = p.readValueAsTree();
            // String a = node.get("a").asText();
            String type = node.get("__typename").asText();
            switch (type) {
                case "A":
                    return om.readValue(node.toString(), A.class);
                case "B":
                    return om.readValue(node.toString(), A.class);
                default:
                    return null;

            }
        }
    }

    public static void main(String[] args) throws JsonParseException, JsonMappingException, IOException {
        SimpleModule deserialization = new SimpleModule();
        deserialization.addDeserializer(I.class, new IDeserialize());

        ObjectMapper om = new ObjectMapper();
        om.registerModule(deserialization);

        // ArrayList<I> a = om.readValue(json, new TypeReference<ArrayList<I>>() {

        // });

        List<I> a = om.readValue(json, om.getTypeFactory().constructCollectionType(List.class, I.class));

        // for (Object o : a) {
        // A x = (A) o;
        // System.out.println(x.val);
        // }
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