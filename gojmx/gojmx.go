/*
 * Copyright 2021 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package gojmx

import (
	"context"
	"time"

	"github.com/apache/thrift/lib/go/thrift"

	"github.com/newrelic/nrjmx/gojmx/internal/nrprotocol"
)

const (
	// pingTimeout specifies how long we wait for a ping response.
	pingTimeout = 10 * time.Second

	// nrJMXExitTimeout specifies how long we wait for nrjmx process to exit.
	nrJMXExitTimeout = 5 * time.Second

	unknownNRJMXVersion = "<unknown>"
)

// Client to connect with a JMX endpoint.
type Client struct {
	// jmxService is the thrift implementation to communicate with nrjmx subprocess.
	jmxService   nrprotocol.JMXService
	nrJMXProcess *process
	ctx          context.Context
	version      string
}

// NewClient returns a JMX client.
func NewClient(ctx context.Context) *Client {
	return &Client{
		ctx:          ctx,
		version:      unknownNRJMXVersion,
		nrJMXProcess: newProcess(ctx),
	}
}

// Open will create the connection the the JMX endpoint.
func (c *Client) Open(config *JMXConfig) (client *Client, err error) {
	c.nrJMXProcess, err = newProcess(c.ctx).start()
	if err != nil {
		return c, err
	}

	c.jmxService, err = c.configureJMXServiceClient()
	if err != nil {
		c.nrJMXProcess.waitExit(nrJMXExitTimeout)
		return c, err
	}

	c.version, err = c.ping(pingTimeout)
	if err != nil {
		c.nrJMXProcess.waitExit(nrJMXExitTimeout)
		return c, err
	}

	return c, c.connect(config)
}

// IsClientRunning returns if the nrjmx client is running.
func (c *Client) IsRunning() bool {
	if c.nrJMXProcess == nil {
		return false
	}

	return c.nrJMXProcess.state.IsRunning()
}

// checkNRJMXProccessError will check if the nrjmx subprocess returned any error.
func (c *Client) checkNRJMXProccessError() error {
	if c.nrJMXProcess == nil {
		return errProcessNotRunning
	}
	return c.nrJMXProcess.error()
}

// QueryMBeanNames returns all the mbeans that match the glob pattern DOMAIN:BEAN.
// e.g *:* or jboss.as:subsystem=remoting,configuration=endpoint
func (c *Client) QueryMBeanNames(mBeanGlobPattern string) ([]string, error) {
	if err := c.checkNRJMXProccessError(); err != nil {
		return nil, err
	}
	result, err := c.jmxService.QueryMBeanNames(c.ctx, mBeanGlobPattern)

	return result, c.handleError(err)
}

// GetMBeanAttributeNames returns all the available JMX attribute names for a given mBeanName.
func (c *Client) GetMBeanAttributeNames(mBeanName string) ([]string, error) {
	if err := c.checkNRJMXProccessError(); err != nil {
		return nil, err
	}
	result, err := c.jmxService.GetMBeanAttributeNames(c.ctx, mBeanName)
	return result, c.handleError(err)
}

// GetMBeanAttributes returns the JMX attribute values.
func (c *Client) GetMBeanAttributes(mBeanName string, mBeanAttrName ...string) ([]*AttributeResponse, error) {
	if err := c.checkNRJMXProccessError(); err != nil {
		return nil, err
	}

	result, err := c.jmxService.GetMBeanAttributes(c.ctx, mBeanName, mBeanAttrName)
	return toAttributeResponseList(result), c.handleError(err)
}

// Close will stop the connection with the JMX endpoint.
func (c *Client) Close() error {
	if err := c.checkNRJMXProccessError(); err != nil {
		return err
	}
	c.jmxService.Disconnect(c.ctx)
	if waitErr := c.nrJMXProcess.waitExit(nrJMXExitTimeout); waitErr != nil {
		return waitErr
	}
	return nil
}

// GetClientVersion returns nrjmx version.
func (c *Client) GetClientVersion() string {
	return c.version
}

// QueryMBeanAttributes performs all calls necessary for retrieving all MBeanAttrs values for the mBeanNamePattern:
// 1. QueryMBeanNames
// 2. GetMBeanAttributeNames
// 3. GetMBeanAttributes
// If an error occur it checks if it's a collection error (it can recover) or a connection error (that blocks all the collection).
func (c *Client) QueryMBeanAttributes(mBeanNamePattern string, mBeanAttrName ...string) ([]*AttributeResponse, error) {
	if err := c.checkNRJMXProccessError(); err != nil {
		return nil, err
	}

	result, err := c.jmxService.QueryMBeanAttributes(c.ctx, mBeanNamePattern, mBeanAttrName)
	return toAttributeResponseList(result), c.handleError(err)
}

// GetInternalStats returns the nrjmx internal query statistics for troubleshooting.
// Internal statistics must be enabled using JMXConfig.EnableInternalStats flag.
// Additionally you can set a maximum size for the collected stats using JMXConfig.MaxInternalStatsSize. (default: 100000)
// Each time you retrieve GetInternalStats, the internal stats will be cleaned.
func (c *Client) GetInternalStats() ([]*InternalStat, error) {
	if err := c.checkNRJMXProccessError(); err != nil {
		return nil, err
	}
	result, err := c.jmxService.GetInternalStats(c.ctx)

	return toInternalStatList(result), c.handleError(err)
}

// connect will pass the JMXConfig to nrjmx subprocess and establish the
// connection with the JMX endpoint.
func (c *Client) connect(config *JMXConfig) (err error) {
	if err = c.checkNRJMXProccessError(); err != nil {
		return err
	}
	err = c.jmxService.Connect(c.ctx, config.convertToProtocol())

	return c.handleError(err)
}

// ping will test the communication with nrjmx subprocess.
func (c *Client) ping(timeout time.Duration) (string, error) {
	ctx, cancel := context.WithCancel(c.ctx)
	defer cancel()

	done := make(chan string, 1)
	go func() {
		for ctx.Err() == nil {
			version, err := c.jmxService.GetClientVersion(ctx)
			if err != nil {
				continue
			}
			done <- version
			break
		}
	}()

	select {
	case <-time.After(timeout):
		return "<nil>", errPingTimeout
	case err, open := <-c.nrJMXProcess.state.ErrorC():
		if err == nil && !open {
			return "<nil>", errProcessNotRunning
		}
		return "<nil>", err
	case version := <-done:
		return version, nil
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
func (c *Client) handleError(err error) error {
	if _, ok := err.(thrift.TTransportException); ok {
		// TTransportException means that interprocess communication
		// failed, and it cannot be restored. We make sure nrJMX subprocess stops.
		return c.nrJMXProcess.waitExit(nrJMXExitTimeout)
	} else if jmxErr, ok := err.(*nrprotocol.JMXError); ok {
		return (*JMXError)(jmxErr)
	} else if jmxConnErr, ok := err.(*nrprotocol.JMXConnectionError); ok {
		return (*JMXConnectionError)(jmxConnErr)
	}
	return err
}
