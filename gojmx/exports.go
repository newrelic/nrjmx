/*
 * Copyright 2021 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package gojmx

import (
	"fmt"
	"github.com/newrelic/nrjmx/gojmx/internal/nrprotocol"
	"unsafe"
)

/*
 * We want to keep the generated thrift structures internal to avoid having a heavy API.
 * Here we place things that we want to export.
 */

// JMXConfig exports internal nrprotocol.JMXConfig.
type JMXConfig nrprotocol.JMXConfig

// JMXAttribute exports internal nrprotocol.JMXAttribute.
type JMXAttribute nrprotocol.JMXAttribute

func (j *JMXAttribute) String() string {
	if j == nil {
		return "<nil>"
	}
	return fmt.Sprintf("JMXAttribute(%+v)", *j)
}

func toJMXAttributeList(in []*nrprotocol.JMXAttribute) []*JMXAttribute {
	return *(*[]*JMXAttribute)(unsafe.Pointer(&in))
}

// GetValue extracts the value from nrprotocol.JMXAttribute.
func (j *JMXAttribute) GetValue() interface{} {
	switch (*j).ValueType {
	case ValueTypeBool:
		return j.BoolValue
	case ValueTypeString:
		return j.StringValue
	case ValueTypeDouble:
		return j.DoubleValue
	case ValueTypeInt:
		return j.IntValue
	default:
		panic(fmt.Sprintf("unkown value type: %v", j.ValueType))
	}
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

// ValueType exports internal nrprotocol.ValueType.
type ValueType nrprotocol.ValueType

var (
	// ValueTypeBool JMXAttribute of bool value
	ValueTypeBool = nrprotocol.ValueType_BOOL
	// ValueTypeString JMXAttribute of string value
	ValueTypeString = nrprotocol.ValueType_STRING
	// ValueTypeDouble JMXAttribute of double value
	ValueTypeDouble = nrprotocol.ValueType_DOUBLE
	// ValueTypeInt JMXAttribute of int value
	ValueTypeInt = nrprotocol.ValueType_INT
)
