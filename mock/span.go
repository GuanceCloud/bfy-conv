package mock

import (
	"context"
	"fmt"
	"github.com/GuanceCloud/bfy-conv/gen-go/apptrace"
	"github.com/GuanceCloud/bfy-conv/gen-go/span"
	"github.com/apache/thrift/lib/go/thrift"
	"math"
	"math/rand"
	"time"
)

func GetKafkaSpanByte() []byte {
	//msg := spanMessage()
	//msg, _ := hex.DecodeString(des)
	rpc := "gateway"
	endpoint := "urus-uw/gateway/request"
	remoteAddr := "140.2.1.53:31003"
	other := "other-string"
	header := "x-forwarded-prefix,/uw-non/busi-biz-ns-pt/eurus-uw;x-b3-traceid,1d5b4afa2dd865a4;x-b3-spanid,fe4830ce06e18229;x-b3-parentspanid,1d5b4afa2dd865a4;x-b3-sampled,1;host,140.2.1.53:31003"
	mill := time.Now().UnixMilli()
	span := &span.TSpan{
		AgentId:                "xytb-uw-pt-7d876c8d-7vjg6",
		ApplicationName:        "urus-uw-pt-1",
		AgentStartTime:         mill,
		TransactionId:          getTranID(),
		Appkey:                 "e246fc1cfb4c428eb744a76cf996d10c",
		SpanId:                 GetRandomWithAll(),
		ParentSpanId:           0,
		StartTime:              mill,
		Elapsed:                100,
		RPC:                    &rpc,
		ServiceType:            1121,
		EndPoint:               &endpoint,
		RemoteAddr:             &remoteAddr,
		Annotations:            nil,
		Flag:                   0,
		Err:                    nil,
		SpanEventList:          getTSpanEventList(), //todo
		ParentApplicationName:  &other,
		ParentApplicationType:  nil,
		AcceptorHost:           nil,
		ApiId:                  nil,
		ExceptionInfo:          nil,
		ApplicationServiceType: nil,
		LoggingTransactionInfo: nil,
		HttpPara:               nil,
		HttpMethod:             nil,
		HttpRequestHeader:      &header,
		HttpRequestUserAgent:   nil,
		HttpRequestBody:        &other,
		HttpResponseBody:       nil,
		Retcode:                nil,
		HttpRequestUID:         nil,
		HttpRequestTID:         nil,
		PagentId:               nil,
		Apidesc:                &other,
		HttpResponseHeader:     &other,
		UserId:                 &other,
		SessionId:              &other,
		AppId:                  "abc",
		Tenant:                 "",
		ThreadId:               nil,
		ThreadName:             &other,
		HasNextCall:            nil,
	}
	transport := thrift.NewTMemoryBuffer()
	strict := true
	protocol := thrift.NewTCompactProtocolConf(transport, &thrift.TConfiguration{
		MaxMessageSize:     1024 * 20,
		MaxFrameSize:       0,
		TBinaryStrictRead:  &strict,
		TBinaryStrictWrite: &strict,
	})
	//	protocol := thrift.NewTBinaryProtocolFactoryDefault().GetProtocol(transport)
	if err := span.Write(context.Background(), protocol); err != nil {
		// 处理错误
		fmt.Printf("span  err=%v", err)
		return nil
	}
	serializedData := transport.Buffer.Bytes() // 序列化后的数据
	buf := make([]byte, 0)
	buf = append(buf, 239, 16, 0, 40)
	buf = append(buf, serializedData...)
	fmt.Println()
	for i := 0; i < 20; i++ {
		fmt.Printf(" %d, ", buf[i])
	}
	fmt.Println()

	return buf
}

func GetAppSpan() []byte {
	//msg := spanMessage()
	//msg, _ := hex.DecodeString(des)
	rpc := "gateway"
	endpoint := "urus-uw/gateway/request"
	remoteAddr := "140.2.1.53:31003"
	other := "other-string"
	header := "x-forwarded-prefix,/uw-non/busi-biz-ns-pt/eurus-uw;x-b3-traceid,1d5b4afa2dd865a4;x-b3-spanid,fe4830ce06e18229;x-b3-parentspanid,1d5b4afa2dd865a4;x-b3-sampled,1;host,140.2.1.53:31003"
	mill := time.Now().UnixMilli()
	span := &apptrace.TSpan{
		AgentId:                "xytb-uw-pt-7d876c8d-7vjg6",
		ApplicationName:        "urus-uw-pt-1",
		AgentStartTime:         mill,
		TransactionId:          getTranID(),
		Appkey:                 "e246fc1cfb4c428eb744a76cf996d10c",
		SpanId:                 GetRandomWithAll(),
		ParentSpanId:           0,
		StartTime:              mill,
		Elapsed:                100,
		RPC:                    &rpc,
		ServiceType:            1121,
		EndPoint:               &endpoint,
		RemoteAddr:             &remoteAddr,
		Annotations:            nil,
		Flag:                   0,
		Err:                    nil,
		ParentApplicationName:  &other,
		ParentApplicationType:  nil,
		AcceptorHost:           nil,
		ApiId:                  nil,
		ExceptionInfo:          nil,
		ApplicationServiceType: nil,
		LoggingTransactionInfo: nil,
		HttpPara:               nil,
		HttpMethod:             nil,
		HttpRequestHeader:      &header,
		HttpRequestUserAgent:   nil,
		HttpRequestBody:        &other,
		HttpResponseBody:       nil,
		Retcode:                nil,
		HttpRequestUID:         nil,
		HttpRequestTID:         nil,
		PagentId:               nil,
		Apidesc:                &other,
		HttpResponseHeader:     &other,
		UserId:                 &other,
		SessionId:              &other,
		AppId:                  "abc",
		Tenant:                 "",
		ThreadId:               nil,
		ThreadName:             &other,
		HasNextCall:            nil,
	}
	transport := thrift.NewTMemoryBuffer()
	protocol := thrift.NewTBinaryProtocolFactoryDefault().GetProtocol(transport)
	if err := span.Write(context.Background(), protocol); err != nil {
		// 处理错误
		fmt.Printf("span  err=%v", err)
		return nil
	}
	serializedData := transport.Buffer.Bytes() // 序列化后的数据
	buf := make([]byte, 0)
	buf = append(buf, 239, 16, 0, 40)
	buf = append(buf, serializedData...)
	fmt.Println()
	for i := 0; i < 20; i++ {
		fmt.Printf(" %d, ", buf[i])
	}
	fmt.Println()

	return buf
}

func getTranID() []byte {
	buf := make([]byte, 0)
	buf = append(buf, 0x00, 0x1a)
	buf = append(buf, []byte("tyjygl-pt-uap")...)
	buf = append(buf, 0xfc, 0xf8, 0xf8, 0x92, 0xf1)
	buf = append(buf, []byte("0Tmesos-b9978739-070d-4056-83e5-a27c49d69e86")...)
	buf = append(buf, 0xaf, 0xae, 0xe9, 0x06)

	return buf
}

func getTSpanEventList() []*span.TSpanEvent {
	tlist := make([]*span.TSpanEvent, 0)
	for i := 0; i < 5; i++ {
		spanID := GetRandomWithAll()
		rpc := "gateway"
		endpoint := "urus-uw/gateway/request"
		remoteAddr := "140.2.1.53:31003"
		other := "other-string"
		event := &span.TSpanEvent{
			SpanId:             &spanID,
			Sequence:           0,
			StartElapsed:       2000,
			EndElapsed:         1000,
			RPC:                &rpc,
			ServiceType:        1121,
			EndPoint:           &endpoint,
			Annotations:        nil,
			Depth:              3,
			NextSpanId:         0,
			DestinationId:      &endpoint,
			ApiId:              nil,
			ExceptionInfo:      nil,
			ExceptionClassName: &other,
			AsyncId:            nil,
			NextAsyncId:        nil,
			AsyncSequence:      nil,
			ApiInfo:            "",
			LineNumber:         nil,
			Sql:                nil,
			Retcode:            nil,
			RequestHeaders:     &remoteAddr,
			RequestBody:        nil,
			ResponseBody:       nil,
			URL:                nil,
			Method:             nil,
			Arguments:          nil,
		}
		tlist = append(tlist, event)
	}

	return tlist
}

func GetRandomWithAll() int64 {
	rand.Seed(time.Now().UnixNano())
	return int64(rand.Intn(math.MaxInt))
}
