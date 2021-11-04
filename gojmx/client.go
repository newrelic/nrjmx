package gojmx

import (
	"context"
	"fmt"

	"time"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/cciutea/gojmx/generated/jmx"
)

func NewJMXServiceClient(ctx context.Context) (client *JMXClient, err error) {
	jmxProcess, err := startJMXProcess(ctx)
	if err != nil {
		return
	}

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTJSONProtocolFactory()

	var transportFactory thrift.TTransportFactory
	transportFactory = thrift.NewTTransportFactory()

	// transportFactory = thrift.NewTBufferedTransportFactory(8192)

	// transportFactory = thrift.NewTFramedTransportFactory(transportFactory)

	var transport thrift.TTransport
	transport = thrift.NewStreamTransport(jmxProcess.Stdout, jmxProcess.Stdin)
	transport, err = transportFactory.GetTransport(transport)
	if err != nil {
		return nil, err
	}

	iprot := protocolFactory.GetProtocol(transport)
	oprot := protocolFactory.GetProtocol(transport)
	client = &JMXClient{
		JMXService: jmx.NewJMXServiceClient(thrift.NewTStandardClient(iprot, oprot)),
		jmxProcess: *jmxProcess,
		ctx:        ctx,
	}
	return
}

type JMXClient struct {
	jmx.JMXService
	jmxProcess jmxProcess
	ctx        context.Context
}

func (j *JMXClient) Close(timeout time.Duration) error {
	j.Disconnect(j.ctx)
	return j.jmxProcess.stop(timeout)
}

type JMXAttributeValueConverter struct {
	*jmx.JMXAttributeValue
}

func (j *JMXAttributeValueConverter) GetValue() interface{} {
	switch j.ValueType {
	case jmx.ValueType_BOOL:
		return j.GetBoolValue()
	case jmx.ValueType_STRING:
		return j.GetStringValue()
	case jmx.ValueType_DOUBLE:
		return j.GetDoubleValue()
	case jmx.ValueType_INT:
		return j.GetIntValue()
	default:
		panic(fmt.Sprintf("unkown value type: %v", j.ValueType))
	}
}
