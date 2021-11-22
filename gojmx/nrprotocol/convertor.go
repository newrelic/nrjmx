package nrprotocol

import "fmt"

func (j *JMXAttribute) GetValue() interface{} {
	switch j.ValueType {
	case ValueType_BOOL:
		return j.GetBoolValue()
	case ValueType_STRING:
		return j.GetStringValue()
	case ValueType_DOUBLE:
		return j.GetDoubleValue()
	case ValueType_INT:
		return j.GetIntValue()
	default:
		panic(fmt.Sprintf("unkown value type: %v", j.ValueType))
	}
}
