FROM maven:3.6-jdk-8 as builder

ADD . .

RUN mvn package

FROM openjdk:8-jre

EXPOSE 4567
EXPOSE 7199

COPY --from=builder target/test-server.jar /
COPY --from=builder truststore /
COPY --from=builder keystore /
COPY --from=builder jmxremote.password /usr/local/openjdk-8/lib/management/jmxremote.password
COPY --from=builder jmxremote.access /usr/local/openjdk-8/lib/management/jmxremote.access

RUN chmod 400 /usr/local/openjdk-8/lib/management/jmxremote.password

ENV JAVA_OPTS -Dcom.sun.management.jmxremote.port=7199 \
      -Dcom.sun.management.jmxremote.authenticate=false \
      -Dcom.sun.management.jmxremote.ssl=false \
      -Dcom.sun.management.jmxremote=true \
      -Dcom.sun.management.jmxremote.rmi.port=7199 \
      -Djava.rmi.server.hostname=localhost

CMD ["/bin/bash", "-c", "java ${JAVA_OPTS} -jar /test-server.jar"]