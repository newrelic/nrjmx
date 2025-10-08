#!/usr/bin/env sh
set -e
#
#
#
# Sign FIPS RPM's & DEB's with proper digest algorithms
#
#
#
# Function to start gpg-agent if not running
start_gpg_agent() {
    if ! pgrep -x "gpg-agent" > /dev/null
    then
        echo "Starting gpg-agent..."
        eval $(gpg-agent --daemon)
    else
        echo "gpg-agent is already running."
    fi
}

# Ensure gpg-agent is running
start_gpg_agent

# Ensure GPG configuration directory exists
mkdir -p ~/.gnupg
chmod 700 ~/.gnupg

# Sign FIPS RPM's with proper digest configuration
echo "===> Create .rpmmacros for FIPS-compliant RPM signing"
echo "%_gpg_name ${GPG_MAIL}" >> ~/.rpmmacros
echo "%_signature gpg" >> ~/.rpmmacros
echo "%_gpg_path /root/.gnupg" >> ~/.rpmmacros
echo "%_gpgbin /usr/bin/gpg" >> ~/.rpmmacros

echo "%__gpg_sign_cmd %{__gpg} gpg --no-verbose --no-armor --passphrase ${GPG_PASSPHRASE} --no-secmem-warning --digest-algo sha256 -u "%{_gpg_name}" -sbo %{__signature_filename} %{__plaintext_filename}" >> ~/.rpmmacros

# FIPS-specific digest settings
echo "%_binary_filedigest_algorithm 8" >> ~/.rpmmacros
echo "%_source_filedigest_algorithm 8" >> ~/.rpmmacros
echo "%_binary_payload w9.gzdio" >> ~/.rpmmacros

echo "===> Importing GPG private key from GHA secrets..."
printf %s ${GPG_PRIVATE_KEY_BASE64} | base64 -d | gpg --batch --import -

echo "===> Importing GPG signature, needed from Goreleaser to verify signature"
gpg --export -a ${GPG_MAIL} > /tmp/RPM-GPG-KEY-${GPG_MAIL}
rpm --import /tmp/RPM-GPG-KEY-${GPG_MAIL}

cd dist

# Only sign FIPS packages
for rpm_file in $(find -name "*fips*.rpm"); do
  if [ -f "$rpm_file" ]; then
    echo "===> Signing FIPS package: $rpm_file"
    
    # Check if package already has proper digests
    echo "===> Checking current digest status of $rpm_file"
    rpm -Kv $rpm_file
    
    # Sign the package using batch mode for FIPS
    sign_rpm.exp $rpm_file ${GPG_PASSPHRASE} batch
    
    echo "===> Post-signing verification of $rpm_file"
    rpm -Kv $rpm_file
  fi
done
