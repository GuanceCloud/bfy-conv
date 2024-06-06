package parseV2

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"testing"
	"time"
)

func Test_parseSQL(t *testing.T) {
	sql := SQL{
		ID:          "eventid12345",
		Appid:       "appid12345",
		AppSysId:    "sysid12345",
		TrxId:       "xxx",
		TraceId:     "60000000000000011111222222",
		SpanId:      "4564564565456456",
		Group:       "78787",
		AgentId:     "java-agent-id",
		DbHost:      "10.200.14.188",
		DbType:      "sql",
		Db:          "0",
		Ts:          time.Now().UnixMilli(),
		Status:      0,
		Tolerated:   0,
		Frustrated:  0,
		Dur:         1,
		Outputs:     "1.1",
		BindValue:   "10,20",
		SqlErr:      0,
		Err:         "",
		ErrGroup:    "",
		ExceptionId: "",
		UserId:      "userid12345",
		SessionId:   "xxxxxxxsession",
		Ip:          "10.200.14.9",
		Province:    "",
		City:        "",
		IsOtel:      true,
	}

	localCache = NewCache()
	InitAppFilter(map[string]string{"LKL": "appid12345"})

	bts, _ := json.Marshal(sql)
	msg := &sarama.ConsumerMessage{
		Headers:        nil,
		Timestamp:      time.Time{},
		BlockTimestamp: time.Time{},
		Key:            nil,
		Value:          bts,
		Topic:          "sql",
		Partition:      0,
		Offset:         10,
	}
	pts, c := parseSQL(msg)
	t.Log(c)
	t.Logf("sql point=%s", pts[0].LineProto())
	/*
		bfy-sql-logging,
		agent_id=java-agent-id,
		app_id=appid12345,
		db=0,
		db_host=10.200.14.188,
		dbhost=10.200.14.188,
		dbtype=sql,
		event_id=eventid12345,
		group_id=78787,
		outputs=1.1,
		p_time=2024-06-06\ 11:28:41.76,
		project_id=LKL,
		span_id=4564564565456456,
		trace_id=60000000000000011111222222
		frustrated=0i,
		message="{dbhost:10.200.14.188,dbtype:sql,db:0,sql:appid12345-sql-78787,outputs:1.1,bind_value:10,20,dur:1}",
		status=ok,
		tolerated=0i
		1717644521759000000
	*/
}
