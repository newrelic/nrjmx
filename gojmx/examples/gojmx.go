package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cciutea/gojmx"
	"github.com/cciutea/gojmx/generated/jmx"
)

func main() {
	var defaultCtx = context.Background()

	client, err := gojmx.NewJMXServiceClient(defaultCtx)
	if err != nil {
		panic(err)
	}

	config := &jmx.JMXConfig{
		ConnURL:            "service:jmx:remote+https://localhost:9993",
		Username:           "admin",
		Password:           "Admin.123",
		KeyStore:           "./../tests/jboss/key/jboss.keystore2",
		KeyStorePassword:   "password",
		TrustStore:         "./../tests/jboss/key/jboss.truststore",
		TrustStorePassword: "password",
	}

	// config := &jmx.JMXConfig{
	// 	Hostname: "localhost",
	// 	Port:     9010,
	// 	// ConnURL:  "service:jmx:remote+https://localhost:9010",
	// 	// Username:           "admin",
	// 	// Password:           "Admin.123",
	// 	// KeyStore:           "./../tests/jboss/key/jboss.keystore",
	// 	// KeyStorePassword:   "password",
	// 	// TrustStore:         "./../tests/jboss/key/jboss.truststore",
	// 	// TrustStorePassword: "password",
	// }

	_, err = client.Connect(defaultCtx, config)
	if err != nil {
		panic(err)
	}

	result, err := client.QueryMbean(defaultCtx, "jboss.as.expr:subsystem=remoting,configuration=endpoint")
	if err != nil {
		panic(err)
	}

	for _, r := range result {
		c := gojmx.JMXAttributeValueConverter{r.GetValue()}
		fmt.Print(r.Attribute + ":")
		fmt.Print(c.GetValue())
		fmt.Print(",")
	}

	client.Close(1 * time.Second)

}
