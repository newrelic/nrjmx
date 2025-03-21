ARG UBUNTU_VERSION=16.04

# Use Ubuntu as the base image
FROM ubuntu:${UBUNTU_VERSION}

ARG GO_VERSION=1.23.2
ARG ARCH='amd64'
ARG GH_VERSION='2.61.0'
ARG DOCKER_VERSION='24.0.5'
ARG BUILDX_VERSION='0.11.2'
ARG MAVEN_VERSION='3.6.3'

# Install dependencies
RUN apt-get update && apt-get install -y \
    gnupg-agent \
    unzip \
    zip \
    curl \
    wget \
    expect \
    git \
    tar \
    gcc \
    jq \
    g++ \
    gnupg2 \
    gnupg-agent \
    debsigs \
    rpm \
    build-essential \
    software-properties-common \
    python-software-properties \
    gcc-arm-linux-gnueabi \
    dpkg-sig \
    openssl \
    libssl-dev \
    ca-certificates \
    gcc-aarch64-linux-gnu

# Install Go with FIPS mode enabled (BoringCrypto)
RUN curl -sSL https://golang.org/dl/go${GO_VERSION}.linux-${ARCH}.tar.gz -o go${GO_VERSION}.linux-${ARCH}.tar.gz && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-${ARCH}.tar.gz && \
    rm go${GO_VERSION}.linux-${ARCH}.tar.gz

# Set Go environment variables
ENV PATH="/usr/local/go/bin:/go/bin:${PATH}"
ENV GOPATH="/go"
ENV GOEXPERIMENT=boringcrypto
ENV GOFLAGS="-buildvcs=false"

# Install Docker
RUN curl -fsSL https://download.docker.com/linux/static/stable/x86_64/docker-${DOCKER_VERSION}.tgz -o docker.tgz && \
    tar --extract --file docker.tgz --strip-components 1 --directory /usr/local/bin/ && \
    rm docker.tgz

# Install Buildx
RUN mkdir -p ~/.docker/cli-plugins && \
    curl -sSL https://github.com/docker/buildx/releases/download/v${BUILDX_VERSION}/buildx-v${BUILDX_VERSION}.linux-${ARCH} -o ~/.docker/cli-plugins/docker-buildx && \
    chmod +x ~/.docker/cli-plugins/docker-buildx

# Since the user does not match the owners of the repo "git rev-parse --is-inside-work-tree" fails and goreleaser does not populate projectName
# https://stackoverflow.com/questions/72978485/git-submodule-update-failed-with-fatal-detected-dubious-ownership-in-repositor
RUN git config --global --add safe.directory '*'

# Copy signing scripts
COPY ./nix/sign.sh ./nix/sign_deb.exp ./nix/sign_rpm.exp ./nix/sign_tar.exp /usr/local/bin/

# Make sure signing scripts are executable
RUN chmod +x /usr/local/bin/sign.sh /usr/local/bin/sign_deb.exp /usr/local/bin/sign_rpm.exp /usr/local/bin/sign_tar.exp

# Install GitHub CLI
RUN curl -L https://github.com/cli/cli/releases/download/v${GH_VERSION}/gh_${GH_VERSION}_linux_amd64.deb -o gh_${GH_VERSION}_linux_amd64.deb && \
    dpkg -i gh_${GH_VERSION}_linux_amd64.deb && \
    rm gh_${GH_VERSION}_linux_amd64.deb

# Add Java repository and install both Java 8 and 11
RUN add-apt-repository -y ppa:openjdk-r/ppa && \
    apt-get update && \
    apt-get install -y openjdk-8-jdk openjdk-11-jdk && \
    ln -s /usr/lib/jvm/java-11-openjdk-amd64/ /usr/local/openjdk-11

# Setup FIPS compliant OpenSSL
RUN mkdir -p /etc/ssl/fips-enabled

# Configure OpenSSL for FIPS mode
RUN echo "openssl_conf = openssl_init" > /etc/ssl/openssl.cnf && \
    echo "" >> /etc/ssl/openssl.cnf && \
    echo "[openssl_init]" >> /etc/ssl/openssl.cnf && \
    echo "providers = provider_sect" >> /etc/ssl/openssl.cnf && \
    echo "" >> /etc/ssl/openssl.cnf && \
    echo "[provider_sect]" >> /etc/ssl/openssl.cnf && \
    echo "default = default_sect" >> /etc/ssl/openssl.cnf && \
    echo "fips = fips_sect" >> /etc/ssl/openssl.cnf && \
    echo "" >> /etc/ssl/openssl.cnf && \
    echo "[default_sect]" >> /etc/ssl/openssl.cnf && \
    echo "activate = 1" >> /etc/ssl/openssl.cnf && \
    echo "" >> /etc/ssl/openssl.cnf && \
    echo "[fips_sect]" >> /etc/ssl/openssl.cnf && \
    echo "activate = 1" >> /etc/ssl/openssl.cnf

# Set environment variable for FIPS mode
ENV OPENSSL_FORCE_FIPS_MODE=1

# Install Maven with FIPS-compliant settings
RUN mkdir -p /usr/share/maven && \
    curl -fsSL https://archive.apache.org/dist/maven/maven-3/${MAVEN_VERSION}/binaries/apache-maven-${MAVEN_VERSION}-bin.tar.gz | \
    tar -xzC /usr/share/maven --strip-components=1 && \
    ln -s /usr/share/maven/bin/mvn /usr/bin/mvn

# Configure Maven to use FIPS-compliant encryption
RUN mkdir -p /root/.m2 && \
    echo "<settings>" > /root/.m2/settings.xml && \
    echo "  <profiles>" >> /root/.m2/settings.xml && \
    echo "    <profile>" >> /root/.m2/settings.xml && \
    echo "      <id>fips</id>" >> /root/.m2/settings.xml && \
    echo "      <properties>" >> /root/.m2/settings.xml && \
    echo "        <https.protocols>TLSv1.2</https.protocols>" >> /root/.m2/settings.xml && \
    echo "        <jdk.tls.client.protocols>TLSv1.2</jdk.tls.client.protocols>" >> /root/.m2/settings.xml && \
    echo "      </properties>" >> /root/.m2/settings.xml && \
    echo "    </profile>" >> /root/.m2/settings.xml && \
    echo "  </profiles>" >> /root/.m2/settings.xml && \
    echo "  <activeProfiles>" >> /root/.m2/settings.xml && \
    echo "    <activeProfile>fips</activeProfile>" >> /root/.m2/settings.xml && \
    echo "  </activeProfiles>" >> /root/.m2/settings.xml && \
    echo "</settings>" >> /root/.m2/settings.xml

ENV MAVEN_HOME /usr/share/maven
ENV PATH $MAVEN_HOME/bin:$PATH

# Verify installations
RUN java -version && \
    /usr/local/openjdk-11/bin/java -version && \
    mvn --version && \
    go version && \
    gh --version

WORKDIR /src/nrjmx