package gojmx

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/apache/thrift/lib/go/thrift"

	"github.com/newrelic/nrjmx/gojmx/nrprotocol"
)

const pingTimeout = 1000 * time.Millisecond

var (
	ErrAlreadyStarted = errors.New("nrjmx subprocess already started")
	ErrNotRunning     = errors.New("nrjmx subprocess is not running")
)

type JMXClient struct {
	jmxService nrprotocol.JMXService
	jmxProcess *jmxProcess
	isRunning  bool
	ctx        context.Context
}

func NewJMXClient(ctx context.Context) *JMXClient {
	return &JMXClient{
		ctx: ctx,
	}
}

func (j *JMXClient) InitStandardIO() (*JMXClient, error) {
	if j.isRunning {
		return j, ErrAlreadyStarted
	}

	jmxProcess, err := startJMXProcess(j.ctx)
	if err != nil {
		jmxProcess.stop() // TODO: Handle err
		return j, err
	}

	transport := thrift.NewStreamTransport(jmxProcess.Stdout, jmxProcess.Stdin)
	jmxServiceClient, err := j.configureJMXServiceClient(transport)

	if err != nil {
		jmxProcess.stop() // TODO: Handle err
		return j, err
	}

	j.jmxProcess = jmxProcess
	j.jmxService = jmxServiceClient

	err = j.ping(pingTimeout)
	if err != nil {
		jmxProcess.stop() // TODO: Handle err
		return j, err
	}
	j.isRunning = true

	return j, nil
}

func (j *JMXClient) configureJMXServiceClient(transport thrift.TTransport) (*nrprotocol.JMXServiceClient, error) {
	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTCompactProtocolFactory()

	var transportFactory thrift.TTransportFactory
	transportFactory = thrift.NewTBufferedTransportFactory(8192)
	transportFactory = thrift.NewTFramedTransportFactory(transportFactory)

	transport, err := transportFactory.GetTransport(transport)
	if err != nil {
		return nil, err
	}

	iprot := protocolFactory.GetProtocol(transport)
	oprot := protocolFactory.GetProtocol(transport)
	jmxServiceClient := nrprotocol.NewJMXServiceClient(thrift.NewTStandardClient(iprot, oprot))
	return jmxServiceClient, err
}

//Connect(ctx context.Context, config *JMXConfig) (err error)
//Disconnect(ctx context.Context) (err error)
// Parameters:
//  - BeanName
//QueryMbean(ctx context.Context, beanName string) (r []*JMXAttribute, err error)
//GetLogs(ctx context.Context) (r []*LogMessage, err error)

func (j *JMXClient) ping(timeout time.Duration) error {
	ctx, cancel := context.WithCancel(j.ctx)
	defer cancel()
	done := make(chan struct{}, 1)
	go func() {
		for ctx.Err() == nil {
			err := j.jmxService.Ping(ctx)
			if err != nil {
				continue
			}
			done <- struct{}{}
			break
		}
	}()
	select {
	case <-time.After(timeout):
		return fmt.Errorf("ping timeout")
	case err := <-j.jmxProcess.errCh:
		return err
	case <-done:
		return nil
	}
}

func (j *JMXClient) checkState() error {
	if !j.isRunning {
		return ErrNotRunning
	}
	if j.jmxProcess != nil {
		err := j.jmxProcess.Error()
		if err != nil {
			return err
		}
	}
	return nil
}

func (j *JMXClient) Connect(config *nrprotocol.JMXConfig, timeout int64) error {
	if err := j.checkState(); err != nil {
		return err
	}
	return j.jmxService.Connect(j.ctx, config, timeout)
}

func (j *JMXClient) GetMBeanNames(mbean string, timeout int64) ([]string, error) {
	if err := j.checkState(); err != nil {
		return nil, err
	}
	result, err := j.jmxService.GetMBeanNames(j.ctx, mbean, timeout)
	if err == io.EOF {
		err2 := j.jmxProcess.WaitExitError(5 * time.Second)
		if err2 != nil {
			return nil, err2
		}
	}
	return result, err
}

func (j *JMXClient) GetMBeanAttrNames(mbean string, timeout int64) ([]string, error) {
	if err := j.checkState(); err != nil {
		return nil, err
	}
	return j.jmxService.GetMBeanAttrNames(j.ctx, mbean, timeout)
}

func (j *JMXClient) GetMBeanAttr(mBeanName, mBeanAttrName string, timeout int64) (*nrprotocol.JMXAttribute, error) {
	if err := j.checkState(); err != nil {
		return nil, err
	}
	return j.jmxService.GetMBeanAttr(j.ctx, mBeanName, mBeanAttrName, timeout)
}

func (j *JMXClient) Disconnect() error {
	if err := j.checkState(); err != nil {
		return err
	}
	defer func() {
		//j.jmxProcess.stop()
		//j = NewJMXClient(j.ctx)
	}()
	return j.jmxService.Disconnect(j.ctx)
}

func (j *JMXClient) InitTCP(startSubprocess bool) (*JMXClient, error) {
	if j.isRunning {
		return j, ErrAlreadyStarted
	}

	if startSubprocess {
		jmxProcess, err := startJMXProcess(j.ctx)
		if err != nil {
			jmxProcess.stop() // TODO: Handle err
			return j, err
		}
		j.jmxProcess = jmxProcess
	}

	transport, err := thrift.NewTSocket("localhost:9090")
	if err != nil {
		if startSubprocess {
			j.jmxProcess.stop()
		}
		return j, err
	}

	err = transport.Open()
	if err != nil {
		if startSubprocess {
			j.jmxProcess.stop()
		}
		return j, err
	}
	jmxServiceClient, err := j.configureJMXServiceClient(transport)
	if err != nil {
		if startSubprocess {
			j.jmxProcess.stop()
		}
		return j, err
	}

	j.jmxService = jmxServiceClient

	if startSubprocess {
		err = j.ping(pingTimeout)
		if err != nil {
			j.jmxProcess.stop() // TODO: Handle err
			return j, err
		}
	}

	j.isRunning = true
	return j, nil
}

func (j *JMXClient) Close() error {
	j.jmxService.Disconnect(j.ctx)
	return j.jmxProcess.stop()
}
