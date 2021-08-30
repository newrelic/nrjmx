namespace java org.newrelic.nrjmx.v2.jmx


 struct JMXConfig {
  1: string connURL
  2: string hostname,
  3: i32 port,
  4: string uriPath,
  5: string username,
  6: string password,
  7: string keyStore,
  8: string keyStorePassword,
  9: string trustStore,
  10: string trustStorePassword,
  11: bool isRemote,
  12: bool isJBossStandaloneMode
}

enum ValueType {
  STRING = 1,
  DOUBLE = 2,
  INT    = 3,
  BOOL   = 4
}

struct JMXAttributeValue {
  1: ValueType valueType,
  2: string stringValue,
  3: optional double doubleValue,
  4: optional i64 intValue,
  5: optional bool boolValue
}

struct JMXAttribute {
  1: string attribute
  2: JMXAttributeValue value
}

exception JMXError {
  1: optional i32 code,
  2: string message
}

exception JMXConnectionError {
  1: i32 code,
  2: string message
}

enum JMXLoggerMessageLevel {
  DEBUG   = 1,
  INFO    = 2,
  WARNING = 3,
  ERROR   = 4
}

struct JMXLoggerMessage {
    1: string message
    2: JMXLoggerMessageLevel level
}

service JMXService {

    bool connect(1:JMXConfig config) throws (1:JMXConnectionError connErr, 2:JMXError jmxErr),

    void disconnect() throws (1:JMXError err),

    list<JMXAttribute> queryMbean(1:string beanName) throws (1:JMXConnectionError connErr, 2:JMXError jmxErr),

    list<JMXLoggerMessage> getLogs()
}