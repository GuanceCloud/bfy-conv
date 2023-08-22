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
	"time"
)

var log *logger.Logger
var HeaderKey = "x-b3-traceid"

func SetLogger(slog *logger.Logger) {
	log = slog
}

/*
spanchunk 和span 关联关系是 spanid

start time 毫秒
agentStart  没意义

*/

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

func parseTSpan(buf []byte) (*span.TSpan, error) {
	transport := &thrift.TMemoryBuffer{
		Buffer: bytes.NewBuffer(buf),
	}
	strict := false
	protocol := thrift.NewTCompactProtocolConf(transport, &thrift.TConfiguration{
		MaxMessageSize:     1024 * 20,
		MaxFrameSize:       0,
		TBinaryStrictRead:  &strict,
		TBinaryStrictWrite: &strict,
	})
	tSpan := span.NewTSpan()
	ctx := context.Background()
	err := tSpan.Read(ctx, protocol)
	return tSpan, err
	/*deserializer := thrift.NewTDeserializer()
	tSpan := span.NewTSpan()
	ctx := context.Background()
	if err := deserializer.Read(ctx, tSpan, buf); err != nil {
		log.Errorf("deserializer tSpan err=%v", err)
		return nil, err
	} else {
		return tSpan, err
	}*/
}

func tSpanToPoint(tSpan *span.TSpan, traceid string, xid string) []*point.Point {
	pts := make([]*point.Point, 0)
	for _, event := range tSpan.SpanEventList {
		eventPt := ptdecodeEvent(event)
		if eventPt == nil {
			continue
		}
		// eventPt.AddTag()
		eventPt.Add([]byte("trace_id"), traceid)
		eventPt.Add([]byte("parent_id"), strconv.FormatInt(tSpan.SpanId, 10))
		eventPt.Add([]byte("start"), (tSpan.StartTime+int64(event.StartElapsed))*1e3)
		eventPt.AddTag([]byte("service"), []byte(tSpan.ApplicationName))
		eventPt.AddTag([]byte("transactionId"), []byte(xid))
		eventPt.SetTime(time.UnixMilli(tSpan.StartTime + int64(event.StartElapsed)))

		pts = append(pts, eventPt)
	}

	pt := &point.Point{}
	pt.SetName("kafka-bfy")
	pt.Add([]byte("span_id"), strconv.FormatInt(tSpan.SpanId, 10))
	pt.Add([]byte("trace_id"), traceid)
	pid := tSpan.ParentSpanId
	if pid < 0 {
		pid = 0
	}
	pt.Add([]byte("parent_id"), strconv.FormatInt(pid, 10))
	pt.Add([]byte("start"), tSpan.StartTime*1e3)
	pt.Add([]byte("duration"), tSpan.Elapsed*1e3)
	pt.Add([]byte("resource"), *tSpan.RPC)

	pt.AddTag([]byte("service"), []byte(tSpan.ApplicationName))
	pt.AddTag([]byte("service_name"), []byte(serviceName(tSpan.ServiceType)))
	pt.AddTag([]byte("operation"), []byte(*tSpan.RPC))
	pt.AddTag([]byte("source_type"), []byte(sourceType(tSpan.ServiceType)))
	pt.AddTag([]byte("transactionId"), []byte(xid))
	pt.AddTag([]byte("original_type"), []byte("Span"))
	if tSpan.ExceptionInfo != nil && tSpan.Err != nil && *tSpan.Err != 0 {
		pt.AddTag([]byte("status"), []byte("error"))
		pt.Add([]byte("exception"), *tSpan.ExceptionInfo)
	} else {
		pt.AddTag([]byte("status"), []byte("ok"))
	}

	pt.AddTag([]byte("span_type"), []byte("entry"))
	pt.AddTag([]byte("source"), []byte("byf-kafka"))
	pt.AddTag([]byte("service_type"), []byte("byf-tspan"))

	pt.SetTime(time.UnixMilli(tSpan.StartTime))
	pt.AddTag([]byte("event_count"), []byte(strconv.Itoa(len(tSpan.SpanEventList))))
	tSpan.SpanEventList = make([]*span.TSpanEvent, 0) // 防止重复数据太多
	jsonBody, err := json.Marshal(tSpan)
	if err == nil {
		pt.Add([]byte("message"), string(jsonBody))
	}
	pts = append(pts, pt)

	return pts
}

func ptdecodeEvent(event *span.TSpanEvent) *point.Point {
	pt := &point.Point{}
	pt.SetName("kafka-bfy")
	pt.Add([]byte("span_id"), strconv.FormatInt(GetRandomWithAll(), 10))
	d := (event.StartElapsed + event.EndElapsed) * 1e3
	if d < 0 {
		d = 1000
	}
	pt.Add([]byte("duration"), d)
	resource := ""
	if st, ok := ServiceTypeMap[event.ServiceType]; ok {
		resource = st.Name
		switch {
		case st.IsQueue:
			pt.AddTag([]byte("source_type"), []byte("db"))
			if event.IsSetSql() {
				pt.AddTag([]byte("db.host"), []byte(event.Sql.Dbhost))
				pt.AddTag([]byte("db.type"), []byte(event.Sql.Dbtype))
				pt.AddTag([]byte("db.status"), []byte(event.Sql.Status))
			}
		case st.IsIncludeDestinationID == 1:
			pt.AddTag([]byte("source_type"), []byte("db"))
		case st.IsRecordStatistics == 1:
			pt.AddTag([]byte("source_type"), []byte("record"))
		case st.IsInternalMethod == 1:
			pt.AddTag([]byte("source_type"), []byte("Internal"))
			if event.IsSetURL() {
				pt.AddTag([]byte("url"), []byte(*event.URL))
			}
		case st.IsRpcClient == 1:
			if event.IsSetRPC() {
				pt.AddTag([]byte("rpc"), []byte(*event.RPC))
			}
		case st.IsTerminal == 1:
			pt.AddTag([]byte("source_type"), []byte("Terminal"))
		case st.IsUser == 1:
			pt.AddTag([]byte("source_type"), []byte("user"))

		case st.IsUnknown == 1:
			pt.AddTag([]byte("source_type"), []byte("unknown"))
		default:
			pt.AddTag([]byte("source_type"), []byte("unknown"))
		}
	} else {
		return nil
	}

	pt.Add([]byte("resource"), resource)
	pt.AddTag([]byte("source"), []byte("byf-kafka"))
	pt.AddTag([]byte("operation"), []byte(resource))
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

	/*	deserializer := thrift.NewTDeserializer()
		tSpanChunk := span.NewTSpanChunk()
		ctx := context.Background()
		if err := deserializer.Read(ctx, tSpanChunk, buf); err != nil {
			log.Errorf("deserializer TSpanChunk err=%v", err)
			return nil, err
		} else {
			return tSpanChunk, err
		}*/
}

func tSpanChunkToPoint(tSpanChunk *span.TSpanChunk, traceID string, transactionID string) (pts []*point.Point) {
	if tSpanChunk == nil {
		return
	}

	if tSpanChunk.SpanEventList == nil || len(tSpanChunk.SpanEventList) == 0 {
		return
	}
	/*	pt := &point.Point{}
		pt.SetName("kafka-bfy")
		pt.Add([]byte("span_id"), strconv.FormatInt(tSpanChunk.SpanId, 10))
		pt.Add([]byte("trace_id"), traceID)

		pt.Add([]byte("parent_id"), "0")
		if tSpanChunk.StartTime != nil {
			pt.Add([]byte("start"), *tSpanChunk.StartTime)
		}
		//if tSpanChunk.AgentStartTime != 0 {
		//	pt.Add([]byte("start"), tSpanChunk.AgentStartTime)
		//}

		pt.Add([]byte("duration"), 1000)
		if tSpanChunk.EndPoint != nil {
			pt.Add([]byte("resource"), *tSpanChunk.EndPoint)
			pt.AddTag([]byte("operation"), []byte(*tSpanChunk.EndPoint))
		}

		pt.AddTag([]byte("service"), []byte(tSpanChunk.ApplicationName))
		pt.AddTag([]byte("service_name"), []byte(serviceName(tSpanChunk.ServiceType)))

		pt.AddTag([]byte("source_type"), []byte(sourceType(tSpanChunk.ServiceType)))
		pt.AddTag([]byte("transactionId"), []byte(transactionID))
		pt.AddTag([]byte("original_type"), []byte("Span"))

		pt.AddTag([]byte("status"), []byte("ok"))

		pt.AddTag([]byte("span_type"), []byte("entry"))
		pt.AddTag([]byte("source"), []byte("byf-kafka"))
		pt.AddTag([]byte("service_type"), []byte("byf-tspanchunk"))
		jsonBody, err := json.Marshal(tSpanChunk)
		if err == nil {
			pt.Add([]byte("message"), string(jsonBody))
		}



	*/
	startTime := time.Now().UnixMicro()
	if tSpanChunk.StartTime != nil {
		//pt.SetTime(time.UnixMilli(*tSpanChunk.StartTime))
		log.Warnf("tspanchunk starttime is null")
		startTime = *tSpanChunk.StartTime * 1e3
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
		eventPt.AddTag([]byte("service"), []byte(tSpanChunk.ApplicationName))
		eventPt.AddTag([]byte("transactionId"), []byte(transactionID))
		pts = append(pts, eventPt)
	}
	return pts
}
