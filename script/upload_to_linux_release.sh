#!/usr/bin/env sh
set -e

#
#
#
# Upload binary to release page
#
#
#

INTEGRATION_PATH=$1
TAG=$2
SEMVER=`echo "${TAG}" | cut -c 2-`

echo "===> Uploading ${INTEGRATION_PATH}_${SEMVER}-1_amd64.deb to ${TAG}"
hub release edit -a "${INTEGRATION_PATH}_${SEMVER}-1_amd64.deb" -m "${TAG}" "${TAG}"

echo "===> Uploading ${INTEGRATION_PATH}-${SEMVER}-1.x86_64.rpm to ${TAG}"
hub release edit -a "${INTEGRATION_PATH}-${SEMVER}-1.x86_64.rpm" -m "${TAG}" "${TAG}"

echo "===> Uploading ${INTEGRATION_PATH}-${SEMVER}-noarch.ja to ${TAG}"
hub release edit -a "${INTEGRATION_PATH}-${SEMVER}-noarch.jar" -m "${TAG}" "${TAG}"

echo "===> Uploading ${INTEGRATION_PATH}-${SEMVER}-jlink.zip to ${TAG}"
hub release edit -a "${INTEGRATION_PATH}-${SEMVER}-jlink.zip" -m "${TAG}" "${TAG}"

echo "===> Uploading ${INTEGRATION_PATH}-${SEMVER}.tar.gz to ${TAG}"
hub release edit -a "${INTEGRATION_PATH}-${SEMVER}.tar.gz" -m "${TAG}" "${TAG}"

echo "===> Uploading ${INTEGRATION_PATH}-${SEMVER}-jlink.tar.gz to ${TAG}"
hub release edit -a "${INTEGRATION_PATH}-${SEMVER}-jlink.tar.gz" -m "${TAG}" "${TAG}"
