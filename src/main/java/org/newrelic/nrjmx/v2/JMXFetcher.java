package org.newrelic.nrjmx.v2;

/*
 * Copyright 2021 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import org.apache.commons.lang3.exception.ExceptionUtils;
import org.newrelic.nrjmx.Application;
import org.newrelic.nrjmx.v2.nrprotocol.*;

import javax.management.*;
import javax.management.openmbean.CompositeData;
import javax.management.remote.JMXConnector;
import javax.management.remote.JMXConnectorFactory;
import javax.management.remote.JMXServiceURL;
import javax.rmi.ssl.SslRMIClientSocketFactory;
import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.math.BigDecimal;
import java.util.*;
import java.util.concurrent.*;
import java.util.stream.Collectors;

/**
 * JMXFetcher class executes requests to an JMX endpoint.
 */
public class JMXFetcher {
    public static final String defaultURIPath = "jmxrmi";

    /* ExecutorService is required to run JMX requests with timeout. */
    private final ExecutorService executor;

    /* MBeanServerConnection is the connection to JMX endpoint. */
    private MBeanServerConnection connection;

    public JMXFetcher(ExecutorService executor) {
        this.executor = executor;
    }

    /**
     * connect performs the connection to the JMX endpoint.
     *
     * @param jmxConfig JMX configuration.
     * @param timeoutMs long timeout for the request in milliseconds
     * @throws JMXError           JMX related Exception
     * @throws JMXConnectionError JMX connection related exception
     */
    public void connect(JMXConfig jmxConfig, long timeoutMs) throws JMXError, JMXConnectionError {
        withTimeout(executor.submit((Callable<Void>) () -> {
            connect(jmxConfig);
            return null;
        }), timeoutMs);
    }

    /**
     * connect performs the connection to the JMX endpoint.
     *
     * @param jmxConfig JMX configuration.
     * @throws JMXConnectionError JMX connection related exception
     */
    public void connect(JMXConfig jmxConfig) throws JMXConnectionError {
        String connectionString = buildConnectionString(jmxConfig);
        Map<String, Object> connectionEnv = buildConnectionEnvConfig(jmxConfig);

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

    /**
     * getMBeanNames returns all founded mBeans that match the provided pattern.
     *
     * @param mBeanGlobPattern String glob pattern DOMAIN:BEAN e.g *:* or jboss.as:subsystem=remoting,configuration=endpoint
     * @param timeoutMs        long timeout for the request in milliseconds
     * @return List<String> containing all mBean names that were found
     * @throws JMXError           JMX related Exception
     * @throws JMXConnectionError JMX connection related exception
     */
    public List<String> getMBeanNames(String mBeanGlobPattern, long timeoutMs) throws JMXError, JMXConnectionError {
        return this.withTimeout(
                executor.submit(() -> this.getMBeanNames(mBeanGlobPattern)),
                timeoutMs
        );
    }

    /**
     * getMBeanNames returns all founded mBeans that match the provided pattern.
     *
     * @param mBeanGlobPattern String glob pattern DOMAIN:BEAN e.g *:* or jboss.as:subsystem=remoting,configuration=endpoint
     * @return List<String> containing all mBean names that were found
     * @throws JMXError JMX related Exception
     */
    public List<String> getMBeanNames(String mBeanGlobPattern) throws JMXError {
        ObjectName objectName = this.getObjectName(mBeanGlobPattern);
        try {
            return getConnection().queryMBeans(objectName, null)
                    .stream()
                    .map(ObjectInstance::getObjectName)
                    .map(ObjectName::toString)
                    .collect(Collectors.toList());
        } catch (IOException ioe) {
            throw new JMXError()
                    .setMessage("can't get beans for query: " + mBeanGlobPattern)
                    .setCauseMessage(ioe.getMessage())
                    .setStacktrace(ExceptionUtils.getStackTrace(ioe));
        }
    }

    /**
     * getMBeanAttrNames returns all the available JMX attribute names for a given mBeanName.
     *
     * @param mBeanName of which we want to retrieve attributes
     * @param timeoutMs long timeout for the request in milliseconds
     * @return List<String> containing all mBean attribute names that were found
     * @throws JMXError           JMX related Exception
     * @throws JMXConnectionError JMX connection related exception
     */
    public List<String> getMBeanAttrNames(String mBeanName, long timeoutMs) throws JMXError, JMXConnectionError {
        return this.withTimeout(
                executor.submit(() -> this.getMBeanAttrNames(mBeanName)),
                timeoutMs
        );
    }

    /**
     * getMBeanAttrNames returns all the available JMX attribute names for a given mBeanName.
     *
     * @param mBeanName of which we want to retrieve attributes
     * @return List<String> containing all mBean attribute names that were found
     * @throws JMXError JMX related Exception
     */
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

    /**
     * getMBeanAttr returns the attribute value for an mBeanName.
     *
     * @param mBeanName of which we want to retrieve the attribute value
     * @param attrName  of which we want to retrieve the attribute value
     * @param timeoutMs long timeout for the request in milliseconds
     * @return JMXAttribute representing the mBean attribute value
     * @throws JMXError           JMX related Exception
     * @throws JMXConnectionError JMX connection related exception
     */
    public List<JMXAttribute> getMBeanAttrs(String mBeanName, String attrName, long timeoutMs) throws JMXError, JMXConnectionError {
        return this.withTimeout(
                executor.submit(() -> this.getMBeanAttrs(mBeanName, attrName)),
                timeoutMs
        );
    }

    /**
     * getMBeanAttr returns the attribute value for an mBeanName.
     *
     * @param mBeanName of which we want to retrieve the attribute value
     * @param attrName  of which we want to retrieve the attribute value
     * @return JMXAttribute representing the mBean attribute value
     * @throws JMXError JMX related Exception
     */
    public List<JMXAttribute> getMBeanAttrs(String mBeanName, String attrName) throws JMXError {
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

    /**
     * getObjectName returns the ObjectName for an mBeanName required on performing JMX requests.
     *
     * @param mBeanName to build the ObjectName
     * @return ObjectName for the mBeanName
     * @throws JMXError JMX related Exception
     */
    private ObjectName getObjectName(String mBeanName) throws JMXError {
        try {
            return new ObjectName(mBeanName);
        } catch (MalformedObjectNameException me) {
            throw new JMXError()
                    .setMessage("cannot parse MBean glob pattern: '" + mBeanName + "', valid: 'DOMAIN:BEAN'")
                    .setCauseMessage(me.getMessage())
                    .setStacktrace(ExceptionUtils.getStackTrace(me));
        }
    }

    /**
     * withTimeout executes a task with timeout.
     *
     * @param future    is the task that has to be executed
     * @param timeoutMs timeout in milliseconds after which we terminate the task
     * @param <T>       Generic type for the task
     * @return Task result
     * @throws JMXError           JMX related Exception
     * @throws JMXConnectionError JMX connection related exception
     */
    private <T> T withTimeout(Future<T> future, long timeoutMs) throws JMXError, JMXConnectionError {
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
            } else if (e.getCause() instanceof JMXConnectionError) {
                throw (JMXConnectionError) e.getCause();
            }
            throw new JMXError()
                    .setMessage("failed to execute operation, error: " + e.getMessage())
                    .setStacktrace(ExceptionUtils.getStackTrace(e));
        }
    }

    /**
     * parseValue converts the received value from JMX into an JMXAttribute object.
     *
     * @param mBeanAttributeName of the value
     * @param value              that has to be converted
     * @return JMXAttribute containing the mBeanAttributeName and the converted value.
     * @throws JMXError JMX related Exception
     */
    private List<JMXAttribute> parseValue(String mBeanAttributeName, Object value) throws JMXError {
        JMXAttribute attr = new JMXAttribute();
        attr.attribute = mBeanAttributeName;

        if (value == null) {
            throw new JMXError()
                    .setMessage("found a null value for bean: " + mBeanAttributeName);
        } else if (value instanceof java.lang.Double) {
            attr.doubleValue = parseDouble((Double) value);
            attr.valueType = ValueType.DOUBLE;
            return Arrays.asList(attr);
        } else if (value instanceof java.lang.Float) {
            attr.doubleValue = parseFloatToDouble((Float) value);
            attr.valueType = ValueType.DOUBLE;
            return Arrays.asList(attr);
        } else if (value instanceof Number) {
            attr.intValue = ((Number) value).longValue();
            attr.valueType = ValueType.INT;
            return Arrays.asList(attr);
        } else if (value instanceof String) {
            attr.stringValue = (String) value;
            attr.valueType = ValueType.STRING;
            return Arrays.asList(attr);
        } else if (value instanceof Boolean) {
            attr.boolValue = (Boolean) value;
            attr.valueType = ValueType.BOOL;
            return Arrays.asList(attr);
        } else if (value instanceof CompositeData) {
            List<JMXAttribute> result = new ArrayList<>();
            CompositeData cdata = (CompositeData) value;
            Set<String> fieldKeys = cdata.getCompositeType().keySet();

            for (String field : fieldKeys) {
                if (field.length() < 1)
                    continue;

                String fieldKey = field.substring(0, 1).toUpperCase() + field.substring(1);
                result.addAll(parseValue(String.format("%s.%s", mBeanAttributeName, fieldKey), cdata.get(field)));
            }
            return result;
        } else {
            throw new JMXError()
                    .setMessage("unsuported data type (" + value.getClass() + ") for bean " + mBeanAttributeName);
        }
    }

    /**
     * getConnection returns the connection the the JMX endpoint.
     *
     * @return MBeanServerConnection the connection to the JMX endpoint
     * @throws JMXError JMX related Exception
     */
    private MBeanServerConnection getConnection() throws JMXError {
        if (this.connection == null) {
            throw new JMXError()
                    .setMessage("connection to JMX endpoint is not established");
        }
        return this.connection;
    }

    /**
     * parseDouble ensures the value has the expected format.
     * We do not support NaN, Infinity, or -Infinity as they come back from
     * JMX. So we parse them out to 0, Max Double, and Min Double respectively.
     *
     * @param value to be parsed
     * @return Double parsed value.
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
     * parseFloatToDouble ensures the value has the expected format.
     * We do not support NaN, Infinity, or -Infinity as they come back from
     * JMX. So we parse them out to 0, Max Double, and Min Double respectively.
     *
     * @param value to be parsed
     * @return Double parsed value.
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

    /**
     * buildConnectionString is used to build the connection URL using the JMXConfig.
     *
     * @param jmxConfig JMX configuration.
     * @return String containing the connection URL.
     */
    public static String buildConnectionString(JMXConfig jmxConfig) {
        if (jmxConfig.connectionURL != null && !jmxConfig.connectionURL.equals("")) {
            return jmxConfig.connectionURL;
        }
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
        if (uriPath == null) {
            uriPath = defaultURIPath;
        }
        if (jmxConfig.isRemote) {
            if (defaultURIPath.equals(uriPath)) {
                uriPath = "";
            } else {
                uriPath = uriPath.concat("/");
            }
            String remoteProtocol = "remote";
            if (jmxConfig.isJBossStandaloneMode) {
                remoteProtocol += (jmxConfig.useSSL) ? "+https" : "+http";
            }
            return String.format("service:jmx:%s://%s:%s%s", remoteProtocol, jmxConfig.hostname, jmxConfig.port, uriPath);
        }
        return String.format("service:jmx:rmi:///jndi/rmi://%s:%s/%s", jmxConfig.hostname, jmxConfig.port, uriPath);
    }

    /**
     * buildConnectionEnvConfig creates a Map containing the environment options required for JMX.
     * based on received JMXConfig
     *
     * @param jmxConfig JMX configuration.
     * @return Map<String, Object> containing the environment options required for JMX
     */
    private static Map<String, Object> buildConnectionEnvConfig(JMXConfig jmxConfig) {
        Map<String, Object> connectionEnv = new HashMap<>();
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
        return connectionEnv;
    }

    public String getVersion() {
        try {
            InputStream inputStream = getClass().getClassLoader().getResourceAsStream("version");
            InputStreamReader inputStreamReader = new InputStreamReader(Optional.ofNullable(inputStream).orElseThrow(IOException::new));
            try (BufferedReader reader = new BufferedReader(inputStreamReader)) {
                return reader.readLine();
            } finally {
                inputStream.close();
                inputStreamReader.close();
            }
        } catch (Exception e) {
            return "unknown";
        }
    }
}
