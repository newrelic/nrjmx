[![Community Project header](https://github.com/newrelic/opensource-website/raw/master/src/images/categories/Community_Project.png)](https://opensource.newrelic.com/oss-category/#community-project)

# New Relic JMX fetcher

A tool for extracting data out of any application exposing a JMX interface.

## Installation

### Using a package manager

Common package managers can be used for this purpose: *yum, apt, zypper*.

For example: `yum install nrjmx`

### Using the tarball

You can download and decompress the Java executable from the [downloads page](http://download.newrelic.com/infrastructure_agent/binaries/linux/noarch/).

### New Relic JMX integration (`nri-jmx`)

`nrjmx` is *not* bundled within the `nri-jmx` package. Still, it's declared as a dependency.

So if you have `nrjmx` already installed while installing `nri-jmx`, the installed version is kept. Otherwise, `nri-jmx` will try to get the latest `nrjmx` release.

## Getting started

nrjmx expects the connection parameters to the JMX interface as command line arguments.

```bash
$ ./bin/nrjmx -hostname 127.0.0.1 -port 7199 -username user -password pwd
```

nrjmx reads lines from the standard input which should contain object name patterns for which we want to fetch their attributes. For each line, it gets the beans matching the pattern and outputs a JSON with all the attributes found.

For example, if you want to fetch some beans from Cassandra JMX metrics, you can execute:

```bash
$ echo
"org.apache.cassandra.metrics:type=Table,keyspace=*,scope=*,name=ReadLatency" | java -jar target/nrjmx-0.0.1-SNAPSHOT-jar-with-dependencies.jar -hostname 127.0.0.1 -port 7199 -username user -password pwd
```

## Usage

Additional options are listed below.

## Custom protocols

JMX allows use of custom protocols to communicate with the application. To use a custom protocol you have to include the custom connectors in the `nrjmx` classpath.

By default, nrjmx will include the sub-folder connectors in its class path. If this folder does not exist create it under the folder where you have nrjmx installed.

For example, to add support for JBoss, create a folder `connectors` under the default (Linux) library path `/usr/lib/nrjmx/` (`/usr/lib/nrjmx/connectors`) and copy the custom connector `jar` (`$JBOSS_HOME/bin/client/jboss-cli-client.jar`) into it. You can now execute JMX queries against JBoss.

### Remote URL connection

If you want to use a remoting-jmx URL you can use the flag `-remote`. In this case it uses the remoting connection URL: `service:jmx:remote://host:port` instead of `service:jmx:rmi:///jndi/rmi://host:port/jmxrmi`.

This sets a URI ready for JBoss Domain mode.

> You will need to add support for the custom JBoss protocol. See the previous section `Custom protocols`.

#### JBoss standalone mode

This is supported via `-remoteJBossStandalone` and sets a connection URL to `service:jmx:remote+http://host:port`.

Example of usage with remoting:

```bash
$ ./bin/nrjmx -hostname 127.0.0.1 -port 7199 -username user -password pwd -remote
```
> You will need to add support for the custom JBoss protocol. See the previous section `Custom protocols`.

### Non-Standard JMX Service URI

If your JMX provider uses a non-standard JMX service URI path (default path is `jmxrmi`), you can use the flag `-uriPath` to specify the path portion (without `/` prefix).

For example:

- A default URI path could be like: `service:jmx:rmi:///jndi/rmi://localhost:1689/jmxrmi` (path is last path of the URI without the prefix `/`)
- ForgeRock OpenDJ uses a JMX service URI like: `service:jmx:rmi:///jndi/rmi://localhost:1689/org.opends.server.protocols.jmx.client-unknown`

To extract data from this application:

```bash
$ ./bin/nrjmx -hostname localhost -port 1689 -uriPath "org.opends.server.protocols.jmx.client-unknown" -username user -password pwd
```

### Troubleshooting

If you are having difficulties with `nrjmx` to get data out of your JMX service, we provide a CLI tool (`jmxterm`) to help you [troubleshoot](./TROUBLESHOOT.md).

## Building

nrjmx uses Gradle for generating the binaries:

```bash
$ ./gradlew build
```

This creates the modularised `nrjmx.jar` file under the `./build/libs` directory as well as a `tar` and `zip` under the `build/distributions` directory. It can also build RPM & DEB packages.

```bash
$ ./gradlew buildDeb  # Debian package
$ ./gradlew buildRpm  # RPM package
$ ./gradlew distTar   # Tar ball
$ ./gradlew distZip   # ZIP package
```

## Support

New Relic hosts and moderates an online forum where customers can interact with New Relic employees as well as other customers to get help and share best practices. Like all official New Relic open source projects, there's a related Community topic in the New Relic Explorers Hub. You can find this project's topic/threads here:

https://discuss.newrelic.com/c/support-products-agents/new-relic-infrastructure

## Using testservers for testing

It builds a test service that introduces some monitoring.

### Build

`./gradlew :test-server-jdk8:build` will generate both a runnable jar file as well as copy the appropriate files from `src/docker` to a location where it can be used to build a container. 
Replace `test-server-jdk8` with `test-server-jdk11` to build a runnable jar with JDK11


### Run (with JMX enabled)

The project itself does not run the container. 
Instead this happens when `./gradlew :test` is executed. 
The containers are built from within testcontainers and then executed.

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
./gradlew :test-server-jdk8:install
docker build test-server-jdk8/build/install/test-server-jdk8
```

(Replace with `test-server-jdk11` as appropriate).

### REST API:

In the port `4567`:

* `POST /cat`
    * BODY: `{"name":"Isidoro"}` would register in JMX a cat named Isidoro to the registry name

* `PUT /clear`
    * Will clear all the cats from JMX

### Example JMX test

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

### Generating keys

```
keytool -genkeypair -dname "cn=server, ou=nrjmx, o=NR, c=US" -keystore serverkeystore -keyalg RSA -alias serverkey -validity 180 -storepass serverpass -keypass serverpass
keytool -exportcert -keystore serverkeystore -alias serverkey -storepass serverpass -file server.cer
keytool -import -v -trustcacerts -alias serverkey -file server.cer -keystore clienttruststore  -storepass clienttrustpass -noprompt
keytool -genkeypair -dname "cn=client, ou=test, o=NR, c=US" -keystore clientkeystore -keyalg RSA -alias clientkey -validity 180 -storepass clientkeystore -keypass clientkeystore
keytool -exportcert -keystore clientkeystore -alias clientkey -storepass clientkeystore -file client.cer
keytool -import -v -trustcacerts -alias clientkey -file client.cer -keystore servertruststore -storepass servertrustpass -noprompt
```

## Contributing
We encourage your contributions to improve New Relic JMX fetcher! Keep in mind when you submit your pull request, you'll need to sign the CLA via the click-through using CLA-Assistant. You only have to sign the CLA one time per project.
If you have any questions, or to execute our corporate CLA, required if your contribution is on behalf of a company,  please drop us an email at opensource@newrelic.com.

## License
New Relic JMX fetcher is licensed under the [Apache 2.0](http://apache.org/licenses/LICENSE-2.0.txt) License.
