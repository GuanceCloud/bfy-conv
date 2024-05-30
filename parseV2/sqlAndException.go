package parseV2

import (
	"encoding/json"
	"github.com/GuanceCloud/cliutils/point"
	"github.com/IBM/sarama"
	"time"
)

/*
1 {
2   "id": "73639cc85e524ecc920795bb06277b6f",
3   "appid": "shop_38080_ot",
4   "appsysid": "dc6b396d-cbb8-44bd-859f-93ac4267225b",
5   "trxid": "shop_38080_ot^1715392942205^shop38080_ot^531",
6   "spanid": "3711720381322726698",
7   "group": "1028",
8   "agentid": "shop38080_ot",
9   "dbhost": "172.17.102.109:3306",
10   "dbtype": "MYSQL",
11   "db": "iwc_shop",
12   "ts": 1715406816589,
13   "status": 0,
14   "tolerated": 0,
15   "frustrated": 0,
16   "dur": 0,
17   "outputs": "1,1",
18   "bind_value": "1, 0, 0",
19   "sqlerr": 0,
20   "is_otel": true
}
*/

type SQL struct {
	ID          string `json:"id"`           // id
	Appid       string `json:"appid"`        // 应用id
	AppSysId    string `json:"appsysid"`     //应用系统id
	TrxId       string `json:"trxid"`        // 链路id
	TraceId     string `json:"trace_id"`     // -
	SpanId      string `json:"spanid"`       // spanid
	Group       string `json:"group"`        // 分组
	AgentId     string `json:"agentid"`      // 探针id
	DbHost      string `json:"dbhost"`       // 数据库地址
	DbType      string `json:"dbtype"`       // 数据库类型
	Db          string `json:"db"`           // 数据库名称
	Ts          int64  `json:"ts"`           // 时间戳
	Status      int    `json:"status"`       // 状态
	Tolerated   int    `json:"tolerated"`    // 缓慢数
	Frustrated  int    `json:"frustrated"`   // 极慢数
	Dur         int    `json:"dur"`          // 耗时
	Outputs     string `json:"outputs"`      // 输出
	BindValue   string `json:"bind_value"`   // 绑定的数据
	SqlErr      int    `json:"sqlerr"`       // 错误数量
	Err         string `json:"err"`          // 错误信息
	ErrGroup    string `json:"err_group"`    // 错误分组
	ExceptionId string `json:"exception_id"` // 异常id
	UserId      string `json:"user_id"`      // 用户id
	SessionId   string `json:"session_id"`   // 会话id
	Ip          string `json:"ip"`           // IP地址
	Province    string `json:"province"`     // 省
	City        string `json:"city"`         // 市
	IsOtel      bool   `json:"is_otel"`      // 是否ot
}

func parseSQL(msg *sarama.ConsumerMessage) (pts []*point.Point, category point.Category) {
	// id 与 event中的event_id一致
	sql := &SQL{}
	err := json.Unmarshal(msg.Value, sql)
	if err != nil {
		return
	}
	if !sql.IsOtel {
		return
	}
	// 过滤
	projectID := projectFilter(sql.Appid)
	if projectID == "" {
		return
	}

	opts := point.DefaultLoggingOptions()
	opts = append(opts, point.WithTime(time.UnixMilli(sql.Ts)))
	sqlStr := sqlGetFromCache(sql.Appid, sql.Group)
	// 组装行协议 logging
	var kvs point.KVs
	kvs = kvs.AddTag("group_id", sql.Group).
		AddTag("span_id", sql.SpanId).
		AddTag("event_id", sql.ID).
		AddTag("trace_id", sql.TraceId).
		AddTag("agent_id", sql.AgentId).
		AddTag("sql_template", sqlStr).
		AddTag("db", sql.Db).
		AddTag(ProjectKey, projectID).
		AddTag("db_host", sql.DbHost).
		Add("message", string(msg.Value), false, false)

	if sql.Outputs != "" {
		kvs = kvs.AddTag("outputs", sql.Outputs)
	}

	pt := point.NewPointV2("bfy-sql-logging", kvs, opts...)
	pts = append(pts, pt)
	return pts, point.Logging
}

/*
{
  "ts": 1715418618775,
  "name": "java.net.ConnectException",
  "type": 0,
  "modelid": [
    "7fabc61c-d5c9-47ef-8f1e-e0b3a34f38a9"
  ],
  "appid": "jvm_106",
  "appsysid": "Other",
   "message": "Failed to connect to /172.17.102.106:28080",
   "method": "connectSocket",
   "class": "java.net.ConnectException",
   "interface": "org.apache.catalina.connector.CoyoteInputStream.read(byte[] b, int off, int len)",
   "url": "/manage/login.action",
   "agentid": "jvm_106",
   "trxid": "jvm_106^1715418587165^jvm_106^1",
   "spanid": "7666444056924685434",
   "pspanid": "5199990252656806391",
   "exception_id": "0b47a8e889bae3ead28d19dcd29aac1df59effc4",
   "group": "a4a3dfb31a9cdceb810fa657d3d8d082e4276362",
   "depth": 0,
   "is_otel": true
 }
*/

type exception struct {
	Ts          int64    `json:"ts"`          // 时间戳 毫秒
	Name        string   `json:"name"`        // 异常名称
	Etype       int      `json:"type"`        // 异常类型
	ModelId     []string `json:"modelid"`     // 异常模型id
	AppId       string   `json:"appid"`       // 应用id
	AppSysId    string   `json:"appsysid"`    // 应用系统id
	Message     string   `json:"message"`     // 异常信息
	Method      string   `json:"method"`      // 方法
	Class       string   `json:"class"`       // 异常类
	Interface   string   `json:"interface"`   // 接口
	Url         string   `json:"url"`         // 请求url
	AgentId     string   `json:"agentId"`     // 探针id
	PagenId     string   `json:"pagenid"`     // 父探针id
	TrxId       string   `json:"trxid"`       // 链路id
	TraceId     string   `json:"trace_id"`    // 链路id
	SpanId      string   `json:"spanid"`      // span
	PspanId     string   `json:"pspanid"`     // 父span
	ExceptionId string   `json:"exceptionid"` // 错误id
	UserId      string   `json:"userid"`      // 用户id
	SessionId   string   `json:"sessionid"`   // 会话id
	IP          string   `json:"ip"`          // ip地址
	Group       string   `json:"group"`       // 异常分组
	Depth       int      `json:"depth"`       // 异常深度
	IsOtel      bool     `json:"is_otel"`     // 是否otel探针
	EventId     string   `json:"event_id"`    // 事件
}

// parseException: 通过thrift反序列化，后生成日志发送到中心。
func parseException(msg *sarama.ConsumerMessage) (pts []*point.Point, category point.Category) {
	e := &exception{}
	err := json.Unmarshal(msg.Value, e)
	if err != nil {
		log.Errorf("")
		return
	}
	if !e.IsOtel {
		return
	}
	// 过滤
	projectID := projectFilter(e.AppId)
	if projectID == "" {
		return
	}
	opts := point.DefaultLoggingOptions()
	opts = append(opts, point.WithTime(time.UnixMilli(e.Ts)))
	var kvs point.KVs
	kvs = kvs.AddTag("app_id", e.AppId).
		AddTag("name", e.Name).
		AddTag("method", e.Method).
		AddTag("class", e.Class).
		AddTag("interface", e.Interface).
		AddTag("url", e.Url).
		AddTag(ProjectKey, projectID).
		AddTag("event_id", e.EventId).
		Add("message", string(msg.Value), false, false)
	pt := point.NewPointV2("bfy-exception", kvs, opts...)
	pts = append(pts, pt)
	return pts, point.Logging
}
