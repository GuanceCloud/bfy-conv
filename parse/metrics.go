package parse

import (
	"bytes"
	"context"
	"fmt"
	"github.com/GuanceCloud/bfy-conv/gen-go/server"
	"github.com/GuanceCloud/cliutils/point"
	"github.com/apache/thrift/lib/go/thrift"
	"time"
)

/*
appId 取值 从 AgentInfo 中对应的值。

*/

func parseAgentStatBatch(buf []byte) (*server.TAgentStatBatch, error) {
	transport := &thrift.TMemoryBuffer{
		Buffer: bytes.NewBuffer(buf),
	}

	protocol := thrift.NewTCompactProtocolConf(transport, &thrift.TConfiguration{})
	batch := server.NewTAgentStatBatch()
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	err := batch.Read(ctx, protocol)
	return batch, err
}

func statBatchToPoints(batch *server.TAgentStatBatch) (pts []*point.Point) {
	pts = make([]*point.Point, 0)
	appID := batch.GetAppId()
	projectKey := "project"
	projectVal := ""
	if appFilter != nil {
		filter := false
		// 过滤 app 名称， 通过之后增加tag：project="project_name"
		for pName, appNames := range appFilter.Projects {
			for _, name := range appNames {
				if name == appID {
					projectVal = pName
					filter = true
					break
				}
			}
		}
		if !filter {
			log.Debugf("del applicationName %s", appID)
			return
		}
	}
	// todo  appID 过滤
	ip := findIPFromRedis(batch.GetAppId(), batch.GetAgentId())
	// todo 添加 IP
	agentID := batch.GetAgentId()

	for _, stat := range batch.AgentStats {
		opts := point.DefaultMetricOptions()
		opts = append(opts, point.WithTime(time.UnixMilli(stat.GetTimestamp())))
		cpuLoad := stat.GetCpuLoad()
		if cpuLoad != nil {
			var cpukv point.KVs
			cpukv = cpukv.Add("SystemCpuLoad", cpuLoad.GetSystemCpuLoad(), false, false).
				Add("JvmCpuLoad", cpuLoad.GetJvmCpuLoad(), false, false).
				AddTag("app_id", appID).
				AddTag("ip", ip).
				AddTag(projectKey, projectVal).
				AddTag("agent_id", agentID)
			pt := point.NewPointV2("agentStats-cpu", cpukv, opts...)
			pts = append(pts, pt)
		}

		gc := stat.GetGc()
		if gc != nil {
			var gckvs point.KVs
			gckvs = gckvs.AddTag("app_id", appID).
				AddTag("agent_id", agentID).
				AddTag("ip", ip).
				AddTag(projectKey, projectVal).
				Add("JvmMemoryHeapUsed", gc.GetJvmMemoryHeapUsed(), false, false).
				Add("JvmMemoryHeapMax", gc.GetJvmMemoryHeapMax(), false, false).
				Add("JvmMemoryNonHeapUsed", gc.GetJvmMemoryNonHeapUsed(), false, false).
				Add("JvmMemoryNonHeapMax", gc.GetJvmMemoryNonHeapMax(), false, false).
				Add("JvmGcOldCount", gc.GetJvmGcOldCount(), false, false).
				Add("JvmGcOldTime", gc.GetJvmGcOldTime(), false, false).
				Add("JvmMemoryNonHeapCommitted", gc.GetJvmMemoryNonHeapCommitted(), false, false).
				Add("TotalPhysicalMemory", gc.GetTotalPhysicalMemory(), false, false)

			if gc.GetJvmGcDetailed() != nil {
				detailed := gc.GetJvmGcDetailed()
				gckvs = gckvs.Add("GcNewCount", detailed.GetJvmGcNewCount(), false, false).
					Add("PoolCodeCacheUsage", detailed.GetJvmPoolCodeCacheUsage(), false, false).
					Add("GcNewTime", detailed.GetJvmGcNewTime(), false, false).
					Add("PoolCodeCacheMax", detailed.GetJvmPoolCodeCacheMax(), false, false).
					Add("PoolCodeCacheUsed", detailed.GetJvmPoolCodeCacheUsed(), false, false).
					Add("PoolCodeCacheCommitted", detailed.GetJvmPoolCodeCacheCommitted(), false, false).
					Add("PoolCodeCacheInit", detailed.GetJvmPoolCodeCacheInit(), false, false).
					Add("PoolNewGenUsage", detailed.GetJvmPoolNewGenUsage(), false, false).
					Add("PoolNewGenMax", detailed.GetJvmPoolNewGenMax(), false, false).
					Add("PoolNewGenUsed", detailed.GetJvmPoolNewGenUsed(), false, false).
					Add("PoolNewGenCommitted", detailed.GetJvmPoolNewGenCommitted(), false, false).
					Add("PoolNewGenInit", detailed.GetJvmPoolNewGenInit(), false, false).
					Add("PoolOldGenUsage", detailed.GetJvmPoolOldGenUsage(), false, false).
					Add("PoolOldGenMax", detailed.GetJvmPoolOldGenMax(), false, false).
					Add("PoolOldGenUsed", detailed.GetJvmPoolOldGenUsed(), false, false).
					Add("PoolOldGenCommitted", detailed.GetJvmPoolOldGenCommitted(), false, false).
					Add("PoolOldGenInit", detailed.GetJvmPoolOldGenInit(), false, false).
					Add("PoolSurvivorSpaceUsage", detailed.GetJvmPoolSurvivorSpaceUsage(), false, false).
					Add("PoolSurvivorSpaceMax", detailed.GetJvmPoolSurvivorSpaceMax(), false, false).
					Add("PoolSurvivorSpaceUsed", detailed.GetJvmPoolSurvivorSpaceUsed(), false, false).
					Add("PoolSurvivorSpaceCommitted", detailed.GetJvmPoolSurvivorSpaceCommitted(), false, false).
					Add("PoolSurvivorSpaceInit", detailed.GetJvmPoolSurvivorSpaceInit(), false, false).
					Add("PoolPermGenUsage", detailed.GetJvmPoolPermGenUsage(), false, false).
					Add("PoolPermGenMax", detailed.GetJvmPoolPermGenMax(), false, false).
					Add("PoolPermGenUsed", detailed.GetJvmPoolPermGenUsed(), false, false).
					Add("PoolPermGenCommitted", detailed.GetJvmPoolPermGenCommitted(), false, false).
					Add("PoolPermGenInit", detailed.GetJvmPoolPermGenInit(), false, false).
					Add("PoolMetaspaceUsage", detailed.GetJvmPoolMetaspaceUsage(), false, false).
					Add("PoolMetaspaceMax", detailed.GetJvmPoolMetaspaceMax(), false, false).
					Add("PoolMetaspaceUsed", detailed.GetJvmPoolMetaspaceUsed(), false, false).
					Add("PoolMetaspaceCommitted", detailed.GetJvmPoolMetaspaceCommitted(), false, false).
					Add("PoolMetaspaceInit", detailed.GetJvmPoolMetaspaceInit(), false, false)
			}
			gckvs = gckvs.Add("JvmGcOldCountNew", gc.GetJvmGcOldCountNew(), false, false).
				Add("JvmGcOldCountNew", gc.GetJvmGcOldCountNew(), false, false).
				Add("ThreadNum", gc.GetThreadNum(), false, false)
			pt := point.NewPointV2("agentStats-gc", gckvs, opts...)
			pts = append(pts, pt)
		}
		// trace:= stat.GetActiveTrace() dk 不支持该指标
	}

	return pts
}

func parseAgentInfo(buf []byte) (*server.TAgentInfo, error) {
	transport := &thrift.TMemoryBuffer{
		Buffer: bytes.NewBuffer(buf),
	}

	protocol := thrift.NewTCompactProtocolConf(transport, &thrift.TConfiguration{})
	info := server.NewTAgentInfo()
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	err := info.Read(ctx, protocol)
	if info != nil {
		err = storeIPToRedis(info.GetAppId(), info.GetAgentId(), info.GetIP())
	}
	return info, err
}

func parseAgentEvent(buf []byte) (*server.TAgentEvent, error) {
	transport := &thrift.TMemoryBuffer{
		Buffer: bytes.NewBuffer(buf),
	}

	protocol := thrift.NewTCompactProtocolConf(transport, &thrift.TConfiguration{})
	batch := server.NewTAgentEvent()
	ctx := context.Background()
	err := batch.Read(ctx, protocol)
	return batch, err
}

func storeIPToRedis(appID, agentID string, ip string) error {
	if pool != nil {
		if agentID != "" && appID != "" && ip != "" {
			RedigoSet(agentID+"|"+appID, ip)
			return nil
		}
	} else {
		return fmt.Errorf("redis conn is nil")
	}
	return nil
}

func findIPFromRedis(appID, agentID string) string {
	if pool != nil {
		return RedigoGet(agentID + "|" + appID)
	}
	return ""
}
