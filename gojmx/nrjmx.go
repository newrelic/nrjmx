package gojmx

import (
	"context"

	"time"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/newrelic/nrjmx/gojmx/nrprotocol"
)

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

func (j *JMXClient) Query(timeout time.Duration, mbean string) ([]*nrprotocol.JMXAttribute, error) {
	queryCtx, cancel := context.WithTimeout(j.ctx, timeout)
	cancel()
	return j.QueryMbean(queryCtx, mbean)
}

func (j *JMXClient) Close(timeout time.Duration) error {
	j.Disconnect(j.ctx)
	return j.jmxProcess.stop(timeout)
}
