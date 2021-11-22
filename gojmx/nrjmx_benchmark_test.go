package gojmx

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/newrelic/nrjmx/gojmx/nrprotocol"
	"github.com/testcontainers/testcontainers-go"
)

var c testcontainers.Container
var ctx context.Context
var jmxClient *JMXClient

func init() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	prjDir = filepath.Join(path, "../..")
	keystorePath = filepath.Join(prjDir, "test-server", "keystore")
	truststorepath = filepath.Join(prjDir, "test-server", "truststore")

	os.Setenv("NR_JMX_TOOL", filepath.Join(prjDir, "bin", "nrjmx"))

	ctx = context.Background()

	c, err = runJMXServiceContainer(ctx)
	if err != nil {
		panic(err)
	}

	var data []map[string]interface{}
	for i := 0; i < 10; i++ {
		data = append(data, map[string]interface{}{
			"name":        fmt.Sprintf("%s-%d", "tomas", i),
			"doubleValue": 1.2,
			"floatValue":  2.2,
			"numberValue": 3,
			"boolValue":   true,
		})
	}

	// Populate the JMX Server with mbeans
	addMBeansBatch(ctx, c, data)

	// err = jmx.Open("localhost", "7199", "admin1234", "Password1!")
	// if err != nil {
	// 	panic(err)
	// }

	// 5	 491660541 ns/op

	//61	 110436625 ns/op
	//69	  81057981 ns/op
	//64	  90962180 ns/op
	//60	  87890616 ns/op
	jmxClient, _ = NewJMXServiceClient(ctx)
	config := &nrprotocol.JMXConfig{
		Hostname: "localhost",
		Port:     7199,
		UriPath:  "jmxrmi",
		// ConnectionURL: "service:jmx:remote+http://localhost:9990/",
		// Username:      "admin1234",
		// Password:      "Password1!",
	}

	_, err = jmxClient.Connect(ctx, config)
	if err != nil {
		panic(err)
	}
}

func BenchmarkTest2(b *testing.B) {

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		jmxClient.QueryMbean(ctx, "*:*")
		// assert.NoError(b, err)
		// fmt.Println(result)
	}

	// defer c.Terminate(context.Background())
}

// func BenchmarkTest2(b *testing.B) {

// 	// run the Fib function b.N times
// 	for n := 0; n < b.N; n++ {
// 		jmx.Query("*:*", 10000)
// 		// assert.NoError(b, err)
// 		// fmt.Println(result)
// 	}
// }
