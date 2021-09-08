DOCKER_BIN 		?= docker
MAVEN_BIN       ?= mvn

DOCKER_CMD 		?= $(DOCKER_BIN) run --rm -t \
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
	@docker build -t nrjmx_builder .

.PHONY : build
build:
	@($(MAVEN_BIN) clean package -DskipTests -P \!deb,\!rpm,\!tarball,\!test)

.PHONY : package
package:
	@($(MAVEN_BIN) clean -DskipTests package)

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

.PHONY : release/sign
release/sign:
	@echo "=== [release/sign] signing packages"
	@bash $(CURDIR)/sign.sh

.PHONY : release/publish
release/publish:
	@echo "=== [release/publish] publishing artifacts"
	@bash $(CURDIR)/upload_artifacts_gh.sh

.PHONY : release-linux
release-linux: package release/sign release/publish
	@echo "=== [release-linux] full pre-release cycle complete for nix"

.PHONY : ci/release
ci/release: deps
		@($(DOCKER_CMD) make release-linux)