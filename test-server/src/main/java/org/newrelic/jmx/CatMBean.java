package org.newrelic.jmx;

public interface CatMBean {
    String getName();
    Double getDoubleValue();
    Float getFloatValue();
    Boolean getBoolValue();
    Number getNumberValue();
}
