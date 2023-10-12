package jvmparse

import (
	"bytes"
	"context"
	"github.com/GuanceCloud/bfy-conv/gen-go/server"
	"github.com/GuanceCloud/cliutils/point"
	"github.com/apache/thrift/lib/go/thrift"
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

	for _, stat := range batch.AgentStats {
		cpuLoad := stat.GetCpuLoad()
		if cpuLoad != nil {
			var cpukv point.KVs
			cpukv.Add([]byte("SystemCpuLoad"), cpuLoad.GetSystemCpuLoad(), false, false)
			cpukv.Add([]byte("JvmCpuLoad"), cpuLoad.GetJvmCpuLoad(), false, false)
			pt := point.NewPointV2([]byte("agentStats-cpu"), cpukv, point.DefaultMetricOptions()...)
			pts = append(pts, pt)
		}

		gc := stat.GetGc()
		if gc != nil {
			var gckvs point.KVs
			gckvs.AddTag([]byte("app_id"), []byte(appID))
			gckvs.AddTag([]byte("agent_id"), []byte(agentID))
			gckvs.Add([]byte("JvmMemoryHeapUsed"), gc.GetJvmMemoryHeapUsed(), false, false)
			gckvs.Add([]byte("JvmMemoryHeapMax"), gc.GetJvmMemoryHeapMax(), false, false)
			gckvs.Add([]byte("JvmMemoryNonHeapUsed"), gc.GetJvmMemoryNonHeapUsed(), false, false)
			gckvs.Add([]byte("JvmMemoryNonHeapMax"), gc.GetJvmMemoryNonHeapMax(), false, false)

			gckvs.Add([]byte("JvmGcOldCount"), gc.GetJvmGcOldCount(), false, false)
			gckvs.Add([]byte("JvmMemoryNonHeapCommitted"), gc.GetJvmMemoryNonHeapCommitted(), false, false)
			gckvs.Add([]byte("TotalPhysicalMemory"), gc.GetTotalPhysicalMemory(), false, false)

			if gc.GetJdbcConnNum() != 0 {
				gckvs.Add([]byte("JdbcConnNum"), gc.GetJdbcConnNum(), false, false)
			}
			gckvs.Add([]byte("JvmGcOldCountNew"), gc.GetJvmGcOldCountNew(), false, false)
			gckvs.Add([]byte("JvmGcOldCountNew"), gc.GetJvmGcOldCountNew(), false, false)
			gckvs.Add([]byte("ThreadNum"), gc.GetThreadNum(), false, false)
			pt := point.NewPointV2([]byte("agentStats-gc"), gckvs, point.DefaultMetricOptions()...)
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
