DOCKER_CMD 		?= $(DOCKER_BIN) run --rm -t \
					--name "nrjmx-builder" \
					-v $(HOME)/.docker/:/root/.docker/ \
					-v /var/run/docker.sock:/var/run/docker.sock \
					-v $(CURDIR):/src/nrjmx \
					-w /src/nrjmx \
					-e GITHUB_TOKEN \
					-e TAG \
					-e GPG_MAIL \
					-e GPG_PASSPHRASE \
					-e GPG_PRIVATE_KEY_BASE64 \
					nrjmx_builder

.PHONY : deps
deps:
	@($(DOCKER_BIN) build -t nrjmx_builder ./build/.)

.PHONY : ci/build
ci/build: deps
	@($(DOCKER_CMD) make build)

.PHONY : ci/package
ci/package: deps
	@($(DOCKER_CMD) make package)

.PHONY : ci/test
ci/test: deps
	@($(DOCKER_CMD) make test)

.PHONY : ci/release
ci/release: deps
	@($(DOCKER_CMD) make release)

.PHONY : ci/go-test
ci/go-test: deps go-test-utils
	@($(DOCKER_CMD) make go-test -o go-test-utils)

.PHONY : ci/go-test-jdk11
ci/go-test-jdk11: deps go-test-utils
	@($(DOCKER_CMD) /bin/bash -c 'export PATH=/usr/local/openjdk-11/bin:$$PATH; java -version; make go-test -o go-test-utils')

TRACKED_GEN_DIR=src/main/java/org/newrelic/nrjmx/v2/nrprotocol \
				gojmx/internal/nrprotocol
.PHONY : ci/check-gen-code
ci/check-gen-code: validate-thrift-version code-gen
	@echo "Checking the generated code..." ; \
	if [ `git status --porcelain $(TRACKED_GEN_DIR) | wc -l` -gt 0 ]; then \
		echo "Code generator produced different code, make sure you pushed the latest changes!"; \
		git --no-pager diff $(TRACKED_GEN_DIR); \
		exit 1;	\
	fi ; \
	echo "Success!"

.PHONY : validate-thrift-version
validate-thrift-version: deps
	@printf '\n------------------------------------------------------\n'
	@printf 'Validating thrift version\n'
	@($(DOCKER_CMD) build/validate_thrift_version.sh)

.PHONY: ci/docker/publish
ci/docker/publish: code-gen-utils
	@printf '\n------------------------------------------------------\n'
	@printf 'Publishing docker image\n'
	@($(DOCKER_BIN) push ohaiops/nrjmx-code-generator:$(THRIFT_VERSION))
	@($(DOCKER_BIN) push ohaiops/nrjmx-code-generator:latest)

.PHONY: ci/snyk-test-java
ci/snyk-test-java:
	@docker run --rm -t \
			--name "nrjmx-snyk-test-java" \
			-v $(CURDIR):/src/nrjmx \
			-w /src/nrjmx \
			-e SNYK_TOKEN \
			snyk/snyk:maven-3-jdk-11 snyk test --severity-threshold=high

.PHONY: ci/snyk-test-go
ci/snyk-test-go:
	@docker run --rm -t \
			--name "nrjmx-snyk-test-go" \
			-v $(CURDIR):/src/nrjmx \
			-w /src/nrjmx/gojmx \
			-e SNYK_TOKEN \
			-e GOFLAGS="-buildvcs=false" \
			snyk/snyk:golang snyk test --severity-threshold=high
