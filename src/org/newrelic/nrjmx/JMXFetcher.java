package org.newrelic.nrjmx;

import java.io.IOException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.Map;
import java.util.Set;
import java.util.logging.Logger;

import javax.management.AttributeNotFoundException;
import javax.management.InstanceNotFoundException;
import javax.management.IntrospectionException;
import javax.management.MBeanAttributeInfo;
import javax.management.MBeanException;
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

public class JMXFetcher {
    private static final Logger logger = Logger.getLogger("nrjmx");

    private MBeanServerConnection connection;
    private Map<String, Object> result = new HashMap<>();;


    public class ConnectionError extends Exception {
        public ConnectionError(String message) { super(message); }
    };

    public class QueryError extends Exception {
        public QueryError(String message) { super(message); }
    };

    public class ValueError extends Exception {
        public ValueError(String message) { super(message); }
    };

    public JMXFetcher(String hostname, int port, String username, String password) throws ConnectionError {
        String connectionString = String.format("service:jmx:rmi:///jndi/rmi://%s:%s/jmxrmi", hostname, port);
        Map<String, String[]> env = new HashMap<>();
        if (username != "") {
          env.put(JMXConnector.CREDENTIALS, new String[] { username, password });
        }

        try {
            JMXServiceURL address = new JMXServiceURL(connectionString);
            JMXConnector connector = JMXConnectorFactory.connect(address, env);
            connection = connector.getMBeanServerConnection();
        } catch (IOException e) {
            throw new ConnectionError("Can't connect to JMX server: " + connectionString);
        }
    }

    public Set<ObjectInstance> query(String beanName) throws QueryError {
        ObjectName queryObject;

        try {
            queryObject = new ObjectName(beanName);
        } catch (MalformedObjectNameException e) {
            throw new QueryError("Can't parse bean name " + beanName);
        }

        Set<ObjectInstance> beanInstances;
        try {
            beanInstances = connection.queryMBeans(queryObject, null);
        } catch (IOException e) {
            throw new QueryError("Can't get beans for query " + beanName);
        }
        
        return beanInstances;
    }

    public void queryAttributes(ObjectInstance instance) throws QueryError {
        ObjectName objectName = instance.getObjectName();
        MBeanInfo info;

        try {
            info = connection.getMBeanInfo(objectName);
        } catch (InstanceNotFoundException | IntrospectionException | ReflectionException | IOException e) {
            throw new QueryError("Can't find bean " + objectName.toString());
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
            } catch (AttributeNotFoundException | InstanceNotFoundException | MBeanException | ReflectionException | IOException e) {
                logger.warning("Can't get attribute " + attrName + " for bean " + objectName.toString());
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
        result = new HashMap<String, Object>();
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
