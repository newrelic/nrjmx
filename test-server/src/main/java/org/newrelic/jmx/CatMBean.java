package org.newrelic.jmx;

import java.util.Date;

public interface CatMBean {
    String getName();

    Double getDoubleValue();

    Float getFloatValue();

    Boolean getBoolValue();

    Number getNumberValue();

    Date getDateValue();
}
