package org.newrelic.jmx;

import java.io.NotSerializableException;

public interface ExceptionalCatMBean {
    ExceptionalCat.NotSerializable getNotSerializable() throws NotSerializableException;
    Double getDoubleValue();
}
