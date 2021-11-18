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
	@cd $(GOMODULE_DIR); go clean -testcache; go test -v -timeout 60s github.com/newrelic/nrjmx