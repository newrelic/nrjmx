FROM maven:3.6-jdk-8

ARG GH_VERSION='2.0.0'
ARG GO_VERSION='1.18.5'

RUN apt-get update && \
    apt-get install -y \
    make \
    rpm \
    gnupg2 \
    gpg-agent \
    debsigs \
    unzip \
    zip \
    gcc \
    openjdk-11-jdk

RUN ln -s /usr/lib/jvm/java-11-openjdk-amd64/ /usr/local/openjdk-11

ENV PATH="/usr/local/go/bin:$PATH"

# Install golang.
RUN curl -L "https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz" -o go.tar.gz && tar -C /usr/local -xzf go.tar.gz && export PATH=${PATH}:/usr/local/go/bin

RUN curl -L https://github.com/cli/cli/releases/download/v${GH_VERSION}/gh_${GH_VERSION}_linux_amd64.deb -o gh_${GH_VERSION}_linux_amd64.deb
RUN dpkg -i gh_${GH_VERSION}_linux_amd64.deb
