package gojmx

import "github.com/newrelic/nrjmx/gojmx/internal/nrprotocol"

/*
 * We want to keep the generated thrift structures internal to avoid having a heavy API.
 * Here we place things that we want to export.
 */

// JMXConfig exports internal nrprotocol.JMXConfig.
type JMXConfig nrprotocol.JMXConfig

// JMXAttribute exports internaln nrprotocol.JMXAttribute.
type JMXAttribute nrprotocol.JMXAttribute

func (j *JMXConfig) toProtocol() *nrprotocol.JMXConfig {
	return (*nrprotocol.JMXConfig)(j)
}

// JMXError exports nrprotocol.JMXError.
type JMXError *nrprotocol.JMXError

// IsJMXError tries to convert the error to exported JMXError.
func IsJMXError(err error) (JMXError, bool) {
	e, ok := err.(*nrprotocol.JMXError)
	return e, ok
}

// JMXConnectionError exports nrprotocol.JMXConnectionError.
type JMXConnectionError *nrprotocol.JMXConnectionError

// IsJMXConnectionError tries to convert the error to exported JMXConnectionError.
func IsJMXConnectionError(err error) (JMXConnectionError, bool) {
	e, ok := err.(*nrprotocol.JMXConnectionError)
	return e, ok
}
