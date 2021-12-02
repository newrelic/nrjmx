package main

import (
	"context"
	"fmt"

	"github.com/newrelic/nrjmx/gojmx"
	"github.com/newrelic/nrjmx/gojmx/nrprotocol"
)

func main() {
	ctx := context.Background()

	client, err := gojmx.NewJMXServiceClient(ctx)
	if err != nil {
		panic(err)
	}

	err = client.Connect(ctx, &nrprotocol.JMXConfig{
		Hostname: "localhost",
		Port:     7199,
		UriPath:  "jmxrmi",
	})
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	result, err := client.QueryMbean(ctx, "*:*")
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
