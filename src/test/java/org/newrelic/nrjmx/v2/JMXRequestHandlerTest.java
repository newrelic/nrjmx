package org.newrelic.nrjmx.v2;

import org.junit.Assert;
import org.junit.Test;

import java.rmi.ConnectException;
import java.rmi.MarshalException;
import java.util.InputMismatchException;

public class JMXRequestHandlerTest {

    @Test(timeout = 20_000)
    public void testJMXRequestHandler() {
        JMXRequestHandler handler = new JMXRequestHandler(e -> {
            if ("java.util.InputMismatchException".equals(e.getClass().getName())) {
                throw new ConnectException(e.getMessage());
            }
        });
        try {
            handler.exec(() -> {
                throw new InputMismatchException();
            });
        } catch (Exception e) {
            Assert.assertEquals("Exception wrapping didn't work as expected.", "java.rmi.ConnectException", e.getClass().getName());
        }

        handler = new JMXRequestHandler(e -> {
            if ("java.util.InputMismatchException".equals(e.getClass().getName())) {
                throw new ConnectException(e.getMessage());
            }
        });

        try {
            handler.exec(() -> {
                throw new MarshalException("");
            });
        } catch (Exception e) {
            Assert.assertEquals("Exception wrapping should be skipped.", "java.rmi.MarshalException", e.getClass().getName());
        }
    }
}
