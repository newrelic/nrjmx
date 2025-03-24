/*
 * Copyright 2021 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package gojmx

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/newrelic/nrjmx/gojmx/internal/testutils"
	gopsutil "github.com/shirou/gopsutil/v3/process"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var timeStamp = time.Date(2022, time.January, 1, 01, 23, 45, 0, time.Local).UnixNano() / 1000000

func init() {
	_ = os.Setenv("NR_JMX_TOOL", filepath.Join(testutils.PrjDir, "bin", "nrjmx"))
	//_ = os.Setenv("NRIA_NRJMX_DEBUG", "true")
}

func Test_Query_Success_LargeAmountOfData(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	var data []map[string]interface{}

	name := strings.Repeat("tomas", 100)

	for i := 0; i < 1500; i++ {
		data = append(data, map[string]interface{}{
			"name":        fmt.Sprintf("%s-%d", name, i),
			"doubleValue": 1.2,
			"floatValue":  2.2,
			"numberValue": 3,
			"boolValue":   true,
			"dateValue":   timeStamp,
		})
	}

	// Populate the JMX Server with mbeans
	resp, err := testutils.AddMBeansBatch(ctx, container, data)
	assert.NoError(t, err)
	assert.Equal(t, "ok!\n", string(resp))

	defer testutils.CleanMBeans(ctx, container)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
	require.NoError(t, err)

	// THEN JMX connection can be oppened
	config := &JMXConfig{
		Hostname: jmxHost,
		Port:     int32(jmxPort.Int()),
	}
	client, err := NewClient(ctx).Open(config)
	assert.NoError(t, err)
	defer assertCloseClientNoError(t, client)

	// AND query returns at least 5Mb of data.
	mBeanNames, err := client.QueryMBeanNames("test:type=Cat,*")
	assert.NoError(t, err)

	var result []*AttributeResponse

	for _, mBeanName := range mBeanNames {
		mBeanAttrNames, err := client.GetMBeanAttributeNames(mBeanName)
		assert.NoError(t, err)

		for _, mBeanAttrName := range mBeanAttrNames {
			jmxAttrs, err := client.GetMBeanAttributes(mBeanName, mBeanAttrName)
			assert.NoError(t, err)
			result = append(result, jmxAttrs...)
		}

	}
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(fmt.Sprintf("%v", result)), 5*1024*1024)
}

func Test_Query_Success(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	// Populate the JMX Server with mbeans
	resp, err := testutils.AddMBeans(ctx, container, map[string]interface{}{
		"name":        "tomas",
		"doubleValue": 1.2,
		"floatValue":  2.2222222,
		"numberValue": 3,
		"boolValue":   true,
		"dateValue":   timeStamp,
	})
	assert.NoError(t, err)
	assert.Equal(t, "ok!\n", string(resp))

	defer testutils.CleanMBeans(ctx, container)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
	require.NoError(t, err)

	// THEN JMX connection can be oppened
	config := &JMXConfig{
		Hostname:         jmxHost,
		Port:             int32(jmxPort.Int()),
		RequestTimeoutMs: testutils.DefaultTimeoutMs,
	}

	client, err := NewClient(ctx).Open(config)
	assert.NoError(t, err)
	defer assertCloseClientNoError(t, client)

	// AND Query returns expected data
	expectedMBeanNames := []string{
		"test:type=Cat,name=tomas",
	}

	actualMBeanNames, err := client.QueryMBeanNames("test:type=Cat,*")
	require.NoError(t, err)
	require.ElementsMatch(t, expectedMBeanNames, actualMBeanNames)

	expectedMBeanAttrNames := []string{
		"BoolValue", "FloatValue", "NumberValue", "DoubleValue", "Name", "DateValue",
	}
	actualMBeanAttrNames, err := client.GetMBeanAttributeNames("test:name=tomas,type=Cat")
	require.NoError(t, err)
	require.ElementsMatch(t, expectedMBeanAttrNames, actualMBeanAttrNames)

	// AND Query returns expected data
	expected := []*AttributeResponse{
		{
			Name: "test:type=Cat,name=tomas,attr=NumberValue",

			ResponseType: ResponseTypeInt,
			IntValue:     3,
		},
		{
			Name: "test:type=Cat,name=tomas,attr=FloatValue",

			ResponseType: ResponseTypeDouble,
			DoubleValue:  2.222222,
		},
		{
			Name: "test:type=Cat,name=tomas,attr=DateValue",

			ResponseType: ResponseTypeString,
			StringValue:  "Jan 1, 2022, 1:23:45 AM",
		},
		{
			Name: "test:type=Cat,name=tomas,attr=BoolValue",

			ResponseType: ResponseTypeBool,
			BoolValue:    true,
		},
		{
			Name: "test:type=Cat,name=tomas,attr=DoubleValue",

			ResponseType: ResponseTypeDouble,
			DoubleValue:  1.2,
		},
		{
			Name: "test:type=Cat,name=tomas,attr=Name",

			ResponseType: ResponseTypeString,
			StringValue:  "tomas",
		},
	}

	var actual []*AttributeResponse
	for _, mBeanAttrName := range expectedMBeanAttrNames {
		jmxAttrs, err := client.GetMBeanAttributes("test:type=Cat,name=tomas", mBeanAttrName)
		assert.NoError(t, err)
		actual = append(actual, jmxAttrs...)
	}
	assert.ElementsMatch(t, expected, actual)
}

func Test_Query_Exception_Success(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	// Populate the JMX Server with mbeans
	resp, err := testutils.AddMBeansWithException(ctx, container, map[string]interface{}{
		"name":        "tomas",
		"doubleValue": 1.2,
	})

	assert.NoError(t, err)
	assert.Equal(t, "ok!\n", string(resp))

	defer testutils.CleanMBeans(ctx, container)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
	require.NoError(t, err)

	// THEN JMX connection can be opened
	config := &JMXConfig{
		Hostname:         jmxHost,
		Port:             int32(jmxPort.Int()),
		RequestTimeoutMs: testutils.DefaultTimeoutMs,
	}

	client, err := NewClient(ctx).Open(config)
	assert.NoError(t, err)
	defer assertCloseClientNoError(t, client)

	actualMBeans, err := client.QueryMBeanAttributes("test:type=ExceptionalCat,*")
	require.NoError(t, err)

	// AND Query returns expected data
	expected := []*AttributeResponse{
		{
			Name: "test:type=ExceptionalCat,name=tomas,attr=DoubleValue",

			ResponseType: ResponseTypeDouble,
			DoubleValue:  1.2,
		},
		{
			Name:         "test:type=ExceptionalCat,name=tomas,attr=NotSerializable",
			StatusMsg:    "can't get attribute, error: 'can't get attribute: NotSerializable for bean: test:type=ExceptionalCat,name=tomas: ', cause: 'error unmarshalling return; nested exception is: \n\tjava.io.WriteAbortedException: writing aborted; java.io.NotSerializableException: org.newrelic.jmx.ExceptionalCat$NotSerializable', stacktrace: ''",
			ResponseType: ResponseTypeErr,
		},
	}

	assert.ElementsMatch(t, expected, actualMBeans)
}

func Test_QueryMBean_Success(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	// Populate the JMX Server with mbeans
	resp, err := testutils.AddMBeans(ctx, container, map[string]interface{}{
		"name":        "tomas",
		"doubleValue": 1.2,
		"floatValue":  2.2222222,
		"numberValue": 3,
		"boolValue":   true,
	})
	assert.NoError(t, err)
	assert.Equal(t, "ok!\n", string(resp))

	defer testutils.CleanMBeans(ctx, container)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
	require.NoError(t, err)

	// THEN JMX connection can be oppened
	config := &JMXConfig{
		Hostname:         jmxHost,
		Port:             int32(jmxPort.Int()),
		RequestTimeoutMs: testutils.DefaultTimeoutMs,
	}

	client, err := NewClient(ctx).Open(config)
	assert.NoError(t, err)
	defer assertCloseClientNoError(t, client)

	testCases := []struct {
		name       string
		query      string
		attributes []string
		expected   []*AttributeResponse
	}{
		{
			name:  "With_attributes_specified",
			query: "test:type=Cat,name=*",
			attributes: []string{
				"NumberValue",
				"BoolValue",
			},
			expected: []*AttributeResponse{
				{
					Name:         "test:type=Cat,name=tomas,attr=NumberValue",
					ResponseType: ResponseTypeInt,
					IntValue:     3,
				},
				{
					Name: "test:type=Cat,name=tomas,attr=BoolValue",

					ResponseType: ResponseTypeBool,
					BoolValue:    true,
				},
			},
		},
		{
			name:       "Without_attributes_specified",
			query:      "test:type=Cat,name=tomas",
			attributes: nil,
			expected: []*AttributeResponse{
				{
					Name: "test:type=Cat,name=tomas,attr=FloatValue",

					ResponseType: ResponseTypeDouble,
					DoubleValue:  2.222222,
				},
				{
					Name: "test:type=Cat,name=tomas,attr=NumberValue",

					ResponseType: ResponseTypeInt,
					IntValue:     3,
				},
				{
					Name: "test:type=Cat,name=tomas,attr=BoolValue",

					ResponseType: ResponseTypeBool,
					BoolValue:    true,
				},
				{
					Name: "test:type=Cat,name=tomas,attr=DoubleValue",

					ResponseType: ResponseTypeDouble,
					DoubleValue:  1.2,
				},
				{
					Name: "test:type=Cat,name=tomas,attr=Name",

					ResponseType: ResponseTypeString,
					StringValue:  "tomas",
				},
				{
					Name:         "test:type=Cat,name=tomas,attr=DateValue",
					ResponseType: ResponseTypeErr,
					StatusMsg:    `can't parse attribute, error: 'found a null value for bean: test:type=Cat,name=tomas,attr=DateValue', cause: 'null', stacktrace: 'null'`,
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// AND Query returns expected data
			actualResponse, err := client.QueryMBeanAttributes(testCase.query, testCase.attributes...)
			require.NoError(t, err)

			assert.ElementsMatch(t, testCase.expected, actualResponse)
		})
	}
}

func Test_Query_CompositeData(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	// Populate the JMX Server with mbeans
	resp, err := testutils.AddMCompositeDataBeans(ctx, container, map[string]interface{}{
		"name":        "tomas",
		"doubleValue": 1.2,
	})
	assert.NoError(t, err)
	assert.Equal(t, "ok!\n", string(resp))

	defer testutils.CleanMBeans(ctx, container)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
	require.NoError(t, err)

	// THEN JMX connection can be oppened
	config := &JMXConfig{
		Hostname:         jmxHost,
		Port:             int32(jmxPort.Int()),
		RequestTimeoutMs: testutils.DefaultTimeoutMs,
	}

	client, err := NewClient(ctx).Open(config)
	assert.NoError(t, err)
	defer assertCloseClientNoError(t, client)

	// AND Query returns expected data
	expectedMBeanNames := []string{
		"test:type=CompositeDataCat,name=tomas",
	}
	actualMBeanNames, err := client.QueryMBeanNames("test:type=CompositeDataCat,*")
	require.NoError(t, err)
	require.ElementsMatch(t, expectedMBeanNames, actualMBeanNames)

	expectedMBeanAttrNames := []string{
		"CatInfo",
	}
	actualMBeanAttrNames, err := client.GetMBeanAttributeNames("test:name=tomas,type=CompositeDataCat")
	require.NoError(t, err)
	require.ElementsMatch(t, expectedMBeanAttrNames, actualMBeanAttrNames)

	// AND Query returns expected data
	expected := []*AttributeResponse{
		{
			Name: "test:type=CompositeDataCat,name=tomas,attr=CatInfo.Double",

			ResponseType: ResponseTypeDouble,
			DoubleValue:  1.2,
		},
		{
			Name: "test:type=CompositeDataCat,name=tomas,attr=CatInfo.Name",

			ResponseType: ResponseTypeString,
			StringValue:  "tomas",
		},
	}

	actual, err := client.GetMBeanAttributes("test:type=CompositeDataCat,name=tomas", "CatInfo")
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected, actual)
}

func Test_Query_Timeout(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
	require.NoError(t, err)

	// THEN JMX connection fails
	config := &JMXConfig{
		Hostname:         jmxHost,
		Port:             int32(jmxPort.Int()),
		RequestTimeoutMs: 1,
	}
	client, err := NewClient(ctx).Open(config)
	assert.NotNil(t, client)
	assert.Error(t, err)
	defer assertCloseClientNoError(t, client)

	// AND Query returns expected error
	actual, err := client.GetMBeanAttributeNames("*:*")
	assert.Nil(t, actual)
	assert.Error(t, err)
}

func Test_ConnectionURL_Success(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	// Populate the JMX Server with mbeans
	resp, err := testutils.AddMBeans(ctx, container, map[string]interface{}{
		"name":        "tomas",
		"doubleValue": 1.2,
		"floatValue":  2.2,
		"numberValue": 3,
		"boolValue":   true,
	})
	assert.NoError(t, err)
	assert.Equal(t, "ok!\n", string(resp))
	defer testutils.CleanMBeans(ctx, container)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
	require.NoError(t, err)

	// THEN JMX connection can be oppened
	config := &JMXConfig{
		ConnectionURL:    fmt.Sprintf("service:jmx:rmi:///jndi/rmi://%s:%s/jmxrmi", jmxHost, jmxPort.Port()),
		RequestTimeoutMs: testutils.DefaultTimeoutMs,
	}
	client, err := NewClient(ctx).Open(config)
	assert.NoError(t, err)
	defer assertCloseClientNoError(t, client)

	// AND Query returns expected data
	actual, err := client.GetMBeanAttributes("test:type=Cat,name=tomas", "FloatValue")
	assert.NoError(t, err)

	expected := []*AttributeResponse{
		{
			Name:         "test:type=Cat,name=tomas,attr=FloatValue",
			ResponseType: ResponseTypeDouble,
			DoubleValue:  2.2,
		},
	}

	assert.Equal(t, expected, actual)
}

func Test_JavaNotInstalledError(t *testing.T) {
	// GIVEN a wrong Java Home
	os.Setenv("NRIA_JAVA_HOME", "/wrong/path")
	defer os.Unsetenv("NRIA_JAVA_HOME")

	ctx := context.Background()

	config := &JMXConfig{
		RequestTimeoutMs: testutils.DefaultTimeoutMs,
	}
	// THEN connect fails with expected error
	client, err := NewClient(ctx).Open(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "/wrong/path/bin/java")

	// AND Query fails with expected error
	actual, err := client.QueryMBeanNames("test:type=Cat,*")
	assert.Nil(t, actual)
	assert.ErrorIs(t, err, errProcessNotRunning)
}

func Test_WrongMBeanFormatError(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
	require.NoError(t, err)

	// THEN JMX connection can be oppened
	config := &JMXConfig{
		ConnectionURL:    fmt.Sprintf("service:jmx:rmi:///jndi/rmi://%s:%s/jmxrmi", jmxHost, jmxPort.Port()),
		RequestTimeoutMs: testutils.DefaultTimeoutMs,
	}

	client, err := NewClient(ctx).Open(config)
	assert.NoError(t, err)
	defer assertCloseClientNoError(t, client)

	// AND Query returns expected error
	actual, err := client.QueryMBeanNames("wrong_format")
	assert.Nil(t, actual)
	assert.EqualError(t, err, `jmx error: cannot parse MBean glob pattern: 'wrong_format', valid: 'DOMAIN:BEAN', cause: Key properties cannot be empty, stacktrace: `)
}

func Test_Wrong_Connection(t *testing.T) {
	ctx := context.Background()

	// GIVEN a wrong hostname and port
	config := &JMXConfig{
		Hostname:         "localhost",
		Port:             1234,
		RequestTimeoutMs: testutils.DefaultTimeoutMs,
	}

	// THEN open fails with expected error
	client, err := NewClient(ctx).Open(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Connection refused to host: localhost;")
	defer assertCloseClientNoError(t, client)

	// AND query returns expected error
	assert.Contains(t, err.Error(), "Connection refused to host: localhost;") // TODO: fix this, doesn't return the correct error

	actual, err := client.QueryMBeanNames("test:type=Cat,*")
	assert.Nil(t, actual)
	assert.Errorf(t, err, "connection to JMX endpoint is not established")
}

func Test_SSLQuery_Success(t *testing.T) {
	ctx := context.Background()

	// GIVEN an SSL JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainerSSL(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	// Populate the JMX Server with mbeans
	resp, err := testutils.AddMBeans(ctx, container, map[string]interface{}{
		"name":        "tomas",
		"doubleValue": 1.2,
		"floatValue":  2.222222,
		"numberValue": 3,
		"boolValue":   true,
		"dateValue":   1641429296000,
	})
	assert.NoError(t, err)
	assert.Equal(t, "ok!\n", string(resp))
	defer testutils.CleanMBeans(ctx, container)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
	require.NoError(t, err)

	// THEN SSL JMX connection can be opened
	config := &JMXConfig{
		Hostname:           jmxHost,
		Port:               int32(jmxPort.Int()),
		Username:           testutils.JmxUsername,
		Password:           testutils.JmxPassword,
		KeyStore:           testutils.KeystorePath,
		KeyStorePassword:   testutils.KeystorePassword,
		TrustStore:         testutils.TruststorePath,
		TrustStorePassword: testutils.TruststorePassword,
		RequestTimeoutMs:   testutils.DefaultTimeoutMs,
	}
	client, err := NewClient(ctx).Open(config)
	assert.NoError(t, err)
	defer assertCloseClientNoError(t, client)

	// AND Query returns expected data
	expectedMBeanNames := []string{
		"test:type=Cat,name=tomas",
	}
	actualMBeanNames, err := client.QueryMBeanNames("test:type=Cat,*")
	require.NoError(t, err)
	require.ElementsMatch(t, expectedMBeanNames, actualMBeanNames)

	expectedMBeanAttrNames := []string{
		"BoolValue", "FloatValue", "NumberValue", "DoubleValue", "Name", "DateValue",
	}
	actualMBeanAttrNames, err := client.GetMBeanAttributeNames("test:name=tomas,type=Cat")
	require.NoError(t, err)
	require.ElementsMatch(t, expectedMBeanAttrNames, actualMBeanAttrNames)
}

func Test_Wrong_Credentials(t *testing.T) {
	ctx := context.Background()

	// GIVEN an SSL JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainerSSL(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
	require.NoError(t, err)

	// WHEN wrong jmx username and password is provided
	config := &JMXConfig{
		Hostname:           jmxHost,
		Port:               int32(jmxPort.Int()),
		Username:           "wrong_username",
		Password:           "wrong_password",
		KeyStore:           testutils.KeystorePath,
		KeyStorePassword:   testutils.KeystorePassword,
		TrustStore:         testutils.TruststorePath,
		TrustStorePassword: testutils.TruststorePassword,
		RequestTimeoutMs:   testutils.DefaultTimeoutMs,
	}

	// THEN open fails with expected error
	client, err := NewClient(ctx).Open(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Authentication failed! Invalid username or password")
	defer assertCloseClientNoError(t, client)

	// AND Query returns expected error
	actual, err := client.QueryMBeanNames("test:type=Cat,*")
	assert.Nil(t, actual)
	assert.Errorf(t, err, "connection to JMX endpoint is not established")
}

func Test_Wrong_Certificate_password(t *testing.T) {
	ctx := context.Background()

	// GIVEN an SSL JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainerSSL(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
	require.NoError(t, err)

	// WHEN wrong jmx username and password is provided
	config := &JMXConfig{
		Hostname:           jmxHost,
		Port:               int32(jmxPort.Int()),
		Username:           testutils.JmxUsername,
		Password:           testutils.JmxPassword,
		KeyStore:           testutils.KeystorePath,
		KeyStorePassword:   "wrong_password",
		TrustStore:         testutils.TruststorePath,
		TrustStorePassword: testutils.TruststorePassword,
		RequestTimeoutMs:   testutils.DefaultTimeoutMs,
	}

	// THEN open returns expected error
	client, err := NewClient(ctx).Open(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SSLContext")
	defer assertCloseClientNoError(t, client)

	// AND Query returns expected error
	actual, err := client.QueryMBeanNames("test:type=Cat,*")
	assert.Nil(t, actual)
	assert.Errorf(t, err, "connection to JMX endpoint is not established")
}

func Test_Connector_Success(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JBoss Server with JMX exposed running inside a container
	container, err := testutils.RunJbossStandaloneJMXContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	// Install the connector
	dstFile := filepath.Join(testutils.PrjDir, "/bin/jboss-client.jar")
	err = testutils.CopyFileFromContainer(ctx, container, "/opt/jboss/wildfly/bin/client/jboss-client.jar", dstFile)
	assert.NoError(t, err)

	defer os.Remove(dstFile)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.JbossJMXPort)
	require.NoError(t, err)

	// THEN JMX connection can be opened
	config := &JMXConfig{
		Hostname:              jmxHost,
		Port:                  int32(jmxPort.Int()),
		Username:              testutils.JbossJMXUsername,
		Password:              testutils.JbossJMXPassword,
		IsJBossStandaloneMode: true,
		IsRemote:              true,
		RequestTimeoutMs:      testutils.DefaultTimeoutMs,
	}
	client, err := NewClient(ctx).Open(config)
	assert.NoError(t, err)
	defer assertCloseClientNoError(t, client)

	// AND Query returns expected data
	expectedMbeanNames := []string{
		"jboss.as:subsystem=remoting,configuration=endpoint",
	}
	actualMbeanNames, err := client.QueryMBeanNames("jboss.as:subsystem=remoting,configuration=endpoint")
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedMbeanNames, actualMbeanNames)

	expectedMBeanAttrNames := []string{
		"authRealm",
		"authenticationRetries",
		"authorizeId",
		"bufferRegionSize",
		"heartbeatInterval",
		"maxInboundChannels",
		"maxInboundMessageSize",
		"maxInboundMessages",
		"maxOutboundChannels",
		"maxOutboundMessageSize",
		"maxOutboundMessages",
		"receiveBufferSize",
		"receiveWindowSize",
		"saslProtocol",
		"sendBufferSize",
		"serverName",
		"transmitWindowSize",
		"worker",
	}
	actualMBeanAttrNames, err := client.GetMBeanAttributeNames("jboss.as:subsystem=remoting,configuration=endpoint")
	assert.ElementsMatch(t, expectedMBeanAttrNames, actualMBeanAttrNames)

	expected := []*AttributeResponse{
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=authRealm",
			StatusMsg:    "can't parse attribute, error: 'found a null value for bean: jboss.as:subsystem=remoting,configuration=endpoint,attr=authRealm', cause: 'null', stacktrace: 'null'",
			ResponseType: ResponseTypeErr,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=authenticationRetries",
			ResponseType: ResponseTypeInt,
			IntValue:     3,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=authorizeId",
			StatusMsg:    "can't parse attribute, error: 'found a null value for bean: jboss.as:subsystem=remoting,configuration=endpoint,attr=authorizeId', cause: 'null', stacktrace: 'null'",
			ResponseType: ResponseTypeErr,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=bufferRegionSize",
			StatusMsg:    "can't parse attribute, error: 'found a null value for bean: jboss.as:subsystem=remoting,configuration=endpoint,attr=bufferRegionSize', cause: 'null', stacktrace: 'null'",
			ResponseType: ResponseTypeErr,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=heartbeatInterval",
			ResponseType: ResponseTypeInt,
			IntValue:     60000,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=maxInboundChannels",
			ResponseType: ResponseTypeInt,
			IntValue:     40,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=maxInboundMessageSize",
			ResponseType: ResponseTypeInt,
			IntValue:     9223372036854775807,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=maxInboundMessages",
			ResponseType: ResponseTypeInt,
			IntValue:     80,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=maxOutboundChannels",
			ResponseType: ResponseTypeInt,
			IntValue:     40,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=maxOutboundMessageSize",
			ResponseType: ResponseTypeInt,
			IntValue:     9223372036854775807,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=maxOutboundMessages",
			ResponseType: ResponseTypeInt,
			IntValue:     65535,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=receiveBufferSize",
			ResponseType: ResponseTypeInt,
			IntValue:     8192,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=serverName",
			StatusMsg:    "can't parse attribute, error: 'found a null value for bean: jboss.as:subsystem=remoting,configuration=endpoint,attr=serverName', cause: 'null', stacktrace: 'null'",
			ResponseType: ResponseTypeErr,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=receiveWindowSize",
			ResponseType: ResponseTypeInt,
			IntValue:     131072,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=saslProtocol",
			ResponseType: ResponseTypeString,
			StringValue:  "remote",
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=sendBufferSize",
			ResponseType: ResponseTypeInt,
			IntValue:     8192,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=transmitWindowSize",
			ResponseType: ResponseTypeInt,
			IntValue:     131072,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=worker",
			ResponseType: ResponseTypeString,
			StringValue:  "default",
		},
	}

	var actual []*AttributeResponse
	for _, mBeanAttrName := range expectedMBeanAttrNames {
		jmxAttrs, err := client.GetMBeanAttributes("jboss.as:subsystem=remoting,configuration=endpoint", mBeanAttrName)
		if err != nil {
			continue
		}
		actual = append(actual, jmxAttrs...)
	}

	assert.ElementsMatch(t, expected, actual)
}

func TestClientClose(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
	require.NoError(t, err)

	// THEN JMX connection can be opened.
	config := &JMXConfig{
		Hostname:         jmxHost,
		Port:             int32(jmxPort.Int()),
		RequestTimeoutMs: testutils.DefaultTimeoutMs,
	}

	client, err := NewClient(ctx).Open(config)
	assert.NoError(t, err)

	// WHEN closing the client there is no error
	assertCloseClientNoError(t, client)

	// AND Query returns expected error
	actual, err := client.QueryMBeanNames("*:*")
	assert.Nil(t, actual)
	assert.ErrorIs(t, err, errProcessNotRunning)

	assert.True(t, client.nrJMXProcess.getOSProcessState().Success())
}

func TestGetInternalStats(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	var data []map[string]interface{}

	name := strings.Repeat("tomas", 100)

	for i := 0; i < 1500; i++ {
		data = append(data, map[string]interface{}{
			"name":        fmt.Sprintf("%s-%d", name, i),
			"doubleValue": 1.2,
			"floatValue":  2.2,
			"numberValue": 3,
			"boolValue":   true,
			"dateValue":   timeStamp,
		})
	}

	// Populate the JMX Server with mbeans
	resp, err := testutils.AddMBeansBatch(ctx, container, data)
	assert.NoError(t, err)
	assert.Equal(t, "ok!\n", string(resp))

	defer testutils.CleanMBeans(ctx, container)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
	require.NoError(t, err)

	// THEN JMX connection can be oppened
	config := &JMXConfig{
		Hostname:             jmxHost,
		Port:                 int32(jmxPort.Int()),
		EnableInternalStats:  true,
		MaxInternalStatsSize: 3000, // We expect 3002 stats. With MaxInternalStatsSize we test the limit.
	}
	client, err := NewClient(ctx).Open(config)
	assert.NoError(t, err)
	defer assertCloseClientNoError(t, client)

	_, err = client.QueryMBeanAttributes("test:type=Cat,*")
	assert.NoError(t, err)

	// AND query generates the expected internal stats
	internalStats, err := client.GetInternalStats()
	assert.NoError(t, err)

	assert.Len(t, internalStats, 3000)
	assert.Regexp(t, "StatType: 'getMBeanInfo', MBean: 'test:type=Cat,name=(.+?)', Attributes: '[]', TotalObjCount: 6, StartTimeMs: [0-9]+,  Duration: [0-9]+\\.[0-9]+ms, Successful: true", internalStats[0].String())
	assert.Regexp(t, "TotalMs: '[0-9]+\\.[0-9]{3}', TotalObjects: 18000, TotalAttr: 9000, TotalCalls: 3000, TotalSuccessful: 3000", internalStats.String())

	// AND internal stats get cleaned
	internalStats, err = client.GetInternalStats()
	assert.NoError(t, err)
	assert.Equal(t, "TotalMs: '0.000', TotalObjects: 0, TotalAttr: 0, TotalCalls: 0, TotalSuccessful: 0", internalStats.String())
}

func TestGetEmptyInternalStats(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
	require.NoError(t, err)

	// THEN JMX connection can be oppened
	config := &JMXConfig{
		Hostname:            jmxHost,
		Port:                int32(jmxPort.Int()),
		EnableInternalStats: false,
	}

	client, err := NewClient(ctx).Open(config)
	assert.NoError(t, err)
	defer assertCloseClientNoError(t, client)

	_, err = client.QueryMBeanAttributes("test:type=Cat,*")
	assert.NoError(t, err)

	// AND query doesn't generate any internal stats.
	internalStats, err := client.GetInternalStats()

	assert.Error(t, err)
	assert.Equal(t, "<nil>", internalStats.String())
}

func TestConnectionRecovers(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainer(ctx)
	require.NoError(t, err)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
	require.NoError(t, err)

	// THEN JMX connection can be opened.
	config := &JMXConfig{
		Hostname:            jmxHost,
		Port:                int32(jmxPort.Int()),
		RequestTimeoutMs:    5000,
		EnableInternalStats: true,
	}

	query := "java.lang:type=*"

	client, err := NewClient(ctx).Open(config)
	assert.NoError(t, err)
	defer assertCloseClientNoError(t, client)

	res, err := client.QueryMBeanNames(query)
	assert.NoError(t, err)
	assert.NotEmpty(t, res)

	assert.NoError(t, container.Terminate(ctx))

	res, err = client.QueryMBeanNames(query)
	assert.Nil(t, res)
	assert.Error(t, err)

	_, ok := IsJMXClientError(err)
	assert.False(t, ok)
	assert.True(t, client.IsRunning())

	assert.Eventually(t, func() bool {
		container, err = testutils.RunJMXServiceContainer(ctx)
		return err == nil
	}, 200*time.Second, 50*time.Millisecond,
		"didn't managed to restart the container")

	defer container.Terminate(ctx)

	assert.Eventually(t, func() bool {
		_, err = client.QueryMBeanNames(query)
		return err == nil
	}, 20*time.Second, 50*time.Millisecond,
		"didn't managed to recover connection")
}

func TestProcessExits(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	// Populate the JMX Server with mbeans
	resp, err := testutils.AddMBeans(ctx, container, map[string]interface{}{
		"name":        "tomas",
		"doubleValue": 1.2,
		// test-server can delay the answer for doubleValue by specifying a millisecond timeout value.
		// we need a delay in processing to simulate when the java process is stuck, and we want to terminate it.
		"timeout": 60000,
	})

	assert.NoError(t, err)
	assert.Equal(t, "ok!\n", string(resp))

	defer testutils.CleanMBeans(ctx, container)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
	require.NoError(t, err)

	// Run gojmx library as a child process.
	cmd := testutils.NrJMXAsSubprocess(ctx, jmxHost, jmxPort.Port())

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Start()

	// Call wait() to avoid defunct process.
	go cmd.Wait()

	// The subprocess will write the gojmx version to stdout. At that point we know that java process is initialized.
	waitToStart := make(chan string)
	go func() {
		r := bufio.NewReader(&stdout)
		assert.Eventually(t, func() bool {
			_, err := r.ReadString('\n')
			return err == nil
		}, 5*time.Second, 50*time.Millisecond,
			"didn't managed to read the gojmx version",
		)

		close(waitToStart)
	}()

	// For troubleshooting purposes.
	defer func() {
		stdoutBytes, _ := io.ReadAll(&stdout)
		fmt.Println(fmt.Sprintf("[DEBUG] Stdout: '%s'", stdoutBytes))

		stderrBytes, _ := io.ReadAll(&stderr)
		fmt.Println(fmt.Sprintf("[DEBUG] Stderr: '%s'", stderrBytes))
	}()

	<-waitToStart

	// The process tree is much larger, we want to identify the main go process and the java process.
	subProcess, err := gopsutil.NewProcess(int32(cmd.Process.Pid))
	require.NoError(t, err)

	mainProcess := findChildProcessByName(t, subProcess, "main")

	javaProcess := findChildProcessByName(t, mainProcess, "java")

	// WHEN main process is terminated.
	mainProcess.Kill()

	// THEN Java child also terminates.
	assert.Eventually(t, func() bool {
		up, err := javaProcess.IsRunning()
		if err != nil {
			return false
		}
		// assert is not running
		return !up
	}, 5*time.Second, 50*time.Millisecond,
		"java subprocess was not properly terminated")
}

// findChildProcessByName will search for a java process in a given process tree.
func findChildProcessByName(t *testing.T, parentProcess *gopsutil.Process, childName string) *gopsutil.Process {
	var err error
	var children []*gopsutil.Process

	require.Eventuallyf(t, func() bool {
		children, err = parentProcess.Children()
		if err != nil {
			return false
		}

		return len(children) > 0
	}, 5*time.Second, 50*time.Millisecond,
		"couldn't found the %s subprocess", childName)

	// We check only the first child since we don't start multiple processes in parallel.
	name, err := children[0].Name()
	require.NoError(t, err)
	if name == childName {
		return children[0]
	}

	return findChildProcessByName(t, children[0], childName)
}

func assertCloseClientNoError(t *testing.T, client *Client) {
	assert.NotNil(t, client)
	assert.NoError(t, client.Close())
}

func assertCloseClientError(t *testing.T, client *Client) {
	assert.NotNil(t, client)
	assert.Error(t, client.Close())
}
