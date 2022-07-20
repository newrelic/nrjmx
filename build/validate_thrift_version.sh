#!/usr/bin/env bash

JAVA_THRIFT_VERSION=$( mvn dependency:list | grep -P "^\[INFO\]\s+org.apache.thrift:libthrift:jar:" | grep -Eo "[0-9]+\.[0-9]+\.[0-9]+" )
DOCKER_THRIFT_VERSION=$( grep "THRIFT_VERSION=" commons/Dockerfile | grep -Eo "[0-9]+\.[0-9]+\.[0-9]+" )
cd gojmx
GO_THRIFT_VERSION=$( go mod graph | grep  -E "^github.com/newrelic/nrjmx/gojmx github.com/apache/thrift" | grep -Eo "[0-9]+\.[0-9]+\.[0-9]+" )

if [[ "${JAVA_THRIFT_VERSION}" == "${GO_THRIFT_VERSION}" && "${GO_THRIFT_VERSION}" == "${DOCKER_THRIFT_VERSION}" ]];then
  exit 0
fi

echo "Different thrift versions found:"
echo "Java: ${JAVA_THRIFT_VERSION}"
echo "Go: ${GO_THRIFT_VERSION}"
echo "Docker: ${DOCKER_THRIFT_VERSION}"
echo ""
echo "Follow the instructions in: https://github.com/newrelic/nrjmx/blob/master/DEVELOP_V2.md#updating-thrift-version"
exit 1