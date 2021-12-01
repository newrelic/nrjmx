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
		jmxService: nrprotocol.NewJMXServiceClient(thrift.NewTStandardClient(iprot, oprot)),
		jmxProcess: *jmxProcess,
		ctx:        ctx,
	}
	err = client.jmxService.Ping(ctx)
	return
}

type JMXClient struct {
	jmxService nrprotocol.JMXService
	jmxProcess jmxProcess
	ctx        context.Context
}


//Connect(ctx context.Context, config *JMXConfig) (err error)
//Disconnect(ctx context.Context) (err error)
// Parameters:
//  - BeanName
//QueryMbean(ctx context.Context, beanName string) (r []*JMXAttribute, err error)
//GetLogs(ctx context.Context) (r []*LogMessage, err error)

func (j *JMXClient) Connect(config *nrprotocol.JMXConfig) error {
	err := j.jmxProcess.Error()
	if err != nil {
		return err
	}
	return j.jmxService.Connect(j.ctx, config)
}

func (j *JMXClient) QueryMbean(beanName string) ([]*nrprotocol.JMXAttribute, error) {
	err := j.jmxProcess.Error()
	if err != nil {
		return nil, err
	}
	return j.jmxService.QueryMbean(j.ctx, beanName)
}

func (j *JMXClient) Close(timeout time.Duration) error {
	//j.Disconnect(j.ctx)
	return j.jmxProcess.stop(timeout)
}

func (j *JMXClient) Disconnect() error {
	err := j.jmxProcess.Error()
	if err != nil {
		return err
	}
	return nil
}
