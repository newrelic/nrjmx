package gojmx

import (
	"github.com/newrelic/nrjmx/gojmx/internal/nrprotocol"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_JMXAttribute_GetValue(t *testing.T) {

	testCases := []struct {
		name     string
		jmxAttr  *JMXAttribute
		expected interface{}
	}{
		{
			name: "Double Value",
			jmxAttr: &JMXAttribute{
				Attribute:   "test:type=Cat,name=tomas,attr=FloatValue",
				ValueType:   nrprotocol.ValueType_DOUBLE,
				DoubleValue: 2.222222,
			},
			expected: 2.222222,
		},
		{
			name: "Number Value",
			jmxAttr: &JMXAttribute{
				Attribute: "test:type=Cat,name=tomas,attr=NumberValue",
				ValueType: nrprotocol.ValueType_INT,
				IntValue:  3,
			},
			expected: int64(3),
		},
		{
			name: "Bool Value",
			jmxAttr: &JMXAttribute{
				Attribute: "test:type=Cat,name=tomas,attr=BoolValue",
				ValueType: nrprotocol.ValueType_BOOL,
				BoolValue: true,
			},
			expected: true,
		},
		{
			name: "Double Value",
			jmxAttr: &JMXAttribute{
				Attribute:   "test:type=Cat,name=tomas,attr=DoubleValue",
				ValueType:   nrprotocol.ValueType_DOUBLE,
				DoubleValue: 1.2,
			},
			expected: 1.2,
		},
		{
			name: "String Value",
			jmxAttr: &JMXAttribute{
				Attribute:   "test:type=Cat,name=tomas,attr=Name",
				ValueType:   nrprotocol.ValueType_STRING,
				StringValue: "tomas",
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
