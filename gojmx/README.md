
`gojmx` is a go module that allows fetching data from a JMX
endpoint.

# Installing
You can use gojmx in your project by running `go get` command:

    go get github.com/newrelic/nrjmx/gojmx

Next, import the dependency into your application:

```go
import "github.com/newrelic/nrjmx/gojmx"
```

This go module will call nrjmx library in order to fetch the metrics.

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
// JMX Client configuration.
config := &gojmx.JMXConfig{
    Hostname:         "localhost",
    Port:             7199,
    RequestTimeoutMs: 10000,
}

// Connect to JMX endpoint.
client, err := gojmx.NewClient(context.Background()).Open(config)
handleError(err)

// Get the mBean names.
mBeanNames, err := client.QueryMBeanNames("java.lang:type=*")
handleError(err)

// Get the Attribute names for each mBeanName.
for _, mBeanName := range mBeanNames {
    mBeanAttrNames, err := client.GetMBeanAttributeNames(mBeanName)
    handleError(err)

    // Get the attribute value for each mBeanName and mBeanAttributeName.
    jmxAttrs, err := client.GetMBeanAttributes(mBeanName, mBeanAttrNames...)
    if err != nil {
        fmt.Println(err)
        continue
    }
    for _, attr := range jmxAttrs {
        if attr.ResponseType == gojmx.ResponseTypeErr {
            fmt.Println(attr.StatusMsg)
            continue
        }
        printAttr(attr)
    }
}

// Or use QueryMBean call which wraps all the necessary requests to get the values for an MBeanNamePattern.
response, err := client.QueryMBeanAttributes("java.lang:type=*")
handleError(err)
for _, attr := range response {
    if attr.ResponseType == gojmx.ResponseTypeErr {
        fmt.Println(attr.StatusMsg)
        continue
    }
    printAttr(attr)
}
```

You can find the full example in the examples directory.

# Custom connectors
JMX allows the use of custom connectors to communicate with the application. In order to use a custom connector, you have to include the custom connectors in the nrjmx classpath.

By default, the sub-folder connectors is in the classpath. If this folder does not exist, create it under the folder where nrjmx is installed.

For example, to add support for JBoss, create a folder named connectors under the default (Linux) library path /usr/lib/nrjmx/ (/usr/lib/nrjmx/connectors/) and copy the custom connector jar ($JBOSS_HOME/bin/client/jboss-cli-client.jar) into it. You can now execute JMX queries against JBoss.
