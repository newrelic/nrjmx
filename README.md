# NR-JMX
New Relic's JMX fetcher, a simple tool for extracting data out of any application exposing a JMX interface.

## Installation

### Via package manager

Usual package managers could be used for this purpose: *yum, apt, zypper*.

Ie: `yum install nrjmx`

### Via tarball

You can download and decompress the Java executable as well from the [downloads url](http://download.newrelic.com/infrastructure_agent/binaries/linux/noarch/)

### nri-jmx relation

`nrjmx` is *not* bundled within the `nri-jmx` package. But, it's declared as a dependency. 

So while installing `nri-jmx` if you have `nrjmx` already installed it keeps the installed version, otherwise it'll try to get the latest `nrjmx` release.


## Custom Build
NR-JMX uses Maven for generating the binaries:

```bash
$ mvn package
```

This will create the `nrjmx.jar` file under the `./bin/` directory. Copy
`bin/nrjmx` & `bin/nrjmx.jar` files to your preferred location. Both files must
be located under the same folder.

It will also create DEB and RPM packages to automatically install NR-JMX. If you
want to skip DEB and RPM packages (e.g. because your development machine does not
provide the required tools), you can disable the `deb` and `rpm` Maven profiles from
the command line:

```bash
mvn clean package -P \!deb,\!rpm,\!tarball,\!test
```

## Configuring java version

`Note: nrjmx is targetted to work with Java 8`

After installation, nrjmx will use the default java version installed on the environment (the one set in the $PATH).
You can configure a different java version for nrjmx by adding one of the following environment variables at the end of the `/etc/environment` file:

### e.g.:
`JAVA_HOME=/usr/lib/jvm/jdk1.x.yz`

or

`NRIA_JAVA_HOME=/usr/lib/jvm/jdk1.x.yz`

The standard way is using JAVA_HOME, but the NRIA_JAVA_HOME takes precedence over JAVA_HOME in case your JAVA_HOME is already set and the version does not fit the NRJMX requirements.

Then pass the environment variable through the infra agent by appending the following lines to newrelic-infra.yml config file:
```
passthrough_environment:
    - JAVA_HOME
```
(replace JAVA_HOME by NRIA_JAVA_HOME if you have used the later property instead the standard one)

## Usage
The applicaton just expects the connection parameters to the JMX interface.

```bash
$ ./bin/nrjmx -hostname 127.0.0.1 -port 7199 -username user -password pwd
```

The tool will read lines from the standard input which should contain object
name patterns for which we want to fetch their attributes. For each line, it
will get the beans matching the pattern and output a JSON which all the
attributes found.

For instance, if you want to fetch some beans from Cassandra JMX metrics, you
could execute:

```bash
$ echo
"org.apache.cassandra.metrics:type=Table,keyspace=*,scope=*,name=ReadLatency" | java -jar target/nrjmx-0.0.1-SNAPSHOT-jar-with-dependencies.jar -hostname 127.0.0.1 -port 7199 -username user -password pwd
```

## Custom protocols

JMX allows use of custom protocols to communicate with the application. In order to use a custom protocol you have to include the custom connectors in the nrjmx classpath.
By default nrjmx will include the sub-folder connectors in it's class path. If this folder does not exist create it under the fodler where you have nrjmx installed.

For example, to add support for JBoss, create a folder `connectors` under the default (Linux) library path `/usr/lib/nrjmx/` (`/usr/lib/nrjmx/connectors`) and copy the custom connector `jar` (`$JBOSS_HOME/bin/client/jboss-cli-client.jar`) into it. You can now execute JMX queries against JBoss.

### Remote URL connection

If you want to use a remoting-jmx URL you can use the flag `-remote`. In this case it will use the remoting connection URL: `service:jmx:remote://host:port` instead of `service:jmx:rmi:///jndi/rmi://host:port/jmxrmi`

This sets URI ready for JBoss Domain mode.

Note: you will need to add support for the custom JBoss protocol. See the previous section `Custom protocols`.

#### JBoss Standalone mode

This is supported via `-remoteJBossStandalone` and will set connection URL to `service:jmx:remote+http://host:port`.

Example of usage with remoting:
```bash
$ ./bin/nrjmx -hostname 127.0.0.1 -port 7199 -username user -password pwd -remote
```
Note: you will need to add support for the custom JBoss protocol. See the previous section `Custom protocols`.

### Non-Standard JMX Service URI 

If your JMX provider uses a non-standard JMX service URI path (default path is `jmxrmi`), you can use the flag `-uriPath` to specify the path portion (without `/` prefix).

For example:

- A default URI path could be like: `service:jmx:rmi:///jndi/rmi://localhost:1689/jmxrmi` (path is last path of the URI without the prefix `/`)
- ForgeRock OpenDJ uses a JMX service URI like: `service:jmx:rmi:///jndi/rmi://localhost:1689/org.opends.server.protocols.jmx.client-unknown`

To extract data from this application:
```bash
$ ./bin/nrjmx -hostname localhost -port 1689 -uriPath "org.opends.server.protocols.jmx.client-unknown" -username user -password pwd
```
