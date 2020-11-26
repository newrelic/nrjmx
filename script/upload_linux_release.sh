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

DEB_FILE=${INTEGRATION_PATH}_${SEMVER}-1_amd64.deb
echo "===> Uploading ${DEB_FILE} to ${TAG}"
hub release edit -a "${DEB_FILE}" -m "${TAG}" "${TAG}"

RPM_FILE=${INTEGRATION_PATH}-${SEMVER}-1.x86_64.rpm
echo "===> Uploading ${RPM_FILE} to ${TAG}"
hub release edit -a "${RPM_FILE}" -m "${TAG}" "${TAG}"

NOARCH_JAR=${INTEGRATION_PATH}-${SEMVER}-noarch.jar
echo "===> Uploading ${NOARCH_JAR} to ${TAG}"
hub release edit -a "${NOARCH_JAR}" -m "${TAG}" "${TAG}"

TAR_FILE=${INTEGRATION_PATH}_linux_${SEMVER}_amd64.tar.gz
echo "===> Uploading ${TAR_FILE} to ${TAG}"
hub release edit -a "${TAR_FILE}" -m "${TAG}" "${TAG}"
