package org.newrelic.nrjmx.v2;

import org.junit.Assert;
import org.junit.Test;

import java.rmi.MarshalException;
import java.util.Arrays;
import java.util.HashSet;
import java.util.InputMismatchException;

public class JMXRequestHandlerTest {

    @Test(timeout = 20_000)
    public void testJMXRequestHandler() {
        JMXRequestHandler handler = new JMXRequestHandler();
        try {
            handler.exec(() -> {
                throw new InputMismatchException();
            }, new HashSet<>(Arrays.asList("java.util.InputMismatchException")));
        } catch (Exception e) {
            Assert.assertEquals("Exception wrapping didn't work as expected.", "java.rmi.ConnectException", e.getClass().getName());
        }

        try {
            handler.exec(() -> {
                throw new MarshalException("");
            }, new HashSet<>(Arrays.asList("java.util.InputMismatchException")));
        } catch (Exception e) {
            Assert.assertEquals("Exception wrapping should be skipped.", "java.rmi.MarshalException", e.getClass().getName());
        }
    }
}
