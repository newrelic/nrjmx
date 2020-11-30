JAVA_VERSION 	?= jdk-11.0.9_11.1
DOCKER_BIN 		?= docker
DOCKER_CMD 		?= $(DOCKER_BIN) run --rm -it -v ${HOME}/.gradle:/root/.gradle -v $(CURDIR):/src/nrjmx -w /src/nrjmx adoptopenjdk/openjdk11:$(JAVA_VERSION)-centos

GRADLE_BIN		?= $(CURDIR)/gradlew
GRADLE_FLAGS	+=  --warn
GRADLE_FLAGS	+=  --stacktrace
#GRADLE_FLAGS	+=  --warning-mode all

.PHONY : package
package :
	@($(GRADLE_BIN) clean package $(GRADLE_FLAGS))

.PHONY : package/linux
package/linux :
	@($(GRADLE_BIN) clean package-linux $(GRADLE_FLAGS))

.PHONY : package/windows
package/windows :
	@($(GRADLE_BIN) clean package-windows $(GRADLE_FLAGS))

.PHONY : build
build : GRADLE_FLAGS += --info
build :
	@($(GRADLE_BIN) clean build $(GRADLE_FLAGS))

.PHONY : ci/build
ci/build :
	@($(DOCKER_CMD) make build)

.PHONY : ci/package
ci/package :
	@($(DOCKER_CMD) make package)

.PHONY : ci/package/linux
ci/package/linux :
	@($(DOCKER_CMD) make package/linux)

.PHONY : ci/package/windows
ci/package/windows :
	@($(DOCKER_CMD) make package/windows)

.PHONY : test
test:
	@($(GRADLE_BIN) clean test)
