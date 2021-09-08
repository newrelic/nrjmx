DOCKER_BIN 		?= docker
DOCKER_CMD 		?= $(DOCKER_BIN) run --rm -it -v $(HOME)/.docker/:/root/.docker/ -v /var/run/docker.sock:/var/run/docker.sock -v $(CURDIR):/src/nrjmx -w /src/nrjmx nrjmx_builder
MAVEN_BIN       ?= mvn

.PHONY : deps
deps:
	@docker build -t nrjmx_builder .

.PHONY : build
build:
	@($(MAVEN_BIN) clean package -DskipTests -P \!deb,\!rpm,\!tarball,\!test)

.PHONY : package
package:
	@($(MAVEN_BIN) clean package)

.PHONY : test
test:
	@($(MAVEN_BIN) clean test -P test)

.PHONY : ci/build
ci/build: deps
	@($(DOCKER_CMD) make build)

.PHONY : ci/package
ci/package: deps
	@($(DOCKER_CMD) make package)

.PHONY : ci/test
ci/test: deps
	@($(DOCKER_CMD) make test)