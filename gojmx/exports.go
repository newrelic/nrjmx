package gojmx

import (
	"fmt"
	"github.com/newrelic/nrjmx/gojmx/internal/nrprotocol"
)

/*
 * We want to keep the generated thrift structures internal to avoid having a heavy API.
 * Here we place things that we want to export.
 */

// JMXConfig exports internal nrprotocol.JMXConfig.
type JMXConfig nrprotocol.JMXConfig

// JMXAttribute exports internal nrprotocol.JMXAttribute.
type JMXAttribute nrprotocol.JMXAttribute

// GetValue extracts the value from nrprotocol.JMXAttribute.
func (j *JMXAttribute) GetValue() interface{} {
	switch j.ValueType {
	case nrprotocol.ValueType_BOOL:
		return j.BoolValue
	case nrprotocol.ValueType_STRING:
		return j.StringValue
	case nrprotocol.ValueType_DOUBLE:
		return j.DoubleValue
	case nrprotocol.ValueType_INT:
		return j.IntValue
	default:
		panic(fmt.Sprintf("unkown value type: %v", j.ValueType))
	}
}

// ConvertJMXAttributeArray converts a list of nrprotocol.JMXAttribute into a list of JMXAttribute.
func ConvertJMXAttributeArray(attrs []*nrprotocol.JMXAttribute) (result []*JMXAttribute) {
	result = make([]*JMXAttribute, len(attrs))
	for i, attr := range attrs {
		result[i] = (*JMXAttribute)(attr)
	}
	return
}

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
