.PHONY : build
build:
	@($(MAVEN_BIN) clean package -DskipTests -P \!deb,\!rpm,\!tarball,\!test)

.PHONY : test
test:
	@($(MAVEN_BIN) clean test -P test)

CUR_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
GOMODULE_DIR:=$(CUR_DIR)/src/go/

go-test: godeps build
	@echo $(GOMODULE_DIR)
	@cd $(GOMODULE_DIR); go clean -testcache; go test -v -timeout 300s github.com/newrelic/nrjmx

DOCKER_THRIFT=$(DOCKER_BIN) run --rm -t \
					--name "nrjmx-code-generator" \
					-v $(CURDIR):/src/nrjmx \
					-w /src/nrjmx \
					nrjmx-code-generator

.PHONY : code-gen-deps
code-gen-deps:
		@($(DOCKER_BIN) build -t nrjmx-code-generator ./commons/.)

.PHONY : code-gen
code-gen: code-gen-deps
	@($(DOCKER_THRIFT) thrift -r --out src/main/java/ --gen java ./commons/nrjmx.thrift)
	@($(DOCKER_THRIFT) thrift -r --out src/go/ --gen go:package=protocol ./commons/nrjmx.thrift)


