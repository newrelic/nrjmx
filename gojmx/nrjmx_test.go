package gojmx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/newrelic/nrjmx/gojmx/nrprotocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	testServerHost                 = "localhost"
	testServerPort                 = "4567"
	testServerJMXPort              = "7199"
	jbossJMXPort                   = "9990"
	jbossJMXUsername               = "admin1234"
	jbossJMXPassword               = "Password1!"
	testServerAddDataEndpoint      = "/cat"
	testServerAddDataBatchEndpoint = "/cat_batch"
	testServerCleanDataEndpoint    = "/clear"
	keystorePassword               = "password"
	truststorePassword             = "password"
	jmxUsername                    = "testuser"
	jmxPassword                    = "testpassword"
	defaultTimeoutMs 			   = 5000
)

var prjDir, keystorePath, truststorepath string

func init() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	prjDir = filepath.Join(path, "..")
	keystorePath = filepath.Join(prjDir, "test-server", "keystore")
	truststorepath = filepath.Join(prjDir, "test-server", "truststore")

	os.Setenv("NR_JMX_TOOL", filepath.Join(prjDir, "bin", "nrjmx"))
}

func Test_Query_Success_LargeAmountOfData(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := runJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	data := []map[string]interface{}{}

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
	resp, err := addMBeansBatch(ctx, container, data)
	assert.NoError(t, err)
	assert.Equal(t, "ok!\n", string(resp))

	defer cleanMBeans(ctx, container)

	jmxPort, err := container.MappedPort(ctx, testServerJMXPort)
	require.NoError(t, err)
	jmxHost, err := container.Host(ctx)
	require.NoError(t, err)

	// THEN JMX connection can be oppened
	client, err := NewJMXServiceClient(ctx)
	assert.NoError(t, err)

	config := &nrprotocol.JMXConfig{
		Hostname: jmxHost,
		Port:     int32(jmxPort.Int()),
		UriPath:  "jmxrmi",
	}

	err = client.Connect(ctx, config, -1)
	assert.NoError(t, err)
	defer client.Disconnect(ctx)

	// AND query returns at least 5Mb of data.
	result, err := client.QueryMbean(ctx, "test:type=Cat,*", -1)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(fmt.Sprintf("%v", result)), 5*1024*1024)
}

func Test_Query_Success(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := runJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	// Populate the JMX Server with mbeans
	resp, err := addMBeans(ctx, container, map[string]interface{}{
		"name":        "tomas",
		"doubleValue": 1.2,
		"floatValue":  2.2222222,
		"numberValue": 3,
		"boolValue":   true,
	})
	assert.NoError(t, err)
	assert.Equal(t, "ok!\n", string(resp))

	defer cleanMBeans(ctx, container)

	jmxPort, err := container.MappedPort(ctx, testServerJMXPort)
	require.NoError(t, err)
	jmxHost, err := container.Host(ctx)
	require.NoError(t, err)

	// THEN JMX connection can be oppened
	client, err := NewJMXServiceClient(ctx)
	assert.NoError(t, err)

	config := &nrprotocol.JMXConfig{
		Hostname: jmxHost,
		Port:     int32(jmxPort.Int()),
		UriPath:  "jmxrmi",
	}

	err = client.Connect(ctx, config, defaultTimeoutMs)
	defer client.Disconnect(ctx)
	assert.NoError(t, err)

	// AND Query returns expected data
	actual, err := client.QueryMbean(ctx, "test:type=Cat,*", defaultTimeoutMs)
	assert.NoError(t, err)

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

	assert.Equal(t, expected, actual)
}

func Test_Query_Timeout(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := runJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	jmxPort, err := container.MappedPort(ctx, testServerJMXPort)
	require.NoError(t, err)
	jmxHost, err := container.Host(ctx)
	require.NoError(t, err)

	// THEN JMX connection can be oppened
	client, err := NewJMXServiceClient(ctx)
	assert.NoError(t, err)

	config := &nrprotocol.JMXConfig{
		Hostname: jmxHost,
		Port:     int32(jmxPort.Int()),
		UriPath:  "jmxrmi",
	}

	err = client.Connect(ctx, config, defaultTimeoutMs)
	defer client.Disconnect(ctx)
	assert.NoError(t, err)

	// AND Query returns expected data
	actual, err := client.Query(1, "*:*")
	assert.Nil(t, actual)

	assert.Error(t, err)
}

func Test_URL_Success(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := runJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	// Populate the JMX Server with mbeans
	resp, err := addMBeans(ctx, container, map[string]interface{}{
		"name":        "tomas",
		"doubleValue": 1.2,
		"floatValue":  2.2,
		"numberValue": 3,
		"boolValue":   true,
	})
	assert.NoError(t, err)
	assert.Equal(t, "ok!\n", string(resp))
	defer cleanMBeans(ctx, container)

	jmxPort, err := container.MappedPort(ctx, testServerJMXPort)
	require.NoError(t, err)
	jmxHost, err := container.Host(ctx)
	require.NoError(t, err)

	// THEN JMX connection can be oppened
	client, err := NewJMXServiceClient(ctx)
	assert.NoError(t, err)

	config := &nrprotocol.JMXConfig{
		ConnectionURL: fmt.Sprintf("service:jmx:rmi:///jndi/rmi://%s:%s/jmxrmi", jmxHost, jmxPort.Port()),
	}

	err = client.Connect(ctx, config, defaultTimeoutMs)
	defer client.Disconnect(ctx)
	assert.NoError(t, err)

	// AND Query returns expected data
	actual, err := client.QueryMbean(ctx, "test:type=Cat,*", defaultTimeoutMs)
	assert.NoError(t, err)

	expected := []*nrprotocol.JMXAttribute{
		{
			Attribute:   "test:type=Cat,name=tomas,attr=FloatValue",
			ValueType:   nrprotocol.ValueType_DOUBLE,
			DoubleValue: 2.2,
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
			Attribute:   "test:type=Cat,name=tomas,attr=DoubleValue",
			ValueType:   nrprotocol.ValueType_DOUBLE,
			DoubleValue: 1.2,
		},
		{
			Attribute:   "test:type=Cat,name=tomas,attr=Name",
			ValueType:   nrprotocol.ValueType_STRING,
			StringValue: "tomas",
		},
	}

	assert.Equal(t, expected, actual)
}

func Test_JavaNotInstalled(t *testing.T) {
	// GIVEN a wrong Java Home
	os.Setenv("NRIA_JAVA_HOME", "/wrong/path")
	defer os.Unsetenv("NRIA_JAVA_HOME")

	ctx := context.Background()
	client, err := NewJMXServiceClient(ctx)
	assert.NoError(t, err)

	config := &nrprotocol.JMXConfig{}

	// THEN connect fails with expected error
	err = client.Connect(ctx, config, defaultTimeoutMs)
	defer client.Disconnect(ctx)
	assert.EqualError(t, err, "EOF") // TODO: this error message should be fixed

	// AND Query fails with expected error
	actual, err := client.QueryMbean(ctx, "test:type=Cat,*", defaultTimeoutMs)
	assert.Nil(t, actual)
	assert.EqualError(t, err, "write |1: broken pipe") // TODO: this error message should be fixed
}

func Test_WrongMbeanFormat(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := runJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	jmxPort, err := container.MappedPort(ctx, testServerJMXPort)
	require.NoError(t, err)
	jmxHost, err := container.Host(ctx)
	require.NoError(t, err)

	// THEN JMX connection can be oppened
	client, err := NewJMXServiceClient(ctx)
	assert.NoError(t, err)

	config := &nrprotocol.JMXConfig{
		ConnectionURL: fmt.Sprintf("service:jmx:rmi:///jndi/rmi://%s:%s/jmxrmi", jmxHost, jmxPort.Port()),
	}

	err = client.Connect(ctx, config, defaultTimeoutMs)
	defer client.Disconnect(ctx)
	assert.NoError(t, err)

	// AND Query returns expected error
	actual, err := client.QueryMbean(ctx, "wrong_format", defaultTimeoutMs)
	assert.Nil(t, actual)

	assert.EqualError(t, err, "cannot parse MBean glob pattern, valid: 'DOMAIN:BEAN'") //TODO: return the correct error from java to match this message.
}

func Test_Wrong_Connection(t *testing.T) {
	ctx := context.Background()

	client, err := NewJMXServiceClient(ctx)
	assert.NoError(t, err)

	// GIVEN a wrong hostname and port
	config := &nrprotocol.JMXConfig{
		Hostname: "localhost",
		Port:     1234,
		UriPath:  "jmxrmi",
	}

	// WHEN connecting expected error is returned
	err = client.Connect(ctx, config, defaultTimeoutMs)
	defer client.Disconnect(ctx)
	assert.Contains(t, err.Error(), "Connection refused to host: localhost;")

	// AND query returns expected error
	actual, err := client.QueryMbean(ctx, "test:type=Cat,*", defaultTimeoutMs)
	assert.Nil(t, actual)
	assert.Contains(t, err.Error(), "Connection refused to host: localhost;") // TODO: fix this, doesn't return the correct error
}

func Test_SSLQuery_Success(t *testing.T) {
	ctx := context.Background()

	// GIVEN an SSL JMX Server running inside a container
	container, err := runJMXServiceContainerSSL(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	// Populate the JMX Server with mbeans
	resp, err := addMBeans(ctx, container, map[string]interface{}{
		"name":        "tomas",
		"doubleValue": 1.2,
		"floatValue":  2.222222,
		"numberValue": 3,
		"boolValue":   true,
	})
	assert.NoError(t, err)
	assert.Equal(t, "ok!\n", string(resp))
	defer cleanMBeans(ctx, container)

	jmxPort, err := container.MappedPort(ctx, testServerJMXPort)
	require.NoError(t, err)
	jmxHost, err := container.Host(ctx)
	require.NoError(t, err)

	// THEN SSL JMX connection can be oppened
	client, err := NewJMXServiceClient(ctx)
	assert.NoError(t, err)

	config := &nrprotocol.JMXConfig{
		Hostname:           jmxHost,
		Port:               int32(jmxPort.Int()),
		UriPath:            "jmxrmi",
		Username:           jmxUsername,
		Password:           jmxPassword,
		KeyStore:           keystorePath,
		KeyStorePassword:   keystorePassword,
		TrustStore:         truststorepath,
		TrustStorePassword: truststorePassword,
	}

	err = client.Connect(ctx, config, defaultTimeoutMs)
	defer client.Disconnect(ctx)
	assert.NoError(t, err)

	// AND Query returns expected data
	actual, err := client.QueryMbean(ctx, "test:type=Cat,*", defaultTimeoutMs)
	assert.NoError(t, err)

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

	assert.Equal(t, expected, actual)
}

func Test_Wrong_Credentials(t *testing.T) {
	ctx := context.Background()

	// GIVEN an SSL JMX Server running inside a container
	container, err := runJMXServiceContainerSSL(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	jmxPort, err := container.MappedPort(ctx, testServerJMXPort)
	require.NoError(t, err)
	jmxHost, err := container.Host(ctx)
	require.NoError(t, err)

	// WHEN wrong jmx username and password is provided
	client, err := NewJMXServiceClient(ctx)
	assert.NoError(t, err)

	config := &nrprotocol.JMXConfig{
		Hostname:           jmxHost,
		Port:               int32(jmxPort.Int()),
		UriPath:            "jmxrmi",
		Username:           "wrong_username",
		Password:           "wrong_password",
		KeyStore:           keystorePath,
		KeyStorePassword:   keystorePassword,
		TrustStore:         truststorepath,
		TrustStorePassword: truststorePassword,
	}

	// THEN connect fails with expected error
	err = client.Connect(ctx, config, defaultTimeoutMs)
	defer client.Disconnect(ctx)
	assert.Contains(t, err.Error(), "Authentication failed! Invalid username or password")

	// AND Query returns expected error
	actual, err := client.QueryMbean(ctx, "test:type=Cat,*", defaultTimeoutMs)
	assert.Nil(t, actual)
	assert.Contains(t, err.Error(), "Authentication failed! Invalid username or password") // TODO: fix this in java tool, as it doesn't return the correct error
}

func Test_Wrong_Certificate_password(t *testing.T) {
	ctx := context.Background()

	// GIVEN an SSL JMX Server running inside a container
	container, err := runJMXServiceContainerSSL(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	jmxPort, err := container.MappedPort(ctx, testServerJMXPort)
	require.NoError(t, err)
	jmxHost, err := container.Host(ctx)
	require.NoError(t, err)

	// WHEN wrong jmx username and password is provided
	client, err := NewJMXServiceClient(ctx)
	assert.NoError(t, err)

	config := &nrprotocol.JMXConfig{
		Hostname:           jmxHost,
		Port:               int32(jmxPort.Int()),
		UriPath:            "jmxrmi",
		Username:           jmxUsername,
		Password:           jmxPassword,
		KeyStore:           keystorePath,
		KeyStorePassword:   "wrong_password",
		TrustStore:         truststorepath,
		TrustStorePassword: truststorePassword,
	}

	// THEN Connect returns expected error
	err = client.Connect(ctx, config, defaultTimeoutMs)
	defer client.Disconnect(ctx)
	assert.Contains(t, err.Error(), "SSLContext") // TODO: improve this error from java

	// AND Query returns expected error
	actual, err := client.QueryMbean(ctx, "test:type=Cat,*", defaultTimeoutMs)
	assert.Nil(t, actual)
	assert.Contains(t, err.Error(), "SSLContext") // TODO: improve this error from java
}

func Test_Connector_Success(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JBoss Server with JMX exposed running inside a container
	container, err := runJbossStandaloneJMXContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	// Install the connector
	dstFile := filepath.Join(prjDir, "/bin/jboss-client.jar")
	err = copyFileFromContainer(ctx, container.GetContainerID(), "/opt/jboss/wildfly/bin/client/jboss-client.jar", dstFile)
	assert.NoError(t, err)

	defer os.Remove(dstFile)

	jmxPort, err := container.MappedPort(ctx, jbossJMXPort)
	require.NoError(t, err)
	jmxHost, err := container.Host(ctx)
	require.NoError(t, err)

	// THEN JMX connection can be oppened
	client, err := NewJMXServiceClient(ctx)
	assert.NoError(t, err)

	config := &nrprotocol.JMXConfig{
		Hostname:              jmxHost,
		Port:                  int32(jmxPort.Int()),
		Username:              jbossJMXUsername,
		Password:              jbossJMXPassword,
		IsJBossStandaloneMode: true,
		IsRemote:              true,
	}

	err = client.Connect(ctx, config, defaultTimeoutMs)
	defer client.Disconnect(ctx)
	assert.NoError(t, err)

	// AND Query returns expected data
	actual, err := client.QueryMbean(ctx, "jboss.as:subsystem=remoting,configuration=endpoint", defaultTimeoutMs)
	assert.NoError(t, err)

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

	assert.Equal(t, expected, actual)
}

// runJMXServiceContainer will start a container running test-server with JMX.
func runJMXServiceContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image: "test-server:latest",
		ExposedPorts: []string{
			fmt.Sprintf("%[1]s:%[1]s", testServerPort),
			fmt.Sprintf("%[1]s:%[1]s", testServerJMXPort),
		},
		Env: map[string]string{
			"JAVA_OPTS": "-Dcom.sun.management.jmxremote.port=" + testServerJMXPort + " " +
				"-Dcom.sun.management.jmxremote.authenticate=false " +
				"-Dcom.sun.management.jmxremote.ssl=false " +
				"-Dcom.sun.management.jmxremote=true " +
				"-Dcom.sun.management.jmxremote.rmi.port=" + testServerJMXPort + " " +
				"-Djava.rmi.server.hostname=localhost",
		},

		WaitingFor: wait.ForListeningPort(testServerPort),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, err
	}

	container.StartLogProducer(ctx)
	container.FollowOutput(&TestLogConsumer{})
	return container, err
}

// runJMXServiceContainerSSL will start a container running test-server configured with SSL JMX.
func runJMXServiceContainerSSL(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image: "test-server:latest",
		ExposedPorts: []string{
			fmt.Sprintf("%[1]s:%[1]s", testServerPort),
			fmt.Sprintf("%[1]s:%[1]s", testServerJMXPort),
		},
		Env: map[string]string{
			"JAVA_OPTS": "-Dcom.sun.management.jmxremote.port=" + testServerJMXPort + " " +
				"-Dcom.sun.management.jmxremote.authenticate=true " +
				"-Dcom.sun.management.jmxremote.ssl=true " +
				"-Dcom.sun.management.jmxremote.ssl.need.client.auth=true " +
				"-Dcom.sun.management.jmxremote.registry.ssl=true " +
				"-Dcom.sun.management.jmxremote=true " +
				"-Dcom.sun.management.jmxremote.rmi.port=" + testServerJMXPort + " " +
				"-Djava.rmi.server.hostname=0.0.0.0 " +
				"-Djavax.net.ssl.keyStore=/keystore  " +
				"-Djavax.net.ssl.keyStorePassword=password " +
				"-Djavax.net.ssl.trustStore=/truststore " +
				"-Djavax.net.ssl.trustStorePassword=password",
		},
		WaitingFor: wait.ForListeningPort(testServerPort),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, err
	}

	container.StartLogProducer(ctx)
	container.FollowOutput(&TestLogConsumer{})
	return container, err
}

// getContainerServiceURL will return the url to the test-server running inside the container.
func getContainerServiceURL(ctx context.Context, container testcontainers.Container, port nat.Port, endpoint string) (string, error) {
	mappedPort, err := container.MappedPort(ctx, port)
	if err != nil {
		return "", err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("http://%s:%s%s", hostIP, mappedPort.Port(), endpoint), nil
}

// cleanMBeans will remove all new added MBeans from test-server.
func cleanMBeans(ctx context.Context, container testcontainers.Container) ([]byte, error) {
	url, err := getContainerServiceURL(ctx, container, testServerPort, testServerCleanDataEndpoint)
	if err != nil {
		return nil, err
	}
	return httpRequest(http.MethodPut, url, nil)
}

// addMBeansBatch will add new MBeans to the test-server.
func addMBeansBatch(ctx context.Context, container testcontainers.Container, body []map[string]interface{}) ([]byte, error) {
	url, err := getContainerServiceURL(ctx, container, testServerPort, testServerAddDataBatchEndpoint)
	if err != nil {
		return nil, err
	}
	json, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return httpRequest(http.MethodPost, url, json)
}

// addMBeans will add new MBeans to the test-server.
func addMBeans(ctx context.Context, container testcontainers.Container, body map[string]interface{}) ([]byte, error) {
	url, err := getContainerServiceURL(ctx, container, testServerPort, testServerAddDataEndpoint)
	if err != nil {
		return nil, err
	}
	json, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return httpRequest(http.MethodPost, url, json)
}

// httpRequest will perform the http request.
func httpRequest(method, url string, body []byte) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("request returned error, status code: %d", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

// TestLogConsumer is used to print container logs to stdout.
type TestLogConsumer struct {
}

func (g *TestLogConsumer) Accept(l testcontainers.Log) {
	fmt.Fprintf(os.Stdout, "[CONTAINER LOG] %s %s\n", time.Now().Format("2006/01/02 15:04:05"), l.Content)
}

// runJbossStandaloneJMXContainer will start a container running a jboss instace with JMX.
func runJbossStandaloneJMXContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image: "test_jboss",

		ExposedPorts: []string{
			fmt.Sprintf("%[1]s:%[1]s", jbossJMXPort),
		},

		WaitingFor: wait.ForListeningPort(jbossJMXPort),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, err
	}

	container.StartLogProducer(ctx)
	container.FollowOutput(&TestLogConsumer{})
	return container, err
}

func copyFileFromContainer(ctx context.Context, containerID, srcPath, dstPath string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	reader, containerPathStat, err := cli.CopyFromContainer(ctx, containerID, srcPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	if !containerPathStat.Mode.IsRegular() {
		return fmt.Errorf("src is not a regular file: %s", srcPath)
	}

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dstPath, b, 0644)
}