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

ZIP_FILE=${INTEGRATION_PATH}_windows_${SEMVER}_amd64.zip
echo "===> Uploading ${ZIP_FILE} to ${TAG}"
hub release edit -a "${ZIP_FILE}" -m "${TAG}" "${TAG}"
