.PHONY : validate-thrift-version
validate-thrift-version:
	@printf '\n------------------------------------------------------\n'
	@printf 'Validating thrift version\n'
	@build/validate_thrift_version.sh

.PHONY : build
build: validate-thrift-version
	@($(MAVEN_BIN) clean package -DskipTests -P \!deb,\!rpm,\!tarball,\!test,\!tarball)

.PHONY : test
test:
	@($(MAVEN_BIN) clean test -P test)

CUR_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
GOMODULE_DIR:=$(CUR_DIR)/gojmx/

.PHONY : go-test-utils
go-test-utils:
	@($(DOCKER_BIN) build -t test-server $(CUR_DIR)/test-server/.)
	@($(DOCKER_BIN) build -t test_jboss -f $(CUR_DIR)/jboss.dockerfile $(CUR_DIR))

.PHONY : go-test
go-test: go-test-utils build
	@echo $(GOMODULE_DIR)
	@cd $(GOMODULE_DIR); \
	go vet .; \
	go clean -testcache; \
	go test -v -timeout 300s github.com/newrelic/nrjmx/gojmx

DOCKER_THRIFT=$(DOCKER_BIN) run --rm -t \
					--name "nrjmx-code-generator" \
					-v $(CURDIR):/src/nrjmx \
					-w /src/nrjmx \
					ohaiops/nrjmx-code-generator:latest

.PHONY : code-gen-utils
code-gen-utils:
		@($(DOCKER_BIN) build \
						-t ohaiops/nrjmx-code-generator:$(THRIFT_VERSION) \
						-t ohaiops/nrjmx-code-generator:latest \
						--build-arg THRIFT_VERSION=$(THRIFT_VERSION) ./commons/.)


.PHONY: docker/publish
docker/publish: code-gen-utils
	@printf '\n------------------------------------------------------\n'
	@printf 'Publishing docker image\n'
	@($(DOCKER_BIN) push ohaiops/nrjmx-code-generator:$(THRIFT_VERSION))
	@($(DOCKER_BIN) push ohaiops/nrjmx-code-generator:latest)


.PHONY : code-gen
code-gen: validate-thrift-version
	rm -rf src/main/java/org/newrelic/nrjmx/v2/nrprotocol
	rm -rf gojmx/internal/nrprotocol
	@($(DOCKER_THRIFT) thrift -r --out src/main/java/ --gen java:generated_annotations=suppress ./commons/nrjmx.thrift)
	@($(DOCKER_THRIFT) thrift -r --out gojmx/internal/ --gen go:package_prefix=github.com/newrelic/nrjmx/gojmx/internal/,package=nrprotocol ./commons/nrjmx.thrift)



