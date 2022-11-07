package org.newrelic.nrjmx.v2;

import org.junit.Test;

public class JMXRequestHandlerTest {

    @Test(timeout = 20_000)
    public void testGetConnectionURL_RMI_With_URIPath() {
        JMXRequestHandler handler = new JMXRequestHandler();
        try {
            handler.exec(() -> {
                throw new Exception();
            });
        } catch (Exception e) {
        }
    }
}
