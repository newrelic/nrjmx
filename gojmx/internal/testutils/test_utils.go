/*
 * Copyright 2021 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package testutils

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	TestServerPort                     = "4567"
	TestServerJMXPort                  = "7199"
	JbossJMXPort                       = "9990"
	JbossJMXUsername                   = "admin1234"
	JbossJMXPassword                   = "Password1!"
	TestServerAddDataEndpoint          = "/cat"
	TestServerAddDataBatchEndpoint     = "/cat_batch"
	TestServerAddCompositeDataEndpoint = "/composite_data_cat"
	TestServerCleanDataEndpoint        = "/clear"
	KeystorePassword                   = "password"
	TruststorePassword                 = "password"
	JmxUsername                        = "testuser"
	JmxPassword                        = "testpassword"
	DefaultTimeoutMs                   = 5000
)

var (
	PrjDir         string
	KeystorePath   string
	TruststorePath string
)

func init() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Configure tests to point to the project's nrjmx build instead of the regular installation.
	PrjDir = filepath.Join(path, "..")
	KeystorePath = filepath.Join(PrjDir, "test-server", "keystore")
	TruststorePath = filepath.Join(PrjDir, "test-server", "truststore")
}

// RunJMXServiceContainer will start a container running test-server with JMX.
func RunJMXServiceContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image: "test-server:latest",
		ExposedPorts: []string{
			fmt.Sprintf("%[1]s:%[1]s", TestServerPort),
			fmt.Sprintf("%[1]s:%[1]s", TestServerJMXPort),
		},
		Env: map[string]string{
			"JAVA_OPTS": "-Dcom.sun.management.jmxremote.port=" + TestServerJMXPort + " " +
				"-Dcom.sun.management.jmxremote.authenticate=false " +
				"-Dcom.sun.management.jmxremote.ssl=false " +
				"-Dcom.sun.management.jmxremote=true " +
				"-Dcom.sun.management.jmxremote.rmi.port=" + TestServerJMXPort + " " +
				"-Djava.rmi.server.hostname=localhost",
		},

		WaitingFor: wait.ForListeningPort(TestServerPort),
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

// RunJMXServiceContainerSSL will start a container running test-server configured with SSL JMX.
func RunJMXServiceContainerSSL(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image: "test-server:latest",
		ExposedPorts: []string{
			fmt.Sprintf("%[1]s:%[1]s", TestServerPort),
			fmt.Sprintf("%[1]s:%[1]s", TestServerJMXPort),
		},
		Env: map[string]string{
			"JAVA_OPTS": "-Dcom.sun.management.jmxremote.port=" + TestServerJMXPort + " " +
				"-Dcom.sun.management.jmxremote.authenticate=true " +
				"-Dcom.sun.management.jmxremote.ssl=true " +
				"-Dcom.sun.management.jmxremote.ssl.need.client.auth=true " +
				"-Dcom.sun.management.jmxremote.registry.ssl=true " +
				"-Dcom.sun.management.jmxremote=true " +
				"-Dcom.sun.management.jmxremote.rmi.port=" + TestServerJMXPort + " " +
				"-Djava.rmi.server.hostname=0.0.0.0 " +
				"-Djavax.net.ssl.keyStore=/keystore  " +
				"-Djavax.net.ssl.keyStorePassword=password " +
				"-Djavax.net.ssl.trustStore=/truststore " +
				"-Djavax.net.ssl.trustStorePassword=password",
		},
		WaitingFor: wait.ForListeningPort(TestServerPort),
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

// GetContainerServiceURL will return the url to the test-server running inside the container.
func GetContainerServiceURL(ctx context.Context, container testcontainers.Container, port nat.Port, endpoint string) (string, error) {
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

// CleanMBeans will remove all new added MBeans from test-server.
func CleanMBeans(ctx context.Context, container testcontainers.Container) ([]byte, error) {
	url, err := GetContainerServiceURL(ctx, container, TestServerPort, TestServerCleanDataEndpoint)
	if err != nil {
		return nil, err
	}
	return DoHttpRequest(http.MethodPut, url, nil)
}

// AddMBeansBatch will add new MBeans to the test-server.
func AddMBeansBatch(ctx context.Context, container testcontainers.Container, body []map[string]interface{}) ([]byte, error) {
	return addMBeans(ctx, container, body, TestServerAddDataBatchEndpoint)
}
// AddMBeans will add new MBeans to the test-server.
func AddMBeans(ctx context.Context, container testcontainers.Container, body map[string]interface{}) ([]byte, error) {
	return addMBeans(ctx, container, body, TestServerAddDataEndpoint)
}

// AddMBeans will add new MBeans to the test-server.
func AddMCompositeDataBeans(ctx context.Context, container testcontainers.Container, body map[string]interface{}) ([]byte, error) {
	return addMBeans(ctx, container, body, TestServerAddCompositeDataEndpoint)
}

// addMBeans will add new MBeans to the test-server.
func addMBeans(ctx context.Context, container testcontainers.Container, body interface{}, endpointPath string) ([]byte, error) {
	url, err := GetContainerServiceURL(ctx, container, TestServerPort, endpointPath)
	if err != nil {
		return nil, err
	}
	json, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return DoHttpRequest(http.MethodPost, url, json)
}

// TestLogConsumer is used to print container logs to stdout.
type TestLogConsumer struct {
}

func (g *TestLogConsumer) Accept(l testcontainers.Log) {
	fmt.Fprintf(os.Stdout, "[CONTAINER LOG] %s %s\n", time.Now().Format("2006/01/02 15:04:05"), l.Content)
}

// RunJbossStandaloneJMXContainer will start a container running a jboss instace with JMX.
func RunJbossStandaloneJMXContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image: "test_jboss",

		ExposedPorts: []string{
			fmt.Sprintf("%[1]s:%[1]s", JbossJMXPort),
		},

		WaitingFor: wait.ForListeningPort(JbossJMXPort),
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

// CopyFileFromContainer will copy a file from a given docker container.
func CopyFileFromContainer(ctx context.Context, containerID, srcPath, dstPath string) error {
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

// DoHttpRequest will perform the http request.
func DoHttpRequest(method, url string, body []byte) ([]byte, error) {
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

// ReadFirstLine will return just the first line of the file.
func ReadFirstLine(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	output := scanner.Text()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

// GetContainerMappedPort returns the hostname and the port for a given container.
func GetContainerMappedPort(ctx context.Context, container testcontainers.Container, targetPort nat.Port) (host string, port nat.Port, err error) {
	host, err = container.Host(ctx)
	if err != nil {
		return
	}

	port, err = container.MappedPort(ctx, targetPort)
	if err != nil {
		return
	}
	return
}

// ReadPidFile will read a given file and extract the pid form it.
func ReadPidFile(fileName string) (int32, error) {
	line, err := ReadFirstLine(fileName)
	if err != nil {
		return -1, err
	}
	pid, err := strconv.Atoi(line)
	if err != nil {
		return -1, err
	}
	return int32(pid), nil
}