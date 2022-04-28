package org.newrelic.jmx;

import javax.management.MBeanRegistration;
import javax.management.MBeanServer;
import javax.management.ObjectName;
import java.util.List;
import java.util.Map;

public class ComplexDataCat implements ComplexDataCatMBean, MBeanRegistration {

    private String name;
    private List<Double> list;
    private Map<String, Double> map;

    public ComplexDataCat(String name, List<Double> catAttrs, Map<String, Double> catProps) {
        this.name = name;
        this.list = catAttrs;
        this.map = catProps;
    }

    @Override
    public List<Double> getList() {
        return list;
    }

    @Override
    public Map<String, Double> getMap() {
        return map;
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
        return new ObjectName("test:type=ComplexDataCat,name=" + this.name);
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
