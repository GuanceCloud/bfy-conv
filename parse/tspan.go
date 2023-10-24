package parse

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/GuanceCloud/bfy-conv/gen-go/span"
	"github.com/GuanceCloud/cliutils/point"
	"github.com/apache/thrift/lib/go/thrift"
	"strconv"
	"time"
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
	ctx := context.Background()
	err := tSpan.Read(ctx, protocol)
	return tSpan, err
}

func tSpanToPoint(tSpan *span.TSpan, traceid string, xid string) []*point.Point {
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
	for _, event := range tSpan.SpanEventList {
		eventPt := ptdecodeEvent(event)
		if eventPt == nil {
			continue
		}
		// eventPt.AddTag()
		eventPt.Add([]byte("trace_id"), traceid)
		eventPt.Add([]byte("parent_id"), strconv.FormatInt(tSpan.SpanId, 10))
		eventPt.Add([]byte("start"), (tSpan.StartTime+int64(event.StartElapsed))*1e3)
		if eventPt.GetTag([]byte("service")) == nil {
			eventPt.AddTag([]byte("service"), []byte(tSpan.ApplicationName))
		}
		if projectVal != "" {
			eventPt.AddTag([]byte(projectKey), []byte(projectVal))
		}
		eventPt.AddTag([]byte("transactionId"), []byte(xid))
		eventPt.SetTime(time.UnixMilli(tSpan.StartTime + int64(event.StartElapsed)))

		pts = append(pts, eventPt)
	}

	pt := &point.Point{}
	pt.SetName("kafka-bfy")
	pt.Add([]byte("span_id"), strconv.FormatInt(tSpan.SpanId, 10))
	pt.Add([]byte("trace_id"), traceid)
	pid := tSpan.ParentSpanId
	if pid == 0 {
		pid = 0
	}
	pt.Add([]byte("parent_id"), strconv.FormatInt(pid, 10))
	pt.Add([]byte("start"), tSpan.StartTime*1e3)
	pt.Add([]byte("duration"), tSpan.Elapsed*1e3)
	if tSpan.IsSetRPC() {
		pt.Add([]byte("resource"), *tSpan.RPC)
		pt.AddTag([]byte("operation"), []byte(*tSpan.RPC))
	} else {
		pt.Add([]byte("resource"), "unknown")
		pt.AddTag([]byte("operation"), []byte("unknown"))
	}
	pt.AddTag([]byte(projectKey), []byte(projectVal))
	pt.AddTag([]byte("service"), []byte(tSpan.ApplicationName))
	pt.AddTag([]byte("service_name"), []byte(serviceName(tSpan.ServiceType)))
	pt.AddTag([]byte("source_type"), []byte(sourceType(tSpan.ServiceType)))
	pt.AddTag([]byte("transactionId"), []byte(xid))
	pt.AddTag([]byte("original_type"), []byte("Span"))
	if tSpan.ExceptionInfo != nil && tSpan.Err != nil && *tSpan.Err != 0 {
		pt.AddTag([]byte("status"), []byte("error"))
		pt.Add([]byte("exception"), *tSpan.ExceptionInfo)
	} else {
		pt.AddTag([]byte("status"), []byte("ok"))
	}

	// requestBody 和 responseBody Headers 没有放进去时因为其中有敏感信息
	if tSpan.IsSetHttpMethod() {
		pt.AddTag([]byte("http_method"), []byte(*tSpan.HttpMethod))
	}
	if tSpan.IsSetHttpRequestTID() {
		pt.AddTag([]byte("http_request_tid"), []byte(*tSpan.HttpRequestTID))
	}
	if tSpan.IsSetRetcode() {
		pt.AddTag([]byte("http_status_code"), []byte(strconv.Itoa(int(*tSpan.Retcode))))
	}
	pt.AddTag([]byte("span_type"), []byte("entry"))
	pt.AddTag([]byte("service_type"), []byte("bfy-tspan"))

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
