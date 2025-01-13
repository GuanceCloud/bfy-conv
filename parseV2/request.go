package parseV2

import (
	"encoding/json"
	"github.com/GuanceCloud/bfy-conv/utils"
	"github.com/GuanceCloud/cliutils/point"
	"github.com/IBM/sarama"
	"strconv"
	"time"
)

type BizData struct {
	Cip   string `json:"cip"`
	Cid   string `json:"cid"`
	Cname string `json:"cname"`
}

type RequestData struct {
	Model           string  `json:"model,omitempty"`             // 模型名称
	Group           string  `json:"group,omitempty"`             // 分组
	Trxid           string  `json:"trxid,omitempty"`             // xid
	TraceID         string  `json:"trace_id,omitempty"`          // 链路id
	SpanID          string  `json:"span_id,omitempty"`           // span id
	PSpanID         string  `json:"pspan_id,omitempty"`          // 上级span id
	IPAddr          string  `json:"ip_addr,omitempty"`           // ip 地址
	Host            string  `json:"host,omitempty"`              // 域名
	Path            string  `json:"path,omitempty"`              // 请求路径
	AgentID         string  `json:"agent_id,omitempty"`          // 探针 id
	AgentIP         string  `json:"agent_ip,omitempty"`          // 探针 IP
	RootAppID       string  `json:"root_appid,omitempty"`        // 根节点 appid
	PAppSysID       string  `json:"pappsysid,omitempty"`         // 上级应用系统id
	PAppID          string  `json:"pappid,omitempty"`            // 上级应用id
	PAppType        int     `json:"papp_type,omitempty"`         // 上级应用类型
	RetCode         int     `json:"ret_code,omitempty"`          // 返回状态码
	Method          string  `json:"method,omitempty"`            // 请求方法
	ModelID         string  `json:"modelid,omitempty"`           // 模型id
	URL             string  `json:"url,omitempty"`               // 请求url 请求参数
	ServiceType     int     `json:"service_type,omitempty"`      // 服务类型 从serviceType表中取值
	APPServiceType  int     `json:"app_service_type,omitempty"`  // 应用类型 从serviceType表中取值
	AppID           string  `json:"appid,omitempty"`             // 应用id
	AppSysID        string  `json:"appsysid,omitempty"`          // 应用系统id
	PAgentID        string  `json:"pagent_id,omitempty"`         // 上级探针id
	PAgentIP        string  `json:"pagent_ip,omitempty"`         // 上级应用 IP
	Status          int     `json:"status,omitempty"`            // 状态
	Err4xx          int     `json:"err_4xx,omitempty"`           // 400~499 错误次数
	Err5xx          int     `json:"err_5xx,omitempty"`           // 500~ 错误次数
	Tag             string  `json:"tag,omitempty"`               // 标签
	Code            string  `json:"code,omitempty"`              // 业务失败吗
	Dur             int     `json:"dur,omitempty"`               // 响应时间
	Header          string  `json:"header,omitempty"`            // 请求头
	Body            string  `json:"body,omitempty"`              // 请求体
	ResHeader       string  `json:"res_header,omitempty"`        // 请求头
	ResBody         string  `json:"res_body,omitempty"`          // 响应体
	UEventModel     string  `json:"uevent_model,omitempty"`      // 用户行为模型
	UEventID        string  `json:"uevent_id,omitempty"`         // 用户行为id
	UserID          string  `json:"user_id,omitempty"`           // 用户id
	SessionID       string  `json:"session_id,omitempty"`        // 会话id
	Province        string  `json:"province,omitempty"`          // 省
	City            string  `json:"city,omitempty"`              // 市
	BizData         BizData `json:"biz_data,omitempty"`          // 业务数据
	PageID          string  `json:"page_id,omitempty"`           // 页面id
	PageGroup       string  `json:"page_group,omitempty"`        // 页面分组
	APIID           int     `json:"api_id,omitempty"`            // 请求接口id
	Exception       int     `json:"exception,omitempty"`         // 是否错误
	Type            string  `json:"type,omitempty"`              // 请求类型
	Ts              int64   `json:"ts,omitempty"`                // 时间戳
	Browser         string  `json:"browser,omitempty"`           // 浏览器信息
	Device          string  `json:"device,omitempty"`            // 设备信息
	OS              string  `json:"os,omitempty"`                // 系统信息
	OSVersion       string  `json:"os_version,omitempty"`        // 系统版本
	OSVersionNumber int     `json:"os_version_number,omitempty"` // 系统版本号
	IsOtel          bool    `json:"is_otel,omitempty"`           // 是否ot协议
	EventCID        string  `json:"event_cid,omitempty"`         // 是否ot协议
}

func request(msg *sarama.ConsumerMessage) (pts []*point.Point, category point.Category) {
	req := &RequestData{}
	err := json.Unmarshal(msg.Value, req)
	if err != nil {
		return
	}
	if !req.IsOtel {
		// 不是ot协议的退出
		return
	}
	// 过滤
	projectID := projectFilter(req.AppID)
	if projectID == "" {
		return
	}

	opts := point.CommonLoggingOptions()
	opts = append(opts, point.WithTime(time.UnixMilli(req.Ts)))

	var kvs point.KVs
	kvs = kvs.Add("trace_id", req.TraceID, false, false).
		Add("span_id", req.SpanID, false, false).
		Add("parent_id", req.PSpanID, false, false).
		AddTag("service", req.AppID).
		Add("resource", req.Path, false, false).
		Add("operation", req.Path, false, false).
		Add("start", req.Ts*1e3, false, false).
		Add("duration", req.Dur*1e3, false, false).
		AddTag("status", GetStatus(req.Status)).
		AddTag("service_type", "").
		AddTag("agent_id", req.AgentID).
		AddTag("agent_ip", req.AgentIP).
		AddTag("pappid", req.PAppID).
		AddTag("http_method", req.Method).
		AddTag("http_url", req.URL).
		AddTag("rpc_route", req.URL).
		AddTag("http_status_code", strconv.Itoa(req.RetCode)).
		AddTag("span_type", "entry").
		AddTag("user_id", req.UserID).
		AddTag("session_id", req.SessionID).
		AddTag(ProjectKey, projectID).
		AddTag("service_type", "bfy-tspan").
		AddTag("source", "kafka").
		AddTag("source_type", utils.GetSourceType(int16(req.ServiceType))).
		Add("message", string(msg.Value), false, false)
	pt := point.NewPointV2("bfy", kvs, opts...)
	pts = append(pts, pt)
	return pts, point.Tracing
}

func GetStatus(status int) string {
	switch status {
	case 0:
		return "ok"
	default:
		return "error"
	}
}
