FROM jboss/wildfly:18.0.0.Final

# Add management user
RUN /opt/jboss/wildfly/bin/add-user.sh admin1234 Password1! --silent

# Expose ports
EXPOSE 8080 9990 9999

# Create custom standalone configuration
USER root
RUN yum update -y && yum install -y curl && yum clean all
USER jboss

# Add JMX configuration via CLI script
RUN echo 'embed-server --server-config=standalone.xml --std-out=echo' > /tmp/jmx-config.cli && \
    echo '/system-property=jboss.bind.address.management:add(value="0.0.0.0")' >> /tmp/jmx-config.cli && \
    echo '/system-property=java.rmi.server.hostname:add(value="0.0.0.0")' >> /tmp/jmx-config.cli && \
    echo '/system-property=com.sun.management.jmxremote:add(value="true")' >> /tmp/jmx-config.cli && \
    echo '/system-property=com.sun.management.jmxremote.port:add(value="9999")' >> /tmp/jmx-config.cli && \
    echo '/system-property=com.sun.management.jmxremote.rmi.port:add(value="9999")' >> /tmp/jmx-config.cli && \
    echo '/system-property=com.sun.management.jmxremote.authenticate:add(value="false")' >> /tmp/jmx-config.cli && \
    echo '/system-property=com.sun.management.jmxremote.ssl:add(value="false")' >> /tmp/jmx-config.cli && \
    echo '/system-property=com.sun.management.jmxremote.local.only:add(value="false")' >> /tmp/jmx-config.cli && \
    echo 'stop-embedded-server' >> /tmp/jmx-config.cli

# Apply the CLI configuration
RUN /opt/jboss/wildfly/bin/jboss-cli.sh --file=/tmp/jmx-config.cli

# Health check with longer timeout for Wildfly
HEALTHCHECK --interval=15s --timeout=10s --start-period=120s --retries=10 \
    CMD curl -f http://localhost:9990/management --connect-timeout 5 || exit 1

# Environment variables for runtime configuration
ENV JBOSS_OPTS="-Djboss.bind.address.management=0.0.0.0 -Djava.rmi.server.hostname=0.0.0.0"

# Start with comprehensive binding
CMD ["/opt/jboss/wildfly/bin/standalone.sh", "-b", "0.0.0.0", "-bmanagement", "0.0.0.0"]