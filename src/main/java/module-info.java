module org.newrelic.nrjmx {
  exports org.newrelic.nrjmx;

  requires com.google.code.gson;
  requires commons.cli;
  requires java.logging;
  requires java.management;
  requires java.rmi;

  // the following modules are not necessary for nrjmx, but we add them due to our jmxterm hack
  requires java.desktop;
}
