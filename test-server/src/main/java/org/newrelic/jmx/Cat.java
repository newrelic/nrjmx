package org.newrelic.jmx;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.management.MBeanRegistration;
import javax.management.MBeanServer;
import javax.management.ObjectName;
import java.util.Date;

public class Cat implements CatMBean, MBeanRegistration {
    private String name;
    private Double doubleValue;
    private Float floatValue;
    private Boolean boolValue;
    private Number numberValue;
    private Integer timeout;
    private long dateValue;

    public Cat(String name, Double doubleValue, Float floatValue, Boolean boolValue, Number numberValue, Integer timeout, long dateValue) {
        this.name = name;
        this.doubleValue = doubleValue;
        this.floatValue = floatValue;
        this.boolValue = boolValue;
        this.numberValue = numberValue;
        this.timeout = timeout;
        this.dateValue = dateValue;
    }

    private void delay() {
        final Logger log = LoggerFactory.getLogger("org.newrelic");

        if (this.timeout == null) {
            return;
        }
        try {
            log.info("delaying");
            Thread.sleep(this.timeout);
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
    }

    @Override
    public String getName() {
        return name;
    }

    @Override
    public Date getDateValue() {
        if (this.dateValue == 0) {
            return null;
        }
        return new Date(this.dateValue);
    }

    @Override
    public Double getDoubleValue() {
        delay();
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
