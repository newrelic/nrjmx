# NR-JMX
New Relic's JMX fetcher, a simple tool for extracting data out of any application exposing a JMX interface.

## Build
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
mvn clean package -P \!deb,\!rpm
```

Note: nrjmx is targetted to work with Java 7 ( [pom.xml](https://github.com/newrelic/nrjmx/blob/master/pom.xml#L217-L219) )

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

## Updates
### 1.0.5
- Add `--url` command line parameter to allow the end-user to use atypical service urls. Defaults to `service:jmx:rmi:///jndi/rmi://%s:%s/jmxrmi` which is `service:jmx:rmi:///jndi/rmi://<hostname>:<port>/jmxrmi`
  - Use caution when quoting the string!