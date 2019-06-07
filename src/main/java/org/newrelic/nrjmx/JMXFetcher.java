package org.newrelic.nrjmx;

import java.io.IOException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.Map;
import java.util.Properties;
import java.util.Set;
import java.util.logging.Logger;

import javax.management.Attribute;
import javax.management.InstanceNotFoundException;
import javax.management.IntrospectionException;
import javax.management.MBeanAttributeInfo;
import javax.management.MBeanInfo;
import javax.management.MBeanServerConnection;
import javax.management.MalformedObjectNameException;
import javax.management.ObjectInstance;
import javax.management.ObjectName;
import javax.management.ReflectionException;
import javax.management.openmbean.CompositeData;
import javax.management.remote.JMXConnector;
import javax.management.remote.JMXConnectorFactory;
import javax.management.remote.JMXServiceURL;
import javax.rmi.ssl.SslRMIClientSocketFactory;

public class JMXFetcher {
    private static final Logger logger = Logger.getLogger("nrjmx");

    private MBeanServerConnection connection;
    private Map<String, Object> result = new HashMap<>();

    public class ConnectionError extends Exception {
        public ConnectionError(String message, Exception cause) {
            super(message, cause);
        }
    }

    public class QueryError extends Exception {
        public QueryError(String message, Exception cause) {
            super(message, cause);
        }
    }

    public class ValueError extends Exception {
        public ValueError(String message) {
            super(message);
        }
    }

    public JMXFetcher(String hostname, int port, String uriPath, String username, String password , String keyStore, String keyStorePassword, String trustStore, String trustStorePassword, boolean isRemote) throws ConnectionError {
        String connectionString = String.format("service:jmx:rmi:///jndi/rmi://%s:%s/%s", hostname, port, uriPath);
        if (isRemote) {
            connectionString = String.format("service:jmx:remoting-jmx://%s:%s", hostname, port);
        }

        Map<String, Object> env = new HashMap<>();
        if (!"".equals(username)) {
            env.put(JMXConnector.CREDENTIALS, new String[]{username, password});
        }

        if (!"".equals(keyStore) && !"".equals(trustStore)) {
            Properties p = System.getProperties();
            p.put("javax.net.ssl.keyStore", keyStore);
            p.put("javax.net.ssl.keyStorePassword", keyStorePassword);
            p.put("javax.net.ssl.trustStore", trustStore);
            p.put("javax.net.ssl.trustStorePassword", trustStorePassword);
            env.put("com.sun.jndi.rmi.factory.socket", new SslRMIClientSocketFactory());
        }

        try {
            JMXServiceURL address = new JMXServiceURL(connectionString);
            JMXConnector connector = JMXConnectorFactory.connect(address, env);
            connection = connector.getMBeanServerConnection();
        } catch (IOException e) {
            throw new ConnectionError("Can't connect to JMX server: " + connectionString, e);
        }
    }

    public Set<ObjectInstance> query(String beanName) throws QueryError {
        ObjectName queryObject;

        try {
            queryObject = new ObjectName(beanName);
        } catch (MalformedObjectNameException e) {
            throw new QueryError("Can't parse bean name " + beanName, e);
        }

        Set<ObjectInstance> beanInstances;
        try {
            beanInstances = connection.queryMBeans(queryObject, null);
        } catch (IOException e) {
            throw new QueryError("Can't get beans for query " + beanName, e);
        }

        return beanInstances;
    }

    public void queryAttributes(ObjectInstance instance) throws QueryError {
        ObjectName objectName = instance.getObjectName();
        MBeanInfo info;

        try {
            info = connection.getMBeanInfo(objectName);
        } catch (InstanceNotFoundException | IntrospectionException | ReflectionException | IOException e) {
            throw new QueryError("Can't find bean " + objectName.toString(), e);
        }

        MBeanAttributeInfo[] attrInfo = info.getAttributes();

        for (MBeanAttributeInfo attr : attrInfo) {
            if (!attr.isReadable()) {
                continue;
            }

            String attrName = attr.getName();
            Object value;

            try {
                value = connection.getAttribute(objectName, attrName);
                if (value instanceof javax.management.Attribute) {
                	Attribute jmxAttr = (Attribute) value;
                	value = jmxAttr.getValue();
                }
            } catch (Exception e) {
                logger.warning("Can't get attribute " + attrName + " for bean " + objectName.toString() + ": " + e.getMessage());
                continue;
            }

            String name = String.format("%s,attr=%s", objectName.toString(), attrName);
            try {
                parseValue(name, value);
            } catch (ValueError e) {
                logger.fine(e.getMessage());
            }
        }
    }

    public Map<String, Object> popResults() {
        Map<String, Object> out = result;
        result = new HashMap<>();
        return out;
    }

    private void parseValue(String name, Object value) throws ValueError {
        if (value == null) {
            throw new ValueError("Found a null value for bean " + name);
        } else if (value instanceof java.lang.Double) {
            Double ddata = parseDouble((Double) value);
            result.put(name, ddata);
        } else if (value instanceof Number || value instanceof String || value instanceof Boolean) {
            result.put(name, value);
        } else if (value instanceof CompositeData) {
            CompositeData cdata = (CompositeData) value;
            Set<String> fieldKeys = cdata.getCompositeType().keySet();

            for (String field : fieldKeys) {
                if (field.length() < 1) continue;

                String fieldKey = field.substring(0, 1).toUpperCase() + field.substring(1);
                parseValue(String.format("%s.%s", name, fieldKey), cdata.get(field));
            }
        } else if (value instanceof HashMap) {
            // TODO: Process hashmaps
            logger.fine("HashMaps are not supported yet: " + name);
        } else if (value instanceof ArrayList || value.getClass().isArray()) {
            // TODO: Process arrays
            logger.fine("Arrays are not supported yet: " + name);
        } else {
            throw new ValueError("Unsuported data type (" + value.getClass() + ") for bean " + name);
        }
    }

    /**
     * XXX: JSON does not support NaN, Infinity, or -Infinity as they come back from JMX.
     * So we parse them out to 0, Max Double, and Min Double respectively.
     */
    private Double parseDouble(Double value) {
        if (value.isNaN()) {
            return 0.0;
        } else if (value == Double.NEGATIVE_INFINITY) {
            return Double.MIN_VALUE;
        } else if (value == Double.POSITIVE_INFINITY) {
            return Double.MAX_VALUE;
        }

        return value;
    }
}
