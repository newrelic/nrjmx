namespace java org.newrelic.nrjmx.v2.nrprotocol

struct JMXConfig {
  1: string connectionURL
  2: string hostname,
  3: i32 port,
  4: optional string uriPath,
  5: string username,
  6: string password,
  7: string keyStore,
  8: string keyStorePassword,
  9: string trustStore,
  10: string trustStorePassword,
  11: bool isRemote,
  12: bool isJBossStandaloneMode
  13: bool useSSL
  14: i64 requestTimeoutMs
  15: bool verbose
}

enum ResponseType {
  STRING = 1,
  DOUBLE = 2,
  INT    = 3,
  BOOL   = 4,
  ERROR  = 5,
}

struct AttributeResponse {
  1: string statusMsg,
  2: string name,
  3: ResponseType responseType,
  4: string stringValue,
  5: double doubleValue,
  6: i64 intValue,
  7: bool boolValue
}

exception JMXError {
  1: string message,
  2: string causeMessage
  3: string stacktrace
}

exception JMXConnectionError {
  1: string message
}

service JMXService {
    void connect(1:JMXConfig config) throws (1:JMXConnectionError connErr, 2:JMXError jmxErr),

    void disconnect() throws (1:JMXError err),

    string getClientVersion() throws (1:JMXError err),

    list<string> queryMBeanNames(1:string mBeanNamePattern) throws (1:JMXConnectionError connErr, 2:JMXError jmxErr),

    list<string> getMBeanAttributeNames(1:string mBeanName) throws (1:JMXConnectionError connErr, 2:JMXError jmxErr),

    list<AttributeResponse> getMBeanAttributes(1:string mBeanName, 2:list<string> attributes) throws (1:JMXConnectionError connErr, 2:JMXError jmxErr),

    list<AttributeResponse> queryMBeanAttributes(1:string mBeanNamePattern) throws (1:JMXConnectionError connErr, 2:JMXError jmxErr)
}