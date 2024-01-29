package parse

import (
	"github.com/GuanceCloud/bfy-conv/mock"
	"testing"
)

func Test_tSpanToPoint(t *testing.T) {
	buf := mock.GetKafkaSpanByte()
	tspan, err := parseTSpan(buf[4:])
	if err != nil {
		t.Log(err)
		return
	}
	appFilter = nil
	pts := tSpanToPoint(tspan, "trace_id", "xxxx")
	if len(pts) == 0 {
		t.Errorf("pts len ==0")
		return
	}
	for _, pt := range pts {
		t.Logf("point=%s", pt.LineProto())
	}
}

/*

kafka-bfy,agentId=xytb-uw-pt-7d876c8d-7vjg6,
event_count=5,operation=gateway,
original_type=Span,
process_time=2023-12-21\ 14:16:28.205,rpc_route=gateway,service=urus-uw-pt-1,service_name=UNDERTOW_METHOD,
service_type=bfy-tspan,source_type=web,span_type=entry,status=ok,transactionId=xxxx duration=100000i,
message="{\"agentId\":\"xytb-uw-pt-7d876c8d-7vjg6\",\"applicationName\":\"urus-uw-pt-1\",\"agentStartTime\":1703139388205,\"transactionId\":\"ABp0eWp5Z2wtcHQtdWFw/Pj4kvEwVG1lc29zLWI5OTc4NzM5LTA3MGQtNDA1Ni04M2U1LWEyN2M0OWQ2OWU4Nq+u6QY=\",\"appkey\":\"e246fc1cfb4c428eb744a76cf996d10c\",\"spanId\":6186247787118501400,\"parentSpanId\":0,\"startTime\":1703139388205,\"elapsed\":100,\"rpc\":\"gateway\",\"serviceType\":1121,\"endPoint\":\"urus-uw/gateway/request\",\"remoteAddr\":\"140.2.1.53:31003\",\"flag\":0,\"parentApplicationName\":\"other-string\",\"httpRequestHeader\":\"x-forwarded-prefix,/uw-non/busi-biz-ns-pt/eurus-uw;x-b3-traceid,1d5b4afa2dd865a4;x-b3-spanid,fe4830ce06e18229;x-b3-parentspanid,1d5b4afa2dd865a4;x-b3-sampled,1;host,140.2.1.53:31003\",\"httpRequestBody\":\"other-string\",\"apidesc\":\"other-string\",\"httpResponseHeader\":\"other-string\",\"userId\":\"other-string\",\"sessionId\":\"other-string\",\"appId\":\"abc\",\"tenant\":\"\",\"threadName\":\"other-string\"}",
parent_id="0",resource="gateway",span_id="6186247787118501400",start=1703139388205000i,trace_id="trace_id" 1703139388205000000
*/
