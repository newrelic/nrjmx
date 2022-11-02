package org.newrelic.jmx;

import javax.management.MBeanServer;
import javax.management.ObjectName;
import java.io.IOException;
import java.io.NotSerializableException;

public class ExceptionalCat extends Cat {
    private String mBeanName;
    private String exception;

    public ExceptionalCat(String mBeanName, String name, Double doubleValue, Float floatValue, Boolean boolValue, Number numberValue, Integer timeout, long dateValue) {
        super(name, doubleValue, floatValue, boolValue, numberValue, timeout, dateValue);
        this.mBeanName = mBeanName;
    }

    @Override
    public ObjectName preRegister(MBeanServer server, ObjectName name) throws Exception {
        return new ObjectName(this.mBeanName);
    }

    public String getException() throws NotSerializableException {
        throw new NotSerializableException();
    }
}
