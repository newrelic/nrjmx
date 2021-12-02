package gojmx

import (
	"context"

	"time"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/newrelic/nrjmx/gojmx/nrprotocol"
)

func NewHttpJMXServiceClient(ctx context.Context, startSubprocess bool) (client *JMXClient, err error) {
	var jmxProcess = &jmxProcess{}
	if startSubprocess {
		jmxProcess, err = startJMXProcess(ctx)
		if err != nil {
			return
		}
	}

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTCompactProtocolFactory()

	var transportFactory thrift.TTransportFactory

	transportFactory = thrift.NewTBufferedTransportFactory(8192)

	transportFactory = thrift.NewTFramedTransportFactory(transportFactory)

	var transport thrift.TTransport

	transport, err = thrift.NewTSocket("localhost:9090")
	if err != nil {
		return nil, err
	}

	err = transport.Open()
	if err !=nil {
		return nil, err
	}

	transport, err = transportFactory.GetTransport(transport)

	iprot := protocolFactory.GetProtocol(transport)
	oprot := protocolFactory.GetProtocol(transport)
	client = &JMXClient{
		JMXService: nrprotocol.NewJMXServiceClient(thrift.NewTStandardClient(iprot, oprot)),
		jmxProcess: *jmxProcess,
		ctx:        ctx,
	}
	return
}

func NewJMXServiceClient(ctx context.Context) (client *JMXClient, err error) {
	jmxProcess, err := startJMXProcess(ctx)
	if err != nil {
		return
	}

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTCompactProtocolFactory()

	var transportFactory thrift.TTransportFactory

	transportFactory = thrift.NewTBufferedTransportFactory(8192)

	transportFactory = thrift.NewTFramedTransportFactory(transportFactory)

	var transport thrift.TTransport

	transport = thrift.NewStreamTransport(jmxProcess.Stdout, jmxProcess.Stdin)
	transport, err = transportFactory.GetTransport(transport)
	if err != nil {
		return nil, err
	}

	iprot := protocolFactory.GetProtocol(transport)
	oprot := protocolFactory.GetProtocol(transport)
	client = &JMXClient{
		JMXService: nrprotocol.NewJMXServiceClient(thrift.NewTStandardClient(iprot, oprot)),
		jmxProcess: *jmxProcess,
		ctx:        ctx,
	}
	return
}

type JMXClient struct {
	nrprotocol.JMXService
	jmxProcess jmxProcess
	ctx        context.Context
}

func (j *JMXClient) Query(timeout int64, mbean string) ([]*nrprotocol.JMXAttribute, error) {
	return j.QueryMbean(j.ctx, mbean, timeout)
}

func (j *JMXClient) Close(timeout time.Duration) error {
	j.Disconnect(j.ctx)
	return j.jmxProcess.stop(timeout)
}
