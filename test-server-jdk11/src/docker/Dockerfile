#FROM maven:3.6-jdk-8 as builder
#
#ADD ../.. .
#
#RUN mvn package

FROM openjdk:11

EXPOSE 4567
EXPOSE 7199

COPY bin /application/bin/
COPY lib /application/lib/
COPY server*store /
#RUN chmod +x /application/bin/test-server
ENTRYPOINT exec /application/bin/test-server-jdk11