package gojmx

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/apache/thrift/lib/go/thrift"

	"github.com/newrelic/nrjmx/gojmx/nrprotocol"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStderrBuffer(t *testing.T) {
	buff := NewStderrBuffer(4)
	n, err := buff.WriteString("12345")
	assert.NoError(t, err)
	assert.Equal(t, 4, n)
	n, err = buff.WriteString("67")
	assert.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.Equal(t, "4567", buff.String())
	buff2 := NewStderrBuffer(5)
	n1, err := buff2.WriteString("12")
	assert.NoError(t, err)
	assert.Equal(t, 2, n1)
	n1, err = buff2.WriteString("3456")
	assert.NoError(t, err)
	assert.Equal(t, "23456", buff2.String())
}

func TestJMXServiceSubprocessStops(t *testing.T) {
	if os.Getenv("SHOULD_RUN_EXIT") == "1" {
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
		client, err := NewJMXClient(ctx).Open()
		assert.NoError(t, err)

		config := &nrprotocol.JMXConfig{
			Hostname: jmxHost,
			Port:     int32(jmxPort.Int()),
			UriPath:  thrift.StringPtr("jmxrmi"),
		}

		err = client.Connect(config, defaultTimeoutMs)
		assert.NoError(t, err)
		f, err := os.OpenFile(os.Getenv("TMP_FILE"), os.O_WRONLY|os.O_TRUNC, 0644)
		defer f.Close()
		fmt.Fprintf(f, "%d\n", client.jmxProcess.cmd.Process.Pid)
		<-time.After(60 * time.Second)
	}

	tmpfile, err := ioutil.TempFile("", "nrjmxtest")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	cmd := exec.Command(os.Args[0], "-test.run=TestJMXServiceSubprocessStops")
	cmd.Env = append(os.Environ(), "SHOULD_RUN_EXIT=1", "TMP_FILE="+tmpfile.Name())
	assert.NoError(t, err)

	err = cmd.Start()
	assert.NoError(t, err)

	var pid int32
	assert.Eventually(t, func() bool {
		line, err := ReadFirstLine(tmpfile.Name())
		if err != nil {
			return false
		}
		npid, err := strconv.Atoi(line)
		if err != nil {
			return false
		}
		pid = int32(npid)
		return true
	}, 5*time.Second, 50*time.Millisecond)

	p, err := process.NewProcess(pid)
	assert.NoError(t, err)
	ch, err := p.Children()
	assert.NoError(t, err)

	err = cmd.Process.Kill()
	assert.NoError(t, err)
	// check nrjmx pid does not exist anymore
	assert.Eventually(t, func() bool {
		up, err := ch[0].IsRunning()
		if err != nil {
			return false
		}
		// assert is not running
		return !up
	}, 5*time.Second, 50*time.Millisecond)
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
