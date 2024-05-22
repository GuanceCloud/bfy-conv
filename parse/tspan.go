package parse

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/GuanceCloud/bfy-conv/utils"
	"strconv"
	"strings"
	"time"

	"github.com/GuanceCloud/bfy-conv/gen-go/span"
	"github.com/GuanceCloud/cliutils/point"
	"github.com/apache/thrift/lib/go/thrift"
)

func parseTSpan(buf []byte) (*span.TSpan, error) {
	transport := &thrift.TMemoryBuffer{
		Buffer: bytes.NewBuffer(buf),
	}
	strict := false
	protocol := thrift.NewTCompactProtocolConf(transport, &thrift.TConfiguration{
		MaxMessageSize:     0,
		MaxFrameSize:       0,
		TBinaryStrictRead:  &strict,
		TBinaryStrictWrite: &strict,
	})
	tSpan := span.NewTSpan()
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	err := tSpan.Read(ctx, protocol)

	return tSpan, err
}

func tSpanToPoint(tSpan *span.TSpan, xid string) []*point.Point {
	pts := make([]*point.Point, 0)
	appName := tSpan.ApplicationName
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
			return pts
		}
	}
	traceID, spanID := getTraceIDAndSpanIDFromW3C(tSpan.GetTraceparent())
	if traceID == "" {
		traceID = xid
	}
	if spanID == "" {
		spanID = strconv.FormatInt(tSpan.SpanId, 10)
	}
	for _, event := range tSpan.SpanEventList {
		eventPt := ptdecodeEvent(event)
		if eventPt == nil {
			continue
		}
		// eventPt.AddTag()
		eventPt.Add("trace_id", traceID)
		eventPt.Add("parent_id", spanID)
		eventPt.Add("start", (tSpan.StartTime+int64(event.StartElapsed))*1e3)
		if eventPt.GetTag("service") == "" {
			eventPt.AddTag("service", tSpan.ApplicationName)
		}

		eventPt.AddTag("source_type", utils.SourceType(event.ServiceType))
		eventPt.AddTag("service_type", "bfy-tspan")
		eventPt.AddTag("process_time", time.Now().Format("2006-01-02 15:04:05.000"))
		if projectVal != "" {
			eventPt.AddTag(projectKey, projectVal)
		}
		eventPt.AddTag("transactionId", xid)
		eventPt.SetTime(time.UnixMilli(tSpan.StartTime + int64(event.StartElapsed)))

		pts = append(pts, eventPt)
	}

	pt := &point.Point{}
	pt.SetName("kafka-bfy")
	pt.Add("span_id", spanID)
	pt.Add("trace_id", traceID)
	pid := tSpan.ParentSpanId
	if pid == 0 {
		pid = 0
	}
	pt.Add("parent_id", parentIDToDK(tSpan.GetParentId()))
	pt.Add("start", tSpan.StartTime*1e3)
	pt.Add("duration", tSpan.Elapsed*1e3)
	if tSpan.IsSetRPC() {
		rpc := tSpan.GetRPC()
		pt.Add("resource", rpc)
		pt.AddTag("operation", rpc)
		index := strings.Index(rpc, "?")
		if index != -1 {
			route := rpc[:index]
			pt.AddTag("rpc_route", route)
		} else {
			pt.AddTag("rpc_route", rpc)
		}
	} else {
		pt.Add("resource", "unknown")
		pt.AddTag("operation", "unknown")
	}
	pt.AddTag("agentId", tSpan.GetAgentId())
	pt.AddTag(projectKey, projectVal)
	pt.AddTag("service", tSpan.ApplicationName)
	pt.AddTag("service_name", utils.ServiceName(tSpan.ServiceType))
	pt.AddTag("source_type", utils.SourceType(tSpan.ServiceType))
	pt.AddTag("transactionId", xid)
	pt.AddTag("original_type", "Span")
	if tSpan.ExceptionInfo != nil && tSpan.Err != nil && *tSpan.Err != 0 {
		pt.AddTag("status", "error")
		pt.Add("exception", *tSpan.ExceptionInfo)
	} else {
		pt.AddTag("status", "ok")
	}
	if tSpan.GetTracestate() != "" {
		states := getTraceState(tSpan.GetTracestate())
		for k, v := range states {
			pt.AddTag(k, v)
		}
	}
	// requestBody 和 responseBody Headers 没有放进去时因为其中有敏感信息
	if tSpan.IsSetHttpMethod() {
		pt.AddTag("http_method", *tSpan.HttpMethod)
	}
	if tSpan.IsSetHttpRequestTID() {
		pt.AddTag("http_request_tid", *tSpan.HttpRequestTID)
	}
	if tSpan.IsSetRetcode() {
		pt.AddTag("http_status_code", strconv.Itoa(int(*tSpan.Retcode)))
	}
	pt.AddTag("span_type", "entry")
	pt.AddTag("service_type", "bfy-tspan")
	pt.AddTag("process_time", time.Now().Format("2006-01-02 15:04:05.000"))

	pt.SetTime(time.UnixMilli(tSpan.StartTime))
	pt.AddTag("event_count", strconv.Itoa(len(tSpan.SpanEventList)))
	tSpan.SpanEventList = make([]*span.TSpanEvent, 0) // 防止重复数据太多
	jsonBody, err := json.Marshal(tSpan)
	if err == nil {
		pt.Add("message", string(jsonBody))
	}
	pts = append(pts, pt)

	return pts
}

func getTraceState(tracestate string) map[string]string {
	traceTags := make(map[string]string)
	kvs := strings.Split(tracestate, ",")
	for _, kv := range kvs {
		if len(kv) == 0 {
			continue
		}
		strs := strings.Split(kv, "=")
		if len(strs) == 2 {
			traceTags[strs[0]] = strs[1]
		}
	}
	return traceTags
}

func isSample(traceparent string) bool {
	strs := strings.Split(traceparent, "-")
	if len(strs) != 4 {
		return true
	}
	return strs[3] == "01"
}

func getTraceIDAndSpanIDFromW3C(traceparent string) (string, string) {
	strs := strings.Split(traceparent, "-")
	if len(strs) != 4 {
		return "", ""
	}

	return strs[1], strs[2]
}

func parentIDToDK(parentID string) string {
	if parentID == "0000000000000000" {
		return "0"
	}
	return parentID
}
