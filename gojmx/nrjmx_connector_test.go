package gojmx

// const (
// 	jbossJMXPort = "9990"
// )

// func Test_Connector_Success(t *testing.T) {
// 	ctx := context.Background()

// 	// GIVEN a JBoss Server with JMX exposed running inside a container
// 	container, err := runJbossStandaloneJMXContainer(ctx)
// 	require.NoError(t, err)
// 	defer container.Terminate(ctx)

// 	// Install the connector
// 	dstFile := filepath.Join(prjDir, "/bin/jboss-client.jar")
// 	err = copyFileFromContainer(ctx, container.GetContainerID(), "/opt/jboss/wildfly/bin/client/jboss-client.jar", dstFile)
// 	assert.NoError(t, err)

// 	defer os.Remove(dstFile)

// 	// THEN JMX connection can be oppened
// 	jmxPort, err := container.MappedPort(ctx, jbossJMXPort)
// 	require.NoError(t, err)
// 	jmxHost, err := container.Host(ctx)
// 	require.NoError(t, err)

// 	client, err := NewJMXServiceClient(ctx)
// 	assert.NoError(t, err)

// 	config := &nrprotocol.JMXConfig{
// 		Hostname:              jmxHost,
// 		Port:                  int32(jmxPort.Int()),
// 		Username:              "admin1234",
// 		Password:              "Password1!",
// 		IsJBossStandaloneMode: true,
// 		IsRemote:              true,
// 	}

// 	_, err = client.Connect(ctx, config)
// 	defer client.Disconnect(ctx)
// 	assert.NoError(t, err)

// 	// AND Query returns expected data
// 	actual, err := client.QueryMbean(ctx, "jboss.as:subsystem=remoting,configuration=endpoint")
// 	assert.NoError(t, err)

// 	expected := []*nrprotocol.JMXAttribute{
// 		{
// 			Attribute: "jboss.as:subsystem=remoting,configuration=endpoint,attr=authenticationRetries",
// 			ValueType: nrprotocol.ValueType_INT,
// 			IntValue:  3,
// 		},
// 		{
// 			Attribute: "jboss.as:subsystem=remoting,configuration=endpoint,attr=heartbeatInterval",
// 			ValueType: nrprotocol.ValueType_INT,
// 			IntValue:  60000,
// 		},
// 		{
// 			Attribute: "jboss.as:subsystem=remoting,configuration=endpoint,attr=maxInboundChannels",
// 			ValueType: nrprotocol.ValueType_INT,
// 			IntValue:  40,
// 		},
// 		{
// 			Attribute: "jboss.as:subsystem=remoting,configuration=endpoint,attr=maxInboundMessageSize",
// 			ValueType: nrprotocol.ValueType_INT,
// 			IntValue:  9223372036854775807,
// 		},
// 		{
// 			Attribute: "jboss.as:subsystem=remoting,configuration=endpoint,attr=maxInboundMessages",
// 			ValueType: nrprotocol.ValueType_INT,
// 			IntValue:  80,
// 		},
// 		{
// 			Attribute: "jboss.as:subsystem=remoting,configuration=endpoint,attr=maxOutboundChannels",
// 			ValueType: nrprotocol.ValueType_INT,
// 			IntValue:  40,
// 		},
// 		{
// 			Attribute: "jboss.as:subsystem=remoting,configuration=endpoint,attr=maxOutboundMessageSize",
// 			ValueType: nrprotocol.ValueType_INT,
// 			IntValue:  9223372036854775807,
// 		},
// 		{
// 			Attribute: "jboss.as:subsystem=remoting,configuration=endpoint,attr=maxOutboundMessages",
// 			ValueType: nrprotocol.ValueType_INT,
// 			IntValue:  65535,
// 		},
// 		{
// 			Attribute: "jboss.as:subsystem=remoting,configuration=endpoint,attr=receiveBufferSize",
// 			ValueType: nrprotocol.ValueType_INT,
// 			IntValue:  8192,
// 		},
// 		{
// 			Attribute: "jboss.as:subsystem=remoting,configuration=endpoint,attr=receiveWindowSize",
// 			ValueType: nrprotocol.ValueType_INT,
// 			IntValue:  131072,
// 		},
// 		{
// 			Attribute:   "jboss.as:subsystem=remoting,configuration=endpoint,attr=saslProtocol",
// 			ValueType:   nrprotocol.ValueType_STRING,
// 			StringValue: "remote",
// 		},
// 		{
// 			Attribute: "jboss.as:subsystem=remoting,configuration=endpoint,attr=sendBufferSize",
// 			ValueType: nrprotocol.ValueType_INT,
// 			IntValue:  8192,
// 		},
// 		{
// 			Attribute: "jboss.as:subsystem=remoting,configuration=endpoint,attr=transmitWindowSize",
// 			ValueType: nrprotocol.ValueType_INT,
// 			IntValue:  131072,
// 		},
// 		{
// 			Attribute:   "jboss.as:subsystem=remoting,configuration=endpoint,attr=worker",
// 			ValueType:   nrprotocol.ValueType_STRING,
// 			StringValue: "default",
// 		},
// 	}

// 	assert.Equal(t, expected, actual)
// }

// // runJbossStandaloneJMXContainer will start a container running a jboss instace with JMX.
// func runJbossStandaloneJMXContainer(ctx context.Context) (testcontainers.Container, error) {
// 	req := testcontainers.ContainerRequest{
// 		Image: "test_jboss",

// 		ExposedPorts: []string{
// 			fmt.Sprintf("%[1]s:%[1]s", jbossJMXPort),
// 		},

// 		WaitingFor: wait.ForListeningPort(jbossJMXPort),
// 	}

// 	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
// 		ContainerRequest: req,
// 		Started:          true,
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	container.StartLogProducer(ctx)
// 	container.FollowOutput(&TestLogConsumer{})
// 	return container, err
// }

// func copyFileFromContainer(ctx context.Context, containerID, srcPath, dstPath string) error {
// 	cli, err := client.NewClientWithOpts(client.FromEnv)
// 	if err != nil {
// 		return err
// 	}
// 	reader, containerPathStat, err := cli.CopyFromContainer(ctx, containerID, srcPath)
// 	if err != nil {
// 		return err
// 	}
// 	defer reader.Close()

// 	if !containerPathStat.Mode.IsRegular() {
// 		return fmt.Errorf("src is not a regular file: %s", srcPath)
// 	}

// 	b, err := ioutil.ReadAll(reader)
// 	if err != nil {
// 		return err
// 	}
// 	return ioutil.WriteFile(dstPath, b, 0644)
// }
