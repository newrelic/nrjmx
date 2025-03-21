PROJECT_WORKSPACE	?= $(CURDIR)
INCLUDE_BUILD_DIR	?= $(PROJECT_WORKSPACE)/build
GO_VERSION 		?= $(shell grep '^go ' go.mod | awk '{print $$2}')

DOCKER_BIN 		?= docker
MAVEN_BIN       ?= mvn

TAG				?= v0.0.0

# FIPS-compliant options - This definition will be propagated to included files
MAVEN_FIPS_OPTS = -Dhttps.protocols=TLSv1.2 -Djdk.tls.client.protocols=TLSv1.2 -Djavax.net.ssl.keyStoreType=PKCS12 -Dcom.sun.net.ssl.checkRevocation=true -Dssl.TrustManagerFactory.algorithm=PKIX
export MAVEN_FIPS_OPTS

export GOEXPERIMENT=boringcrypto

include $(INCLUDE_BUILD_DIR)/build.mk
include $(INCLUDE_BUILD_DIR)/ci.mk
include $(INCLUDE_BUILD_DIR)/release.mk


DOCKER_ENV_VARS = -e MAVEN_FIPS_OPTS -e GOEXPERIMENT

DOCKER_CMD 		?= $(DOCKER_BIN) run --rm -t \
					--name "nrjmx-builder" \
					-v $(HOME)/.docker/:/root/.docker/ \
					-v /var/run/docker.sock:/var/run/docker.sock \
					-v $(CURDIR):/src/nrjmx \
					-w /src/nrjmx \
					$(DOCKER_ENV_VARS) \
					-e GITHUB_TOKEN \
					-e TAG \
					-e GPG_MAIL \
					-e GPG_PASSPHRASE \
					-e GPG_PRIVATE_KEY_BASE64 \
					nrjmx_builder