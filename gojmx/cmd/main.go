package main

import (
	"context"
	"fmt"
	"github.com/newrelic/nrjmx/gojmx"
	"github.com/newrelic/nrjmx/gojmx/internal/nrprotocol"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

var prjDir, keystorePath, truststorepath string

func init() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	prjDir = filepath.Join(path, "../..")
	keystorePath = filepath.Join(prjDir, "test-server", "keystore")
	truststorepath = filepath.Join(prjDir, "test-server", "truststore")

	os.Setenv("NR_JMX_TOOL", filepath.Join(prjDir, "bin", "nrjmx"))
}

func main() {
	// THEN JMX connection can be oppened
	client, err := gojmx.NewClient(context.Background()).Init()
	//if err != nil {
	//	panic(err)
	//}
	config := &nrprotocol.JMXConfig{
		Hostname: "localhost",
		Port:     7199,
	}

	err = client.Connect(config, 10000)

	if err != nil {
		panic(err)
	}
	//defer client.Disconnect()
	//if err != nil {
	//	panic(err)
	//}

	client.WriteJunk()
	//client.Init()
	result, err := client.GetMBeanNames("*:*", 20000)
	fmt.Println(result, err)


	//client, err = client.Init()



	result, err = client.GetMBeanNames("*:*", 20000)
	//if err != nil {
	//	panic(err)
	//}
	fmt.Println(result, err)
	result, err = client.GetMBeanNames("*:*", 20000)
	//if err != nil {
	//	panic(err)
	//}
	fmt.Println(result, reflect.TypeOf(err))
	client.Init()

	result, err = client.GetMBeanNames("*:*", 20000)

	//if err != nil {
	//	panic(err)
	//}
	fmt.Println(result, err)
	fmt.Println("here")

	result, err = client.GetMBeanNames("*:*", 20000)
	fmt.Println(result, err)
	//err := jmx.Open("localhost", "7199", "", "")
	//if err != nil {
	//	panic(err)
	//}
	//
	//result, err := jmx.Query("test:type=Cat,*", 600000)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(result)
	time.Sleep(1 * time.Hour)

}
