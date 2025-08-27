# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## Unreleased

## v2.10.0 - 2025-08-26

### 🚀 Enhancements
- Upgraded golang version to v1.24.6
- Switched to Windows 2025 runners

### 🐞 Bug fixes
- Fix publishing FIPS compliant nrjmx packages
- Added name and version to pom files

## v2.8.0 - 2025-04-23

### 🚀 Enhancements
- Add FIPS compliance packages for nrjmx

## v2.7.1 - 2025-04-02

### 🐞 Bug fixes
- Upgraded golang.org/x/net to v0.35.0

## v1.5.3 - 2020-09-14

### 🚀 Enhancements
- Removed `org.yaml` unused bundled dependency. This will reduce the size of the JAR file.
- Upgraded `jmxterm` to v1.0.2 (fixes a vulnerability in bundled jar)

## v1.5.2 - 2019-11-18

### 🐞 Bug fixes
- Install `jmxterm` in `/usr/lib/nrjmx` for deb packages.

## v1.5.1 - 2019-11-15

### 🐞 Bug fixes
- Install `jmxterm` in `/usr/bin` for deb packages.

## v1.5.0 - 2019-11-15

### 🚀 Enhancements
- Build JMXFetcher from full connection URL
- Set debug log entry on nice log lvl to be shown only for verbose mode
- Clean up cmd execution log entries
- Support custom connectors
- Increase test timeout to build on slow boxes
- Java version file
- Windows build
- Include `jmxterm` for troubleshooting queries within mvn packaging for tarball, deb, and rpm

## v1.4.1 - 2019-10-01

### 🐞 Bug fixes
- Fixed issue when parsing float NaN values.

## v1.4.0 - 2019-09-18

### 🐞 Bug fixes
- Upgrade project target to Java 1.8 and allow using a different Java version than the default one by configuring JAVA_HOME or NRIA_JAVA_HOME environment variables.

## v1.3.1 - 2019-06-17

### 🚀 Enhancements
- (Linux-only) tar.gz packaging as an alternative to the current package managers

## v1.3.0 - 2019-06-17

### 🚀 Enhancements
- Non standard (`jmxrmi`) URI path support via `-uriPath` argument.
- JBoss remoting v3 support for JBoss Domain-mode as default and Standalone-mode optionally.

## v1.2.1 - 2019-06-04

### 🐞 Bug fixes
- Fixed SSL connection with keyStore and trustStore 

## v1.1.2 - 2019-03-18

### 🚀 Enhancements
- Added remote argument for JMX remote connections

## v1.0.2 - 2018-09-12

### 🚀 Enhancements
- Catch all exceptions

## v1.0.0 - 2017-07-21

### 🚀 Enhancements
- Initial release
