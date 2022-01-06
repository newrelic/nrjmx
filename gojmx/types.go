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

// JMXConfig keeps the JMX connection settings.
type JMXConfig nrprotocol.JMXConfig

func (j *JMXConfig) convertToProtocol() *nrprotocol.JMXConfig {
	return (*nrprotocol.JMXConfig)(j)
}

// JMXAttribute keeps the JMX MBean query response.
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

// GetValue extracts the value from JMXAttribute based on type.
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

// JMXError is reported when a JMX query fails.
type JMXError nrprotocol.JMXError

func (j *JMXError) String() string {
	if j == nil {
		return "<nil>"
	}
	return fmt.Sprintf("jmx error: %q, cause: %q, stacktrace: %q", j.Message, j.CauseMessage, j.Stacktrace)
}

func (j *JMXError) Error() string {
	return j.String()
}

// IsJMXError asserts if the error is JMXError.
func IsJMXError(err error) (*JMXError, bool) {
	if e, ok := err.(*JMXError); ok {
		return e, ok
	}
	return nil, false
}

// JMXConnectionError is returned when there is a JMX connection error or a nrjmx process error.
type JMXConnectionError nrprotocol.JMXConnectionError

func (e *JMXConnectionError) String() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("connection error: %q", e.Message)
}

func (e *JMXConnectionError) Error() string {
	return e.String()
}

func newJMXConnectionError(message string, args ...interface{}) *JMXConnectionError {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return &JMXConnectionError{
		Message: message,
	}
}

// IsJMXConnectionError tries to convert the error to exported JMXConnectionError.
func IsJMXConnectionError(err error) (*JMXConnectionError, bool) {
	if e, ok := err.(*JMXConnectionError); ok {
		return e, ok
	}
	return nil, false
}

// ValueType specify the type of the value of the JMXAttribute.
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

// ResponseStatus for QueryResponse
type ResponseStatus int

const (
	// QueryResponseStatusSuccess is returned when JMXAttribute was successfully retrieved.
	QueryResponseStatusSuccess ResponseStatus = iota
	// QueryResponseStatusError is returned when gojmx fails to retrieve the JMXAttribute.
	QueryResponseStatusError
)

// QueryAttrResponse wraps the JMXAttribute with the status and status message for the request.
type QueryAttrResponse struct {
	Attribute *JMXAttribute
	Status    ResponseStatus
	StatusMsg string
}

// QueryResponse wraps the JMXAttribute with the status and status message for the request.
type QueryResponse []*QueryAttrResponse

// GetValidAttributes returns all the valid JMXAttribute from the QueryResponse by checking the query status.
func (qr QueryResponse) GetValidAttributes() (result []*JMXAttribute) {
	for _, attr := range qr {
		if attr.Status != QueryResponseStatusSuccess {
			continue
		}
		result = append(result, attr.Attribute)
	}
	return
}
