package parseV2

import (
	"bytes"
	"context"
	"fmt"
	"github.com/GuanceCloud/bfy-conv/gen-go/span"
	"github.com/IBM/sarama"
	"github.com/apache/thrift/lib/go/thrift"
	"time"
)

func metadata(msg *sarama.ConsumerMessage) {
	for _, hander := range msg.Headers {
		if string(hander.Key) == "type" {
			val := string(hander.Value)
			switch val {
			case "sql_metadata":
				transport := &thrift.TMemoryBuffer{
					Buffer: bytes.NewBuffer(msg.Value),
				}
				protocol := thrift.NewTCompactProtocolConf(transport, &thrift.TConfiguration{})
				sqlMeta := span.NewTSqlMetaData()
				ctx, _ := context.WithTimeout(context.Background(), time.Second)
				err := sqlMeta.Read(ctx, protocol)
				if err != nil {
					return
				}
				sqlSetToCache(sqlMeta)
			case "api_metadata":
				transport := &thrift.TMemoryBuffer{
					Buffer: bytes.NewBuffer(msg.Value),
				}
				protocol := thrift.NewTCompactProtocolConf(transport, &thrift.TConfiguration{})
				apiMeta := span.NewTApiMetaData()
				ctx, _ := context.WithTimeout(context.Background(), time.Second)
				err := apiMeta.Read(ctx, protocol)
				if err != nil {
					return
				}
				apiSet(apiMeta)
			}
		}
	}
}

func sqlSetToCache(sql *span.TSqlMetaData) {
	key := fmt.Sprintf("%s-sql-%s", sql.AppId, sql.Hash)
	val := sql.GetTemplate()
	RedigoSet(key, val)
}

func sqlGetFromCache(apiID, hash string) string {
	key := fmt.Sprintf("%s-sql-%s", apiID, hash)
	val := RedigoGet(key)

	if val == "" {
		return key // 如果查询不到，将hash返回。
	}
	return hash
}

func apiSet(api *span.TApiMetaData) {
	key := fmt.Sprintf("%s-%s-%d", api.GetAgentId(), "api", api.GetApiId())
	val := fmt.Sprintf("%s line:%d", api.GetApiInfo(), api.GetLine())
	RedigoSet(key, val)
}

func apiGet(agentId string, apiid int) string {
	key := fmt.Sprintf("%s-%s-%d", agentId, "api", apiid)
	val := RedigoGet(key)
	if val == "" {
		return key
	}
	return val
}
