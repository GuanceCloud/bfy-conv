package parseV2

import (
	"bytes"
	"context"
	"fmt"
	"github.com/GuanceCloud/bfy-conv/gen-go/span"
	"github.com/IBM/sarama"
	"github.com/apache/thrift/lib/go/thrift"
	"strconv"
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
				if projectFilter(sqlMeta.AppId) == "" {
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
				if projectFilter(apiMeta.AppId) == "" {
					return
				}

				apiSet(apiMeta)
			}
		}
	}
}

func sqlSetToCache(sql *span.TSqlMetaData) {
	//key := fmt.Sprintf("%s-sql-%s", sql.AppId, sql.Hash)
	val := sql.GetTemplate()
	//RedigoSet(key, val)

	RedigoHSet(sql.AppId, sql.Hash, val)
}

func sqlGetFromCache(appID, hash string) string {
	//key := fmt.Sprintf("%s-sql-%s", appID, hash)
	//val := RedigoGet(key)
	val := RedigoHGet(appID, hash)
	if val == "" {
		if getFromOld { // 如果获取不到，从旧数据中查询。
			key := fmt.Sprintf("%s-sql-%s", appID, hash)
			val = RedigoGet(key)
			// set
			if val != "" {
				RedigoHSet(appID, hash, val)
				return val
			}
		}
		return hash // 如果查询不到，将hash返回。
	}
	return val
}

func apiSet(api *span.TApiMetaData) {
	//key := fmt.Sprintf("%s-%s-%d", api.GetAgentId(), "api", api.GetApiId())
	val := fmt.Sprintf("%s line:%d", api.GetApiInfo(), api.GetLine())
	//RedigoSet(key, val)
	id := strconv.Itoa(int(api.GetApiId()))
	RedigoHSet(api.GetAgentId(), id, val)
}

func apiGet(agentId string, apiid int) string {
	//key := fmt.Sprintf("%s-%s-%d", agentId, "api", apiid)
	//val := RedigoGet(key)

	id := strconv.Itoa(apiid)
	val := RedigoHGet(agentId, id)

	if val == "" {
		if getFromOld {
			key := fmt.Sprintf("%s-%s-%d", agentId, "api", apiid)
			val = RedigoGet(key)
			if val != "" {
				RedigoHSet(agentId, id, val)
				return val
			}
		}
		return id // 获取不到返回id
	}
	return val
}
