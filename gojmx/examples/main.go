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
	//_ = os.Setenv("NRIA_NRJMX_DEBUG", "true")
}

func main() {
	// JMX Client configuration.
	config := &gojmx.JMXConfig{
		Hostname:        "localhost",
		Port:            7199,
		RequestTimoutMs: 10000,
	}

	// Connect to JMX endpoint.
	client, err := gojmx.NewClient(context.Background()).Open(config)
	handleError(err)

	// Query the mBean names.
	mBeanNames, err := client.GetMBeanNames("java.lang:type=*")
	handleError(err)

	// Query the Attribute names for each mBeanName.
	for _, mBeanName := range mBeanNames {
		mBeanAttrNames, err := client.GetMBeanAttrNames(mBeanName)
		handleError(err)

		for _, mBeanAttrName := range mBeanAttrNames {
			// Query the attribute value for each mBeanName and mBeanAttributeName.
			jmxAttrs, err := client.GetMBeanAttrs(mBeanName, mBeanAttrName)
			if err != nil {
				fmt.Println(err)
				continue
			}
			for _, jmxAttr := range jmxAttrs {
				printAttr(jmxAttr)
			}
		}
	}
}

func printAttr(jmxAttr *gojmx.JMXAttribute) {
	_, _ = fmt.Fprintf(
		os.Stdout,
		"Attribute Name: %s\nAttribute Value: %v\nAttribute Value Type: %v\n\n",
		jmxAttr.Attribute,
		jmxAttr.GetValue(),
		jmxAttr.ValueType,
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
