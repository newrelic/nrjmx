FROM jboss/wildfly:18.0.0.Final

# Add user with FIPS-compliant password (complex with special chars, numbers, mixed case)
RUN /opt/jboss/wildfly/bin/add-user.sh admin1234 Password1! --silent

# Copy FIPS configuration file
COPY fips-jboss.properties /opt/jboss/wildfly/standalone/configuration/

# Configure JBoss/Wildfly for FIPS
RUN echo 'JAVA_OPTS="$JAVA_OPTS -Djava.security.properties=/opt/jboss/wildfly/standalone/configuration/fips-jboss.properties -Djavax.net.ssl.keyStoreType=PKCS12 -Djdk.tls.client.protocols=TLSv1.2 -Dhttps.protocols=TLSv1.2 -Dcom.sun.net.ssl.checkRevocation=true"' >> /opt/jboss/wildfly/bin/standalone.conf

# Start JBoss/Wildfly with FIPS configuration
CMD ["/opt/jboss/wildfly/bin/standalone.sh", "-b", "0.0.0.0", "-bmanagement", "0.0.0.0"]