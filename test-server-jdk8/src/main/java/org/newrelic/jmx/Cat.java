package org.newrelic.jmx;

import javax.management.MBeanRegistration;
import javax.management.MBeanServer;
import javax.management.ObjectName;

public class Cat implements CatMBean, MBeanRegistration {

  private String name;

  public Cat(String name) {
    this.name = name;
  }

  @Override
  public String getName() {
    return name;
  }

  public String toJSON() {
    return "{\"name\":\"" + name + "\"}";
  }

  @Override
  public String toString() {
    return toJSON();
  }

  @Override
  public ObjectName preRegister(MBeanServer server, ObjectName name) throws Exception {
    return new ObjectName("test:type=Cat,name=" + this.name);
  }

  @Override
  public void postRegister(Boolean registrationDone) {}

  @Override
  public void preDeregister() {}

  @Override
  public void postDeregister() {}
}
