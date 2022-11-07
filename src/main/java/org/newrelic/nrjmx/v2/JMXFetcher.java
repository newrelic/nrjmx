/*
 * Copyright 2021 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package org.newrelic.nrjmx.v2;

import org.apache.commons.lang3.exception.ExceptionUtils;
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
import java.rmi.ConnectException;
import java.text.DateFormat;
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

    /* JMXConnector is used to connect to JMX endpoint. */
    private JMXConnector connector;

    /* MBeanServerConnection is the connection to JMX endpoint. */
    private MBeanServerConnection connection;

    /* Date format used for Date type mBeans. */
    private final DateFormat dateFormat = DateFormat.getDateTimeInstance(2, 2, Locale.US);

    /* JMX configuration used to connect to JMX endpoint. */
    private JMXConfig jmxConfig;

    /* InternalStats used for troubleshooting. */
    private InternalStats internalStats;

    private JMXRequestHandler jmxRequestHandler = new JMXRequestHandler();

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
        if (jmxConfig == null) {
            throw new JMXConnectionError("failed to connect to JMX server: configuration not provided");
        }

        this.jmxConfig = jmxConfig;

        if (jmxConfig.enableInternalStats) {
            this.internalStats = new InternalStats(jmxConfig.maxInternalStatsSize);
        }

        String connectionString = buildConnectionString(jmxConfig);
        Map<String, Object> connectionEnv = buildConnectionEnvConfig(jmxConfig);

        InternalStat internalStat = null;
        if (this.internalStats != null) {
            internalStat = internalStats.record("connect");
        }

        try {
            JMXServiceURL address = new JMXServiceURL(connectionString);

            this.connector = JMXConnectorFactory.connect(address, connectionEnv);

            this.connection = connector.getMBeanServerConnection();

            if (internalStat != null) {
                internalStat.setSuccessful(true);
            }
        } catch (Exception e) {
            String message = String.format("can't connect to JMX server: '%s', error: '%s'",
                    connectionString,
                    getErrorMessage(e));
            throw new JMXConnectionError(message);
        } finally {
            if (internalStat != null) {
                InternalStats.setElapsedMs(internalStat);
            }
        }
    }

    /**
     * disconnect from the JMX endpoint.
     *
     * @param timeoutMs long timeout for the request in milliseconds
     * @throws JMXConnectionError JMX connection related exception
     */
    public void disconnect(long timeoutMs) throws JMXError, JMXConnectionError {
        withTimeout(executor.submit((Callable<Void>) () -> {
            disconnect();
            return null;
        }), timeoutMs);
    }

    /**
     * disconnect from the JMX endpoint.
     *
     * @throws JMXConnectionError JMX connection related exception
     */
    public void disconnect() throws JMXConnectionError {
        if (Thread.interrupted()) {
            return;
        }

        if (this.connector == null) {
            throw new JMXConnectionError()
                    .setMessage("cannot disconnect, connection to JMX endpoint is not established");
        }

        InternalStat internalStat = null;
        if (this.internalStats != null) {
            internalStat = internalStats.record("disconnect");
        }

        // Move this to a different variable in case close operation timeouts.
        JMXConnector oldConnector = this.connector;

        // Mark the connector as null in case to allow reconnection.
        this.connector = null;

        try {
            oldConnector.close();
            if (internalStat != null) {
                internalStat.setSuccessful(true);
            }
        } catch (Exception e) {
        } finally {
            if (internalStat != null) {
                InternalStats.setElapsedMs(internalStat);
            }
        }
    }

    /**
     * queryMBeanNames returns all founded mBean names that match the provided pattern.
     *
     * @param mBeanGlobPattern String glob pattern DOMAIN:BEAN e.g *:* or jboss.as:subsystem=remoting,configuration=endpoint
     * @param timeoutMs        long timeout for the request in milliseconds
     * @return List<String> containing all mBean names that were found
     * @throws JMXError           JMX related Exception
     * @throws JMXConnectionError JMX connection related exception
     */
    public List<String> queryMBeanNames(String mBeanGlobPattern, long timeoutMs) throws JMXError, JMXConnectionError {
        return this.withTimeout(
                executor.submit(() -> queryMBeanNames(mBeanGlobPattern)),
                timeoutMs
        );
    }

    /**
     * queryMBeanNames returns all founded mBean names that match the provided pattern.
     *
     * @param mBeanGlobPattern String glob pattern DOMAIN:BEAN e.g *:* or jboss.as:subsystem=remoting,configuration=endpoint
     * @return List<String> containing all mBean names that were found
     * @throws JMXError           JMX related Exception
     * @throws JMXConnectionError JMX connection related exception
     */
    public List<String> queryMBeanNames(String mBeanGlobPattern) throws JMXConnectionError, JMXError {
        ObjectName objectName = getObjectName(mBeanGlobPattern);
        return queryMBeans(objectName)
                .stream()
                .map(ObjectInstance::getObjectName)
                .map(ObjectName::toString)
                .collect(Collectors.toList());
    }

    /**
     * queryMBeans returns all founded mBeans that match the provided pattern.
     *
     * @param objectName ObjectName glob pattern DOMAIN:BEAN e.g *:* or jboss.as:subsystem=remoting,configuration=endpoint
     * @return Set<ObjectInstance>
     * @throws JMXError           JMX related Exception
     * @throws JMXConnectionError JMX connection related exception
     */
    private Set<ObjectInstance> queryMBeans(ObjectName objectName) throws JMXConnectionError, JMXError {
        if (objectName == null) {
            throw new JMXError()
                    .setMessage("can't query MBeans, provided objectName is Null");
        }

        Set<ObjectInstance> result = null;

        InternalStat internalStat = null;
        if (this.internalStats != null) {
            internalStat = internalStats.record("queryMBeans")
                    .setMBean(objectName.toString());
        }

        try {
            result = jmxRequestHandler.exec(() ->
                    getConnection().queryMBeans(objectName, null)
            );


            if (internalStat != null) {
                internalStat.setSuccessful(true);
            }

            return result;
        } catch (JMXConnectionError je) {
            throw je;
        } catch (IOException io) {
            disconnect();

            String message = String.format("problem occurred when talking to the JMX server while querying mBeans, error: '%s'", io.getMessage());
            throw new JMXConnectionError(message);
        } catch (Exception e) {
            throw new JMXError()
                    .setMessage("can't get beans for query: " + objectName)
                    .setCauseMessage(e.getMessage())
                    .setStacktrace(getStackTrace(e));
        } finally {
            if (internalStat != null) {
                InternalStats.setElapsedMs(internalStat);

                if (result != null) {
                    internalStat.setResponseCount(result.size());
                }
            }
        }
    }

    /**
     * getMBeanAttributeNames returns all the available JMX attribute names for a given mBeanName.
     *
     * @param mBeanName of which we want to retrieve attribute names
     * @param timeoutMs long timeout for the request in milliseconds
     * @return List<String> containing all mBean attribute names that were found
     * @throws JMXError           JMX related Exception
     * @throws JMXConnectionError JMX connection related exception
     */
    public List<String> getMBeanAttributeNames(String mBeanName, long timeoutMs) throws JMXError, JMXConnectionError {
        return withTimeout(
                executor.submit(() -> getMBeanAttributeNames(getObjectName(mBeanName))),
                timeoutMs
        );
    }

    /**
     * getMBeanAttributeNames returns all the available JMX attribute names for a given mBeanName.
     *
     * @param objectName of which we want to retrieve attributes names
     * @return List<String> containing all mBean attribute names that were found
     * @throws JMXConnectionError JMX connection related exception
     * @throws JMXError           JMX related Exception
     */
    private List<String> getMBeanAttributeNames(ObjectName objectName) throws JMXConnectionError, JMXError {
        if (objectName == null) {
            throw new JMXError()
                    .setMessage("can't get attribute names, provided objectName is Null");
        }

        MBeanInfo info;

        InternalStat internalStat = null;
        if (this.internalStats != null) {
            internalStat = internalStats.record("getMBeanInfo")
                    .setMBean(objectName.toString());
        }

        try {
            info = jmxRequestHandler.exec(() ->
                    getConnection().getMBeanInfo(objectName)
            );

            if (internalStat != null) {
                internalStat.setSuccessful(true);
            }
        } catch (JMXConnectionError je) {
            throw je;
        } catch (IOException io) {
            disconnect();

            String message = String.format("problem occurred when talking to the JMX server while requesting mBean info, error: '%s'", io.getMessage());
            throw new JMXConnectionError(message);
        } catch (InstanceNotFoundException | IntrospectionException | ReflectionException e) {
            throw new JMXError()
                    .setMessage("can't find mBean: " + objectName)
                    .setCauseMessage(e.getMessage())
                    .setStacktrace(getStackTrace(e));
        } catch (Exception e) {
            throw new JMXConnectionError();
        } finally {
            if (internalStat != null) {
                InternalStats.setElapsedMs(internalStat);
            }
        }

        List<String> result = new ArrayList<>();
        if (info == null) {
            return result;
        }

        for (MBeanAttributeInfo attrInfo : info.getAttributes()) {
            if (internalStat != null) {
                internalStat.setResponseCount(internalStat.responseCount + 1);
            }

            if (attrInfo == null || !attrInfo.isReadable()) {
                continue;
            }
            result.add(attrInfo.getName());
        }

        return result;
    }

    /**
     * getMBeanAttributes returns the attribute values for an mBeanName.
     *
     * @param mBeanName  of which we want to retrieve the attribute values
     * @param attributes List of attribute names that we want to retrieve the values for
     * @param timeoutMs  long timeout for the request in milliseconds
     * @return List<AttributeResponse> representing the mBean attribute values
     * @throws JMXError           JMX related Exception
     * @throws JMXConnectionError JMX connection related exception
     */
    public List<AttributeResponse> getMBeanAttributes(String mBeanName, List<String> attributes, long timeoutMs) throws JMXError, JMXConnectionError {
        List<AttributeResponse> result = new ArrayList<>();
        withTimeout(
                executor.submit((Callable<Void>) () -> {
                    getMBeanAttributes(getObjectName(mBeanName), attributes, result);
                    return null;
                }), timeoutMs
        );
        return result;
    }

    /**
     * getMBeanAttributes fetches the attribute values for an mBeanName.
     *
     * @param objectName of which we want to retrieve the attribute values
     * @param attributes List of attribute names that we want to retrieve the values for
     * @param output     List<AttributeResponse> to add the fetched attribute values
     * @throws JMXError           JMX related Exception
     * @throws JMXConnectionError JMX connection related exception
     */
    private void getMBeanAttributes(ObjectName objectName, List<String> attributes, List<AttributeResponse> output) throws JMXConnectionError, JMXError {
        if (objectName == null) {
            throw new JMXError()
                    .setMessage("can't get attribute value, provided objectName is Null");
        }

        if (output == null) {
            throw new JMXError()
                    .setMessage("can't deserialize attribute value, provided output list is Null");
        }

        if (attributes == null || attributes.size() == 0) {
            attributes = getMBeanAttributeNames(objectName);
        }

        InternalStat internalStat = null;
        if (this.internalStats != null) {
            internalStat = internalStats.record("getAttributes")
                    .setMBean(objectName.toString())
                    .setAttrs(attributes);
        }

        List<Attribute> attrValues = new ArrayList<>();
        AttributeList attributeList;
        try {

            List<String> finalAttributes = attributes;

            attributeList = jmxRequestHandler.exec(() ->
                    getConnection().getAttributes(objectName, finalAttributes.toArray(new String[0]))
            );

            if (internalStat != null) {
                internalStat.setSuccessful(true);
            }
        } catch (JMXConnectionError je) {
            throw je;
        } catch (ConnectException ce) {
            disconnect();

            String message = String.format("problem occurred when talking to the JMX server while requesting attributes, error: '%s'", ce.getMessage());
            throw new JMXConnectionError(message);
        } catch (Exception e) {
            // When running a call for multiple attributes it can fail only because one of them.
            // In that case we try to make a separate call for each one.
            for (String attribute : attributes) {
                String formattedAttrName = formatAttributeName(objectName, attribute);

                try {
                    getMBeanAttribute(objectName, attribute, output);
                } catch (JMXError je) {
                    String statusMessage = String.format("can't get attribute, error: '%s', cause: '%s', stacktrace: '%s'", je.message, je.causeMessage, je.stacktrace);
                    output.add(new AttributeResponse()
                            .setName(formattedAttrName)
                            .setResponseType(ResponseType.ERROR)
                            .setStatusMsg(statusMessage));
                }
            }
            return;
        } finally {
            if (internalStat != null) {
                InternalStats.setElapsedMs(internalStat);
            }
        }

        // Keep a track of requested attributes to report the ones that we fail to retrieve.
        List<String> missingAttrs = new ArrayList<>(attributes);

        try {
            if (attributeList == null) {
                return;
            }

            for (Object value : attributeList) {
                if (internalStat != null) {
                    internalStat.setResponseCount(internalStat.responseCount + 1);
                }

                if (value instanceof Attribute) {
                    Attribute attr = (Attribute) value;
                    attrValues.add(attr);

                    missingAttrs.remove(attr.getName());
                }
            }

            for (Attribute attrValue : attrValues) {
                String formattedAttrName = formatAttributeName(objectName, attrValue.getName());

                try {
                    parseValue(formattedAttrName, attrValue.getValue(), output);
                } catch (JMXError je) {
                    String statusMessage = String.format("can't parse attribute, error: '%s', cause: '%s', stacktrace: '%s'", je.message, je.causeMessage, je.stacktrace);
                    output.add(new AttributeResponse()
                            .setName(formattedAttrName)
                            .setResponseType(ResponseType.ERROR)
                            .setStatusMsg(statusMessage));
                }
            }
        } finally {
            // Report requested attributes that we didn't retrieve.
            for (String attr : missingAttrs) {
                String formattedAttrName = formatAttributeName(objectName, attr);
                output.add(new AttributeResponse()
                        .setName(formattedAttrName)
                        .setResponseType(ResponseType.ERROR)
                        .setStatusMsg("failed to retrieve attribute value from server"));
            }
        }
    }

    /**
     * getMBeanAttribute fetches the attribute value for an mBeanName.
     * CompositeData is handled as multiple values.
     *
     * @param objectName of which we want to retrieve the attribute values
     * @param attribute  of which we want to retrieve the attribute values
     * @param output     List<AttributeResponse> to add the fetched attribute values.
     * @throws JMXError           JMX related Exception
     * @throws JMXConnectionError JMX connection related exception
     */
    private void getMBeanAttribute(ObjectName objectName, String attribute, List<AttributeResponse> output) throws JMXConnectionError, JMXError {
        if (objectName == null) {
            throw new JMXError()
                    .setMessage("can't get attribute value, provided objectName is Null");
        }

        if (attribute == null) {
            throw new JMXError()
                    .setMessage("can't get attribute value, provided attribute name is Null");
        }

        if (output == null) {
            throw new JMXError()
                    .setMessage("can't deserialize attribute value, provided output list is Null");
        }

        Object value;
        try {
            value = jmxRequestHandler.exec(() ->
                    getConnection().getAttribute(objectName, attribute)
            );
            if (value instanceof Attribute) {
                Attribute jmxAttr = (Attribute) value;
                value = jmxAttr.getValue();
            }
        } catch (ConnectException ce) {
            disconnect();

            String message = String.format("can't connect to JMX server, error: '%s'", ce.getMessage());
            throw new JMXConnectionError(message);
        } catch (Exception e) {
            throw new JMXError()
                    .setMessage("can't get attribute: " + attribute + " for bean: " + objectName + ": ")
                    .setCauseMessage(e.getMessage())
                    .setStacktrace(getStackTrace(e));
        }


        String formattedAttrName = formatAttributeName(objectName, attribute);
        parseValue(formattedAttrName, value, output);
    }

    /**
     * queryMBeanAttributes will fetch all the available mBeans attributes for the mBean pattern.
     *
     * @param mBeanGlobPattern String glob pattern DOMAIN:BEAN e.g *:* or jboss.as:subsystem=remoting,configuration=endpoint
     * @param timeoutMs        long timeout for the request in milliseconds
     * @return List<AttributeResponse> containing all the fetched attribute values
     * @throws JMXError           JMX related Exception
     * @throws JMXConnectionError JMX connection related exception
     */
    public List<AttributeResponse> queryMBeanAttributes(String mBeanGlobPattern, List<String> attributes, long timeoutMs) throws JMXError, JMXConnectionError {
        return withTimeout(
                executor.submit(() -> queryMBeanAttributes(mBeanGlobPattern, attributes)),
                timeoutMs
        );
    }

    /**
     * queryMBeanAttributes will fetch all the available mBeans attributes for the mBean pattern.
     *
     * @param mBeanGlobPattern String glob pattern DOMAIN:BEAN e.g *:* or jboss.as:subsystem=remoting,configuration=endpoint
     * @return List<AttributeResponse> containing all the fetched attribute values
     * @throws JMXError           JMX related Exception
     * @throws JMXConnectionError JMX connection related exception
     */
    public List<AttributeResponse> queryMBeanAttributes(String mBeanGlobPattern, List<String> attributes) throws JMXError, JMXConnectionError {
        ObjectName pattern = getObjectName(mBeanGlobPattern);

        Set<ObjectInstance> mBeans;
        if (mBeanGlobPattern.contains("*")) {
            mBeans = queryMBeans(pattern);
        } else {
            mBeans = new HashSet<>(Arrays.asList(new ObjectInstance(pattern, "")));
        }

        List<AttributeResponse> result = new ArrayList<>();

        for (ObjectInstance mBean : mBeans) {
            if (mBean == null) {
                continue;
            }
            ObjectName objectName = mBean.getObjectName();

            getMBeanAttributes(objectName, attributes, result);
        }

        return result;
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
                    .setStacktrace(getStackTrace(me));
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
                    .setMessage("request was interrupted " + e.getMessage())
                    .setCauseMessage(e.getMessage())
                    .setStacktrace(getStackTrace(e));
        } catch (TimeoutException e) {
            throw new JMXError()
                    .setMessage("request timeout exceeded: " + timeoutMs + "ms")
                    .setCauseMessage(e.getMessage())
                    .setStacktrace(getStackTrace(e));
        } catch (ExecutionException e) {
            if (e.getCause() instanceof JMXError) {
                throw (JMXError) e.getCause();
            } else if (e.getCause() instanceof JMXConnectionError) {
                throw (JMXConnectionError) e.getCause();
            }
            throw new JMXError()
                    .setMessage("failed to execute operation, error: " + e.getMessage())
                    .setStacktrace(getStackTrace(e));
        } finally {
            future.cancel(true);
        }
    }

    /**
     * parseValue converts the received value from JMX into an JMXAttribute object.
     *
     * @param mBeanAttributeName of the value
     * @param value              that has to be converted
     * @throws JMXError JMX related Exception
     */
    private void parseValue(String mBeanAttributeName, Object value, List<AttributeResponse> output) throws JMXError {
        if (output == null) {
            output = new ArrayList<>();
        }
        AttributeResponse attr = new AttributeResponse();
        attr.name = mBeanAttributeName;

        if (value == null) {
            throw new JMXError()
                    .setMessage("found a null value for bean: " + mBeanAttributeName);
        } else if (value instanceof java.lang.Double) {
            attr.doubleValue = (Double) value;
            attr.responseType = ResponseType.DOUBLE;
        } else if (value instanceof java.lang.Float) {
            attr.doubleValue = new BigDecimal(value.toString()).doubleValue();
            attr.responseType = ResponseType.DOUBLE;
        } else if (value instanceof Number) {
            attr.intValue = ((Number) value).longValue();
            attr.responseType = ResponseType.INT;
        } else if (value instanceof String) {
            attr.stringValue = (String) value;
            attr.responseType = ResponseType.STRING;
        } else if (value instanceof Boolean) {
            attr.boolValue = (Boolean) value;
            attr.responseType = ResponseType.BOOL;
        } else if (value instanceof java.util.Date) {
            attr.stringValue = dateFormat.format(value);
            attr.responseType = ResponseType.STRING;
        } else if (value instanceof CompositeData) {
            CompositeData cdata = (CompositeData) value;
            Set<String> fieldKeys = cdata.getCompositeType().keySet();
            JMXError jmxError = null;

            for (String field : fieldKeys) {
                if (field.length() < 1) {
                    continue;
                }

                String fieldKey = field.substring(0, 1).toUpperCase() + field.substring(1);
                try {
                    parseValue(String.format("%s.%s", mBeanAttributeName, fieldKey), cdata.get(field), output);
                } catch (JMXError e) {
                    jmxError = e;
                }
            }
            if (output.size() == 0 && jmxError != null) {
                throw jmxError;
            }
            return;
        } else {
            throw new JMXError()
                    .setMessage("unsuported data type (" + value.getClass() + ") for bean " + mBeanAttributeName);
        }
        output.add(attr);
    }

    /**
     * collect the internal stats used for troubleshooting
     *
     * @return List<InternalStat> the collected nrjmx internal query stats.
     */
    public List<InternalStat> getInternalStats() throws JMXError {
        if (internalStats == null) {
            throw new JMXError()
                    .setMessage("internal stats not activated");
        }

        return internalStats.getStats();
    }

    /**
     * getConnection returns the connection the the JMX endpoint.
     *
     * @return MBeanServerConnection the connection to the JMX endpoint
     * @throws JMXConnectionError JMX connection related Exception
     */
    private MBeanServerConnection getConnection() throws JMXConnectionError {
        if (jmxConfig == null) {
            throw new JMXConnectionError("failed to get connection to JMX server: configuration not provided");
        }

        if (this.connector == null) {
            connect(jmxConfig);
        }

        if (this.connection == null) {
            throw new JMXConnectionError()
                    .setMessage("connection to JMX endpoint is not established");
        }
        return this.connection;
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

    private String getStackTrace(Throwable throwable) {
        if (throwable == null || jmxConfig == null || !jmxConfig.verbose) {
            return "";
        }
        return ExceptionUtils.getStackTrace(throwable);
    }

    private String getErrorMessage(Throwable throwable) {
        if (throwable == null) {
            return "NULL";
        }

        String message = throwable.getMessage();
        if (message != null && !message.equals("")) {
            return message;
        }

        message = throwable.getLocalizedMessage();

        if (message != null && !message.equals("")) {
            return message;
        }

        return ExceptionUtils.getStackTrace(throwable);
    }

    private String formatAttributeName(ObjectName objectName, String attribute) {
        return String.format("%s,attr=%s", objectName, attribute);
    }
}
