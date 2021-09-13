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
	@docker build -t nrjmx_builder ./build/.

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


