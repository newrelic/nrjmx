[![Community Project header](https://github.com/newrelic/opensource-website/raw/master/src/images/categories/Community_Project.png)](https://opensource.newrelic.com/oss-category/#community-project)

# New Relic JMX fetcher

A tool for extracting data out of any application exposing a JMX interface.

## Installation

### Via package manager

Usual package managers could be used for this purpose: *yum, apt, zypper*.

Ie: `yum install nrjmx`

### Via tarball

You can download and decompress the Java executable as well from the [downloads url](http://download.newrelic.com/infrastructure_agent/binaries/linux/noarch/)

### New Relic JMX integration (nri-jmx)

`nrjmx` is *not* bundled within the `nri-jmx` package. But, it's declared as a dependency.

So while installing `nri-jmx` if you have `nrjmx` already installed it keeps the installed version, otherwise it'll try to get the latest `nrjmx` release.

## Getting Started

This tool just expects the connection parameters to the JMX interface as command line arguments.

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

## Usage

Additional options are listed below.

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
### Troubleshooting

If you are having difficulties with `nrjmx` to get data out of your JMX service, we provide a CLI tool (jmxterm) to help you [troubleshoot](./TROUBLESHOOT.md).


## Building
mv  
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

## Support

New Relic hosts and moderates an online forum where customers can interact with New Relic employees as well as other customers to get help and share best practices. Like all official New Relic open source projects, there's a related Community topic in the New Relic Explorers Hub. You can find this project's topic/threads here:

https://discuss.newrelic.com/c/support-products-agents/new-relic-infrastructure

## Contributing
We encourage your contributions to improve New Relic JMX fetcher! Keep in mind when you submit your pull request, you'll need to sign the CLA via the click-through using CLA-Assistant. You only have to sign the CLA one time per project.
If you have any questions, or to execute our corporate CLA, required if your contribution is on behalf of a company,  please drop us an email at opensource@newrelic.com.

## License
New Relic JMX fetcher is licensed under the [Apache 2.0](http://apache.org/licenses/LICENSE-2.0.txt) License.
