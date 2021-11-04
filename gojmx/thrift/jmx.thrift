/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

# Thrift Tutorial
# Mark Slee (mcslee@facebook.com)
#
# This file aims to teach you how to use Thrift, in a .thrift file. Neato. The
# first thing to notice is that .thrift files support standard shell comments.
# This lets you make your thrift file executable and include your Thrift build
# step on the top line. And you can place comments like this anywhere you like.
#
# Before running this file, you will need to have installed the thrift compiler
# into /usr/local/bin.

/**
 * The first thing to know about are types. The available types in Thrift are:
 *
 *  bool        Boolean, one byte
 *  i8 (byte)   Signed 8-bit integer
 *  i16         Signed 16-bit integer
 *  i32         Signed 32-bit integer
 *  i64         Signed 64-bit integer
 *  double      64-bit floating point value
 *  string      String
 *  binary      Blob (byte array)
 *  map<t1,t2>  Map from one type to another
 *  list<t1>    Ordered list of one type
 *  set<t1>     Set of unique elements of one type
 *
 * Did you also notice that Thrift supports C style comments?
 */

// Just in case you were wondering... yes. We support simple C comments too.

/**
 * Thrift files can reference other Thrift files to include common struct
 * and service definitions. These are found using the current path, or by
 * searching relative to any paths specified with the -I compiler flag.
 *
 * Included objects are accessed using the name of the .thrift file as a
 * prefix. i.e. shared.SharedObject
 */
// include "shared.thrift"

/**
 * Thrift files can namespace, package, or prefix their output in various
 * target languages.
 */

// namespace cl tutorial
// namespace cpp tutorial
// namespace d tutorial
// namespace dart tutorial
namespace java org.newrelic.nrjmx.v2.jmx
// namespace php tutorial
// namespace perl tutorial
// namespace haxe tutorial
// namespace netcore tutorial
// namespace netstd tutorial

// /**
//  * Thrift lets you do typedefs to get pretty names for your types. Standard
//  * C style here.
//  */
// typedef i32 MyInteger

// /**
//  * Thrift also lets you define constants for use across languages. Complex
//  * types and structs are specified using JSON notation.
//  */
// const i32 INT32CONSTANT = 9853
// const map<string,string> MAPCONSTANT = {'hello':'world', 'goodnight':'moon'}

// /**
//  * You can define enums, which are just 32 bit integers. Values are optional
//  * and start at 1 if not supplied, C style again.
//  */
// enum Operation {
//   ADD = 1,
//   SUBTRACT = 2,
//   MULTIPLY = 3,
//   DIVIDE = 4
// }

// /**
//  * Structs are the basic complex data structures. They are comprised of fields
//  * which each have an integer identifier, a type, a symbolic name, and an
//  * optional default value.
//  *
//  * Fields can be declared "optional", which ensures they will not be included
//  * in the serialized output if they aren't set.  Note that this requires some
//  * manual management in some languages.
//  */
// struct Work {
//   1: i32 num1 = 0,
//   2: i32 num2,
//   3: Operation op,
//   4: optional string comment,
// }

// /**
//  * Structs can also be exceptions, if they are nasty.
//  */
// exception InvalidOperation {
//   1: i32 whatOp,
//   2: string why
// }

// /**
//  * Ahh, now onto the cool part, defining a service. Services just need a name
//  * and can optionally inherit from another service using the extends keyword.
//  */
// service Calculator extends shared.SharedService {

//   /**
//    * A method definition looks like C code. It has a return type, arguments,
//    * and optionally a list of exceptions that it may throw. Note that argument
//    * lists and exception lists are specified using the exact same syntax as
//    * field lists in struct or exception definitions.
//    */

//    void ping(),

//    i32 add(1:i32 num1, 2:i32 num2),

//    i32 calculate(1:i32 logid, 2:Work w) throws (1:InvalidOperation ouch),

//    /**
//     * This method has a oneway modifier. That means the client only makes
//     * a request and does not listen for any response at all. Oneway methods
//     * must be void.
//     */
//    oneway void zip()

// }

/**
 * That just about covers the basics. Take a look in the test/ folder for more
 * detailed examples. After you run this file, your generated code shows up
 * in folders with names gen-<language>. The generated code isn't too scary
 * to look at. It even has pretty indentation.
 */

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
  3: double doubleValue,
  4: i64 intValue,
  5: bool boolValue
}

struct JMXAttribute {
  1: string attribute
  2: JMXAttributeValue value
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

  /**
   * A method definition looks like C code. It has a return type, arguments,
   * and optionally a list of exceptions that it may throw. Note that argument
   * lists and exception lists are specified using the exact same syntax as
   * field lists in struct or exception definitions.
   */

    bool connect(1:JMXConfig config) throws (1:JMXConnectionError connErr, 2:JMXError jmxErr),

    void disconnect() throws (1:JMXError err),

    list<JMXAttribute> queryMbean(1:string beanName) throws (1:JMXConnectionError connErr, 2:JMXError jmxErr),

    list<LogMessage> getLogs()
}