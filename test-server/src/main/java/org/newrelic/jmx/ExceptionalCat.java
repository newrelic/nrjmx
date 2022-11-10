package org.newrelic.jmx;

import javax.management.MBeanRegistration;
import javax.management.MBeanServer;
import javax.management.ObjectName;

public class ExceptionalCat implements ExceptionalCatMBean, MBeanRegistration {

    private String name;

    private Double doubleValue;

    /* Dummy class to mimic NotSerializableException. */
    public static class NotSerializable {
    }

    public ExceptionalCat(String name, Double doubleValue) {
        this.name = name;
        this.doubleValue = doubleValue;
    }

    @Override
    public NotSerializable getNotSerializable() {
        return new NotSerializable();
    }

    @Override
    public Double getDoubleValue() {
        return doubleValue;
    }

    @Override
    public ObjectName preRegister(MBeanServer server, ObjectName name) throws Exception {
        return new ObjectName("test:type=ExceptionalCat,name=" + this.name);
    }

    @Override
    public void postRegister(Boolean registrationDone) {
    }

    @Override
    public void preDeregister() throws Exception {
    }

    @Override
    public void postDeregister() {
    }
}
