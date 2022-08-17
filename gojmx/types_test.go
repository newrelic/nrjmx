/*
 * Copyright 2021 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package gojmx

import (
	"github.com/newrelic/nrjmx/gojmx/internal/nrprotocol"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_JMXAttribute_GetValue(t *testing.T) {
	testCases := []struct {
		name     string
		jmxAttr  *AttributeResponse
		expected interface{}
	}{
		{
			name: "Double Value",
			jmxAttr: &AttributeResponse{
				Name:         "test:type=Cat,name=tomas,attr=FloatValue",
				ResponseType: nrprotocol.ResponseType_DOUBLE,
				DoubleValue:  2.222222,
			},
			expected: 2.222222,
		},
		{
			name: "Number Value",
			jmxAttr: &AttributeResponse{
				Name:         "test:type=Cat,name=tomas,attr=NumberValue",
				ResponseType: nrprotocol.ResponseType_INT,
				IntValue:     3,
			},
			expected: int64(3),
		},
		{
			name: "Bool Value",
			jmxAttr: &AttributeResponse{
				Name:         "test:type=Cat,name=tomas,attr=BoolValue",
				ResponseType: nrprotocol.ResponseType_BOOL,
				BoolValue:    true,
			},
			expected: true,
		},
		{
			name: "Double Value",
			jmxAttr: &AttributeResponse{
				Name:         "test:type=Cat,name=tomas,attr=DoubleValue",
				ResponseType: nrprotocol.ResponseType_DOUBLE,
				DoubleValue:  1.2,
			},
			expected: 1.2,
		},
		{
			name: "String Value",
			jmxAttr: &AttributeResponse{
				Name:         "test:type=Cat,name=tomas,attr=Name",
				ResponseType: nrprotocol.ResponseType_STRING,
				StringValue:  "tomas",
			},
			expected: "tomas",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, testCase.jmxAttr.GetValue())
		})
	}
}

func Test_JMXAttribute_GetValueAsFloat(t *testing.T) {
	testCases := []struct {
		name          string
		jmxAttr       *AttributeResponse
		errorExpected bool
		expected      float64
	}{
		{
			name: "Incorrect String Value",
			jmxAttr: &AttributeResponse{
				Name:         "test:type=Cat,name=tomas,attr=Name",
				ResponseType: nrprotocol.ResponseType_STRING,
				StringValue:  "aaa",
			},
			errorExpected: true,
			expected:      0,
		},
		{
			name: "Double Value",
			jmxAttr: &AttributeResponse{
				Name:         "test:type=Cat,name=tomas,attr=FloatValue",
				ResponseType: nrprotocol.ResponseType_DOUBLE,
				DoubleValue:  2.222222,
			},
			expected: 2.222222,
		},
		{
			name: "Number Value",
			jmxAttr: &AttributeResponse{
				Name:         "test:type=Cat,name=tomas,attr=NumberValue",
				ResponseType: nrprotocol.ResponseType_INT,
				IntValue:     3,
			},
			expected: float64(3),
		},
		{
			name: "Bool Value",
			jmxAttr: &AttributeResponse{
				Name:         "test:type=Cat,name=tomas,attr=BoolValue",
				ResponseType: nrprotocol.ResponseType_BOOL,
				BoolValue:    true,
			},
			expected: 1,
		},
		{
			name: "Double Value",
			jmxAttr: &AttributeResponse{
				Name:         "test:type=Cat,name=tomas,attr=DoubleValue",
				ResponseType: nrprotocol.ResponseType_DOUBLE,
				DoubleValue:  1.2,
			},
			expected: 1.2,
		},
		{
			name: "String Value",
			jmxAttr: &AttributeResponse{
				Name:         "test:type=Cat,name=tomas,attr=Name",
				ResponseType: nrprotocol.ResponseType_STRING,
				StringValue:  "1.2",
			},
			expected: 1.2,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			floatVal, err := testCase.jmxAttr.GetValueAsFloat()
			if testCase.errorExpected {
				assert.Error(t, err)
			}
			assert.Equal(t, testCase.expected, floatVal)
		})
	}
}
