package org.newrelic.jmx

import javax.management.MBeanRegistration
import javax.management.MBeanServer
import javax.management.ObjectName

interface CatMBean {
    val name: String
}

class Cat(override val name: String) : CatMBean, MBeanRegistration {
    override fun preDeregister() {}
    override fun postRegister(registrationDone: Boolean?) {}
    override fun postDeregister() {}

    override fun preRegister(server: MBeanServer?, oname: ObjectName?): ObjectName =
            ObjectName("test:type=Cat,name=$name")

    override fun toString(): String {
        return "Cat{name=$name}"
    }
}
