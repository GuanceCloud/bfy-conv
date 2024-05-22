package parse

import (
	"context"
	"fmt"
	"github.com/GuanceCloud/bfy-conv/gen-go/apptrace"
	"github.com/GuanceCloud/bfy-conv/mock"
	"github.com/GuanceCloud/bfy-conv/utils"
	"github.com/IBM/sarama"
	"github.com/apache/thrift/lib/go/thrift"
	"testing"
	"time"
)

// go test -benchmem -run=^$  -bench ^BenchmarkParseTSpan$ gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/plugins/inputs/ddtrace -memprofile memprofile.out -benchtime=100x

func BenchmarkParseTSpan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := mock.GetKafkaSpanByte()
		tspan, err := parseTSpan(buf[4:])
		if err != nil {
			b.Log(err)
		}
		if tspan != nil {

		}
	}
	/*
		BenchmarkParseTSpan-8   	   12538	     95009 ns/op
		BenchmarkParseTSpan-8   	   12655	     91711 ns/op
	*/
}

func Test_parseTSpanV2(t *testing.T) {
	buf := mock.GetKafkaSpanByte()
	tspan, err := parseTSpan(buf[4:])
	if err != nil {
		t.Log(err)
	}

	t.Logf("tspan = %+v \n", tspan)
	t.Logf("tspan = %+v", tspan.ApiId)
}

func TestParseServiceType(t *testing.T) {
	utils.ParseServiceType()

	fmt.Println(len(utils.ServiceTypeMap))
}

func TestSendMetricToKafka(t *testing.T) {
	load := int64(10)
	cpu := &apptrace.TCpuLoad{
		JvmCpuLoad:    &load,
		SystemCpuLoad: &load,
	}
	gc := &apptrace.TJvmGc{
		Type:                      0,
		JvmMemoryHeapUsed:         100,
		JvmMemoryHeapMax:          1000,
		JvmMemoryNonHeapUsed:      0,
		JvmMemoryNonHeapMax:       0,
		JvmGcOldCount:             5,
		JvmGcOldTime:              0,
		JvmGcDetailed:             nil,
		JvmMemoryNonHeapCommitted: 0,
		TotalPhysicalMemory:       nil,
		TExecuteDfs:               nil,
		TExecuteIostat:            nil,
		JdbcConnNum:               nil,
		ThreadNum:                 nil,
		JvmGcOldCountNew:          nil,
		JvmGcOldTimeNew:           nil,
	}
	timeNow := time.Now().UnixMilli()
	stat := &apptrace.TAgentStatBatch{
		AgentId:        "agentID",
		StartTimestamp: timeNow,
		AppKey:         "aaaaaaaa",
		AppId:          "aaaaaaaaa",
		Tenant:         "aaaaaaaa",
		AgentStats: []*apptrace.TAgentStat{{
			AgentId:         nil,
			StartTimestamp:  nil,
			Timestamp:       nil,
			CollectInterval: nil,
			Gc:              gc,
			CpuLoad:         cpu,
			Transaction:     nil,
			ActiveTrace:     nil,
			Metadata:        nil,
			ThreadCount:     nil,
		},
		},
		ServiceType: nil,
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
	if err := stat.Write(context.Background(), protocol); err != nil {
		// 处理错误
		fmt.Printf("span  err=%v", err)
	}
	serializedData := transport.Buffer.Bytes() // 序列化后的数据
	buf := make([]byte, 0)
	buf = append(buf, 239, 16, 0, 56)
	buf = append(buf, serializedData...)
	fmt.Println()
	for i := 0; i < 20; i++ {
		fmt.Printf(" %d, ", buf[i])
	}
	fmt.Println()

	// 发送到 kafka
	topic := "agentstat"
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
