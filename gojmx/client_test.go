package gojmx

import (
	"context"
	"testing"
	"time"

	"github.com/cciutea/gojmx/generated/jmx"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfullCollection(t *testing.T) {
	var defaultCtx = context.Background()

	client, err := NewJMXServiceClient(defaultCtx)
	assert.NoError(t, err)

	config := &jmx.JMXConfig{
		ConnURL:            "service:jmx:remote+https://localhost:9993",
		Username:           "admin",
		Password:           "Admin.123",
		KeyStore:           "./tests/jboss/key/jboss.keystore",
		KeyStorePassword:   "password",
		TrustStore:         "./tests/jboss/key/jboss.truststore",
		TrustStorePassword: "password",
	}

	ok, err := client.Connect(defaultCtx, config)
	assert.NoError(t, err)
	assert.True(t, ok)

	actual, err := client.QueryMbean(defaultCtx, "jboss.as.expr:subsystem=remoting,configuration=endpoint")
	assert.NoError(t, err)

	expected := []struct {
		attribute string
		valueType jmx.ValueType
		value     interface{}
	}{
		{
			attribute: "jboss.as.expr:subsystem=remoting,configuration=endpoint,attr=authenticationRetries",
			valueType: jmx.ValueType_STRING,
			value:     "3",
		},
		{
			attribute: "jboss.as.expr:subsystem=remoting,configuration=endpoint,attr=heartbeatInterval",
			valueType: jmx.ValueType_STRING,
			value:     "2147483647",
		},
		{
			attribute: "jboss.as.expr:subsystem=remoting,configuration=endpoint,attr=maxInboundChannels",
			valueType: jmx.ValueType_STRING,
			value:     "40",
		},
		{
			attribute: "jboss.as.expr:subsystem=remoting,configuration=endpoint,attr=maxInboundMessageSize",
			valueType: jmx.ValueType_STRING,
			value:     "9223372036854775807",
		},
		{
			attribute: "jboss.as.expr:subsystem=remoting,configuration=endpoint,attr=maxInboundMessages",
			valueType: jmx.ValueType_STRING,
			value:     "80",
		},
		{
			attribute: "jboss.as.expr:subsystem=remoting,configuration=endpoint,attr=maxOutboundChannels",
			valueType: jmx.ValueType_STRING,
			value:     "40",
		},
		{
			attribute: "jboss.as.expr:subsystem=remoting,configuration=endpoint,attr=maxOutboundMessageSize",
			valueType: jmx.ValueType_STRING,
			value:     "9223372036854775807",
		},

		{
			attribute: "jboss.as.expr:subsystem=remoting,configuration=endpoint,attr=maxOutboundMessages",
			valueType: jmx.ValueType_STRING,
			value:     "65535",
		},
		{
			attribute: "jboss.as.expr:subsystem=remoting,configuration=endpoint,attr=receiveBufferSize",
			valueType: jmx.ValueType_STRING,
			value:     "8192",
		},
		{
			attribute: "jboss.as.expr:subsystem=remoting,configuration=endpoint,attr=receiveWindowSize",
			valueType: jmx.ValueType_STRING,
			value:     "131072",
		},
		{
			attribute: "jboss.as.expr:subsystem=remoting,configuration=endpoint,attr=saslProtocol",
			valueType: jmx.ValueType_STRING,
			value:     "remote",
		},
		{
			attribute: "jboss.as.expr:subsystem=remoting,configuration=endpoint,attr=sendBufferSize",
			valueType: jmx.ValueType_STRING,
			value:     "8192",
		},
		{
			attribute: "jboss.as.expr:subsystem=remoting,configuration=endpoint,attr=transmitWindowSize",
			valueType: jmx.ValueType_STRING,
			value:     "131072",
		},
		{
			attribute: "jboss.as.expr:subsystem=remoting,configuration=endpoint,attr=worker",
			valueType: jmx.ValueType_STRING,
			value:     "default",
		},
	}

	assert.Equal(t, len(expected), len(actual))

	for i, expectedAttr := range expected {
		actualAttr := actual[i]
		assert.Equal(t, expectedAttr.attribute, actualAttr.Attribute)

		assert.Equal(t, expectedAttr.valueType, actualAttr.Value.ValueType)

		assert.Equal(t, expectedAttr.value, actualAttr.GetValue())
	}

	assert.NoError(t, client.Disconnect(defaultCtx))

	done := make(chan error, 1)
	go func() {
		done <- client.jmxProcess.cmd.Wait()
	}()

	select {
	case <-time.After(1 * time.Second):
		assert.Fail(t, "timeout waiting for jmx process")
	case err := <-done:
		assert.NoError(t, err)
	}
}

func TestSuccessWrongPassword(t *testing.T) {
	var defaultCtx = context.Background()

	client, err := NewJMXServiceClient(defaultCtx)
	assert.NoError(t, err)

	config := &jmx.JMXConfig{
		ConnURL:            "service:jmx:remote+https://localhost:9993",
		Username:           "admin",
		Password:           "Admin.1234",
		KeyStore:           "./tests/jboss/key/jboss.keystore",
		KeyStorePassword:   "password",
		TrustStore:         "./tests/jboss/key/jboss.truststore",
		TrustStorePassword: "password",
	}

	expectedErr := &jmx.JMXConnectionError{
		Code:    1,
		Message: "Can't connect to JMX server: 'service:jmx:remote+https://localhost:9993', error: 'Authentication failed: all available authentication mechanisms failed:\n   DIGEST-MD5: javax.security.sasl.SaslException: DIGEST-MD5: Server rejected authentication'",
	}

	ok, err := client.Connect(defaultCtx, config)
	assert.Equal(t, err, expectedErr)
	assert.False(t, ok)
}

func TestSuccessWrongCertPassword(t *testing.T) {
	var defaultCtx = context.Background()

	client, err := NewJMXServiceClient(defaultCtx)
	assert.NoError(t, err)

	config := &jmx.JMXConfig{
		ConnURL:            "service:jmx:remote+https://localhost:9993",
		Username:           "admin",
		Password:           "Admin.123",
		KeyStore:           "./tests/jboss/key/jboss.keystore",
		KeyStorePassword:   "password1",
		TrustStore:         "./tests/jboss/key/jboss.truststore",
		TrustStorePassword: "password",
	}

	expectedErr := &jmx.JMXConnectionError{
		Code:    1,
		Message: "Can't connect to JMX server: 'service:jmx:remote+https://localhost:9993', error: 'JBREM000212: Failed to configure SSL context'",
	}

	ok, err := client.Connect(defaultCtx, config)
	assert.Equal(t, err, expectedErr)
	assert.False(t, ok)
}

func TestSuccessWrongURL(t *testing.T) {
	var defaultCtx = context.Background()

	client, err := NewJMXServiceClient(defaultCtx)
	assert.NoError(t, err)

	config := &jmx.JMXConfig{
		ConnURL:            "service:jmx:remote+https://localhost:9994",
		Username:           "admin",
		Password:           "Admin.123",
		KeyStore:           "./tests/jboss/key/jboss.keystore",
		KeyStorePassword:   "password",
		TrustStore:         "./tests/jboss/key/jboss.truststore",
		TrustStorePassword: "password",
	}

	expectedErr := &jmx.JMXConnectionError{
		Code:    1,
		Message: "Can't connect to JMX server: 'service:jmx:remote+https://localhost:9994', error: 'Connection refused'",
	}

	ok, err := client.Connect(defaultCtx, config)
	assert.Equal(t, err, expectedErr)
	assert.False(t, ok)
}
