# JMX test service

It builds a test service that introduces some monitoring.

## Build

`./gradlew :test-server:build` will generate both a runnable jar file as well as copy the appropriate files
from `src/docker` to a location where it can be used to build a container

## Run (with JMX enabled)

The project itself does not run the container. 
Instead this happens when `./gradlew :test` is executed. 
The container is built from within testcontainers and then executed.

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

(This is already done with the test code).

## Building the container manually

If you need to build and run the container manually then you can do:

```
./gradlew :test-server:dockerFile
docker build test-server/build/dockerFiles
```
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

    
