package org.newrelic.jmx;

import javax.management.MBeanException;
import javax.management.openmbean.CompositeData;

public interface CompositeDataCatMBean {
    String NAME = "name";
    String BOOL = "bool";
    String DOUBLE = "double";
    String FLOAT = "float";
    String NUMBER = "number";

    CompositeData getCatInfo() throws MBeanException;
}
