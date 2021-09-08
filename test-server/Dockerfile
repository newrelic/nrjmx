FROM maven:3.6-jdk-8 as builder

ADD . .

RUN mvn package

FROM openjdk:8-jre

EXPOSE 4567
EXPOSE 7199

COPY --from=builder target/test-server.jar /
COPY --from=builder truststore /
COPY --from=builder keystore /

ENTRYPOINT exec java $JAVA_OPTS -jar /test-server.jar