package gojmx

import (
	"context"
	"fmt"
	"github.com/newrelic/nrjmx/gojmx/internal/nrjmx"
	"github.com/newrelic/nrjmx/gojmx/internal/nrprotocol"
	"github.com/newrelic/nrjmx/gojmx/internal/testutils"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func Test_Query_Success_LargeAmountOfData(t *testing.T) {
	ctx := context.Background()
	//
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
		})
	}

	// Populate the JMX Server with mbeans
	resp, err := testutils.AddMBeansBatch(ctx, container, data)
	assert.NoError(t, err)
	assert.Equal(t, "ok!\n", string(resp))

	defer testutils.CleanMBeans(ctx, container)

	jmxHost, jmxPort, err := testutils.GetContainerHostAndPort(ctx, container)
	require.NoError(t, err)

	// THEN JMX connection can be oppened
	config := &JMXConfig{
		Hostname: jmxHost,
		Port:     int32(jmxPort.Int()),
	}
	client, err := NewClient(ctx).Open(config)
	defer client.Close()

	assert.NoError(t, err)

	// AND query returns at least 5Mb of data.
	mBeanNames, err := client.GetMBeanNames("test:type=Cat,*")
	assert.NoError(t, err)

	var result []*nrprotocol.JMXAttribute

	for _, mBeanName := range mBeanNames {
		mBeanAttrNames, err := client.GetMBeanAttrNames(mBeanName)
		assert.NoError(t, err)

		for _, mBeanAttrName := range mBeanAttrNames {
			jmxAttr, err := client.GetMBeanAttr(mBeanName, mBeanAttrName)
			assert.NoError(t, err)
			result = append(result, jmxAttr)
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
	})
	assert.NoError(t, err)
	assert.Equal(t, "ok!\n", string(resp))

	defer testutils.CleanMBeans(ctx, container)

	jmxHost, jmxPort, err := testutils.GetContainerHostAndPort(ctx, container)
	require.NoError(t, err)

	// THEN JMX connection can be oppened
	config := &JMXConfig{
		Hostname:        jmxHost,
		Port:            int32(jmxPort.Int()),
		RequestTimoutMs: testutils.DefaultTimeoutMs,
	}

	client, err := NewClient(ctx).Open(config)
	defer client.Close()
	assert.NoError(t, err)

	// AND Query returns expected data
	expectedMBeanNames := []string{
		"test:type=Cat,name=tomas",
	}
	actualMBeanNames, err := client.GetMBeanNames("test:type=Cat,*")
	require.NoError(t, err)
	require.ElementsMatch(t, expectedMBeanNames, actualMBeanNames)

	expectedMBeanAttrNames := []string{
		"BoolValue", "FloatValue", "NumberValue", "DoubleValue", "Name",
	}
	actualMBeanAttrNames, err := client.GetMBeanAttrNames("test:name=tomas,type=Cat")
	require.NoError(t, err)
	require.ElementsMatch(t, expectedMBeanAttrNames, actualMBeanAttrNames)

	// AND Query returns expected data
	expected := []*nrprotocol.JMXAttribute{
		{
			Attribute: "test:type=Cat,name=tomas,attr=FloatValue",

			ValueType:   nrprotocol.ValueType_DOUBLE,
			DoubleValue: 2.222222,
		},
		{
			Attribute: "test:type=Cat,name=tomas,attr=NumberValue",

			ValueType: nrprotocol.ValueType_INT,
			IntValue:  3,
		},
		{
			Attribute: "test:type=Cat,name=tomas,attr=BoolValue",

			ValueType: nrprotocol.ValueType_BOOL,
			BoolValue: true,
		},
		{
			Attribute: "test:type=Cat,name=tomas,attr=DoubleValue",

			ValueType:   nrprotocol.ValueType_DOUBLE,
			DoubleValue: 1.2,
		},
		{
			Attribute: "test:type=Cat,name=tomas,attr=Name",

			ValueType:   nrprotocol.ValueType_STRING,
			StringValue: "tomas",
		},
	}

	var actual []*nrprotocol.JMXAttribute
	for _, mBeanAttrName := range expectedMBeanAttrNames {
		jmxAttribute, err := client.GetMBeanAttr("test:type=Cat,name=tomas", mBeanAttrName)
		assert.NoError(t, err)
		actual = append(actual, jmxAttribute)
	}

	assert.ElementsMatch(t, expected, actual)
}

func Test_Query_Timeout(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	jmxHost, jmxPort, err := testutils.GetContainerHostAndPort(ctx, container)
	require.NoError(t, err)

	// THEN JMX connection fails
	config := &JMXConfig{
		Hostname:        jmxHost,
		Port:            int32(jmxPort.Int()),
		RequestTimoutMs: 1,
	}
	client, err := NewClient(ctx).Open(config)
	defer client.Close()
	assert.NotNil(t, client)
	assert.Error(t, err)

	// AND Query returns expected error
	actual, err := client.GetMBeanAttrNames("*:*")
	assert.Nil(t, actual)
	assert.Error(t, err)
}

func Test_URL_Success(t *testing.T) {
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

	jmxHost, jmxPort, err := testutils.GetContainerHostAndPort(ctx, container)
	require.NoError(t, err)

	// THEN JMX connection can be oppened
	config := &JMXConfig{
		ConnectionURL:   fmt.Sprintf("service:jmx:rmi:///jndi/rmi://%s:%s/jmxrmi", jmxHost, jmxPort.Port()),
		RequestTimoutMs: testutils.DefaultTimeoutMs,
	}
	client, err := NewClient(ctx).Open(config)
	defer client.Close()

	assert.NoError(t, err)

	// AND Query returns expected data
	actual, err := client.GetMBeanAttr("test:type=Cat,name=tomas", "FloatValue")
	assert.NoError(t, err)

	expected := &nrprotocol.JMXAttribute{
		Attribute:   "test:type=Cat,name=tomas,attr=FloatValue",
		ValueType:   nrprotocol.ValueType_DOUBLE,
		DoubleValue: 2.2,
	}

	assert.Equal(t, expected, actual)
}

func Test_JavaNotInstalled(t *testing.T) {
	// GIVEN a wrong Java Home
	os.Setenv("NRIA_JAVA_HOME", "/wrong/path")
	defer os.Unsetenv("NRIA_JAVA_HOME")

	ctx := context.Background()

	config := &JMXConfig{
		RequestTimoutMs: testutils.DefaultTimeoutMs,
	}
	// THEN connect fails with expected error
	client, err := NewClient(ctx).Open(config)
	assert.Contains(t, err.Error(), "/wrong/path/bin/java")

	// AND Query fails with expected error
	actual, err := client.GetMBeanNames("test:type=Cat,*")
	assert.Nil(t, actual)
	assert.ErrorIs(t, err, nrjmx.ErrProcessNotRunning)
}

func Test_WrongMbeanFormat(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	jmxHost, jmxPort, err := testutils.GetContainerHostAndPort(ctx, container)
	require.NoError(t, err)

	// THEN JMX connection can be oppened
	config := &JMXConfig{
		ConnectionURL:   fmt.Sprintf("service:jmx:rmi:///jndi/rmi://%s:%s/jmxrmi", jmxHost, jmxPort.Port()),
		RequestTimoutMs: testutils.DefaultTimeoutMs,
	}

	client, err := NewClient(ctx).Open(config)
	defer client.Close()
	assert.NoError(t, err)

	// AND Query returns expected error
	actual, err := client.GetMBeanNames("wrong_format")
	assert.Nil(t, actual)

	jmxErr, ok := err.(*nrprotocol.JMXError)
	assert.True(t, ok)
	assert.Equal(t, jmxErr.GetMessage(), "cannot parse MBean glob pattern: 'wrong_format', valid: 'DOMAIN:BEAN'")
}

func Test_Wrong_Connection(t *testing.T) {
	ctx := context.Background()

	// GIVEN a wrong hostname and port
	config := &JMXConfig{
		Hostname:        "localhost",
		Port:            1234,
		RequestTimoutMs: testutils.DefaultTimeoutMs,
	}

	// THEN open fails with expected error
	client, err := NewClient(ctx).Open(config)
	defer client.Close()
	assert.Contains(t, err.Error(), "Connection refused to host: localhost;")

	// AND query returns expected error
	assert.Contains(t, err.Error(), "Connection refused to host: localhost;") // TODO: fix this, doesn't return the correct error

	actual, err := client.GetMBeanNames("test:type=Cat,*")
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
	})
	assert.NoError(t, err)
	assert.Equal(t, "ok!\n", string(resp))
	defer testutils.CleanMBeans(ctx, container)

	jmxHost, jmxPort, err := testutils.GetContainerHostAndPort(ctx, container)
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
		RequestTimoutMs:    testutils.DefaultTimeoutMs,
	}
	client, err := NewClient(ctx).Open(config)
	defer client.Close()
	assert.NoError(t, err)

	// AND Query returns expected data
	expectedMBeanNames := []string{
		"test:type=Cat,name=tomas",
	}
	actualMBeanNames, err := client.GetMBeanNames("test:type=Cat,*")
	require.NoError(t, err)
	require.ElementsMatch(t, expectedMBeanNames, actualMBeanNames)

	expectedMBeanAttrNames := []string{
		"BoolValue", "FloatValue", "NumberValue", "DoubleValue", "Name",
	}
	actualMBeanAttrNames, err := client.GetMBeanAttrNames("test:name=tomas,type=Cat")
	require.NoError(t, err)
	require.ElementsMatch(t, expectedMBeanAttrNames, actualMBeanAttrNames)

	// AND Query returns expected data
	expected := []*nrprotocol.JMXAttribute{
		{
			Attribute: "test:type=Cat,name=tomas,attr=FloatValue",

			ValueType:   nrprotocol.ValueType_DOUBLE,
			DoubleValue: 2.222222,
		},
		{
			Attribute: "test:type=Cat,name=tomas,attr=NumberValue",

			ValueType: nrprotocol.ValueType_INT,
			IntValue:  3,
		},
		{
			Attribute: "test:type=Cat,name=tomas,attr=BoolValue",

			ValueType: nrprotocol.ValueType_BOOL,
			BoolValue: true,
		},
		{
			Attribute: "test:type=Cat,name=tomas,attr=DoubleValue",

			ValueType:   nrprotocol.ValueType_DOUBLE,
			DoubleValue: 1.2,
		},
		{
			Attribute: "test:type=Cat,name=tomas,attr=Name",

			ValueType:   nrprotocol.ValueType_STRING,
			StringValue: "tomas",
		},
	}

	var actual []*nrprotocol.JMXAttribute
	for _, mBeanAttrName := range expectedMBeanAttrNames {
		jmxAttribute, err := client.GetMBeanAttr("test:type=Cat,name=tomas", mBeanAttrName)
		assert.NoError(t, err)
		actual = append(actual, jmxAttribute)
	}

	assert.ElementsMatch(t, expected, actual)
}

func Test_Wrong_Credentials(t *testing.T) {
	ctx := context.Background()

	// GIVEN an SSL JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainerSSL(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	jmxHost, jmxPort, err := testutils.GetContainerHostAndPort(ctx, container)
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
		RequestTimoutMs:    testutils.DefaultTimeoutMs,
	}

	// THEN open fails with expected error
	client, err := NewClient(ctx).Open(config)
	defer client.Close()
	assert.Contains(t, err.Error(), "Authentication failed! Invalid username or password")

	// AND Query returns expected error
	actual, err := client.GetMBeanNames("test:type=Cat,*")
	assert.Nil(t, actual)
	assert.Errorf(t, err, "connection to JMX endpoint is not established")
}

func Test_Wrong_Certificate_password(t *testing.T) {
	ctx := context.Background()

	// GIVEN an SSL JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainerSSL(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	jmxHost, jmxPort, err := testutils.GetContainerHostAndPort(ctx, container)
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
		RequestTimoutMs:    testutils.DefaultTimeoutMs,
	}

	// THEN open returns expected error
	client, err := NewClient(ctx).Open(config)
	defer client.Close()
	assert.Contains(t, err.Error(), "SSLContext") // TODO: improve this error from java

	// AND Query returns expected error
	actual, err := client.GetMBeanNames("test:type=Cat,*")
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
	err = testutils.CopyFileFromContainer(ctx, container.GetContainerID(), "/opt/jboss/wildfly/bin/client/jboss-client.jar", dstFile)
	assert.NoError(t, err)

	defer os.Remove(dstFile)

	jmxHost, jmxPort, err := testutils.GetContainerHostAndPort(ctx, container)
	require.NoError(t, err)

	// THEN JMX connection can be opened
	config := &JMXConfig{
		Hostname:              jmxHost,
		Port:                  int32(jmxPort.Int()),
		Username:              testutils.JbossJMXUsername,
		Password:              testutils.JbossJMXPassword,
		IsJBossStandaloneMode: true,
		IsRemote:              true,
		RequestTimoutMs:       testutils.DefaultTimeoutMs,
	}
	client, err := NewClient(ctx).Open(config)
	defer client.Close()
	assert.NoError(t, err)

	// AND Query returns expected data
	expectedMbeanNames := []string{
		"jboss.as:subsystem=remoting,configuration=endpoint",
	}
	actualMbeanNames, err := client.GetMBeanNames("jboss.as:subsystem=remoting,configuration=endpoint")
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
	actualMBeanAttrNames, err := client.GetMBeanAttrNames("jboss.as:subsystem=remoting,configuration=endpoint")
	assert.ElementsMatch(t, expectedMBeanAttrNames, actualMBeanAttrNames)

	expected := []*nrprotocol.JMXAttribute{
		{
			Attribute: "jboss.as:subsystem=remoting,configuration=endpoint,attr=authenticationRetries",
			ValueType: nrprotocol.ValueType_INT,
			IntValue:  3,
		},
		{
			Attribute: "jboss.as:subsystem=remoting,configuration=endpoint,attr=heartbeatInterval",
			ValueType: nrprotocol.ValueType_INT,
			IntValue:  60000,
		},
		{
			Attribute: "jboss.as:subsystem=remoting,configuration=endpoint,attr=maxInboundChannels",
			ValueType: nrprotocol.ValueType_INT,
			IntValue:  40,
		},
		{
			Attribute: "jboss.as:subsystem=remoting,configuration=endpoint,attr=maxInboundMessageSize",
			ValueType: nrprotocol.ValueType_INT,
			IntValue:  9223372036854775807,
		},
		{
			Attribute: "jboss.as:subsystem=remoting,configuration=endpoint,attr=maxInboundMessages",
			ValueType: nrprotocol.ValueType_INT,
			IntValue:  80,
		},
		{
			Attribute: "jboss.as:subsystem=remoting,configuration=endpoint,attr=maxOutboundChannels",
			ValueType: nrprotocol.ValueType_INT,
			IntValue:  40,
		},
		{
			Attribute: "jboss.as:subsystem=remoting,configuration=endpoint,attr=maxOutboundMessageSize",
			ValueType: nrprotocol.ValueType_INT,
			IntValue:  9223372036854775807,
		},
		{
			Attribute: "jboss.as:subsystem=remoting,configuration=endpoint,attr=maxOutboundMessages",
			ValueType: nrprotocol.ValueType_INT,
			IntValue:  65535,
		},
		{
			Attribute: "jboss.as:subsystem=remoting,configuration=endpoint,attr=receiveBufferSize",
			ValueType: nrprotocol.ValueType_INT,
			IntValue:  8192,
		},
		{
			Attribute: "jboss.as:subsystem=remoting,configuration=endpoint,attr=receiveWindowSize",
			ValueType: nrprotocol.ValueType_INT,
			IntValue:  131072,
		},
		{
			Attribute:   "jboss.as:subsystem=remoting,configuration=endpoint,attr=saslProtocol",
			ValueType:   nrprotocol.ValueType_STRING,
			StringValue: "remote",
		},
		{
			Attribute: "jboss.as:subsystem=remoting,configuration=endpoint,attr=sendBufferSize",
			ValueType: nrprotocol.ValueType_INT,
			IntValue:  8192,
		},
		{
			Attribute: "jboss.as:subsystem=remoting,configuration=endpoint,attr=transmitWindowSize",
			ValueType: nrprotocol.ValueType_INT,
			IntValue:  131072,
		},
		{
			Attribute:   "jboss.as:subsystem=remoting,configuration=endpoint,attr=worker",
			ValueType:   nrprotocol.ValueType_STRING,
			StringValue: "default",
		},
	}

	var actual []*nrprotocol.JMXAttribute
	for _, mBeanAttrName := range expectedMBeanAttrNames {
		jmxAttribute, err := client.GetMBeanAttr("jboss.as:subsystem=remoting,configuration=endpoint", mBeanAttrName)
		if err != nil {
			continue
		}
		actual = append(actual, jmxAttribute)
	}

	assert.ElementsMatch(t, expected, actual)
}

func TestJMXServiceClose(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	jmxHost, jmxPort, err := testutils.GetContainerHostAndPort(ctx, container)
	assert.NoError(t, err)

	// THEN JMX connection can be opened.
	config := &JMXConfig{
		Hostname:        jmxHost,
		Port:            int32(jmxPort.Int()),
		RequestTimoutMs: testutils.DefaultTimeoutMs,
	}

	client, err := NewClient(ctx).Open(config)
	assert.NoError(t, err)
	err = client.Close()
	assert.NoError(t, err)

	// AND Query returns expected error
	actual, err := client.GetMBeanNames("*:*")
	assert.Nil(t, actual)
	assert.ErrorIs(t, err, nrjmx.ErrProcessNotRunning) // TODO: Get valid error message

	assert.True(t, client.nrJMXProcess.GetOSProcessState().Success())
}

func TestJavaProcessExits(t *testing.T) {
	// gojmx starts nrjmx bash script which stats a java process.
	// We want to make sure that if gojmx process dies, java process stops also.
	// To reproduce this scenario, we run the current test twice:
	// - subprocess SHOULD_RUN_EXIT = 1
	// - main test SHOULD_RUN_EXIT unset.
	if os.Getenv("IS_SUBPROCESS") == "1" {
		ctx := context.Background()

		// GIVEN a JMX Server running inside a container
		container, err := testutils.RunJMXServiceContainer(ctx)
		require.NoError(t, err)
		defer container.Terminate(ctx)

		jmxHost, jmxPort, err := testutils.GetContainerHostAndPort(ctx, container)
		require.NoError(t, err)

		// THEN JMX connection can be opened
		config := &JMXConfig{
			Hostname:        jmxHost,
			Port:            int32(jmxPort.Int()),
			RequestTimoutMs: testutils.DefaultTimeoutMs,
		}
		client, err := NewClient(ctx).Open(config)
		require.NoError(t, err)

		f, err := os.OpenFile(os.Getenv("TMP_FILE"), os.O_WRONLY|os.O_TRUNC, 0644)
		require.NoError(t, err)
		defer f.Close()
		_, err = fmt.Fprintf(f, "%d\n", client.nrJMXProcess.GetPID())
		require.NoError(t, err)
		<-time.After(5 * time.Minute)
	}

	// Create a temporary file to communicate get the pid of the subprocess.
	tmpFile, err := ioutil.TempFile("", "nrjmx_pid_test")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Run again the current test function with IS_SUBPROCESS=1
	cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
	cmd.Env = append(os.Environ(), "IS_SUBPROCESS=1", "TMP_FILE="+tmpFile.Name())

	stdErrBuffer := nrjmx.NewDefaultLimitedBuffer()
	stdOutBuffer := nrjmx.NewDefaultLimitedBuffer()
	cmd.Stderr = stdErrBuffer
	cmd.Stdout = stdOutBuffer
	err = cmd.Start()

	ctx, waitCancel := context.WithCancel(context.Background())
	go func() {
		err := cmd.Wait()
		if ctx.Err() == nil && err != nil {
			assert.NoError(t, err)
			panic(fmt.Errorf("stdout: %s\nstderr: %s", stdOutBuffer, stdErrBuffer))
		}
	}()

	// Get the pid from the subprocess.
	var pid int32
	require.Eventually(t, func() bool {
		var err error
		pid, err = testutils.ReadPidFile(tmpFile.Name())
		if err != nil {
			return false
		}
		return true
	}, 30*time.Second, 50*time.Millisecond)

	p, err := process.NewProcess(pid)
	assert.NoError(t, err)
	ch, err := p.Children()
	assert.NoError(t, err)

	// Stop listening for cmd.Wait() error. We want to kill the subprocess, at this point an error is expected.
	waitCancel()
	err = cmd.Process.Kill()
	require.NoError(t, err)

	// Check that java pid does not exist anymore.
	require.Eventually(t, func() bool {
		up, err := ch[0].IsRunning()
		if err != nil {
			return false
		}
		// assert is not running
		return !up
	}, 5*time.Second, 50*time.Millisecond)
}
