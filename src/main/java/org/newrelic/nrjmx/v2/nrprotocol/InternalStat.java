/**
 * Autogenerated by Thrift Compiler (0.21.0)
 *
 * DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING
 *  @generated
 */
package org.newrelic.nrjmx.v2.nrprotocol;

@SuppressWarnings({"cast", "rawtypes", "serial", "unchecked", "unused"})
public class InternalStat implements org.apache.thrift.TBase<InternalStat, InternalStat._Fields>, java.io.Serializable, Cloneable, Comparable<InternalStat> {
  private static final org.apache.thrift.protocol.TStruct STRUCT_DESC = new org.apache.thrift.protocol.TStruct("InternalStat");

  private static final org.apache.thrift.protocol.TField STAT_TYPE_FIELD_DESC = new org.apache.thrift.protocol.TField("statType", org.apache.thrift.protocol.TType.STRING, (short)1);
  private static final org.apache.thrift.protocol.TField M_BEAN_FIELD_DESC = new org.apache.thrift.protocol.TField("mBean", org.apache.thrift.protocol.TType.STRING, (short)2);
  private static final org.apache.thrift.protocol.TField ATTRS_FIELD_DESC = new org.apache.thrift.protocol.TField("attrs", org.apache.thrift.protocol.TType.LIST, (short)3);
  private static final org.apache.thrift.protocol.TField RESPONSE_COUNT_FIELD_DESC = new org.apache.thrift.protocol.TField("responseCount", org.apache.thrift.protocol.TType.I64, (short)4);
  private static final org.apache.thrift.protocol.TField MILLISECONDS_FIELD_DESC = new org.apache.thrift.protocol.TField("milliseconds", org.apache.thrift.protocol.TType.DOUBLE, (short)5);
  private static final org.apache.thrift.protocol.TField START_TIMESTAMP_FIELD_DESC = new org.apache.thrift.protocol.TField("startTimestamp", org.apache.thrift.protocol.TType.I64, (short)6);
  private static final org.apache.thrift.protocol.TField SUCCESSFUL_FIELD_DESC = new org.apache.thrift.protocol.TField("successful", org.apache.thrift.protocol.TType.BOOL, (short)7);

  private static final org.apache.thrift.scheme.SchemeFactory STANDARD_SCHEME_FACTORY = new InternalStatStandardSchemeFactory();
  private static final org.apache.thrift.scheme.SchemeFactory TUPLE_SCHEME_FACTORY = new InternalStatTupleSchemeFactory();

  public @org.apache.thrift.annotation.Nullable java.lang.String statType; // required
  public @org.apache.thrift.annotation.Nullable java.lang.String mBean; // required
  public @org.apache.thrift.annotation.Nullable java.util.List<java.lang.String> attrs; // required
  public long responseCount; // required
  public double milliseconds; // required
  public long startTimestamp; // required
  public boolean successful; // required

  /** The set of fields this struct contains, along with convenience methods for finding and manipulating them. */
  public enum _Fields implements org.apache.thrift.TFieldIdEnum {
    STAT_TYPE((short)1, "statType"),
    M_BEAN((short)2, "mBean"),
    ATTRS((short)3, "attrs"),
    RESPONSE_COUNT((short)4, "responseCount"),
    MILLISECONDS((short)5, "milliseconds"),
    START_TIMESTAMP((short)6, "startTimestamp"),
    SUCCESSFUL((short)7, "successful");

    private static final java.util.Map<java.lang.String, _Fields> byName = new java.util.HashMap<java.lang.String, _Fields>();

    static {
      for (_Fields field : java.util.EnumSet.allOf(_Fields.class)) {
        byName.put(field.getFieldName(), field);
      }
    }

    /**
     * Find the _Fields constant that matches fieldId, or null if its not found.
     */
    @org.apache.thrift.annotation.Nullable
    public static _Fields findByThriftId(int fieldId) {
      switch(fieldId) {
        case 1: // STAT_TYPE
          return STAT_TYPE;
        case 2: // M_BEAN
          return M_BEAN;
        case 3: // ATTRS
          return ATTRS;
        case 4: // RESPONSE_COUNT
          return RESPONSE_COUNT;
        case 5: // MILLISECONDS
          return MILLISECONDS;
        case 6: // START_TIMESTAMP
          return START_TIMESTAMP;
        case 7: // SUCCESSFUL
          return SUCCESSFUL;
        default:
          return null;
      }
    }

    /**
     * Find the _Fields constant that matches fieldId, throwing an exception
     * if it is not found.
     */
    public static _Fields findByThriftIdOrThrow(int fieldId) {
      _Fields fields = findByThriftId(fieldId);
      if (fields == null) throw new java.lang.IllegalArgumentException("Field " + fieldId + " doesn't exist!");
      return fields;
    }

    /**
     * Find the _Fields constant that matches name, or null if its not found.
     */
    @org.apache.thrift.annotation.Nullable
    public static _Fields findByName(java.lang.String name) {
      return byName.get(name);
    }

    private final short _thriftId;
    private final java.lang.String _fieldName;

    _Fields(short thriftId, java.lang.String fieldName) {
      _thriftId = thriftId;
      _fieldName = fieldName;
    }

    @Override
    public short getThriftFieldId() {
      return _thriftId;
    }

    @Override
    public java.lang.String getFieldName() {
      return _fieldName;
    }
  }

  // isset id assignments
  private static final int __RESPONSECOUNT_ISSET_ID = 0;
  private static final int __MILLISECONDS_ISSET_ID = 1;
  private static final int __STARTTIMESTAMP_ISSET_ID = 2;
  private static final int __SUCCESSFUL_ISSET_ID = 3;
  private byte __isset_bitfield = 0;
  public static final java.util.Map<_Fields, org.apache.thrift.meta_data.FieldMetaData> metaDataMap;
  static {
    java.util.Map<_Fields, org.apache.thrift.meta_data.FieldMetaData> tmpMap = new java.util.EnumMap<_Fields, org.apache.thrift.meta_data.FieldMetaData>(_Fields.class);
    tmpMap.put(_Fields.STAT_TYPE, new org.apache.thrift.meta_data.FieldMetaData("statType", org.apache.thrift.TFieldRequirementType.DEFAULT, 
        new org.apache.thrift.meta_data.FieldValueMetaData(org.apache.thrift.protocol.TType.STRING)));
    tmpMap.put(_Fields.M_BEAN, new org.apache.thrift.meta_data.FieldMetaData("mBean", org.apache.thrift.TFieldRequirementType.DEFAULT, 
        new org.apache.thrift.meta_data.FieldValueMetaData(org.apache.thrift.protocol.TType.STRING)));
    tmpMap.put(_Fields.ATTRS, new org.apache.thrift.meta_data.FieldMetaData("attrs", org.apache.thrift.TFieldRequirementType.DEFAULT, 
        new org.apache.thrift.meta_data.ListMetaData(org.apache.thrift.protocol.TType.LIST, 
            new org.apache.thrift.meta_data.FieldValueMetaData(org.apache.thrift.protocol.TType.STRING))));
    tmpMap.put(_Fields.RESPONSE_COUNT, new org.apache.thrift.meta_data.FieldMetaData("responseCount", org.apache.thrift.TFieldRequirementType.DEFAULT, 
        new org.apache.thrift.meta_data.FieldValueMetaData(org.apache.thrift.protocol.TType.I64)));
    tmpMap.put(_Fields.MILLISECONDS, new org.apache.thrift.meta_data.FieldMetaData("milliseconds", org.apache.thrift.TFieldRequirementType.DEFAULT, 
        new org.apache.thrift.meta_data.FieldValueMetaData(org.apache.thrift.protocol.TType.DOUBLE)));
    tmpMap.put(_Fields.START_TIMESTAMP, new org.apache.thrift.meta_data.FieldMetaData("startTimestamp", org.apache.thrift.TFieldRequirementType.DEFAULT, 
        new org.apache.thrift.meta_data.FieldValueMetaData(org.apache.thrift.protocol.TType.I64)));
    tmpMap.put(_Fields.SUCCESSFUL, new org.apache.thrift.meta_data.FieldMetaData("successful", org.apache.thrift.TFieldRequirementType.DEFAULT, 
        new org.apache.thrift.meta_data.FieldValueMetaData(org.apache.thrift.protocol.TType.BOOL)));
    metaDataMap = java.util.Collections.unmodifiableMap(tmpMap);
    org.apache.thrift.meta_data.FieldMetaData.addStructMetaDataMap(InternalStat.class, metaDataMap);
  }

  public InternalStat() {
  }

  public InternalStat(
    java.lang.String statType,
    java.lang.String mBean,
    java.util.List<java.lang.String> attrs,
    long responseCount,
    double milliseconds,
    long startTimestamp,
    boolean successful)
  {
    this();
    this.statType = statType;
    this.mBean = mBean;
    this.attrs = attrs;
    this.responseCount = responseCount;
    setResponseCountIsSet(true);
    this.milliseconds = milliseconds;
    setMillisecondsIsSet(true);
    this.startTimestamp = startTimestamp;
    setStartTimestampIsSet(true);
    this.successful = successful;
    setSuccessfulIsSet(true);
  }

  /**
   * Performs a deep copy on <i>other</i>.
   */
  public InternalStat(InternalStat other) {
    __isset_bitfield = other.__isset_bitfield;
    if (other.isSetStatType()) {
      this.statType = other.statType;
    }
    if (other.isSetMBean()) {
      this.mBean = other.mBean;
    }
    if (other.isSetAttrs()) {
      java.util.List<java.lang.String> __this__attrs = new java.util.ArrayList<java.lang.String>(other.attrs);
      this.attrs = __this__attrs;
    }
    this.responseCount = other.responseCount;
    this.milliseconds = other.milliseconds;
    this.startTimestamp = other.startTimestamp;
    this.successful = other.successful;
  }

  @Override
  public InternalStat deepCopy() {
    return new InternalStat(this);
  }

  @Override
  public void clear() {
    this.statType = null;
    this.mBean = null;
    this.attrs = null;
    setResponseCountIsSet(false);
    this.responseCount = 0;
    setMillisecondsIsSet(false);
    this.milliseconds = 0.0;
    setStartTimestampIsSet(false);
    this.startTimestamp = 0;
    setSuccessfulIsSet(false);
    this.successful = false;
  }

  @org.apache.thrift.annotation.Nullable
  public java.lang.String getStatType() {
    return this.statType;
  }

  public InternalStat setStatType(@org.apache.thrift.annotation.Nullable java.lang.String statType) {
    this.statType = statType;
    return this;
  }

  public void unsetStatType() {
    this.statType = null;
  }

  /** Returns true if field statType is set (has been assigned a value) and false otherwise */
  public boolean isSetStatType() {
    return this.statType != null;
  }

  public void setStatTypeIsSet(boolean value) {
    if (!value) {
      this.statType = null;
    }
  }

  @org.apache.thrift.annotation.Nullable
  public java.lang.String getMBean() {
    return this.mBean;
  }

  public InternalStat setMBean(@org.apache.thrift.annotation.Nullable java.lang.String mBean) {
    this.mBean = mBean;
    return this;
  }

  public void unsetMBean() {
    this.mBean = null;
  }

  /** Returns true if field mBean is set (has been assigned a value) and false otherwise */
  public boolean isSetMBean() {
    return this.mBean != null;
  }

  public void setMBeanIsSet(boolean value) {
    if (!value) {
      this.mBean = null;
    }
  }

  public int getAttrsSize() {
    return (this.attrs == null) ? 0 : this.attrs.size();
  }

  @org.apache.thrift.annotation.Nullable
  public java.util.Iterator<java.lang.String> getAttrsIterator() {
    return (this.attrs == null) ? null : this.attrs.iterator();
  }

  public void addToAttrs(java.lang.String elem) {
    if (this.attrs == null) {
      this.attrs = new java.util.ArrayList<java.lang.String>();
    }
    this.attrs.add(elem);
  }

  @org.apache.thrift.annotation.Nullable
  public java.util.List<java.lang.String> getAttrs() {
    return this.attrs;
  }

  public InternalStat setAttrs(@org.apache.thrift.annotation.Nullable java.util.List<java.lang.String> attrs) {
    this.attrs = attrs;
    return this;
  }

  public void unsetAttrs() {
    this.attrs = null;
  }

  /** Returns true if field attrs is set (has been assigned a value) and false otherwise */
  public boolean isSetAttrs() {
    return this.attrs != null;
  }

  public void setAttrsIsSet(boolean value) {
    if (!value) {
      this.attrs = null;
    }
  }

  public long getResponseCount() {
    return this.responseCount;
  }

  public InternalStat setResponseCount(long responseCount) {
    this.responseCount = responseCount;
    setResponseCountIsSet(true);
    return this;
  }

  public void unsetResponseCount() {
    __isset_bitfield = org.apache.thrift.EncodingUtils.clearBit(__isset_bitfield, __RESPONSECOUNT_ISSET_ID);
  }

  /** Returns true if field responseCount is set (has been assigned a value) and false otherwise */
  public boolean isSetResponseCount() {
    return org.apache.thrift.EncodingUtils.testBit(__isset_bitfield, __RESPONSECOUNT_ISSET_ID);
  }

  public void setResponseCountIsSet(boolean value) {
    __isset_bitfield = org.apache.thrift.EncodingUtils.setBit(__isset_bitfield, __RESPONSECOUNT_ISSET_ID, value);
  }

  public double getMilliseconds() {
    return this.milliseconds;
  }

  public InternalStat setMilliseconds(double milliseconds) {
    this.milliseconds = milliseconds;
    setMillisecondsIsSet(true);
    return this;
  }

  public void unsetMilliseconds() {
    __isset_bitfield = org.apache.thrift.EncodingUtils.clearBit(__isset_bitfield, __MILLISECONDS_ISSET_ID);
  }

  /** Returns true if field milliseconds is set (has been assigned a value) and false otherwise */
  public boolean isSetMilliseconds() {
    return org.apache.thrift.EncodingUtils.testBit(__isset_bitfield, __MILLISECONDS_ISSET_ID);
  }

  public void setMillisecondsIsSet(boolean value) {
    __isset_bitfield = org.apache.thrift.EncodingUtils.setBit(__isset_bitfield, __MILLISECONDS_ISSET_ID, value);
  }

  public long getStartTimestamp() {
    return this.startTimestamp;
  }

  public InternalStat setStartTimestamp(long startTimestamp) {
    this.startTimestamp = startTimestamp;
    setStartTimestampIsSet(true);
    return this;
  }

  public void unsetStartTimestamp() {
    __isset_bitfield = org.apache.thrift.EncodingUtils.clearBit(__isset_bitfield, __STARTTIMESTAMP_ISSET_ID);
  }

  /** Returns true if field startTimestamp is set (has been assigned a value) and false otherwise */
  public boolean isSetStartTimestamp() {
    return org.apache.thrift.EncodingUtils.testBit(__isset_bitfield, __STARTTIMESTAMP_ISSET_ID);
  }

  public void setStartTimestampIsSet(boolean value) {
    __isset_bitfield = org.apache.thrift.EncodingUtils.setBit(__isset_bitfield, __STARTTIMESTAMP_ISSET_ID, value);
  }

  public boolean isSuccessful() {
    return this.successful;
  }

  public InternalStat setSuccessful(boolean successful) {
    this.successful = successful;
    setSuccessfulIsSet(true);
    return this;
  }

  public void unsetSuccessful() {
    __isset_bitfield = org.apache.thrift.EncodingUtils.clearBit(__isset_bitfield, __SUCCESSFUL_ISSET_ID);
  }

  /** Returns true if field successful is set (has been assigned a value) and false otherwise */
  public boolean isSetSuccessful() {
    return org.apache.thrift.EncodingUtils.testBit(__isset_bitfield, __SUCCESSFUL_ISSET_ID);
  }

  public void setSuccessfulIsSet(boolean value) {
    __isset_bitfield = org.apache.thrift.EncodingUtils.setBit(__isset_bitfield, __SUCCESSFUL_ISSET_ID, value);
  }

  @Override
  public void setFieldValue(_Fields field, @org.apache.thrift.annotation.Nullable java.lang.Object value) {
    switch (field) {
    case STAT_TYPE:
      if (value == null) {
        unsetStatType();
      } else {
        setStatType((java.lang.String)value);
      }
      break;

    case M_BEAN:
      if (value == null) {
        unsetMBean();
      } else {
        setMBean((java.lang.String)value);
      }
      break;

    case ATTRS:
      if (value == null) {
        unsetAttrs();
      } else {
        setAttrs((java.util.List<java.lang.String>)value);
      }
      break;

    case RESPONSE_COUNT:
      if (value == null) {
        unsetResponseCount();
      } else {
        setResponseCount((java.lang.Long)value);
      }
      break;

    case MILLISECONDS:
      if (value == null) {
        unsetMilliseconds();
      } else {
        setMilliseconds((java.lang.Double)value);
      }
      break;

    case START_TIMESTAMP:
      if (value == null) {
        unsetStartTimestamp();
      } else {
        setStartTimestamp((java.lang.Long)value);
      }
      break;

    case SUCCESSFUL:
      if (value == null) {
        unsetSuccessful();
      } else {
        setSuccessful((java.lang.Boolean)value);
      }
      break;

    }
  }

  @org.apache.thrift.annotation.Nullable
  @Override
  public java.lang.Object getFieldValue(_Fields field) {
    switch (field) {
    case STAT_TYPE:
      return getStatType();

    case M_BEAN:
      return getMBean();

    case ATTRS:
      return getAttrs();

    case RESPONSE_COUNT:
      return getResponseCount();

    case MILLISECONDS:
      return getMilliseconds();

    case START_TIMESTAMP:
      return getStartTimestamp();

    case SUCCESSFUL:
      return isSuccessful();

    }
    throw new java.lang.IllegalStateException();
  }

  /** Returns true if field corresponding to fieldID is set (has been assigned a value) and false otherwise */
  @Override
  public boolean isSet(_Fields field) {
    if (field == null) {
      throw new java.lang.IllegalArgumentException();
    }

    switch (field) {
    case STAT_TYPE:
      return isSetStatType();
    case M_BEAN:
      return isSetMBean();
    case ATTRS:
      return isSetAttrs();
    case RESPONSE_COUNT:
      return isSetResponseCount();
    case MILLISECONDS:
      return isSetMilliseconds();
    case START_TIMESTAMP:
      return isSetStartTimestamp();
    case SUCCESSFUL:
      return isSetSuccessful();
    }
    throw new java.lang.IllegalStateException();
  }

  @Override
  public boolean equals(java.lang.Object that) {
    if (that instanceof InternalStat)
      return this.equals((InternalStat)that);
    return false;
  }

  public boolean equals(InternalStat that) {
    if (that == null)
      return false;
    if (this == that)
      return true;

    boolean this_present_statType = true && this.isSetStatType();
    boolean that_present_statType = true && that.isSetStatType();
    if (this_present_statType || that_present_statType) {
      if (!(this_present_statType && that_present_statType))
        return false;
      if (!this.statType.equals(that.statType))
        return false;
    }

    boolean this_present_mBean = true && this.isSetMBean();
    boolean that_present_mBean = true && that.isSetMBean();
    if (this_present_mBean || that_present_mBean) {
      if (!(this_present_mBean && that_present_mBean))
        return false;
      if (!this.mBean.equals(that.mBean))
        return false;
    }

    boolean this_present_attrs = true && this.isSetAttrs();
    boolean that_present_attrs = true && that.isSetAttrs();
    if (this_present_attrs || that_present_attrs) {
      if (!(this_present_attrs && that_present_attrs))
        return false;
      if (!this.attrs.equals(that.attrs))
        return false;
    }

    boolean this_present_responseCount = true;
    boolean that_present_responseCount = true;
    if (this_present_responseCount || that_present_responseCount) {
      if (!(this_present_responseCount && that_present_responseCount))
        return false;
      if (this.responseCount != that.responseCount)
        return false;
    }

    boolean this_present_milliseconds = true;
    boolean that_present_milliseconds = true;
    if (this_present_milliseconds || that_present_milliseconds) {
      if (!(this_present_milliseconds && that_present_milliseconds))
        return false;
      if (this.milliseconds != that.milliseconds)
        return false;
    }

    boolean this_present_startTimestamp = true;
    boolean that_present_startTimestamp = true;
    if (this_present_startTimestamp || that_present_startTimestamp) {
      if (!(this_present_startTimestamp && that_present_startTimestamp))
        return false;
      if (this.startTimestamp != that.startTimestamp)
        return false;
    }

    boolean this_present_successful = true;
    boolean that_present_successful = true;
    if (this_present_successful || that_present_successful) {
      if (!(this_present_successful && that_present_successful))
        return false;
      if (this.successful != that.successful)
        return false;
    }

    return true;
  }

  @Override
  public int hashCode() {
    int hashCode = 1;

    hashCode = hashCode * 8191 + ((isSetStatType()) ? 131071 : 524287);
    if (isSetStatType())
      hashCode = hashCode * 8191 + statType.hashCode();

    hashCode = hashCode * 8191 + ((isSetMBean()) ? 131071 : 524287);
    if (isSetMBean())
      hashCode = hashCode * 8191 + mBean.hashCode();

    hashCode = hashCode * 8191 + ((isSetAttrs()) ? 131071 : 524287);
    if (isSetAttrs())
      hashCode = hashCode * 8191 + attrs.hashCode();

    hashCode = hashCode * 8191 + org.apache.thrift.TBaseHelper.hashCode(responseCount);

    hashCode = hashCode * 8191 + org.apache.thrift.TBaseHelper.hashCode(milliseconds);

    hashCode = hashCode * 8191 + org.apache.thrift.TBaseHelper.hashCode(startTimestamp);

    hashCode = hashCode * 8191 + ((successful) ? 131071 : 524287);

    return hashCode;
  }

  @Override
  public int compareTo(InternalStat other) {
    if (!getClass().equals(other.getClass())) {
      return getClass().getName().compareTo(other.getClass().getName());
    }

    int lastComparison = 0;

    lastComparison = java.lang.Boolean.compare(isSetStatType(), other.isSetStatType());
    if (lastComparison != 0) {
      return lastComparison;
    }
    if (isSetStatType()) {
      lastComparison = org.apache.thrift.TBaseHelper.compareTo(this.statType, other.statType);
      if (lastComparison != 0) {
        return lastComparison;
      }
    }
    lastComparison = java.lang.Boolean.compare(isSetMBean(), other.isSetMBean());
    if (lastComparison != 0) {
      return lastComparison;
    }
    if (isSetMBean()) {
      lastComparison = org.apache.thrift.TBaseHelper.compareTo(this.mBean, other.mBean);
      if (lastComparison != 0) {
        return lastComparison;
      }
    }
    lastComparison = java.lang.Boolean.compare(isSetAttrs(), other.isSetAttrs());
    if (lastComparison != 0) {
      return lastComparison;
    }
    if (isSetAttrs()) {
      lastComparison = org.apache.thrift.TBaseHelper.compareTo(this.attrs, other.attrs);
      if (lastComparison != 0) {
        return lastComparison;
      }
    }
    lastComparison = java.lang.Boolean.compare(isSetResponseCount(), other.isSetResponseCount());
    if (lastComparison != 0) {
      return lastComparison;
    }
    if (isSetResponseCount()) {
      lastComparison = org.apache.thrift.TBaseHelper.compareTo(this.responseCount, other.responseCount);
      if (lastComparison != 0) {
        return lastComparison;
      }
    }
    lastComparison = java.lang.Boolean.compare(isSetMilliseconds(), other.isSetMilliseconds());
    if (lastComparison != 0) {
      return lastComparison;
    }
    if (isSetMilliseconds()) {
      lastComparison = org.apache.thrift.TBaseHelper.compareTo(this.milliseconds, other.milliseconds);
      if (lastComparison != 0) {
        return lastComparison;
      }
    }
    lastComparison = java.lang.Boolean.compare(isSetStartTimestamp(), other.isSetStartTimestamp());
    if (lastComparison != 0) {
      return lastComparison;
    }
    if (isSetStartTimestamp()) {
      lastComparison = org.apache.thrift.TBaseHelper.compareTo(this.startTimestamp, other.startTimestamp);
      if (lastComparison != 0) {
        return lastComparison;
      }
    }
    lastComparison = java.lang.Boolean.compare(isSetSuccessful(), other.isSetSuccessful());
    if (lastComparison != 0) {
      return lastComparison;
    }
    if (isSetSuccessful()) {
      lastComparison = org.apache.thrift.TBaseHelper.compareTo(this.successful, other.successful);
      if (lastComparison != 0) {
        return lastComparison;
      }
    }
    return 0;
  }

  @org.apache.thrift.annotation.Nullable
  @Override
  public _Fields fieldForId(int fieldId) {
    return _Fields.findByThriftId(fieldId);
  }

  @Override
  public void read(org.apache.thrift.protocol.TProtocol iprot) throws org.apache.thrift.TException {
    scheme(iprot).read(iprot, this);
  }

  @Override
  public void write(org.apache.thrift.protocol.TProtocol oprot) throws org.apache.thrift.TException {
    scheme(oprot).write(oprot, this);
  }

  @Override
  public java.lang.String toString() {
    java.lang.StringBuilder sb = new java.lang.StringBuilder("InternalStat(");
    boolean first = true;

    sb.append("statType:");
    if (this.statType == null) {
      sb.append("null");
    } else {
      sb.append(this.statType);
    }
    first = false;
    if (!first) sb.append(", ");
    sb.append("mBean:");
    if (this.mBean == null) {
      sb.append("null");
    } else {
      sb.append(this.mBean);
    }
    first = false;
    if (!first) sb.append(", ");
    sb.append("attrs:");
    if (this.attrs == null) {
      sb.append("null");
    } else {
      sb.append(this.attrs);
    }
    first = false;
    if (!first) sb.append(", ");
    sb.append("responseCount:");
    sb.append(this.responseCount);
    first = false;
    if (!first) sb.append(", ");
    sb.append("milliseconds:");
    sb.append(this.milliseconds);
    first = false;
    if (!first) sb.append(", ");
    sb.append("startTimestamp:");
    sb.append(this.startTimestamp);
    first = false;
    if (!first) sb.append(", ");
    sb.append("successful:");
    sb.append(this.successful);
    first = false;
    sb.append(")");
    return sb.toString();
  }

  public void validate() throws org.apache.thrift.TException {
    // check for required fields
    // check for sub-struct validity
  }

  private void writeObject(java.io.ObjectOutputStream out) throws java.io.IOException {
    try {
      write(new org.apache.thrift.protocol.TCompactProtocol(new org.apache.thrift.transport.TIOStreamTransport(out)));
    } catch (org.apache.thrift.TException te) {
      throw new java.io.IOException(te);
    }
  }

  private void readObject(java.io.ObjectInputStream in) throws java.io.IOException, java.lang.ClassNotFoundException {
    try {
      // it doesn't seem like you should have to do this, but java serialization is wacky, and doesn't call the default constructor.
      __isset_bitfield = 0;
      read(new org.apache.thrift.protocol.TCompactProtocol(new org.apache.thrift.transport.TIOStreamTransport(in)));
    } catch (org.apache.thrift.TException te) {
      throw new java.io.IOException(te);
    }
  }

  private static class InternalStatStandardSchemeFactory implements org.apache.thrift.scheme.SchemeFactory {
    @Override
    public InternalStatStandardScheme getScheme() {
      return new InternalStatStandardScheme();
    }
  }

  private static class InternalStatStandardScheme extends org.apache.thrift.scheme.StandardScheme<InternalStat> {

    @Override
    public void read(org.apache.thrift.protocol.TProtocol iprot, InternalStat struct) throws org.apache.thrift.TException {
      org.apache.thrift.protocol.TField schemeField;
      iprot.readStructBegin();
      while (true)
      {
        schemeField = iprot.readFieldBegin();
        if (schemeField.type == org.apache.thrift.protocol.TType.STOP) { 
          break;
        }
        switch (schemeField.id) {
          case 1: // STAT_TYPE
            if (schemeField.type == org.apache.thrift.protocol.TType.STRING) {
              struct.statType = iprot.readString();
              struct.setStatTypeIsSet(true);
            } else { 
              org.apache.thrift.protocol.TProtocolUtil.skip(iprot, schemeField.type);
            }
            break;
          case 2: // M_BEAN
            if (schemeField.type == org.apache.thrift.protocol.TType.STRING) {
              struct.mBean = iprot.readString();
              struct.setMBeanIsSet(true);
            } else { 
              org.apache.thrift.protocol.TProtocolUtil.skip(iprot, schemeField.type);
            }
            break;
          case 3: // ATTRS
            if (schemeField.type == org.apache.thrift.protocol.TType.LIST) {
              {
                org.apache.thrift.protocol.TList _list0 = iprot.readListBegin();
                struct.attrs = new java.util.ArrayList<java.lang.String>(_list0.size);
                @org.apache.thrift.annotation.Nullable java.lang.String _elem1;
                for (int _i2 = 0; _i2 < _list0.size; ++_i2)
                {
                  _elem1 = iprot.readString();
                  struct.attrs.add(_elem1);
                }
                iprot.readListEnd();
              }
              struct.setAttrsIsSet(true);
            } else { 
              org.apache.thrift.protocol.TProtocolUtil.skip(iprot, schemeField.type);
            }
            break;
          case 4: // RESPONSE_COUNT
            if (schemeField.type == org.apache.thrift.protocol.TType.I64) {
              struct.responseCount = iprot.readI64();
              struct.setResponseCountIsSet(true);
            } else { 
              org.apache.thrift.protocol.TProtocolUtil.skip(iprot, schemeField.type);
            }
            break;
          case 5: // MILLISECONDS
            if (schemeField.type == org.apache.thrift.protocol.TType.DOUBLE) {
              struct.milliseconds = iprot.readDouble();
              struct.setMillisecondsIsSet(true);
            } else { 
              org.apache.thrift.protocol.TProtocolUtil.skip(iprot, schemeField.type);
            }
            break;
          case 6: // START_TIMESTAMP
            if (schemeField.type == org.apache.thrift.protocol.TType.I64) {
              struct.startTimestamp = iprot.readI64();
              struct.setStartTimestampIsSet(true);
            } else { 
              org.apache.thrift.protocol.TProtocolUtil.skip(iprot, schemeField.type);
            }
            break;
          case 7: // SUCCESSFUL
            if (schemeField.type == org.apache.thrift.protocol.TType.BOOL) {
              struct.successful = iprot.readBool();
              struct.setSuccessfulIsSet(true);
            } else { 
              org.apache.thrift.protocol.TProtocolUtil.skip(iprot, schemeField.type);
            }
            break;
          default:
            org.apache.thrift.protocol.TProtocolUtil.skip(iprot, schemeField.type);
        }
        iprot.readFieldEnd();
      }
      iprot.readStructEnd();

      // check for required fields of primitive type, which can't be checked in the validate method
      struct.validate();
    }

    @Override
    public void write(org.apache.thrift.protocol.TProtocol oprot, InternalStat struct) throws org.apache.thrift.TException {
      struct.validate();

      oprot.writeStructBegin(STRUCT_DESC);
      if (struct.statType != null) {
        oprot.writeFieldBegin(STAT_TYPE_FIELD_DESC);
        oprot.writeString(struct.statType);
        oprot.writeFieldEnd();
      }
      if (struct.mBean != null) {
        oprot.writeFieldBegin(M_BEAN_FIELD_DESC);
        oprot.writeString(struct.mBean);
        oprot.writeFieldEnd();
      }
      if (struct.attrs != null) {
        oprot.writeFieldBegin(ATTRS_FIELD_DESC);
        {
          oprot.writeListBegin(new org.apache.thrift.protocol.TList(org.apache.thrift.protocol.TType.STRING, struct.attrs.size()));
          for (java.lang.String _iter3 : struct.attrs)
          {
            oprot.writeString(_iter3);
          }
          oprot.writeListEnd();
        }
        oprot.writeFieldEnd();
      }
      oprot.writeFieldBegin(RESPONSE_COUNT_FIELD_DESC);
      oprot.writeI64(struct.responseCount);
      oprot.writeFieldEnd();
      oprot.writeFieldBegin(MILLISECONDS_FIELD_DESC);
      oprot.writeDouble(struct.milliseconds);
      oprot.writeFieldEnd();
      oprot.writeFieldBegin(START_TIMESTAMP_FIELD_DESC);
      oprot.writeI64(struct.startTimestamp);
      oprot.writeFieldEnd();
      oprot.writeFieldBegin(SUCCESSFUL_FIELD_DESC);
      oprot.writeBool(struct.successful);
      oprot.writeFieldEnd();
      oprot.writeFieldStop();
      oprot.writeStructEnd();
    }

  }

  private static class InternalStatTupleSchemeFactory implements org.apache.thrift.scheme.SchemeFactory {
    @Override
    public InternalStatTupleScheme getScheme() {
      return new InternalStatTupleScheme();
    }
  }

  private static class InternalStatTupleScheme extends org.apache.thrift.scheme.TupleScheme<InternalStat> {

    @Override
    public void write(org.apache.thrift.protocol.TProtocol prot, InternalStat struct) throws org.apache.thrift.TException {
      org.apache.thrift.protocol.TTupleProtocol oprot = (org.apache.thrift.protocol.TTupleProtocol) prot;
      java.util.BitSet optionals = new java.util.BitSet();
      if (struct.isSetStatType()) {
        optionals.set(0);
      }
      if (struct.isSetMBean()) {
        optionals.set(1);
      }
      if (struct.isSetAttrs()) {
        optionals.set(2);
      }
      if (struct.isSetResponseCount()) {
        optionals.set(3);
      }
      if (struct.isSetMilliseconds()) {
        optionals.set(4);
      }
      if (struct.isSetStartTimestamp()) {
        optionals.set(5);
      }
      if (struct.isSetSuccessful()) {
        optionals.set(6);
      }
      oprot.writeBitSet(optionals, 7);
      if (struct.isSetStatType()) {
        oprot.writeString(struct.statType);
      }
      if (struct.isSetMBean()) {
        oprot.writeString(struct.mBean);
      }
      if (struct.isSetAttrs()) {
        {
          oprot.writeI32(struct.attrs.size());
          for (java.lang.String _iter4 : struct.attrs)
          {
            oprot.writeString(_iter4);
          }
        }
      }
      if (struct.isSetResponseCount()) {
        oprot.writeI64(struct.responseCount);
      }
      if (struct.isSetMilliseconds()) {
        oprot.writeDouble(struct.milliseconds);
      }
      if (struct.isSetStartTimestamp()) {
        oprot.writeI64(struct.startTimestamp);
      }
      if (struct.isSetSuccessful()) {
        oprot.writeBool(struct.successful);
      }
    }

    @Override
    public void read(org.apache.thrift.protocol.TProtocol prot, InternalStat struct) throws org.apache.thrift.TException {
      org.apache.thrift.protocol.TTupleProtocol iprot = (org.apache.thrift.protocol.TTupleProtocol) prot;
      java.util.BitSet incoming = iprot.readBitSet(7);
      if (incoming.get(0)) {
        struct.statType = iprot.readString();
        struct.setStatTypeIsSet(true);
      }
      if (incoming.get(1)) {
        struct.mBean = iprot.readString();
        struct.setMBeanIsSet(true);
      }
      if (incoming.get(2)) {
        {
          org.apache.thrift.protocol.TList _list5 = iprot.readListBegin(org.apache.thrift.protocol.TType.STRING);
          struct.attrs = new java.util.ArrayList<java.lang.String>(_list5.size);
          @org.apache.thrift.annotation.Nullable java.lang.String _elem6;
          for (int _i7 = 0; _i7 < _list5.size; ++_i7)
          {
            _elem6 = iprot.readString();
            struct.attrs.add(_elem6);
          }
        }
        struct.setAttrsIsSet(true);
      }
      if (incoming.get(3)) {
        struct.responseCount = iprot.readI64();
        struct.setResponseCountIsSet(true);
      }
      if (incoming.get(4)) {
        struct.milliseconds = iprot.readDouble();
        struct.setMillisecondsIsSet(true);
      }
      if (incoming.get(5)) {
        struct.startTimestamp = iprot.readI64();
        struct.setStartTimestampIsSet(true);
      }
      if (incoming.get(6)) {
        struct.successful = iprot.readBool();
        struct.setSuccessfulIsSet(true);
      }
    }
  }

  private static <S extends org.apache.thrift.scheme.IScheme> S scheme(org.apache.thrift.protocol.TProtocol proto) {
    return (org.apache.thrift.scheme.StandardScheme.class.equals(proto.getScheme()) ? STANDARD_SCHEME_FACTORY : TUPLE_SCHEME_FACTORY).getScheme();
  }
}

