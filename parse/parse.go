package parse

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/GuanceCloud/bfy-conv/gen-go/span"
	"github.com/GuanceCloud/bfy-conv/utils"
	"github.com/GuanceCloud/cliutils/logger"
	"github.com/GuanceCloud/cliutils/point"
	"github.com/apache/thrift/lib/go/thrift"
	"strconv"
	"strings"
	"time"
)

var log = logger.DefaultSLogger("bfy")
var appFilter *AppFilter
var SkipSpanChunk = true

func SetLogger(slog *logger.Logger) {
	log = slog
}

type AppFilter struct {
	ProjectName string
	Projects    map[string][]string
}

func InitAppFilter(apps map[string]string) {
	if apps == nil || len(apps) == 0 {
		return
	}
	af := &AppFilter{
		Projects: make(map[string][]string),
	}
	for pname, anames := range apps {
		ns := strings.Split(anames, ",")
		af.Projects[pname] = ns
	}
	appFilter = af
}

// Handle : message to points.
func Handle(message []byte) (pts []*point.Point, category point.Category) {
	pts = make([]*point.Point, 0)
	// 判断类型
	msgType, err := utils.Code(message)
	if err != nil {
		log.Warnf("code err %v", err)
		return
	}
	// 序列化 tspan 或者 tspanchunk
	// 通过对象获取xid
	// 查询redis中获取traceid 如果获取不到 从http header 中获取并放到redis中
	// traceid 替换对象中的的traceID，返回

	switch msgType {
	case 40:
		tSpan, err := parseTSpan(message[4:])
		if err != nil {
			log.Warnf("parse tSpan err=%v", err)
			// 不返回错误，因为tspan 不一定为空
		}
		// log.Debugf("tspan=%s", tSpan.String())
		log.Debugf("TransactionId=%s  AppId=%s  AgentId=%s traceparent=%s tracestate=%s", tSpan.TransactionId, tSpan.AppId, tSpan.AgentId, tSpan.GetTraceparent(), tSpan.GetTracestate())

		xID := utils.Xid(tSpan.TransactionId, tSpan.AppId, tSpan.AgentId)
		pts = tSpanToPoint(tSpan, xID)
		category = point.Tracing
	case 70:
		if SkipSpanChunk {
			return
		}
		tSpanChunk, err := parseTSpanChunk(message[4:])
		if err != nil {
			log.Warnf("parse tSpan err=%v", err)
			//	return pts
		}

		// log.Debugf("tspanChunk=%s", tSpanChunk.String())
		xID := utils.Xid(tSpanChunk.TransactionId, tSpanChunk.AppId, tSpanChunk.AgentId)
		/* tid := getTidFromRedis(xID)

		if tid == "" {
			tid = xID
		}*/

		pts = tSpanChunkToPoint(tSpanChunk, xID)
		category = point.Tracing
	case 56:
		agentStat, err := parseAgentStatBatch(message[4:])
		if err != nil {
			log.Warnf("can parse to AgentStatBatch! err=%v", err)
		}
		if agentStat != nil {
			pts = statBatchToPoints(agentStat)
			category = point.Metric
		}
	case 50:
		agentInfo, err := parseAgentInfo(message[4:])
		if err != nil {
			log.Warnf("can parse to AgentInfo")
		}

		if agentInfo != nil {
			log.Warnf("AgentInfo=%+v", agentInfo)
			return
		}

	case 57:
		AgentEvent, err := parseAgentEvent(message[4:])
		if err != nil {
			log.Warnf("can parse to AgentEvent")
		}

		if AgentEvent != nil {
			log.Warnf("AgentEvent=%+v", AgentEvent)
		}

	default:
		// todo ...
		// 50 AgentInfo
		// 55 AgentStats
		// 56 AgentStatBatch
		// 57 AgentEvent
		// 58 AgentLifeCycle

		log.Debugf("unknown type code=%d", msgType)
	}

	return pts, category
}

func ptdecodeEvent(event *span.TSpanEvent) *point.Point {
	var kvs point.KVs

	kvs = kvs.Add("span_id", strconv.FormatInt(utils.GetRandomWithAll(), 16), false, false)

	d := (event.StartElapsed + event.EndElapsed) * 1e3 // 不乘
	if d < 0 {
		d = 1000
	}
	kvs = kvs.Add("duration", d, false, false)
	resource := ""
	if st, ok := utils.ServiceTypeMap[event.ServiceType]; ok {
		resource = st.Name

		if st.IsQueue {
			kvs = kvs.MustAddTag("source_type", "message_queue")
		}

		if st.IsIncludeDestinationID == 1 {
			kvs = kvs.MustAddTag("source_type", "db")
		}

		if st.IsRecordStatistics == 1 {
			kvs = kvs.MustAddTag("source_type", "custom")
		}

		if st.IsInternalMethod == 1 {
			kvs = kvs.MustAddTag("source_type", "custom")
		}

		if st.IsRpcClient == 1 {
			kvs = kvs.MustAddTag("source_type", "http")
		}

		if st.IsTerminal == 1 {
			kvs = kvs.MustAddTag("service", strings.ToLower(st.TypeDesc)).
				MustAddTag("source_type", "db")
		}

		if st.IsUser == 1 {
			kvs = kvs.MustAddTag("source_type", "custom")
		}

		if st.IsUnknown == 1 {
			kvs = kvs.MustAddTag("source_type", "unknown")
		}

		if kvs.GetTag("source_type") == "" {
			kvs = kvs.MustAddTag("source_type", "unknown")
		}
	} else {
		return nil
	}

	if event.IsSetRPC() {
		rpc := event.GetRPC()
		kvs = kvs.AddTag("resource", rpc).
			AddTag("operation", rpc)
		index := strings.Index(rpc, "?")
		if index != -1 {
			route := rpc[:index]
			kvs = kvs.AddTag("rpc_route", route)
		} else {
			kvs = kvs.AddTag("rpc_route", rpc)
		}
	}
	if event.IsSetURL() {
		kvs = kvs.AddTag("url", *event.URL)
	}
	if event.IsSetSql() {
		kvs = kvs.AddTag("db.host", event.Sql.Dbhost).
			AddTag("db.type", event.Sql.Dbtype).
			AddTag("db.status", event.Sql.Status)
	}

	if event.IsSetDestinationId() {
		kvs = kvs.AddTag("operation", *event.DestinationId)
	} else {
		kvs = kvs.AddTag("operation", resource)
	}

	kvs = kvs.Add("resource", resource, false, false).
		AddTag("source", "byf-kafka")

	if event.IsSetAnnotations() {
		for _, ann := range event.Annotations {
			kvs = kvs.AddTag("key"+strconv.Itoa(int(ann.Key)), (ann.GetValue().String()))
		}
	}
	jsonBody, err := json.Marshal(event)
	if err == nil {
		kvs = kvs.Add("message", string(jsonBody), false, false)
	}

	return point.NewPointV2("kafka-bfy", kvs, point.CommonLoggingOptions()...)
}

func parseTSpanChunk(buf []byte) (*span.TSpanChunk, error) {
	transport := &thrift.TMemoryBuffer{
		Buffer: bytes.NewBuffer(buf),
	}

	protocol := thrift.NewTCompactProtocolConf(transport, &thrift.TConfiguration{})
	tSpanChunk := span.NewTSpanChunk()
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	err := tSpanChunk.Read(ctx, protocol)
	return tSpanChunk, err
}

func tSpanChunkToPoint(tSpanChunk *span.TSpanChunk, transactionID string) (pts []*point.Point) {
	if tSpanChunk == nil {
		return
	}
	appName := tSpanChunk.ApplicationName
	projectKey := "project"
	projectVal := ""
	if appFilter != nil {
		filter := false
		// 过滤 app 名称， 通过之后增加tag：project="project_name"
		for pName, appNames := range appFilter.Projects {
			for _, name := range appNames {
				if name == appName {
					projectVal = pName
					filter = true
					break
				}
			}
		}
		if !filter {
			log.Debugf("del applicationName %s", appName)
			return
		}
	}
	if tSpanChunk.SpanEventList == nil || len(tSpanChunk.SpanEventList) == 0 {
		return
	}
	startTime := time.Now().UnixMicro()
	if tSpanChunk.StartTime != nil {
		//pt.SetTime(time.UnixMilli(*tSpanChunk.StartTime))
		startTime = *tSpanChunk.StartTime * 1e3
	} else {
		log.Warnf("tspanchunk starttime is null")
	}
	traceParent := tSpanChunk.GetTraceparent()
	for _, event := range tSpanChunk.SpanEventList {
		traceID, spanID := getTraceIDAndSpanIDFromW3C(traceParent)
		eventPt := ptdecodeEvent(event)
		if eventPt == nil {
			continue
		}
		// eventPt.AddTag()
		eventPt.Add("trace_id", traceID)
		eventPt.Add("parent_id", parentIDToDK(spanID))
		eventPt.Add("start", startTime+int64(event.StartElapsed)*1e3)

		if eventPt.GetTag("service") == "" {
			eventPt.AddTag("service", tSpanChunk.ApplicationName)
		}
		if projectVal != "" {
			eventPt.AddTag(projectKey, projectVal)
		}
		eventPt.AddTag("span_type", "entry")
		eventPt.AddTag("source", "byf-kafka")
		eventPt.AddTag("service_type", "bfy-tspanchunk")
		eventPt.AddTag("process_time", time.Now().Format("2006-01-02 15:04:05.000"))
		eventPt.AddTag("transactionId", transactionID)
		pts = append(pts, eventPt)
	}
	return pts
}
