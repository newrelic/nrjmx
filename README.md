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

---


# Troubleshooting via jmxterm

If you are having difficulties with `nrjmx` to get data out of your JMX service, this interactive CLI tool could help you out in the process.


## Install

Since `nrjmx` version X.X.X `jmxterm` comes bundled within the `nrjmx` package installer.

For previous versions you can easily get it from https://docs.cyclopsgroup.org/jmxterm

Or via your usual package manager. Ie, for Mac:
```
brew install jmxterm
```


## Connect

As `jmxterm` is an interactive CLI REPL, you could simple launch `jmxterm` and then `open <URL>`. `URL` can be a `<PID>`, `<hostname>:<port>` or full qualified JMX service URL. For example:

```
 open localhost:7199,
 open jmx:service:...
```

There's an option to open connection before getting into the REPL mode via `-l/--url` argument, ie:

```
jmxterm -l localhost:7199
```


### Connectivity issues

If you face connection issues with `nrjmx` you can specify the connection URL in the same way via `-C/--connURL` argument, ie:

```
jmxterm -l service:jmx:rmi:///jndi/rmi://localhost:7199/jmxrmi
nrjmx   -C service:jmx:rmi:///jndi/rmi://localhost:7199/jmxrmi
```


## Verbose

There's a verbose mode that will help with extra information output when troubleshooting in both tools:

```
jmxterm -v/--verbose
nrjmx   -v/--verbose
```


## Query

We are interested on fetching:

- from a JMX *domain*
- some MBean *objects*
- and some MBean *attributes*

JMX MBeans queries are launched in `nrjmx` using a glob pattern matching query against both *DOMAIN* and *BEAN* object. All *readable* bean attributes and values fill be retrieved.

The way to do so on `nrjmx` is concatenating both with the `:` separator, like:  `"DOMAIN_PATTERN:BEAN_PATTERN"`.

For instance to retrieve:
- all beans "*:*"
- all beans of *type* `Foo` "*:type=Foo,*"
- beans from domain Bar and *type* `Foo` "Bar:type=Foo,*"

> If there are issues while fetching these values possible errors like `java.lang.UnsupportedOperationException` will be printed by `nrjmx` into **stderr** output.

You can navigate through your exposed JMX data using `jmxterm` REPL mod.

- `domains` list available domains
- `domain <DOMAIN_NAME>` will set current domain
- `beans` list available MBeans within your current domain

Querying at `jmxterm` uses the same glob fashion mode. Just take into account that this tool divides *DOMAIN* and *BEAN* queries in 2 steps. So you can use

- `domain <PATTERN>` causes subsequent bean queries to run against matched domains, ie `*` unsets the domain
- `bean `<NAME>` sets or retrieves bean from current domain context, `*` unsets the domain
- `get `<ATTRIBUTES>` fetches bean attributes info from current domain & bean context

> A one liner query is possible as well, ie for querying all available attributes: `get -d DOMAIN -b BEAN *`


## Help

Both tools provide extra help via:

```
jmxterm -h/--help
nrjmx   -h/--help
```

