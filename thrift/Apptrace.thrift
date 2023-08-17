namespace java com.szly.apptrace.thrift.dto

enum TJvmGcType {
  UNKNOWN,
  SERIAL,
  PARALLEL,
  CMS,
  G1,
  JRockitThroughput
}

struct TServiceInfo {
  1: optional string          serviceName
  2: optional list<string>    serviceLibs
}

struct TServerMetaData {
  1: optional string              serverInfo
  2: optional list<string>        vmArgs
  10: optional list<TServiceInfo>  serviceInfos
}

struct TJvmInfo {
  1:          i16         version = 0
  2: optional string      vmVersion
  3: optional TJvmGcType  gcType = TJvmGcType.UNKNOWN
}

struct TAgentInfo {
	1: string	hostname
	2: string	ip
	3: string	ports
	4: string	agentId
	5: string	applicationName
	6: i16	    serviceType
	7: i32      pid
	8: string   agentVersion;
	9: string   vmVersion;

	10: i64	    startTimestamp

	11: optional i64     endTimestamp
	12: optional i32     endStatus
	13: string   appkey;
	14: string   osName;
	15: string   osVersion;
	16: string	appId
    17: string  tenant
  	18: optional string  collectionStatus
	20: optional TServerMetaData   serverMetaData

	30: optional TJvmInfo   jvmInfo
}

struct TJvmGc {
    1: TJvmGcType   type = TJvmGcType.UNKNOWN
    2: i64          jvmMemoryHeapUsed
    3: i64          jvmMemoryHeapMax
    4: i64          jvmMemoryNonHeapUsed
    5: i64          jvmMemoryNonHeapMax
    6: i64          jvmGcOldCount
    7: i64          jvmGcOldTime
    8: optional TJvmGcDetailed    jvmGcDetailed
    9: i64          jvmMemoryNonHeapCommitted
    10: optional i64          totalPhysicalMemory
    11: optional list<TExecuteDf>          tExecuteDfs
    12: optional TExecuteIostat         tExecuteIostat
    13: optional i16         jdbcConnNum
    14: optional i32         threadNum
    15: optional i64         jvmGcOldCountNew
    16: optional i64         jvmGcOldTimeNew
}

struct TExecuteDf {
	1: string      fileSystem
	2: i64      size
	3: i64      used
	4: i64      avail
	5: i16      usage
	6: string      mountedOn
}

struct TExecuteIostat {
	1: TExecuteIostatCpu      tExecuteIostatCpu
	2: list<TExecuteIostatDevice>      tExecuteIostatDevices
}

struct TExecuteIostatCpu {
	1: string       userUsage
	2: string       niceUsage
	3: string       systemUsage
	4: string       iowaitUsage
	5: string       stealUsage
	6: string       idleUsage
}

struct TExecuteIostatDevice {
	1: string       device
	2: string       tps
	3: string       kB_read_pers
	4: string       kB_wrtn_pers
	5: string       kB_read
	6: string       kB_wrtn
}

struct TJvmGcDetailed {
    1: optional i64 jvmGcNewCount
    2: optional i64 jvmGcNewTime
    3: optional double jvmPoolCodeCacheUsage
    4: optional i64 jvmPoolCodeCacheMax
    5: optional i64 jvmPoolCodeCacheUsed
    6: optional i64 jvmPoolCodeCacheCommitted
    7: optional i64 jvmPoolCodeCacheInit
    8: optional double jvmPoolNewGenUsage
    9: optional i64 jvmPoolNewGenMax
    10: optional i64 jvmPoolNewGenUsed
    11: optional i64 jvmPoolNewGenCommitted
    12: optional i64 jvmPoolNewGenInit
    13: optional double jvmPoolOldGenUsage
    14: optional i64 jvmPoolOldGenMax
    15: optional i64 jvmPoolOldGenUsed
    16: optional i64 jvmPoolOldGenCommitted
    17: optional i64 jvmPoolOldGenInit
    18: optional double jvmPoolSurvivorSpaceUsage
    19: optional i64 jvmPoolSurvivorSpaceMax
    20: optional i64 jvmPoolSurvivorSpaceUsed
    21: optional i64 jvmPoolSurvivorSpaceCommitted
    22: optional i64 jvmPoolSurvivorSpaceInit
    23: optional double jvmPoolPermGenUsage
    24: optional i64 jvmPoolPermGenMax
    25: optional i64 jvmPoolPermGenUsed
    26: optional i64 jvmPoolPermGenCommitted
    27: optional i64 jvmPoolPermGenInit
    28: optional double jvmPoolMetaspaceUsage
    29: optional i64 jvmPoolMetaspaceMax
    30: optional i64 jvmPoolMetaspaceUsed
    31: optional i64 jvmPoolMetaspaceCommitted
    32: optional i64 jvmPoolMetaspaceInit
}

struct TCpuLoad {
  1: optional i64       jvmCpuLoad
  2: optional i64       systemCpuLoad
}

struct TTransaction {
  2: optional i64     sampledNewCount
  3: optional i64     sampledContinuationCount
  4: optional i64     unsampledNewCount
  5: optional i64     unsampledContinuationCount
}

struct TActiveTraceHistogram {
  1: i16         version = 0
  2: optional i32         histogramSchemaType
  3: optional list<i32>   activeTraceCount
}

struct TActiveTrace {
	1: optional TActiveTraceHistogram   histogram
}

struct TAgentStat {
  1: optional string      agentId
  2: optional i64         startTimestamp
  3: optional i64         timestamp
  4: optional i64         collectInterval
  10: optional TJvmGc     gc
  20: optional TCpuLoad   cpuLoad
  30: optional TTransaction   transaction
  40: optional TActiveTrace   activeTrace
  200: optional string    metadata
  210: optional i32       threadCount
}

struct TAgentStatBatch {
  1: string                   agentId
  2: i64                      startTimestamp
  3: string                   appKey
  4: string                   appId
  5: string                   tenant
  10: list<TAgentStat>        agentStats
  11: optional i16 serviceType
}

struct TAgentLifeCycle {
  1: string appkey
  2: string appId
  3: string tenant
  5: string agentId
  6: i64    startTimestamp
  7: i64    eventTimestamp
  8: string status
  9: i32    duration
  10: optional i16 serviceType
}

struct TAgentEventType {
  1: i32    code
  2: string desc
}

struct TAgentEvent {
  1: string           appkey
  5: string           agentId
  6: i64              eventTimestamp
  7: TAgentEventType  eventType
  8: optional string  eventMessage
}

struct LockedMonitorInfo {
	1: i32 stackDepth
	2: string stackTraceElement
	3: string className
	4: i32 identityHashCode
}

struct ThreadDetailStackTraceElement {
	1: string className
	2: string methodName
	3: string fileName
	4: i32 lineNumber
}

struct TAgentThreadDetail {
	1: string name
	2: string group
	3: i64 cpuTime
	4: string state
	5: optional list<LockedMonitorInfo> ownedMonitors
	6: string waitOn
	7: i32 priority
	8: i64 threadId
	9: optional list<ThreadDetailStackTraceElement> stackTraceElements
}

struct TAgentThread {
  1:  string appkey
  2:  i64 id
  4:  string applicationName
  5:  string agentId
  6:  i64 ts
  7:  bool isDeadLock
  8:  optional string deadLockMessage
  9:  optional list<TAgentThreadDetail> threadDetail
  10: i64 analysisid
  11: string appId
  12: string tenant
}

struct TAgentThreadChunk {
  1:  string appkey
  2:  i64 id
  4:  string applicationName
  5:  string agentId
  6:  i64 ts
  9:  optional list<TAgentThreadDetail> threadDetail
  10: i64 analysisid
  11: string appId
  12: string tenant
}

struct TDbMetaData {
  1: string agentId
  2: string appKey
  4: string dbUrl
  5: string dbHost
  6: string dbName
  8: i16 dbTypeCode
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

struct TAnnotation {
  1: i32 key,
  2: optional TAnnotationValue value
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

struct TExceptionMetaData2Api {
  1: i64 ts
  2: string name
  3: string msg
  4: string method
  5: string exceptionClass
  6: string apiName
  7: string url
  8: string tier
  9: string agent_id
  10: string app_key
  11: string tenant
  12: string appId

  14: string transactionId
  15: string spanId
  16: string pspanId
  17: optional string pagentId
  18: optional i32  exceptionId
  19: optional string userId
  20: optional string sessionId
  21: i64 agentStartTime
  22: optional string ip
}

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

struct TServiceType {
  1: string name
  2: i16  code
  3: string desc
  4: bool isInternalMethod
  5: bool isRpcClient
  6: bool isRecordStatistics
  7: bool isUnknown
  8: bool isUser
  9: bool isTerminal
  10: bool isQueue
  11: bool isIncludeDestinationId
}

struct TServiceTypeChunk {
	1: string app_key
	2: string agent_id
	3: string appId
    4: string tenant
	10: list<TServiceType> serviceTypeList
}

struct TStringMetaData {
  1: string agentId
  2: i64 agentStartTime

  4: i32 stringId
  5: string stringValue;
}

struct TSqlMetaData {
  1: string appkey
  2: string template
  3: string hash
  4: string appId
  5: string tenant
}

struct TSqlMetaData2Api {

    1: string host
    2: string dbhost
    3: string db
    4: string table
    5: string status
    6: string method
    7: i64 dur
    8: string agent_id
    9: i64 startTime
    10: string owner
    11: string url
    12: string clause
    13: string bindValue
    14: string transactionId
    15: string spanId
    16: string appKey
    17: string err
    18: string pspanId
    19: optional string pagentId
    20: string dbtype
    21: i64 agentStartTime
    22: i32 sqlId
    23: string outputs
    24: string appId
    25: string tenant
}

struct TApiMetaData {
  1: string agentId
  2: i64 agentStartTime

  3: optional string appkey

  4: i32 apiId,
  5: string apiInfo,
  6: optional i32 line,
  7: string appId,
  8: string tenant,

  10: optional i32 type,
}

struct TResult {
  1: bool success
  2: optional string message
}
