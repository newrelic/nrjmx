package org.newrelic.jmx;

import javax.management.MBeanException;
import java.util.List;
import java.util.Map;

public interface ComplexDataCatMBean {
    List<Double> getList() throws MBeanException;

    Map<String, Double> getMap() throws MBeanException;
}
