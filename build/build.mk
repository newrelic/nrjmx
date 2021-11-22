.PHONY : build
build:
	@($(MAVEN_BIN) clean package -DskipTests -P \!deb,\!rpm,\!tarball,\!test,\!tarball)

.PHONY : test
test:
	@($(MAVEN_BIN) clean test -P test)

CUR_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
GOMODULE_DIR:=$(CUR_DIR)/gojmx/

go-test: godeps build
	@echo $(GOMODULE_DIR)
	@cd $(GOMODULE_DIR); go clean -testcache; go test -v -timeout 300s github.com/newrelic/nrjmx

DOCKER_THRIFT=$(DOCKER_BIN) run --rm -t \
					--name "nrjmx-code-generator" \
					-v $(CURDIR):/src/nrjmx \
					-w /src/nrjmx \
					cciutea/thrift

.PHONY : code-gen-deps
code-gen-deps:
		@($(DOCKER_BIN) build -t nrjmx-code-generator ./commons/.)

.PHONY : code-gen
code-gen: 
	@($(DOCKER_THRIFT) thrift -r --out src/main/java/ --gen java ./commons/nrjmx.thrift)
	@($(DOCKER_THRIFT) thrift -r --out gojmx/ --gen go:package_prefix=github.com/newrelic/nrjmx/gojmx/,package=nrprotocol ./commons/nrjmx.thrift)

TRACKED_GEN_DIR=src/main/java/nrprotocol \
				gojmx/nrprotocol
.PHONY : check-gen-code
check-gen-code: code-gen
	@echo "Checking the generated code..." ; \
	if [ `git status --porcelain --untracked-files=no $(TRACKED_GEN_DIR) | wc -l` -gt 0 ]; then \
		echo "Code generator produced different code, make sure you pushed the latest changes!"; \
		exit 1;	\
	fi

