#!/bin/bash
set -e
#
#
# Upload target artifacts to GH Release assets
#
#
cd target
for filename in $(find . -name "*.msi" -o -name "*.rpm" -o -name "*.deb" -o -name "*.zip" -o -name "*.tar.gz");do
  echo "===> Uploading to GH $TAG: ${filename}"
      gh release upload $TAG $filename
done