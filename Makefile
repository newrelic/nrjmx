JAVA_VERSION 	?= jdk-11.0.9_11.1
DOCKER_BIN 		?= docker
DOCKER_CMD 		?= $(DOCKER_BIN) run --rm -it -v ${HOME}/.gradle:/root/.gradle -v $(CURDIR):/src/nrjmx -w /src/nrjmx adoptopenjdk/openjdk11:$(JAVA_VERSION)-centos

.PHONY : package
package :
	@($(DOCKER_CMD) ./gradlew clean package --warn --stacktrace)

.PHONY : package/linux
package/linux :
	@($(DOCKER_CMD) ./gradlew clean package-linux --warn --stacktrace)

.PHONY : package/windows
package/windows :
	@($(DOCKER_CMD) ./gradlew clean package-windows --warn --stacktrace)
