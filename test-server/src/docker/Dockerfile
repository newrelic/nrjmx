#FROM maven:3.6-jdk-8 as builder
#
#ADD ../.. .
#
#RUN mvn package

FROM openjdk:8-jre

EXPOSE 4567
EXPOSE 7199

#COPY --from=builder target/test-server.jar /
#COPY --from=builder servertruststore /
#COPY --from=builder serverkeystore /
COPY test-server.jar /
COPY servertruststore /
COPY serverkeystore /

ENTRYPOINT exec java $JAVA_OPTS -jar /test-server.jar