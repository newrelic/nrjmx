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
	"strconv"
)

func main() {
	if len(os.Args) != 3 {
		panic("missing hostname and port")
	}

	hostname := os.Args[1]

	port, err := strconv.ParseInt(os.Args[2], 10, 32)
	if err != nil {
		panic(err)
	}

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	config := &gojmx.JMXConfig{
		Hostname:         hostname,
		Port:             int32(port),
		RequestTimeoutMs: 60000,
	}

	client, err := gojmx.NewClient(ctx).Open(config)
	if err != nil {
		panic(err)
	}

	fmt.Println(client.GetClientVersion())

	result, err := client.QueryMBeanAttributes("*:*")
	if err != nil {
		panic(err)
	}

	fmt.Println(gojmx.FormatJMXAttributes(result))
}
