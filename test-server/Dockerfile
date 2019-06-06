FROM openjdk:8-jre

EXPOSE 4567
EXPOSE 7199
EXPOSE 9100

ADD target/test-server.jar /
ADD servertruststore /
ADD serverkeystore /

ENTRYPOINT exec java $JAVA_OPTS -jar /test-server.jar
