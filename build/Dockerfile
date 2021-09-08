FROM maven:3.6-jdk-8

ARG GH_VERSION='1.1.0'

RUN apt-get update && \
    apt-get install -y \
    make \
    rpm \
    gnupg2 \
    gpg-agent \
    debsigs \
    unzip \
    zip

RUN curl -L https://github.com/cli/cli/releases/download/v${GH_VERSION}/gh_${GH_VERSION}_linux_amd64.deb -o gh_${GH_VERSION}_linux_amd64.deb
RUN dpkg -i gh_${GH_VERSION}_linux_amd64.deb