# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## Next

- Build JMXFetcher from full connection URL
- Set debug log entry on nice log lvl to be shown only for verbose mode
- Clean up cmd execution log entries
- Increase test timeout to build on slow boxes
- Support custom connectors

## 1.4.1 (2019-10-01)
- Fixed issue when parsing float NaN values.

## 1.4.0 (2019-09-18)
- Upgrade project target to Java 1.8 and allow using a different Java version than 
the default one by configuring JAVA_HOME or NRIA_JAVA_HOME environment variables.

## 1.3.1 (2019-06-17)
- (Linux-only) tar.gz packaging as an alternative to the current package managers

## 1.3.0 (2019-06-17)
## Added
- Non standard (`jmxrmi`) URI path support via `-uriPath` argument.
- JBoss remoting v3 support for JBoss Domain-mode as default and Standalone-mode
  optionally.

## 1.2.1 (2019-06-04)
## Fixed
- Fixed SSL connection with keyStore and trustStore 

## 1.1.2 (2019-03-18)
### Added
- Added remote argument for JMX remote connections

## 1.0.2 (2018-09-12)
### Added
- Catch all exceptions

## 1.0.0 (2017-07-21)
### Added
- Initial release
