#!/usr/bin/env sh
set -e

#
#
#
# Upload binary to release page
#
#
#

INTEGRATION=$1
TAG=$2
SEMVER=`echo "${TAG}" | cut -c 2-`

hub release edit -a "${INTEGRATION}.${SEMVER}.msi" -m "${TAG}" "${TAG}"
