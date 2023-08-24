package parse

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/GuanceCloud/bfy-conv/gen-go/span"
	"github.com/GuanceCloud/cliutils/logger"
	"github.com/GuanceCloud/cliutils/point"
	"github.com/apache/thrift/lib/go/thrift"
	"strconv"
	"strings"
	"time"
)

var log *logger.Logger
var HeaderKey = "x-b3-traceid"

func SetLogger(slog *logger.Logger) {
	log = slog
}

// Handle : message to points.
func Handle(message []byte) (pts []*point.Point) {
	// 判断类型
	msgType, err := code(message)
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
		log.Debugf("tspan=%s", tSpan.String())
		log.Debugf("TransactionId=%s  AppId=%s  AgentId=%s", tSpan.TransactionId, tSpan.AppId, tSpan.AgentId)

		xID := xid(tSpan.TransactionId, tSpan.AppId, tSpan.AgentId)
		tid := getTidFromRedis(xID)
		if tid == "" {
			tid = getTidFromHeader(tSpan.GetHttpRequestHeader(), HeaderKey, xID)
		}

		if tid == "" {
			tid = xID
		}
		pts = tSpanToPoint(tSpan, tid, xID)
	case 70:
		tSpanChunk, err := parseTSpanChunk(message[4:])
		if err != nil {
			log.Warnf("parse tSpan err=%v", err)
			//	return pts
		}

		log.Debugf("tspanChunk=%s", tSpanChunk.String())
		xID := xid(tSpanChunk.TransactionId, tSpanChunk.AppId, tSpanChunk.AgentId)
		tid := getTidFromRedis(xID)

		if tid == "" {
			tid = xID
		}

		pts = tSpanChunkToPoint(tSpanChunk, tid, xID)
	default:
		log.Debugf("unknown type code=%d", msgType)
	}

	return pts
}

func ptdecodeEvent(event *span.TSpanEvent) *point.Point {
	pt := &point.Point{}
	pt.SetName("kafka-bfy")
	pt.Add([]byte("span_id"), strconv.FormatInt(GetRandomWithAll(), 10))

	d := (event.StartElapsed + event.EndElapsed) * 1e3 // 不乘
	if d < 0 {
		d = 1000
	}
	pt.Add([]byte("duration"), d)
	resource := ""
	if st, ok := ServiceTypeMap[event.ServiceType]; ok {
		resource = st.Name

		if st.IsQueue {
			pt.MustAddTag([]byte("source_type"), []byte("message_queue"))
		}

		if st.IsIncludeDestinationID == 1 {
			pt.MustAddTag([]byte("source_type"), []byte("db"))
		}

		if st.IsRecordStatistics == 1 {
			pt.MustAddTag([]byte("source_type"), []byte("custom"))
		}

		if st.IsInternalMethod == 1 {
			pt.MustAddTag([]byte("source_type"), []byte("custom"))
		}

		if st.IsRpcClient == 1 {
			pt.MustAddTag([]byte("source_type"), []byte("http"))
		}

		if st.IsTerminal == 1 {
			pt.MustAddTag([]byte("service"), []byte(strings.ToLower(st.TypeDesc)))
			pt.MustAddTag([]byte("source_type"), []byte("db"))
		}

		if st.IsUser == 1 {
			pt.MustAddTag([]byte("source_type"), []byte("custom"))
		}

		if st.IsUnknown == 1 {
			pt.MustAddTag([]byte("source_type"), []byte("unknown"))
		}

		if pt.GetTag([]byte("source_type")) == nil {
			pt.MustAddTag([]byte("source_type"), []byte("unknown"))
		}
		/*		switch {
				case st.IsQueue:
					pt.AddTag([]byte("source_type"), []byte("message_queue"))
				case st.IsIncludeDestinationID == 1:
					pt.AddTag([]byte("source_type"), []byte("db"))
				case st.IsRecordStatistics == 1:
					pt.AddTag([]byte("source_type"), []byte("custom"))
				case st.IsInternalMethod == 1:
					pt.AddTag([]byte("source_type"), []byte("custom"))
				case st.IsRpcClient == 1:
					pt.AddTag([]byte("source_type"), []byte("http"))
				case st.IsTerminal == 1:
					pt.AddTag([]byte("service"), []byte(strings.ToLower(st.TypeDesc)))
					pt.AddTag([]byte("source_type"), []byte("db"))
				case st.IsUser == 1:
					pt.AddTag([]byte("source_type"), []byte("custom"))
				case st.IsUnknown == 1:
					pt.AddTag([]byte("source_type"), []byte("unknown"))
				default:
					//	pt.AddTag([]byte("source_type"), []byte("unknown"))
				}*/

	} else {
		return nil
	}

	if event.IsSetRPC() {
		pt.AddTag([]byte("rpc"), []byte(*event.RPC))
	}
	if event.IsSetURL() {
		pt.AddTag([]byte("url"), []byte(*event.URL))
	}
	if event.IsSetSql() {
		pt.AddTag([]byte("db.host"), []byte(event.Sql.Dbhost))
		pt.AddTag([]byte("db.type"), []byte(event.Sql.Dbtype))
		pt.AddTag([]byte("db.status"), []byte(event.Sql.Status))
	}

	if event.IsSetDestinationId() {
		pt.AddTag([]byte("operation"), []byte(*event.DestinationId))
	} else {
		pt.AddTag([]byte("operation"), []byte(resource))
	}

	pt.Add([]byte("resource"), resource)
	pt.AddTag([]byte("source"), []byte("byf-kafka"))

	if event.IsSetAnnotations() {
		for _, ann := range event.Annotations {
			pt.AddTag([]byte("source"), []byte(ann.GetValue().String()))
		}
	}
	jsonBody, err := json.Marshal(event)
	if err == nil {
		pt.Add([]byte("message"), string(jsonBody))
	}
	return pt
}

func parseTSpanChunk(buf []byte) (*span.TSpanChunk, error) {
	transport := &thrift.TMemoryBuffer{
		Buffer: bytes.NewBuffer(buf),
	}

	protocol := thrift.NewTCompactProtocolConf(transport, &thrift.TConfiguration{})
	tSpanChunk := span.NewTSpanChunk()
	ctx := context.Background()
	err := tSpanChunk.Read(ctx, protocol)
	return tSpanChunk, err
}

func tSpanChunkToPoint(tSpanChunk *span.TSpanChunk, traceID string, transactionID string) (pts []*point.Point) {
	if tSpanChunk == nil {
		return
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

	for _, event := range tSpanChunk.SpanEventList {
		eventPt := ptdecodeEvent(event)
		if eventPt == nil {
			continue
		}
		// eventPt.AddTag()
		eventPt.Add([]byte("trace_id"), traceID)
		eventPt.Add([]byte("parent_id"), strconv.FormatInt(tSpanChunk.SpanId, 10))
		eventPt.Add([]byte("start"), startTime+int64(event.StartElapsed)*1e3)
		//	eventPt.AddTag([]byte("service"), []byte(tSpanChunk.ApplicationName))
		if eventPt.GetTag([]byte("service")) == nil {
			eventPt.AddTag([]byte("service"), []byte(tSpanChunk.ApplicationName))
		}
		eventPt.AddTag([]byte("span_type"), []byte("entry"))
		eventPt.AddTag([]byte("source"), []byte("byf-kafka"))
		eventPt.AddTag([]byte("service_type"), []byte("bfy-tspanchunk"))
		eventPt.AddTag([]byte("transactionId"), []byte(transactionID))
		pts = append(pts, eventPt)
	}
	return pts
}
