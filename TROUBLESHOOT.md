# Troubleshooting via jmxterm

If you are having difficulties with `nrjmx` to get data out of your JMX service, this interactive CLI tool could help you out in the process.


## Install

Since `nrjmx` version 1.5.0 `jmxterm` comes bundled within the `nrjmx` package installer.

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
