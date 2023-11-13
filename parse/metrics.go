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
	ctx := context.Background()
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
			cpukv = cpukv.Add([]byte("SystemCpuLoad"), cpuLoad.GetSystemCpuLoad(), false, false).
				Add([]byte("JvmCpuLoad"), cpuLoad.GetJvmCpuLoad(), false, false).
				AddTag([]byte("app_id"), []byte(appID)).
				AddTag([]byte("ip"), []byte(ip)).
				AddTag([]byte(projectKey), []byte(projectVal)).
				AddTag([]byte("agent_id"), []byte(agentID))
			pt := point.NewPointV2([]byte("agentStats-cpu"), cpukv, opts...)
			pts = append(pts, pt)
		}

		gc := stat.GetGc()
		if gc != nil {
			var gckvs point.KVs
			gckvs = gckvs.AddTag([]byte("app_id"), []byte(appID)).
				AddTag([]byte("agent_id"), []byte(agentID)).
				AddTag([]byte("ip"), []byte(ip)).
				AddTag([]byte(projectKey), []byte(projectVal)).
				Add([]byte("JvmMemoryHeapUsed"), gc.GetJvmMemoryHeapUsed(), false, false).
				Add([]byte("JvmMemoryHeapMax"), gc.GetJvmMemoryHeapMax(), false, false).
				Add([]byte("JvmMemoryNonHeapUsed"), gc.GetJvmMemoryNonHeapUsed(), false, false).
				Add([]byte("JvmMemoryNonHeapMax"), gc.GetJvmMemoryNonHeapMax(), false, false).
				Add([]byte("JvmGcOldCount"), gc.GetJvmGcOldCount(), false, false).
				Add([]byte("JvmGcOldTime"), gc.GetJvmGcOldTime(), false, false).
				Add([]byte("JvmMemoryNonHeapCommitted"), gc.GetJvmMemoryNonHeapCommitted(), false, false).
				Add([]byte("TotalPhysicalMemory"), gc.GetTotalPhysicalMemory(), false, false)

			if gc.GetJvmGcDetailed() != nil {
				detailed := gc.GetJvmGcDetailed()
				gckvs = gckvs.Add([]byte("GcNewCount"), detailed.GetJvmGcNewCount(), false, false).
					Add([]byte("PoolCodeCacheUsage"), detailed.GetJvmPoolCodeCacheUsage(), false, false).
					Add([]byte("GcNewTime"), detailed.GetJvmGcNewTime(), false, false).
					Add([]byte("PoolCodeCacheMax"), detailed.GetJvmPoolCodeCacheMax(), false, false).
					Add([]byte("PoolCodeCacheUsed"), detailed.GetJvmPoolCodeCacheUsed(), false, false).
					Add([]byte("PoolCodeCacheCommitted"), detailed.GetJvmPoolCodeCacheCommitted(), false, false).
					Add([]byte("PoolCodeCacheInit"), detailed.GetJvmPoolCodeCacheInit(), false, false).
					Add([]byte("PoolNewGenUsage"), detailed.GetJvmPoolNewGenUsage(), false, false).
					Add([]byte("PoolNewGenMax"), detailed.GetJvmPoolNewGenMax(), false, false).
					Add([]byte("PoolNewGenUsed"), detailed.GetJvmPoolNewGenUsed(), false, false).
					Add([]byte("PoolNewGenCommitted"), detailed.GetJvmPoolNewGenCommitted(), false, false).
					Add([]byte("PoolNewGenInit"), detailed.GetJvmPoolNewGenInit(), false, false).
					Add([]byte("PoolOldGenUsage"), detailed.GetJvmPoolOldGenUsage(), false, false).
					Add([]byte("PoolOldGenMax"), detailed.GetJvmPoolOldGenMax(), false, false).
					Add([]byte("PoolOldGenUsed"), detailed.GetJvmPoolOldGenUsed(), false, false).
					Add([]byte("PoolOldGenCommitted"), detailed.GetJvmPoolOldGenCommitted(), false, false).
					Add([]byte("PoolOldGenInit"), detailed.GetJvmPoolOldGenInit(), false, false).
					Add([]byte("PoolSurvivorSpaceUsage"), detailed.GetJvmPoolSurvivorSpaceUsage(), false, false).
					Add([]byte("PoolSurvivorSpaceMax"), detailed.GetJvmPoolSurvivorSpaceMax(), false, false).
					Add([]byte("PoolSurvivorSpaceUsed"), detailed.GetJvmPoolSurvivorSpaceUsed(), false, false).
					Add([]byte("PoolSurvivorSpaceCommitted"), detailed.GetJvmPoolSurvivorSpaceCommitted(), false, false).
					Add([]byte("PoolSurvivorSpaceInit"), detailed.GetJvmPoolSurvivorSpaceInit(), false, false).
					Add([]byte("PoolPermGenUsage"), detailed.GetJvmPoolPermGenUsage(), false, false).
					Add([]byte("PoolPermGenMax"), detailed.GetJvmPoolPermGenMax(), false, false).
					Add([]byte("PoolPermGenUsed"), detailed.GetJvmPoolPermGenUsed(), false, false).
					Add([]byte("PoolPermGenCommitted"), detailed.GetJvmPoolPermGenCommitted(), false, false).
					Add([]byte("PoolPermGenInit"), detailed.GetJvmPoolPermGenInit(), false, false).
					Add([]byte("PoolMetaspaceUsage"), detailed.GetJvmPoolMetaspaceUsage(), false, false).
					Add([]byte("PoolMetaspaceMax"), detailed.GetJvmPoolMetaspaceMax(), false, false).
					Add([]byte("PoolMetaspaceUsed"), detailed.GetJvmPoolMetaspaceUsed(), false, false).
					Add([]byte("PoolMetaspaceCommitted"), detailed.GetJvmPoolMetaspaceCommitted(), false, false).
					Add([]byte("PoolMetaspaceInit"), detailed.GetJvmPoolMetaspaceInit(), false, false)
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

func parseAgentInfo(buf []byte) (*server.TAgentInfo, error) {
	transport := &thrift.TMemoryBuffer{
		Buffer: bytes.NewBuffer(buf),
	}

	protocol := thrift.NewTCompactProtocolConf(transport, &thrift.TConfiguration{})
	info := server.NewTAgentInfo()
	ctx := context.Background()
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
