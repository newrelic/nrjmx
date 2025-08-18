FROM jboss/wildfly:18.0.0.Final

# Add management user
RUN /opt/jboss/wildfly/bin/add-user.sh admin1234 Password1! --silent

# Expose ports
EXPOSE 8080 9990 9999

# Add JMX configuration via environment variables in standalone.conf
USER root
RUN echo 'JAVA_OPTS="$JAVA_OPTS -Djboss.bind.address.management=0.0.0.0"' >> /opt/jboss/wildfly/bin/standalone.conf && \
    echo 'JAVA_OPTS="$JAVA_OPTS -Dcom.sun.management.jmxremote"' >> /opt/jboss/wildfly/bin/standalone.conf && \
    echo 'JAVA_OPTS="$JAVA_OPTS -Dcom.sun.management.jmxremote.port=9990"' >> /opt/jboss/wildfly/bin/standalone.conf && \
    echo 'JAVA_OPTS="$JAVA_OPTS -Dcom.sun.management.jmxremote.rmi.port=9990"' >> /opt/jboss/wildfly/bin/standalone.conf && \
    echo 'JAVA_OPTS="$JAVA_OPTS -Dcom.sun.management.jmxremote.authenticate=false"' >> /opt/jboss/wildfly/bin/standalone.conf && \
    echo 'JAVA_OPTS="$JAVA_OPTS -Dcom.sun.management.jmxremote.ssl=false"' >> /opt/jboss/wildfly/bin/standalone.conf && \
    echo 'JAVA_OPTS="$JAVA_OPTS -Djava.rmi.server.hostname=0.0.0.0"' >> /opt/jboss/wildfly/bin/standalone.conf && \
    echo 'JAVA_OPTS="$JAVA_OPTS -Dcom.sun.management.jmxremote.local.only=false"' >> /opt/jboss/wildfly/bin/standalone.conf

USER jboss

# Start with comprehensive binding
CMD ["/opt/jboss/wildfly/bin/standalone.sh", "-b", "0.0.0.0", "-bmanagement", "0.0.0.0"]
