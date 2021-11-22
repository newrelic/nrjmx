
`gojmx` is a go module that allows fetching data from a JMX
endpoint.

# Installing
You can use gojmx in your project by running `go get` command:

    go get -u github.com/newrelic/nrjmx/gojmx

Next, import the dependency into your application:

```go
import "github.com/newrelic/nrjmx/gojmx"
```

This go module will call nrjmx library java application in order to fetch the metrics.

To install nrjmx library on your system, packages are available on this repository
in the release assets or on our package manager repositories stored [here](https://download.newrelic.com/infrastructure_agent/). You can folow this [documentation](https://docs.newrelic.com/docs/infrastructure/install-infrastructure-agent/linux-installation/install-infrastructure-monitoring-agent-linux/#ubuntu-repository) to include our repository into your package manager.

nrjmx library is also available as a [tarball](https://download.newrelic.com/infrastructure_agent/binaries/linux/noarch/).

# Usage
After nrjmx library was installed on the system, you need a
JMX service exposed. For this example, we will use the `test-server` available on the root of this repository:

```bash
docker build -t nrjmx/test-server .
docker run -d -p 7199:7199 nrjmx/test-server
```

```go
func main() {
	ctx := context.Background()

    // Start the nrjmx process
	client, err := gojmx.NewJMXServiceClient(ctx)
	if err != nil {
		panic(err)
	}

    // Connect to the JMX endpoint
	err = client.Connect(ctx, &nrprotocol.JMXConfig{
		Hostname: "localhost",
		Port:     7199,
		UriPath:  "jmxrmi",
	})
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

    // Query for the mbeans. The library also supports
    // using `*` as a wildcard
	result, err := client.QueryMbean(ctx, "*:*")
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
```