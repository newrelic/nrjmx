package org.newrelic.jmx;

import javax.management.*;
import javax.management.openmbean.*;

public class CompositeDataCat implements CompositeDataCatMBean, MBeanRegistration {

    private Cat catInfo;

    public CompositeDataCat(Cat cat) {
        catInfo = cat;
    }

    @Override
    public CompositeData getCatInfo() throws MBeanException {
        try {
            Object[] itemValues = new Object[]{
                    catInfo.getName(),
                    catInfo.getBoolValue(),
                    catInfo.getDoubleValue(),
                    catInfo.getFloatValue(),
                    catInfo.getNumberValue()};
            CompositeDataSupport support = new CompositeDataSupport(getCatDataType(), catAttributeNames(),
                    itemValues);
            return support;
        } catch (OpenDataException e) {
            throw new MBeanException(e, "Error occurred when getting cat information via JMX");
        }
    }

    private static CompositeType getCatDataType() throws OpenDataException {
        CompositeType catDataType = new CompositeType("cat data type",
                "dta type for cat information",
                catAttributeNames(),
                catAttributeDescriptions(),
                catAttributeTypes());
        return catDataType;
    }

    private static String[] catAttributeNames() {

        String[] attributeNames = {
                CompositeDataCat.NAME,
                CompositeDataCat.BOOL,
                CompositeDataCat.DOUBLE,
                CompositeDataCat.FLOAT,
                CompositeDataCat.NUMBER};

        return attributeNames;
    }

    private static String[] catAttributeDescriptions() {
        return catAttributeNames();
    }

    private static OpenType[] catAttributeTypes() {

        OpenType[] attributeTypes = new OpenType[5];

        attributeTypes[0] = SimpleType.STRING;  //for name field
        attributeTypes[1] = SimpleType.BOOLEAN; //bool
        attributeTypes[2] = SimpleType.DOUBLE;  //double
        attributeTypes[3] = SimpleType.FLOAT;   //float
        attributeTypes[4] = SimpleType.INTEGER; //number

        return attributeTypes;
    }


    public String toJSON() {
        return "{\"name\":\"" + catInfo.getName() + "\"}";
    }

    @Override
    public String toString() {
        return toJSON();
    }

    @Override
    public ObjectName preRegister(MBeanServer server, ObjectName name) throws Exception {
        return new ObjectName("test:type=CompositeDataCat,name=" + catInfo.getName());
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