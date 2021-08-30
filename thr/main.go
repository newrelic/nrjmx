package main

import (
	"context"
	"fmt"
	"thr/jmx"
)

func main() {

	var defaultCtx = context.Background()

	jmxProcess, err := StartJMXProcess(defaultCtx)
	if err != nil {
		panic(err)
	}

	client, err := NewJMXServiceClient(jmxProcess)
	if err != nil {
		handleError(jmxProcess, err)
	}

	config := &jmx.JMXConfig{
		// Hostname: "localhost",
		// Port:     9999,
		ConnURL:            "service:jmx:remote+https://0.0.0.0:9993",
		Username:           "admin",
		Password:           "Admin.1234",
		KeyStorePassword:   "password",
		KeyStore:           "/Users/cciutea/workspace/nr/int/nri-jmx/src/github.com/newrelic/nri-jmx/docs/jboss-eap-https-tutorial/key/jboss.keystore",
		TrustStorePassword: "password",
		TrustStore:         "/Users/cciutea/workspace/nr/int/nri-jmx/src/github.com/newrelic/nri-jmx/docs/jboss-eap-https-tutorial/key/jboss.truststore",
	}

	ok, err := client.Connect(defaultCtx, config)

	if e, ok := err.(*jmx.JMXConnectionError); ok {
		handleError(jmxProcess, e)
	}

	fmt.Println("connected: ", ok)
	if err != nil {
		handleError(jmxProcess, err)
	}

	result, err := client.QueryMbean(defaultCtx, "*:*")
	if err != nil {
		handleError(jmxProcess, err)
	}

	for _, attr := range result {
		fmt.Println(fmt.Sprintf("%s=%v", attr.Attribute, getValue(*attr.Value)))
	}

	//[1,"disconnect",1,1,{}]
	err = client.Disconnect(defaultCtx)
	handleError(jmxProcess, err)

	err = jmxProcess.Stop()
	if err != nil {
		panic(err)
	}
}

func handleError(jmxProcess *JMXProcess, err error) {
	if err == nil {
		return
	}
	err2 := jmxProcess.Stop()
	if err2 != nil {
		panic(err2)
	}
	panic(err)
}

func getValue(value jmx.JMXAttributeValue) interface{} {
	switch value.ValueType {
	case jmx.ValueType_DOUBLE:
		return value.GetDoubleValue()
	case jmx.ValueType_STRING:
		return value.GetStringValue()
	case jmx.ValueType_BOOL:
		return value.GetBoolValue()
	case jmx.ValueType_INT:
		return value.GetIntValue()
	default:
		return nil
	}
}
