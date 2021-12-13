package gojmx

import (
	"context"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/newrelic/nrjmx/gojmx/internal/nrjmx"
	"github.com/pkg/errors"
	"time"

	"github.com/newrelic/nrjmx/gojmx/internal/nrprotocol"
)

const (
	// pingTimeout specifies how long we wait for a ping response.
	pingTimeout = 10 * time.Second

	// nrJMXExitTimeout specifies how long we wait for nrjmx process to exit.
	nrJMXExitTimeout = 5 * time.Second
)

// errPingTimeout returned if pingTimeout exceeded.
var errPingTimeout = nrjmx.NewJMXConnectionError("could not establish communication with nrjmx process: ping timeout")

// Client to connect with a JMX endpoint.
type Client struct {
	// jmxService is the thrift implementation to communicate with nrjmx subprocess.
	jmxService   nrprotocol.JMXService
	nrJMXProcess *nrjmx.Process
	ctx          context.Context
}

// NewClient returns a JMX client.
func NewClient(ctx context.Context) *Client {
	return &Client{
		ctx: ctx,
	}
}

// Open will create the connection the the JMX endpoint.
func (c *Client) Open(config *JMXConfig) (client *Client, err error) {
	c.nrJMXProcess, err = nrjmx.NewProcess(c.ctx).Start()
	if err != nil {
		return c, err
	}

	defer func() {
		if err != nil {
			_ = c.nrJMXProcess.Terminate()
		}
	}()

	c.jmxService, err = c.configureJMXServiceClient()
	if err != nil {
		return c, err
	}

	err = c.ping(pingTimeout)
	if err != nil {
		return c, err
	}

	return c, c.connect(config)
}

// GetMBeanNames returns all the mbeans that match the glob pattern DOMAIN:BEAN.
// e.g *:* or jboss.as:subsystem=remoting,configuration=endpoint
func (c *Client) GetMBeanNames(mBeanGlobPattern string) ([]string, error) {
	if err := c.nrJMXProcess.Error(); err != nil {
		return nil, err
	}
	result, err := c.jmxService.GetMBeanNames(c.ctx, mBeanGlobPattern)

	return result, c.handleTransportError(err)
}

// GetMBeanAttrNames returns all the available JMX attribute names for a given mBeanName.
func (c *Client) GetMBeanAttrNames(mBeanName string) ([]string, error) {
	if err := c.nrJMXProcess.Error(); err != nil {
		return nil, err
	}
	result, err := c.jmxService.GetMBeanAttrNames(c.ctx, mBeanName)
	return result, c.handleTransportError(err)
}

// GetMBeanAttrs returns the JMX attribute value.
func (c *Client) GetMBeanAttrs(mBeanName, mBeanAttrName string) ([]*JMXAttribute, error) {
	if err := c.nrJMXProcess.Error(); err != nil {
		return nil, err
	}

	result, err := c.jmxService.GetMBeanAttrs(c.ctx, mBeanName, mBeanAttrName)
	return toJMXAttributeList(result), err
}

// Close will stop the connection with the JMX endpoint.
func (c *Client) Close() error {
	if err := c.nrJMXProcess.Error(); err != nil {
		return err
	}
	err := c.jmxService.Disconnect(c.ctx)
	if waitErr := c.nrJMXProcess.WaitExit(nrJMXExitTimeout); waitErr != nil {
		err = errors.Wrap(err, waitErr.Error())
	}
	return err
}

// GetClientVersion returns nrjmx version.
func (c *Client) GetClientVersion() (version string, err error) {
	if err = c.nrJMXProcess.Error(); err != nil {
		return "<nil>", err
	}
	version, err = c.jmxService.GetClientVersion(c.ctx)

	return version, c.handleTransportError(err)
}

// connect will pass the JMXConfig to nrjmx subprocess and establish the
// connection with the JMX endpoint.
func (c *Client) connect(config *JMXConfig) (err error) {
	if err = c.nrJMXProcess.Error(); err != nil {
		return err
	}
	err = c.jmxService.Connect(c.ctx, config.toProtocol())

	return c.handleTransportError(err)
}

// ping will test the communication with nrjmx subprocess.
func (c *Client) ping(timeout time.Duration) error {
	ctx, cancel := context.WithCancel(c.ctx)
	defer cancel()
	done := make(chan struct{}, 1)
	go func() {
		for ctx.Err() == nil {
			_, err := c.jmxService.GetClientVersion(ctx)
			if err != nil {
				continue
			}
			done <- struct{}{}
			break
		}
	}()
	select {
	case <-time.After(timeout):
		return errPingTimeout
	case err := <-c.nrJMXProcess.ErrorC():
		return err
	case <-done:
		return nil
	}
}

// configureJMXServiceClient will configure the thrift service to communicate via stdin/stdout.
func (c *Client) configureJMXServiceClient() (*nrprotocol.JMXServiceClient, error) {
	var transport thrift.TTransport
	transport = thrift.NewStreamTransport(c.nrJMXProcess.Stdout, c.nrJMXProcess.Stdin)

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTCompactProtocolFactory()

	var transportFactory thrift.TTransportFactory
	transportFactory = thrift.NewTBufferedTransportFactory(8192)
	transportFactory = thrift.NewTFramedTransportFactory(transportFactory)

	transport, err := transportFactory.GetTransport(transport)
	if err != nil {
		return nil, err
	}

	inputProtocol := protocolFactory.GetProtocol(transport)
	outputProtocol := protocolFactory.GetProtocol(transport)
	jmxServiceClient := nrprotocol.NewJMXServiceClient(
		thrift.NewTStandardClient(inputProtocol, outputProtocol),
	)
	return jmxServiceClient, err
}

// handleTransportError will check if the error is TTransportException
// and if required will terminate nrjmx subprocess.
func (c *Client) handleTransportError(err error) error {
	if _, ok := err.(thrift.TTransportException); ok {
		// TTransportException means that interprocess communication
		// failed and it cannot be restored. We make sure nrJMX subprocess stops.
		return c.nrJMXProcess.WaitExit(nrJMXExitTimeout)
	}
	return err
}
