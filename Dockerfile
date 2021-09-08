FROM maven:3.6-jdk-8

RUN apt-get update && \
    apt-get install -y \
    make \
    rpm \
    gnupg2 \
    gpg-agent \
    debsigs \
    unzip \
    zip