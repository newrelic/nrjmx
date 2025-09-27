.PHONY : deps
deps:
	@($(DOCKER_BIN) build -t nrjmx_builder ./build/.)
# A new nrjmx_builder_fips image is created because:
# In ubuntu 16.04 the lastest version of rpm supported is 4.12
# Inorder to build RPM with SHA256 payload digest the rpm version should be 4.14+
# So, a new nrjmx_builder_fips image which use ubuntu 22.04 as base image.
	@($(DOCKER_BIN) build -t nrjmx_builder_fips -f ./build/fips.dockerfile ./build/.)

.PHONY : ci/build
ci/build: deps
	@($(DOCKER_CMD) make build)
	@($(DOCKER_CMD) make build-fips)
# Verify if the build is successful on a ubuntu 22.04 container.
	@($(DOCKER_FIPS_CMD) make build-fips)

.PHONY : ci/package
ci/package: deps
	@($(DOCKER_CMD) make package)

.PHONY : ci/test
ci/test: deps
	@($(DOCKER_CMD) make test)

.PHONY : ci/release
ci/release: deps
# Uploads all fips and non-fips packages expect FIPS RPM
	@($(DOCKER_CMD) make release)
# Uploads FIPS RPM packages
	@($(DOCKER_FIPS_CMD) make release/package-fips/rpm)

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
