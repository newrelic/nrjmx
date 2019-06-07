# JMX test service

It builds a test service that introduces some monitoring.

## Build

`mvn clean package` will generate both a runnable jar file in `target/test-server.jar` as
  well as a container image named `testserver:latest`

## Run (with JMX enabled)

* As a local service:
```
export JAVA_OPTS="-Dcom.sun.management.jmxremote.port=7199 -Dcom.sun.management.jmxremote.authenticate=false -Dcom.sun.management.jmxremote.ssl=false -Dcom.sun.management.jmxremote=true -Dcom.sun.management.jmxremote.rmi.port=7199 -Djava.rmi.server.hostname=localhost"
java -jar target/test-server.jar
```
* As a container:

```
export JAVA_OPTS="-Dcom.sun.management.jmxremote.port=7199 -Dcom.sun.management.jmxremote.authenticate=false -Dcom.sun.management.jmxremote.ssl=false -Dcom.sun.management.jmxremote=true -Dcom.sun.management.jmxremote.rmi.port=7199 -Djava.rmi.server.hostname=localhost"
docker run --rm --name testserver -p 4567:4567 -p 7199:7199 --env JAVA_OPTS="$JAVA_OPTS"  -it testserver
```

It uses the following ports:

* `4567`: HTTP REST port
* `7199`: JMX RMI port

If you want to enable SSL, `JAVA_OPTS` should be:

```
-Dcom.sun.management.jmxremote.authenticate=false
-Dcom.sun.management.jmxremote.ssl=true
-Dcom.sun.management.jmxremote.ssl.need.client.auth=true 
-Dcom.sun.management.jmxremote.registry.ssl=true 
-Djavax.net.ssl.keyStore=/serverkeystore 
-Djavax.net.ssl.keyStorePassword=serverpass 
-Djavax.net.ssl.trustStore=/servertruststore 
-Djavax.net.ssl.trustStorePassword=servertrustpass
```

The client certificates can be found in [../src/test/resources](../src/test/resources) (passwords are `clienttrustpass`
and `clientpass`)

## REST API:

In the port `4567`:

* `POST /cat`
    * BODY: `{"name":"Isidoro"}` would register in JMX a cat named Isidoro to the registry name

* `PUT /clear`
    * Will clear all the cats from JMX

## Example JMX test

```
$ curl -X POST -d '{"name":"Isidoro"}' http://localhost:4567/cat
ok!
$ curl -X POST -d '{"name":"Heathcliff"}' http://localhost:4567/cat
ok!
$ ./nrjmx
test:type=Cat,*
{"test:type\u003dCat,name\u003dIsidoro,attr\u003dName":"Isidoro","test:type\u003dCat,name\u003dHeathcliff,attr\u003dName":"Heathcliff"}
$ curl -X PUT http://localhost:4567/clear
ok!
$ ./nrjmx
test:type=Cat,*
{}
```

    
