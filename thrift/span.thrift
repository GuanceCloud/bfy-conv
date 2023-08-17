struct TSpan {
  1: string agentId
  2: string applicationName
  3: i64 agentStartTime

  // identical to agentId if null
  4: binary  transactionId;
  5: string appkey

  7: i64 spanId
  8: optional i64 parentSpanId = -1

  // span event's startTimestamp
  9: i64 startTime
  10: optional i32 elapsed = 0

  11: optional string rpc

  12: i16 serviceType
  13: optional string endPoint
  14: optional string remoteAddr

  15: optional list<TAnnotation> annotations
  16: optional i16 flag = 0

  17: optional i32 err

  18: optional list<TSpanEvent> spanEventList

  19: optional string parentApplicationName
  20: optional i16 parentApplicationType
  21: optional string acceptorHost

  25: optional i32 apiId;
  26: optional TIntStringValue exceptionInfo;

  30: optional i16 applicationServiceType;
  31: optional byte loggingTransactionInfo;
  32: optional string httpPara
  33: optional string httpMethod
  34: optional string httpRequestHeader
  35: optional string httpRequestUserAgent
  36: optional string httpRequestBody
  37: optional string httpResponseBody
  38: optional i16 retcode
  39: optional string httpRequestUID
  40: optional string httpRequestTID
  41: optional string pagentId
  42: optional string apidesc
  43: optional string httpResponseHeader
  44: optional string userId
  45: optional string sessionId
  46: string appId
  47: string tenant

  50: optional i64 threadId
  51: optional string threadName
  52: optional bool hasNextCall
}

struct TSpanChunk {
  1: string agentId
  2: string applicationName
  3: i64 agentStartTime

  4: i16 serviceType

  // identical to agentId if null
  5: binary  transactionId
  6: string appkey

  8: i64 spanId

  9: optional string endPoint

  10: list<TSpanEvent> spanEventList

  11: optional i16 applicationServiceType

  12: string appId
  13: string tenant

  15: optional i64 threadId
  16: optional string threadName
  17: optional string userId
  18: optional string sessionId
  19: optional i64 startTime
}

struct TSpanEvent {
  7: optional i64 spanId
  8: i32 sequence

  9: i32 startElapsed
  10: optional i32 endElapsed = 0

  11: optional string rpc
  12: i16 serviceType
  13: optional string endPoint

  14: optional list<TAnnotation> annotations

  15: optional i32 depth = -1
  16: optional i64 nextSpanId = -1

  20: optional string destinationId

  25: optional i32 apiId;
  26: optional TIntStringValue exceptionInfo;
  27: optional string exceptionClassName;

  30: optional i32 asyncId;
  31: optional i32 nextAsyncId;
  32: optional i32 asyncSequence;
  33: string apiInfo;
  34: optional i32 lineNumber;

  40: optional TSql sql;

  41: optional i32 retcode;
  51: optional string requestHeaders;
  61: optional string requestBody;
  71: optional string responseBody;
  81: optional string url;
  91: optional string method;
  92: optional string arguments;
}

struct TAnnotation {
  1: i32 key,
  2: optional TAnnotationValue value
}

union TAnnotationValue {
  1: string stringValue
  2: bool boolValue;
  3: i32 intValue;
  4: i64 longValue;
  5: i16 shortValue
  6: double doubleValue;
  7: binary binaryValue;
  8: byte byteValue;
  9: TIntStringValue intStringValue;
  10: TIntStringStringValue intStringStringValue;
}

struct TIntStringValue {
  1: i32 intValue;
  2: optional string stringValue;
}

struct TIntStringStringValue {
  1: i32 intValue;
  2: optional string stringValue1;
  3: optional string stringValue2;
}

struct TSql {
  1: string dbhost
  2: string dbtype
  3: string db

  6: string sqlHash
  7: optional string outputs
  8: optional string bindValue

  10: string status
  11: optional string err
  12: i64 startTime
  13: i64 dur
}