package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/docker/client"
	"github.com/newrelic/infra-integrations-sdk/jmx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	jbossJMXPort = "9990"
)

func Test_Connector_Success(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JBoss Server with JMX exposed running inside a container
	container, err := runJbossStandaloneJMXContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)
	container.GetContainerID()

	// Install the connector
	dstFile := filepath.Join(prjDir, "/bin/jboss-client.jar")
	err = copyFileFromContainer(ctx, container.GetContainerID(), "/opt/jboss/wildfly/bin/client/jboss-client.jar", dstFile)
	assert.NoError(t, err)

	defer os.Remove(dstFile)

	// THEN JMX connection can be oppened
	jmxPort, err := container.MappedPort(ctx, jbossJMXPort)
	require.NoError(t, err)
	jmxHost, err := container.Host(ctx)
	require.NoError(t, err)

	options := []jmx.Option{
		jmx.WithRemoteProtocol(),
		jmx.WithRemoteStandAloneJBoss(),
	}

	err = jmx.Open(jmxHost, jmxPort.Port(), "admin1234", "Password1!", options...)
	defer jmx.Close()
	assert.NoError(t, err)

	// AND Query returns expected data
	result, err := jmx.Query("jboss.as:subsystem=remoting,configuration=endpoint", 10000)
	assert.NoError(t, err)

	val, found := result["jboss.as:subsystem=remoting,configuration=endpoint,attr=heartbeatInterval"]
	assert.True(t, found)
	assert.Equal(t, float64(60000), val)
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
