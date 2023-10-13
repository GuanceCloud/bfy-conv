package jvmparse

import (
	"bytes"
	"context"
	"github.com/GuanceCloud/bfy-conv/gen-go/server"
	"github.com/GuanceCloud/cliutils/point"
	"github.com/apache/thrift/lib/go/thrift"
	"time"
)

func ParseJVMMetrics(buf []byte) (*server.TAgentStatBatch, error) {
	transport := &thrift.TMemoryBuffer{
		Buffer: bytes.NewBuffer(buf),
	}

	protocol := thrift.NewTCompactProtocolConf(transport, &thrift.TConfiguration{})
	batch := server.NewTAgentStatBatch()
	ctx := context.Background()
	err := batch.Read(ctx, protocol)
	return batch, err
}

func StatBatchToPoints(batch *server.TAgentStatBatch) (pts []*point.Point) {
	pts = make([]*point.Point, 0)
	appID := batch.GetAppId()
	agentID := batch.GetAgentId()
	opts := point.DefaultMetricOptions()
	opts = append(opts, point.WithTime(time.UnixMilli(batch.GetStartTimestamp())))

	for _, stat := range batch.AgentStats {
		cpuLoad := stat.GetCpuLoad()
		if cpuLoad != nil {
			var cpukv point.KVs
			cpukv = cpukv.Add([]byte("SystemCpuLoad"), cpuLoad.GetSystemCpuLoad(), false, false).
				Add([]byte("JvmCpuLoad"), cpuLoad.GetJvmCpuLoad(), false, false).
				AddTag([]byte("app_id"), []byte(appID)).
				AddTag([]byte("agent_id"), []byte(agentID))
			pt := point.NewPointV2([]byte("agentStats-cpu"), cpukv, opts...)
			pts = append(pts, pt)
		}

		gc := stat.GetGc()
		if gc != nil {
			var gckvs point.KVs
			gckvs = gckvs.AddTag([]byte("app_id"), []byte(appID)).
				AddTag([]byte("agent_id"), []byte(agentID)).
				Add([]byte("JvmMemoryHeapUsed"), gc.GetJvmMemoryHeapUsed(), false, false).
				Add([]byte("JvmMemoryHeapMax"), gc.GetJvmMemoryHeapMax(), false, false).
				Add([]byte("JvmMemoryNonHeapUsed"), gc.GetJvmMemoryNonHeapUsed(), false, false).
				Add([]byte("JvmMemoryNonHeapMax"), gc.GetJvmMemoryNonHeapMax(), false, false).
				Add([]byte("JvmGcOldCount"), gc.GetJvmGcOldCount(), false, false).
				Add([]byte("JvmMemoryNonHeapCommitted"), gc.GetJvmMemoryNonHeapCommitted(), false, false).
				Add([]byte("TotalPhysicalMemory"), gc.GetTotalPhysicalMemory(), false, false)

			if gc.GetJdbcConnNum() != 0 {
				gckvs = gckvs.Add([]byte("JdbcConnNum"), gc.GetJdbcConnNum(), false, false)
			}
			gckvs = gckvs.Add([]byte("JvmGcOldCountNew"), gc.GetJvmGcOldCountNew(), false, false).
				Add([]byte("JvmGcOldCountNew"), gc.GetJvmGcOldCountNew(), false, false).
				Add([]byte("ThreadNum"), gc.GetThreadNum(), false, false)
			pt := point.NewPointV2([]byte("agentStats-gc"), gckvs, opts...)
			pts = append(pts, pt)
		}
		// trace:= stat.GetActiveTrace() dk 不支持该指标
	}

	return pts
}

func ParseAgentInfo(buf []byte) (*server.TAgentInfo, error) {
	transport := &thrift.TMemoryBuffer{
		Buffer: bytes.NewBuffer(buf),
	}

	protocol := thrift.NewTCompactProtocolConf(transport, &thrift.TConfiguration{})
	batch := server.NewTAgentInfo()
	ctx := context.Background()
	err := batch.Read(ctx, protocol)
	return batch, err
}

func ParseAgentEvent(buf []byte) (*server.TAgentEvent, error) {
	transport := &thrift.TMemoryBuffer{
		Buffer: bytes.NewBuffer(buf),
	}

	protocol := thrift.NewTCompactProtocolConf(transport, &thrift.TConfiguration{})
	batch := server.NewTAgentEvent()
	ctx := context.Background()
	err := batch.Read(ctx, protocol)
	return batch, err
}
