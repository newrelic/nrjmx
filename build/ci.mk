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

TRACKED_GEN_DIR=src/main/java/org/newrelic/nrjmx/v2/nrprotocol \
				gojmx/internal/nrprotocol
.PHONY : ci/check-gen-code
ci/check-gen-code: code-gen
	@echo "Checking the generated code..." ; \
	if [ `git status --porcelain --untracked-files=no $(TRACKED_GEN_DIR) | wc -l` -gt 0 ]; then \
		echo "Code generator produced different code, make sure you pushed the latest changes!"; \
		git --no-pager diff; \
		exit 1;	\
	fi ; \
	echo "Success!"