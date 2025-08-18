FROM maven:3.6-jdk-11 as builder

ADD . .

RUN mvn package

FROM openjdk:11-jre

EXPOSE 4567
EXPOSE 7199

COPY --from=builder target/test-server.jar /
COPY --from=builder truststore /
COPY --from=builder keystore /
COPY --from=builder jmxremote.password /usr/local/openjdk-11/conf/management/jmxremote.password
COPY --from=builder jmxremote.access /usr/local/openjdk-11/conf/management/jmxremote.access

RUN chmod 400 /usr/local/openjdk-11/conf/management/jmxremote.password

# Environment variables for SSL-enabled JMX
ENV JAVA_OPTS_BASE -Dcom.sun.management.jmxremote.port=7199 \
      -Dcom.sun.management.jmxremote.rmi.port=7199 \
      -Djava.rmi.server.hostname=0.0.0.0 \
      -Dcom.sun.management.jmxremote.local.only=false

# Default to non-SSL for regular tests
ENV JAVA_OPTS ${JAVA_OPTS_BASE} \
      -Dcom.sun.management.jmxremote.authenticate=false \
      -Dcom.sun.management.jmxremote.ssl=false

# SSL-specific options (can be overridden at runtime)
ENV JAVA_OPTS_SSL ${JAVA_OPTS_BASE} \
      -Dcom.sun.management.jmxremote.authenticate=true \
      -Dcom.sun.management.jmxremote.ssl=true \
      -Dcom.sun.management.jmxremote.password.file=/usr/local/openjdk-11/conf/management/jmxremote.password \
      -Dcom.sun.management.jmxremote.access.file=/usr/local/openjdk-11/conf/management/jmxremote.access \
      -Djavax.net.ssl.keyStore=/keystore \
      -Djavax.net.ssl.keyStorePassword=password \
      -Djavax.net.ssl.trustStore=/truststore \
      -Djavax.net.ssl.trustStorePassword=password

# Install curl for health check
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*

# Health check
HEALTHCHECK --interval=5s --timeout=3s --start-period=30s --retries=10 \
    CMD curl -f http://localhost:4567/health || exit 1

CMD ["/bin/bash", "-c", "java ${JAVA_OPTS} -jar /test-server.jar"]