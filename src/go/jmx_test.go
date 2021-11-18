package main

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

	"github.com/docker/go-connections/nat"
	"github.com/newrelic/infra-integrations-sdk/jmx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	testServerHost                 = "localhost"
	testServerPort                 = "4567"
	testServerJMXPort              = "7199"
	testServerAddDataEndpoint      = "/cat"
	testServerAddDataBatchEndpoint = "/cat_batch"
	testServerCleanDataEndpoint    = "/clear"
	keystorePassword               = "password"
	truststorePassword             = "password"
	jmxUsername                    = "testuser"
	jmxPassword                    = "testpassword"
)

var prjDir, keystorePath, truststorepath string

func init() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	prjDir = filepath.Join(path, "../..")
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

	for i := 0; i < 2000; i++ {
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

	// THEN JMX connection can be oppened
	jmxPort, err := container.MappedPort(ctx, testServerJMXPort)
	require.NoError(t, err)
	jmxHost, err := container.Host(ctx)
	require.NoError(t, err)

	err = jmx.Open(jmxHost, jmxPort.Port(), "", "")
	defer jmx.Close()
	assert.NoError(t, err)

	result, err := jmx.Query("test:type=Cat,*", 600000)
	assert.NoError(t, err)

	// AND query returns at least 5Mb of data.
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
		"floatValue":  2.2,
		"numberValue": 3,
		"boolValue":   true,
	})
	assert.NoError(t, err)
	assert.Equal(t, "ok!\n", string(resp))

	defer cleanMBeans(ctx, container)

	// THEN JMX connection can be oppened
	jmxPort, err := container.MappedPort(ctx, testServerJMXPort)
	require.NoError(t, err)
	jmxHost, err := container.Host(ctx)
	require.NoError(t, err)

	err = jmx.Open(jmxHost, jmxPort.Port(), "", "")
	defer jmx.Close()
	assert.NoError(t, err)

	// AND Query returns expected data
	result, err := jmx.Query("test:type=Cat,*", 10000)
	assert.NoError(t, err)

	expected := map[string]interface{}{
		"test:type=Cat,name=tomas,attr=Name":        "tomas",
		"test:type=Cat,name=tomas,attr=DoubleValue": 1.2,
		"test:type=Cat,name=tomas,attr=FloatValue":  2.2,
		"test:type=Cat,name=tomas,attr=BoolValue":   true,
		"test:type=Cat,name=tomas,attr=NumberValue": float64(3),
	}

	assert.EqualValues(t, result, expected)
}

func Test_JavaNotInstalled(t *testing.T) {
	// GIVEN a wrong Java Home
	os.Setenv("NRIA_JAVA_HOME", "/wrong/path")
	defer os.Unsetenv("NRIA_JAVA_HOME")

	err := jmx.Open("wrong", "12345", "", "")
	defer jmx.Close()
	assert.NoError(t, err)

	// WHEN connecting
	resp, err := jmx.Query("test:type=Cat,*", 10000)
	assert.Nil(t, resp)

	// Error is returned
	assert.EqualError(t, err, "EOF") // TODO: this error message should be fixed
}

func Test_Query_WrongFormat(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := runJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	// THEN JMX connection can be oppened
	jmxPort, err := container.MappedPort(ctx, testServerJMXPort)
	require.NoError(t, err)
	jmxHost, err := container.Host(ctx)
	require.NoError(t, err)

	err = jmx.Open(jmxHost, jmxPort.Port(), "", "")
	defer jmx.Close()
	assert.NoError(t, err)

	// AND Query returns expected data
	result, err := jmx.Query("wrong_format", 10000)
	assert.Nil(t, result)
	assert.EqualError(t, err, "cannot parse MBean glob pattern, valid: 'DOMAIN:BEAN'")
}

func Test_Wrong_Hostname(t *testing.T) {
	// GIVEN a wrong hostname and port
	err := jmx.Open("wrong", "12345", "", "")
	defer jmx.Close()

	assert.NoError(t, err)
	// WHEN connecting
	resp, err := jmx.Query("test:type=Cat,*", 10000)
	assert.Nil(t, resp)

	// Error is returned
	assert.EqualError(t, err, "jmx endpoint connection error")
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
		"floatValue":  2.2,
		"numberValue": 3,
		"boolValue":   true,
	})
	assert.NoError(t, err)
	assert.Equal(t, "ok!\n", string(resp))
	defer cleanMBeans(ctx, container)

	// THEN SSL JMX connection can be oppened
	options := []jmx.Option{}
	ssl := jmx.WithSSL(keystorePath, keystorePassword, truststorepath, truststorePassword)
	options = append(options, ssl)

	jmxPort, err := container.MappedPort(ctx, testServerJMXPort)
	require.NoError(t, err)
	jmxHost, err := container.Host(ctx)
	require.NoError(t, err)

	err = jmx.Open(jmxHost, jmxPort.Port(), jmxUsername, jmxPassword, options...)
	assert.NoError(t, err)

	// AND Query returns expected data
	result, err := jmx.Query("test:type=Cat,*", 10000)
	defer jmx.Close()

	assert.NoError(t, err)

	expected := map[string]interface{}{
		"test:type=Cat,name=tomas,attr=Name":        "tomas",
		"test:type=Cat,name=tomas,attr=DoubleValue": 1.2,
		"test:type=Cat,name=tomas,attr=FloatValue":  2.2,
		"test:type=Cat,name=tomas,attr=BoolValue":   true,
		"test:type=Cat,name=tomas,attr=NumberValue": float64(3),
	}

	assert.EqualValues(t, result, expected)
}

func Test_Wrong_Credentials(t *testing.T) {
	ctx := context.Background()

	// GIVEN an SSL JMX Server running inside a container
	container, err := runJMXServiceContainerSSL(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	// WHEN wrong jmx username and password is provided
	options := []jmx.Option{}
	ssl := jmx.WithSSL(keystorePath, keystorePassword, truststorepath, truststorePassword)
	options = append(options, ssl)

	jmxPort, err := container.MappedPort(ctx, testServerJMXPort)
	require.NoError(t, err)
	jmxHost, err := container.Host(ctx)
	require.NoError(t, err)

	err = jmx.Open(jmxHost, jmxPort.Port(), "wrongUser", "wrongPassword", options...)
	assert.NoError(t, err)

	// THEN error is returned
	result, err := jmx.Query("test:type=Cat,*", 10000)
	defer jmx.Close()
	assert.Nil(t, result)

	assert.EqualError(t, err, " error running nrjmx: Authentication failed! Invalid username or password\n")
}

func Test_Wrong_Certificate_password(t *testing.T) {
	ctx := context.Background()

	// GIVEN an SSL JMX Server running inside a container
	container, err := runJMXServiceContainerSSL(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	// WHEN wrong wrong certificate password is provided
	options := []jmx.Option{}
	ssl := jmx.WithSSL(keystorePath, "wrong_password", truststorepath, truststorePassword)
	options = append(options, ssl)

	jmxPort, err := container.MappedPort(ctx, testServerJMXPort)
	require.NoError(t, err)
	jmxHost, err := container.Host(ctx)
	require.NoError(t, err)

	err = jmx.Open(jmxHost, jmxPort.Port(), jmxUsername, jmxPassword, options...)
	assert.NoError(t, err)

	// THEN error is returned
	result, err := jmx.Query("test:type=Cat,*", 10000)
	defer jmx.Close()
	assert.Nil(t, result)

	assert.EqualError(t, err, "jmx endpoint connection error")
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
