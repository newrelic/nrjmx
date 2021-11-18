CUR_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
GOMODULE_DIR:=$(CUR_DIR)/src/go/

deps:
	@$(DOCKER_BIN) build -t test-server $(CUR_DIR)/test-server/.

.PHONY : gotest
gotest: deps build
	@echo $(GOMODULE_DIR)
	@cd $(GOMODULE_DIR); go clean -testcache; go test -v -timeout 60s github.com/newrelic/nrjmx
