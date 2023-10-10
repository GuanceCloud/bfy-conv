struct TServerMetaData {
  1: optional string              serverInfo
  2: optional list<string>        vmArgs
  10: optional list<TServiceInfo>  serviceInfos
}

struct TServiceInfo {
  1: optional string          serviceName
  2: optional list<string>    serviceLibs
}

struct TJvmInfo {
  1:          i16         version = 0
  2: optional string      vmVersion
  3: optional TJvmGcType  gcType = TJvmGcType.UNKNOWN
}

enum TJvmGcType {
  UNKNOWN,
  SERIAL,
  PARALLEL,
  CMS,
  G1,
  JRockitThroughput
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

struct TActiveTrace {
	1: optional TActiveTraceHistogram   histogram
}

struct TActiveTraceHistogram {
  1: i16         version = 0
  2: optional i32         histogramSchemaType
  3: optional list<i32>   activeTraceCount
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

struct TAgentEvent {
  1: string           appkey
  5: string           agentId
  6: i64              eventTimestamp
  7: TAgentEventType  eventType
  8: optional string  eventMessage
}

struct TAgentEventType {
  1: i32    code
  2: string desc
}