/*
 * Copyright 2021 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package gojmx

import (
	"context"
	"fmt"
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

// GetMBeanAttrs returns the JMX attribute values.
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

// ResponseStatus for QueryResponse
type ResponseStatus int

const (
	// QueryResponseStatusSuccess is returned when JMXAttribute was successfully retrieved.
	QueryResponseStatusSuccess ResponseStatus = iota
	// QueryResponseStatusError is returned when gojmx fails to retrieve the JMXAttribute.
	QueryResponseStatusError
)

// QueryAttrResponse wraps the JMXAttribute with the status and status message for the request.
type QueryAttrResponse struct {
	Attribute *JMXAttribute
	Status    ResponseStatus
	StatusMsg string
}

type QueryResponse []*QueryAttrResponse

// GetValidAttributes returns all the valid JMXAttribute from the QueryResponse by checking the query status.
func (qr *QueryResponse) GetValidAttributes() (result []*JMXAttribute) {
	if qr == nil {
		return
	}
	for _, attr := range *qr {
		if attr.Status != QueryResponseStatusSuccess {
			continue
		}
		result = append(result, attr.Attribute)
	}
	return
}

// QueryMBean performs all calls necessary for retrieving all MBeanAttrs values for the mBeanNamePattern:
// 1. GetMBeanNames
// 2. GetMBeanAttrNames
// 3. GetMBeanAttrs
// If an error occur it checks if it's a collection error (it can recover) or a connection error (that blocks all the collection).
func (c *Client) QueryMBean(mBeanNamePattern string) (QueryResponse, error) {
	var result QueryResponse

	mBeanNames, err := c.GetMBeanNames(mBeanNamePattern)
	if err != nil {
		return nil, err
	}

	for _, mBeanName := range mBeanNames {
		mBeanAttrNames, err := c.GetMBeanAttrNames(mBeanName)
		if jmxErr, isJMXErr := IsJMXError(err); isJMXErr {
			result = append(result, &QueryAttrResponse{
				Status: QueryResponseStatusError,
				StatusMsg: fmt.Sprintf("error while querying mBean name: '%s', error message: %s, error cause: %s, stacktrace: %q",
					mBeanName,
					jmxErr.Message,
					jmxErr.CauseMessage,
					jmxErr.Stacktrace,
				),
			})
			continue
		} else if err != nil {
			return nil, err
		}

		for _, mBeanAttrName := range mBeanAttrNames {
			jmxAttributes, err := c.GetMBeanAttrs(mBeanName, mBeanAttrName)
			if jmxErr, isJMXErr := IsJMXError(err); isJMXErr {
				result = append(result, &QueryAttrResponse{
					Status: QueryResponseStatusError,
					StatusMsg: fmt.Sprintf("error while querying mBean '%s', attribute: '%s', error message: %s, error cause: %s, stacktrace: %q",
						mBeanName,
						mBeanAttrName,
						jmxErr.Message,
						jmxErr.CauseMessage,
						jmxErr.Stacktrace,
					),
				})
				continue
			} else if err != nil {
				return nil, err
			}

			for _, attr := range jmxAttributes {
				result = append(result, &QueryAttrResponse{
					Attribute: attr,
					Status:    QueryResponseStatusSuccess,
					StatusMsg: "success",
				})
			}
		}
	}
	return result, nil
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
	case err, open := <-c.nrJMXProcess.ErrorC():
		if err == nil && !open {
			return nrjmx.ErrProcessNotRunning
		}
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
