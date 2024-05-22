package parseV2

import (
	"fmt"
	"github.com/GuanceCloud/bfy-conv/pb/callTree"
	"github.com/IBM/sarama"
	"github.com/golang/protobuf/proto"
	"testing"
	"time"
)

func Test_parseCallTree(t *testing.T) {
	// 初始化环境
	InitRedis("", "", "", 0)
	// 初始化数据
	callTree := &callTree.CallTree{Callevents: make([]*callTree.CallEvent, 0)}
	callTree.Callevents = append(callTree.Callevents, newCallEvent("000000000000092f0d8ca3e554b3a735", "174155503488541", 1))
	callTree.Callevents = append(callTree.Callevents, newCallEvent("000000000000092f0d8ca3e554b3a735", "174155503488541", 2))
	callTree.Callevents = append(callTree.Callevents, newCallEvent("000000000000092f0d8ca3e554b3a735", "174155503488541", 3))
	callTree.Callevents = append(callTree.Callevents, newCallEvent("000000000000092f0d8ca3e554b3a735", "174155503488541", 4))
	callTree.Callevents = append(callTree.Callevents, newCallEvent("000000000000092f0d8ca3e554b3a735", "174155503488541", 5))
	callTree.Callevents = append(callTree.Callevents, newCallEvent("000000000000092f0d8ca3e554b3a735", "174155503488541", 5))
	callTree.Callevents = append(callTree.Callevents, newCallEvent("000000000000092f0d8ca3e554b3a735", "174155503488541", 5))
	callTree.Callevents = append(callTree.Callevents, newCallEvent("000000000000092f0d8ca3e554b3a735", "174155503488541", 5))
	callTree.Callevents = append(callTree.Callevents, newCallEvent("000000000000092f0d8ca3e554b3a735", "174155503488541", 5))
	callTree.Callevents = append(callTree.Callevents, newCallEvent("000000000000092f0d8ca3e554b3a735", "174155503488541", 5))
	callTree.Callevents = append(callTree.Callevents, newCallEvent("000000000000092f0d8ca3e554b3a735", "174155503488541", 5))
	callTree.Callevents = append(callTree.Callevents, newCallEvent("000000000000092f0d8ca3e554b3a735", "174155503488541", 5))
	callTree.Callevents = append(callTree.Callevents, newCallEvent("000000000000092f0d8ca3e554b3a735", "174155503488541", 5))
	callTree.Callevents = append(callTree.Callevents, newCallEvent("000000000000092f0d8ca3e554b3a735", "174155503488541", 2))

	bts, err := proto.Marshal(callTree)
	if err != nil {
		t.Errorf("marshal err=%v", err)
		return
	}
	msg := &sarama.ConsumerMessage{
		Headers:        nil,
		Timestamp:      time.Time{},
		BlockTimestamp: time.Time{},
		Key:            nil,
		Value:          bts,
		Topic:          "calltree",
		Partition:      0,
		Offset:         1,
	}
	// 1 调用方法
	pts, category := parseCallTree(msg)
	t.Logf("category =%d", category)
	for _, pt := range pts {
		t.Logf("point=%s", pt.LineProto())
	}
	// 2 发送到kafka
	sendToKafka("dwd_callevents", bts, t)
	Close()
}

func newCallEvent(traceID string, spanID string, depth int32) *callTree.CallEvent {
	event := &callTree.CallEvent{
		Id:                 "000",
		SpanId:             spanID,
		Sequence:           0,
		StartElapsed:       1000,
		EndElapsed:         100,
		Rpc:                "",
		ServiceType:        1010,
		EndPoint:           "",
		Depth:              depth,
		NextSpanId:         0,
		DestinationId:      "",
		ApiId:              991,
		ExceptionClassName: "",
		AsyncId:            0,
		NextAsyncId:        0,
		AsyncSequence:      0,
		ApiInfo:            "",
		LineNumber:         0,
		Retcode:            0,
		RequestHeaders:     "",
		RequestBody:        "",
		ResponseBody:       "",
		Status:             0,
		Url:                "",
		Method:             "",
		Arguments:          "",
		Ps:                 "",
		Tenant:             "",
		Appid:              "client",
		Appsysid:           "",
		AgentId:            "agentid",
		AgentIp:            "",
		Trxid:              "",
		BootTime:           0,
		HasException:       false,
		ExceptionId:        "",
		AppServiceType:     0,
		UserId:             "",
		SessionId:          "",
		Ts:                 1716361416000,
		FromWebAndMobile:   false,
		TraceId:            traceID,
		IsOtel:             false,
		EventCid:           "",
	}

	return event
}

func sendToKafka(topic string, buf []byte, t *testing.T) {
	// 发送到 kafka
	brokerAddr := []string{"10.200.14.226:9092"}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokerAddr, config)
	if err != nil {
		t.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	// Send message to Kafka.
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(buf),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		t.Fatalf("Failed to send message to Kafka: %v", err)
	}

	fmt.Printf("Message sent successfully. Partition: %d, Offset: %d\n", partition, offset)
}
