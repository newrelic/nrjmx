package org.newrelic.jmx;

import javax.management.MBeanRegistration;
import javax.management.MBeanServer;
import javax.management.ObjectName;

public class CustomCat extends Cat {
    private String mBeanName;

    public CustomCat(String mBeanName, String name, Double doubleValue, Float floatValue, Boolean boolValue, Number numberValue, Integer timeout, long dateValue) {
        super(name, doubleValue, floatValue, boolValue, numberValue, timeout, dateValue);
        this.mBeanName = mBeanName;
    }

    @Override
    public ObjectName preRegister(MBeanServer server, ObjectName name) throws Exception {
        return new ObjectName(this.mBeanName);
    }
}
