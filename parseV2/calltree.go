package parseV2

import (
	"encoding/json"
	"github.com/GuanceCloud/bfy-conv/pb/callTree"
	"github.com/GuanceCloud/bfy-conv/utils"
	"github.com/GuanceCloud/cliutils/point"
	"github.com/IBM/sarama"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"strconv"
	"time"
)

func parseCallTree(msg *sarama.ConsumerMessage) (traces []*point.Point, category point.Category) {
	ct := &callTree.CallTree{}
	err := proto.Unmarshal(msg.Value, ct)
	if err != nil {
		log.Errorf("callTree unmarshall err=%v", err)
		return
	}

	for _, event := range ct.GetCallevents() {
		var kvs = point.KVs{}
		opts := point.DefaultLoggingOptions()
		opts = append(opts, point.WithTime(time.UnixMilli(event.Ts)))
		start := (event.GetTs() + int64(event.GetStartElapsed())) * 1e3
		dur := 0
		if event.GetEndElapsed() != 0 {
			dur = int(event.GetEndElapsed()) * 1e3
		}
		spanID := strconv.FormatInt(rand.Int63(), 10) //nolint
		parentID := parentIDFromDepth(event, event.SpanId, traces)
		res := apiGet(event.AgentId, int(event.ApiId))

		kvs = kvs.Add("trace_id", event.TraceId, false, false).
			Add("parent_id", parentID, false, false).
			Add("span_id", spanID, false, false).
			AddTag("service", event.Appid).
			Add("resource", res, false, false).
			Add("operation", res, false, false).
			Add("start", start, false, false).
			Add("duration", dur, false, false).
			AddTag("status", GetStatus(int(event.Status))).
			AddTag("depth", strconv.Itoa(int(event.Depth))).
			AddTag("sequence", strconv.Itoa(int(event.Sequence))).
			AddTag("agent_id", event.AgentId).
			AddTag("agent_ip", event.AgentIp).
			AddTag("span_type", "local").
			AddTag("service_type", "bfy-tspan").
			AddTag("source", "kafka").
			AddTag("source_type", utils.SourceType(int16(event.ServiceType)))

		if event.Method != "" {
			kvs = kvs.AddTag("http_method", event.Method)
		}
		if event.Url != "" {
			kvs = kvs.AddTag("http_url", event.Url)
		}
		if event.Retcode != 0 {
			kvs = kvs.AddTag("http_status_code", strconv.Itoa(int(event.Retcode)))
		}

		bts, err := json.MarshalIndent(event, "", "  ")
		if err == nil {
			kvs = kvs.Add("message", string(bts), false, false)
		}
		pt := point.NewPointV2("bfy", kvs, opts...)

		traces = append(traces, pt)
	}

	return traces, point.Tracing
}

func parentIDFromDepth(event *callTree.CallEvent, parentID string, trace []*point.Point) string {
	switch event.Depth {
	case 1, -1:
		return parentID
	case 0:
		for j := len(trace) - 1; j >= 0; j-- {
			if depth := trace[j].GetTag("depth"); depth == "0" {
				return (trace[j].Get("parent_id")).(string)
			}
		}
	default:
		for j := len(trace) - 1; j >= 0; j-- {
			if depth := trace[j].GetTag("depth"); depth == strconv.Itoa(int(event.Depth)-1) {
				//return trace[j].GetFiledToString(itrace.FieldSpanid)
				return trace[j].Get("span_id").(string)
			}
		}
	}
	return parentID
}

func sourceType(st int32) (resource string, source_type string) {
	resource = "unknown"
	source_type = "custom"
	if st, ok := utils.ServiceTypeMap[int16(st)]; !ok {
		return
	} else {
		resource = st.Name
		if st.IsQueue {
			source_type = "message_queue"
		}

		if st.IsIncludeDestinationID == 1 {
			source_type = "db"
		}

		if st.IsRpcClient == 1 {
			source_type = "http"
		}

		if st.IsTerminal == 1 {
			source_type = "db"
		}
	}
	return
}
