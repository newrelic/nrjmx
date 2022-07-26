/*
 * Copyright 2021 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package gojmx

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/newrelic/nrjmx/gojmx/internal/nrprotocol"
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

// AttributeResponse keeps the JMX MBean query response.
type AttributeResponse nrprotocol.AttributeResponse

func (j *AttributeResponse) String() string {
	if j == nil {
		return "<nil>"
	}
	return fmt.Sprintf("AttributeResponse(%+v)", *j)
}

func toAttributeResponseList(in []*nrprotocol.AttributeResponse) []*AttributeResponse {
	return *(*[]*AttributeResponse)(unsafe.Pointer(&in))
}

// GetValue extracts the value from AttributeResponse based on type.
func (j *AttributeResponse) GetValue() interface{} {
	switch (*j).ResponseType {
	case ResponseTypeBool:
		return j.BoolValue
	case ResponseTypeString:
		return j.StringValue
	case ResponseTypeDouble:
		return j.DoubleValue
	case ResponseTypeInt:
		return j.IntValue
	case ResponseTypeErr:
		return "<nil>"
	default:
		panic(fmt.Sprintf("unkown value type: %v", j.ResponseType))
	}
}

// JMXError is reported when a JMX query fails.
type JMXError nrprotocol.JMXError

func (j *JMXError) String() string {
	if j == nil {
		return "<nil>"
	}
	return fmt.Sprintf("jmx error: %s, cause: %s, stacktrace: %s",
		removeNewLines(j.Message),
		removeNewLines(j.CauseMessage),
		removeNewLines(j.Stacktrace))
}

func removeNewLines(text string) string {
	text = strings.Replace(text, "\n", "\\n", -1)
	text = strings.Replace(text, "\r", " ", -1)
	return text
}

func (j *JMXError) Error() string {
	return j.String()
}

// IsJMXError asserts if the error is JMXError.
func IsJMXError(err error) (*JMXError, bool) {
	if e, ok := err.(*JMXError); err != nil && ok {
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
	return fmt.Sprintf("connection error: %s", removeNewLines(e.Message))
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
	if e, ok := err.(*JMXConnectionError); err != nil && ok {
		return e, ok
	}
	return nil, false
}

// ResponseType specify the type of the value of the AttributeResponse.
type ResponseType nrprotocol.ResponseType

var (
	// ResponseTypeBool AttributeResponse of bool value
	ResponseTypeBool = nrprotocol.ResponseType_BOOL
	// ResponseTypeString AttributeResponse of string value
	ResponseTypeString = nrprotocol.ResponseType_STRING
	// ResponseTypeDouble AttributeResponse of double value
	ResponseTypeDouble = nrprotocol.ResponseType_DOUBLE
	// ResponseTypeInt AttributeResponse of int value
	ResponseTypeInt = nrprotocol.ResponseType_INT
	// ResponseTypeErr AttributeResponse with error
	ResponseTypeErr = nrprotocol.ResponseType_ERROR
)

// InternalStat gathers stats about queries performed by nrjmx.
type InternalStat nrprotocol.InternalStat

func (is *InternalStat) String() string {
	return fmt.Sprintf("StatType: '%s', MBean: '%s', Attributes: '%v', TotalObjCount: %d, StartTimeMs: %d,  Duration: %.3fms, Successful: %t",
		is.StatType,
		is.MBean,
		is.Attrs,
		is.ResponseCount,
		is.StartTimestamp,
		is.Milliseconds,
		is.Successful,
	)
}

func toInternalStatList(in []*nrprotocol.InternalStat) []*InternalStat {
	return *(*[]*InternalStat)(unsafe.Pointer(&in))
}

// JMXClientError is returned when there is an nrjmx process error.
// Those errors require opening a new client.
type JMXClientError struct {
	Message string
}

func (e *JMXClientError) String() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("nrjmx client error: %s", removeNewLines(e.Message))
}

func (e *JMXClientError) Error() string {
	return e.String()
}

func newJMXClientError(message string, args ...interface{}) *JMXClientError {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return &JMXClientError{
		Message: message,
	}
}

// IsJMXClientError checks if the err is JMXJMXClientError.
func IsJMXClientError(err error) (*JMXClientError, bool) {
	if e, ok := err.(*JMXClientError); err != nil && ok {
		return e, ok
	}
	return nil, false
}
