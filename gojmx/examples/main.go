/*
 * Copyright 2021 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"context"
	"fmt"
	"github.com/newrelic/nrjmx/gojmx"
	"os"
	"path/filepath"
)

func init() {
	// Uncomment this when you want use the nrjmx.jar build from the project bin directory.
	_ = os.Setenv("NR_JMX_TOOL", filepath.Join("..", "bin", "nrjmx"))

	// Uncomment this when you want to run both: golang debugger and java debugger.
	_ = os.Setenv("NRIA_NRJMX_DEBUG", "true")
}

func main() {
	// JMX Client configuration.
	config := &gojmx.JMXConfig{
		Hostname:         "localhost",
		Port:             7199,
		RequestTimeoutMs: 100000,
		Verbose:          true,
	}

	// Connect to JMX endpoint.
	client, err := gojmx.NewClient(context.Background()).Open(config)
	handleError(err)

	defer client.Close()

	//// Get the mBean names.
	//mBeanNames, err := client.QueryMBeanNames("java.lang:type=*")
	//handleError(err)
	//
	//// Get the Attribute names for each mBeanName.
	//for _, mBeanName := range mBeanNames {
	//	mBeanAttrNames, err := client.GetMBeanAttributeNames(mBeanName)
	//	handleError(err)
	//
	//	// Get the attribute value for each mBeanName and mBeanAttributeName.
	//	jmxAttrs, err := client.GetMBeanAttributes(mBeanName, mBeanAttrNames...)
	//	if err != nil {
	//		fmt.Println(err)
	//		continue
	//	}
	//	for _, attr := range jmxAttrs {
	//		if attr.ResponseType == gojmx.ResponseTypeErr {
	//			fmt.Println(attr.StatusMsg)
	//			continue
	//		}
	//		printAttr(attr)
	//	}
	//}

	// Or use QueryMBean call which wraps all the necessary requests to get the values for an MBeanNamePattern.
	response, err := client.QueryMBeanAttributes("*:*")
	handleError(err)
	for _, attr := range response {
		if attr.ResponseType == gojmx.ResponseTypeErr {
			fmt.Println(attr.StatusMsg)
			continue
		}
		printAttr(attr)
	}
}

func printAttr(jmxAttr *gojmx.AttributeResponse) {
	_, _ = fmt.Fprintf(
		os.Stdout,
		"Attribute Name: %s\nAttribute Value: %v\nAttribute Value Type: %v\n\n",
		jmxAttr.Name,
		jmxAttr.GetValue(),
		jmxAttr.ResponseType,
	)
}

func handleError(err error) {
	if jmxErr, ok := gojmx.IsJMXError(err); ok {
		fmt.Println("JMXError message:", jmxErr.Message)
		fmt.Println("JMXError stacktrace:", jmxErr.Stacktrace)
		os.Exit(1)
	} else if jmxConnErr, ok := gojmx.IsJMXConnectionError(err); ok {
		fmt.Println("JMXConnectionError message:", jmxConnErr.Message)
		os.Exit(2)
	} else if err != nil {
		panic(err)
	}
}
