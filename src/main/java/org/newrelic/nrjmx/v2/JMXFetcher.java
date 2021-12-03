package org.newrelic.nrjmx.v2;

/*
 * Copyright 2020 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import org.apache.commons.lang3.exception.ExceptionUtils;
import org.newrelic.nrjmx.v2.nrprotocol.*;

import javax.management.*;
import javax.management.openmbean.CompositeData;
import javax.management.remote.JMXConnector;
import javax.management.remote.JMXConnectorFactory;
import javax.management.remote.JMXServiceURL;
import javax.rmi.ssl.SslRMIClientSocketFactory;
import java.io.IOException;
import java.math.BigDecimal;
import java.util.*;
import java.util.concurrent.*;
import java.util.logging.Logger;
import java.util.stream.Collectors;

/**
 * JMXFetcher class reads queries from an InputStream (usually stdin) and sends
 * the results to an OutputStream (usually stdout)
 */
public class JMXFetcher {
    public static final String defaultURIPath = "jmxrmi";
    public static final Boolean defaultJBossModeIsStandalone = false;

    private static final Logger logger = Logger.getLogger("nrjmx");

    private ExecutorService executor;

    private MBeanServerConnection connection;

    private String connectionString;
    private Map<String, Object> connectionEnv = new HashMap<>();

    public JMXFetcher(ExecutorService executor) {
        this.executor = executor;
    }

    public void connect(JMXConfig jmxConfig, long timeoutMs) throws JMXError {
        withTimeout(executor.submit((Callable<Void>) () -> {
            this.connect(jmxConfig);
            return null;
        }), timeoutMs);
    }

    public void connect(JMXConfig jmxConfig) throws JMXConnectionError {

        if (jmxConfig.connectionURL != null && !jmxConfig.connectionURL.equals("")) {
            connectionString = jmxConfig.connectionURL;
        } else {
            // Official doc for remoting v3 is not available, see:
            // - https://developer.jboss.org/thread/196619
            // - http://jbossremoting.jboss.org/documentation/v3.html
            // Some doc on URIS at:
            // -
            // https://github.com/jboss-remoting/jboss-remoting/blob/master/src/main/java/org/jboss/remoting3/EndpointImpl.java#L292-L304
            // - https://stackoverflow.com/questions/42970921/what-is-http-remoting-protocol
            // -
            // http://www.mastertheboss.com/jboss-server/jboss-monitoring/using-jconsole-to-monitor-a-remote-wildfly-server
            String uriPath = jmxConfig.uriPath;
            if (jmxConfig.isRemote) {
                if (defaultURIPath.equals(uriPath)) {
                    uriPath = "";
                } else {
                    uriPath = uriPath.concat("/");
                }
                String remoteProtocol = "remote";
                if (jmxConfig.isJBossStandaloneMode) {
                    remoteProtocol = "remote+http";
                }
                connectionString =
                        String.format("service:jmx:%s://%s:%s%s", remoteProtocol, jmxConfig.hostname, jmxConfig.port, uriPath);
            } else {
                connectionString =
                        String.format("service:jmx:rmi:///jndi/rmi://%s:%s/%s", jmxConfig.hostname, jmxConfig.port, uriPath);
            }
        }

        if (!"".equals(jmxConfig.username)) {
            connectionEnv.put(JMXConnector.CREDENTIALS, new String[]{jmxConfig.username, jmxConfig.password});
        }

        if (!"".equals(jmxConfig.keyStore) && !"".equals(jmxConfig.trustStore)) {
            Properties p = System.getProperties();
            p.put("javax.net.ssl.keyStore", jmxConfig.keyStore);
            p.put("javax.net.ssl.keyStorePassword", jmxConfig.keyStorePassword);
            p.put("javax.net.ssl.trustStore", jmxConfig.trustStore);
            p.put("javax.net.ssl.trustStorePassword", jmxConfig.trustStorePassword);
            connectionEnv.put("com.sun.jndi.rmi.factory.socket", new SslRMIClientSocketFactory());
        }

        try {
            JMXServiceURL address = new JMXServiceURL(connectionString);

            JMXConnector connector = JMXConnectorFactory.connect(address, connectionEnv);

            this.connection = connector.getMBeanServerConnection();
        } catch (Exception e) {
            String message = String.format("Can't connect to JMX server: '%s', error: '%s'", connectionString,
                    e.getMessage());
            throw new JMXConnectionError(1, message);
        }
    }

    public List<String> getMBeanNames(String beanName, long timeoutMs) throws JMXError {
        return this.withTimeout(
                executor.submit(() -> this.getMBeanNames(beanName)),
                timeoutMs
        );
    }

    public List<String> getMBeanNames(String mBeanNamePattern) throws JMXError {
        ObjectName objectName = this.getObjectName(mBeanNamePattern);
        try {
            return getConnection().queryMBeans(objectName, null)
                    .stream()
                    .map(ObjectInstance::getObjectName)
                    .map(ObjectName::toString)
                    .collect(Collectors.toList());
        } catch (IOException ioe) {
            throw new JMXError()
                    .setMessage("can't get beans for query: " + mBeanNamePattern)
                    .setCauseMessage(ioe.getMessage())
                    .setStacktrace(ExceptionUtils.getStackTrace(ioe));
        }
    }

    public List<String> getMBeanAttrNames(String mBeanName, long timeoutMs) throws JMXError {
        return this.withTimeout(
                executor.submit(() -> this.getMBeanAttrNames(mBeanName)),
                timeoutMs
        );
    }

    public List<String> getMBeanAttrNames(String mBeanName) throws JMXError {
        ObjectName objectName = this.getObjectName(mBeanName);
        MBeanInfo info;

        try {
            info = getConnection().getMBeanInfo(objectName);
        } catch (InstanceNotFoundException | IntrospectionException | ReflectionException | IOException e) {
            throw new JMXError()
                    .setMessage("can't find mBean: " + mBeanName)
                    .setCauseMessage(e.getMessage())
                    .setStacktrace(ExceptionUtils.getStackTrace(e));
        }

        return Arrays.stream(info.getAttributes())
                .filter(MBeanAttributeInfo::isReadable)
                .map(MBeanAttributeInfo::getName)
                .collect(Collectors.toList());
    }

    public JMXAttribute getMBeanAttr(String mBeanName, String attrName, long timeoutMs) throws JMXError {
        return this.withTimeout(
                executor.submit(() -> this.getMBeanAttr(mBeanName, attrName)),
                timeoutMs
        );
    }

    public JMXAttribute getMBeanAttr(String mBeanName, String attrName) throws JMXError {
        Object value;
        ObjectName objectName = this.getObjectName(mBeanName);
        try {
            value = getConnection().getAttribute(objectName, attrName);
            if (value instanceof Attribute) {
                Attribute jmxAttr = (Attribute) value;
                value = jmxAttr.getValue();
            }
        } catch (Exception e) {
            throw new JMXError()
                    .setMessage("can't get attribute: " + attrName + " for bean: " + mBeanName + ": ")
                    .setCauseMessage(e.getMessage())
                    .setStacktrace(ExceptionUtils.getStackTrace(e));
        }

        String name = String.format("%s,attr=%s", mBeanName, attrName);
        return parseValue(name, value);
    }

    private ObjectName getObjectName(String mBeanName) throws JMXError {
        try {
            return new ObjectName(mBeanName);
        } catch (MalformedObjectNameException me) {
            throw new JMXError()
                    .setMessage("cannot parse MBean glob pattern: '"+mBeanName+"', valid: 'DOMAIN:BEAN'")
                    .setCauseMessage(me.getMessage())
                    .setStacktrace(ExceptionUtils.getStackTrace(me));
        }
    }

    private <T> T withTimeout(Future<T> future, long timeoutMs) throws JMXError {
        try {
            if (timeoutMs <= 0) {
                return future.get();
            }
            return future.get(timeoutMs, TimeUnit.MILLISECONDS);

        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new JMXError()
                    .setMessage("operation was interrupted " + e.getMessage())
                    .setCauseMessage(e.getMessage())
                    .setStacktrace(ExceptionUtils.getStackTrace(e));
        } catch (TimeoutException e) {
            throw new JMXError()
                    .setMessage("operation timeout exceeded: " + timeoutMs + "ms")
                    .setCauseMessage(e.getMessage())
                    .setStacktrace(ExceptionUtils.getStackTrace(e));
        } catch (ExecutionException e) {
            if (e.getCause() instanceof JMXError) {
                throw (JMXError) e.getCause();
            }
            throw new JMXError()
                    .setMessage("failed to execute operation, error: " + e.getMessage())
                    .setStacktrace(ExceptionUtils.getStackTrace(e));
        }
    }

    private JMXAttribute parseValue(String name, Object value) throws JMXError {
        JMXAttribute attr = new JMXAttribute();
        attr.attribute = name;

        if (value == null) {
            throw new JMXError()
                    .setMessage("found a null value for bean: " + name);
        } else if (value instanceof java.lang.Double) {
            attr.doubleValue = parseDouble((Double) value);
            attr.valueType = ValueType.DOUBLE;
            return attr;
        } else if (value instanceof java.lang.Float) {
            attr.doubleValue = parseFloatToDouble((Float) value);
            attr.valueType = ValueType.DOUBLE;
            return attr;
        } else if (value instanceof Number) {
            attr.intValue = ((Number) value).longValue();
            attr.valueType = ValueType.INT;
            return attr;
        } else if (value instanceof String) {
            attr.stringValue = (String) value;
            attr.valueType = ValueType.STRING;
            return attr;
        } else if (value instanceof Boolean) {
            attr.boolValue = (Boolean) value;
            attr.valueType = ValueType.BOOL;
            return attr;
        } else if (value instanceof CompositeData) {
            CompositeData cdata = (CompositeData) value;
            Set<String> fieldKeys = cdata.getCompositeType().keySet();

            for (String field : fieldKeys) {
                if (field.length() < 1)
                    continue;

                String fieldKey = field.substring(0, 1).toUpperCase() + field.substring(1);
                parseValue(String.format("%s.%s", name, fieldKey), cdata.get(field));
            }
        } else if (value instanceof HashMap) {
            // TODO: Process hashmaps
            // logger.fine("HashMaps are not supported yet: " + name);
            return null;
        } else if (value instanceof ArrayList || value.getClass().isArray()) {
            // TODO: Process arrays
            // logger.fine("Arrays are not supported yet: " + name);
            return null;
        } else {
            throw new JMXError()
                    .setMessage("unsuported data type (" + value.getClass() + ") for bean " + name);
        }
        return null;
    }

    private MBeanServerConnection getConnection() throws JMXError {
        if (this.connection == null) {
            throw new JMXError()
                    .setMessage("connection to JMX endpoint is not established");
        }
        return this.connection;
    }

    /**
     * XXX: JSON does not support NaN, Infinity, or -Infinity as they come back from
     * JMX. So we parse them out to 0, Max Double, and Min Double respectively.
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

    /**
     * XXX: JSON does not support NaN, Infinity, or -Infinity as they come back from
     * JMX. So we parse them out to 0, Max Double, and Min Double respectively.
     */
    private Double parseFloatToDouble(Float value) {
        if (value.isNaN()) {
            return 0.0d;
        } else if (value == Double.NEGATIVE_INFINITY) {
            return Double.MIN_VALUE;
        } else if (value == Double.POSITIVE_INFINITY) {
            return Double.MAX_VALUE;
        }

        return new BigDecimal(value.toString()).doubleValue();
    }

    public boolean StringIsNullOrEmpty(String value) {
        return value == null || value.equals("");
    }
}
