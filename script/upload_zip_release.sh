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

echo "===> Uploading ${INTEGRATION_PATH}-${SEMVER}-jlink.zip to ${TAG}"
hub release edit -a "${INTEGRATION_PATH}-${SEMVER}-jlink.zip" -m "${TAG}" "${TAG}"
