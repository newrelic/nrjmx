PROJECT_WORKSPACE	?= $(CURDIR)
INCLUDE_BUILD_DIR	?= $(PROJECT_WORKSPACE)/build

DOCKER_BIN 		?= docker
MAVEN_BIN       ?= mvn

TAG				?= v0.0.0

include $(INCLUDE_BUILD_DIR)/build.mk
include $(INCLUDE_BUILD_DIR)/ci.mk
include $(INCLUDE_BUILD_DIR)/release.mk
include $(INCLUDE_BUILD_DIR)/release.mk
include $(INCLUDE_BUILD_DIR)/gomodule.mk