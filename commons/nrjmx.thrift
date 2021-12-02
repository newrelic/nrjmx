namespace java org.newrelic.nrjmx.v2.nrprotocol

struct JMXConfig {
  1: string connectionURL
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
  BOOL   = 4,
}

struct JMXAttribute {
  1: string attribute
  2: ValueType valueType,
  3: string stringValue,
  4: double doubleValue,
  5: i64 intValue,
  6: bool boolValue
}

struct LogMessage {
  1: string message
}

exception JMXError {
  1: optional i32 code,
  2: string message
}

exception JMXConnectionError {
  1: i32 code,
  2: string message
}

service JMXService {
    void connect(1:JMXConfig config) throws (1:JMXConnectionError connErr, 2:JMXError jmxErr),

    void disconnect() throws (1:JMXError err),

    list<JMXAttribute> queryMbean(1:string beanName) throws (1:JMXConnectionError connErr, 2:JMXError jmxErr),
}