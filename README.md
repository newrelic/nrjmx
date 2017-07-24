# NR-JMX
New Relic's JMX fetcher, a simple tool for extracting data out of any application exposing a JMX interface.

## Build
NR-JMX uses Maven for generating the binaries:

```bash
$ mvn package
```

This will create the jar file under the /target/ directory.

## Usage
The applicaton just expects the connection parameters to the JMX interface.

```bash
$ java -jar target/nrjmx-0.0.1-SNAPSHOT-jar-with-dependencies.jar -hostname 127.0.0.1 -port 7199 -username user  -password pwd
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
