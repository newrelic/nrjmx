package org.newrelic.jmx;

import javax.management.MBeanRegistration;
import javax.management.MBeanServer;
import javax.management.ObjectName;

public class Cat implements CatMBean, MBeanRegistration {
    private String name;
    private Double doubleValue;
    private Float floatValue;
    private Boolean boolValue;
    private Number numberValue;

    public Cat(String name, Double doubleValue, Float floatValue, Boolean boolValue, Number numberValue) {
        this.name = name;
        this.doubleValue = doubleValue;
        this.floatValue = floatValue;
        this.boolValue = boolValue;
        this.numberValue = numberValue;
    }

    @Override
    public String getName() {
        return name;
    }

    @Override
    public Double getDoubleValue() {
        return doubleValue;
    }

    @Override
    public Float getFloatValue() {
        return floatValue;
    }

    @Override
    public Boolean getBoolValue() {
        return boolValue;
    }


    @Override
    public Number getNumberValue() {
        return numberValue;
    }

    public String toJSON() {
        return "{\"name\":\"" + name + "\"}";
    }

    @Override
    public String toString() {
        return toJSON();
    }

    @Override
    public ObjectName preRegister(MBeanServer server, ObjectName name) throws Exception {
        return new ObjectName("test:type=Cat,name=" + this.name);
    }

    @Override
    public void postRegister(Boolean registrationDone) {
    }

    @Override
    public void preDeregister() {
    }

    @Override
    public void postDeregister() {
    }

}
