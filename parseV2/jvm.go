package parseV2

import (
	"encoding/json"
	"github.com/GuanceCloud/cliutils/point"
	"github.com/IBM/sarama"
	"strconv"
	"time"
)

/*
1 {
2   "appid": "55mockserver10_ot",
3   "appsysid": "Other",
4   "agentid": "55mockserver10_agent_ot",
5   "agent_version": 3004000,
6   "ts": 1715158665322,
7   "jvm_cpu": 6,
8   "sys_cpu": 8495,
9   "gc_type": "PARALLEL",
10   "fgc_count": 0,
11   "fgc_time": 0,
12   "gc_old_count": 3,
13   "gc_old_time": 491,
14   "gc_old_count_new": 0,
15   "gc_old_time_new": 0,
16   "total_physical_memory": 14974713856,
17   "jdbc_conn_num": 0,
18   "thread_num": 0,
19   "heap_max": 14974713856,
20   "heap_used": 68414176,
21   "non_heap_max": 0,
22   "non_heap_used": 109828776,
23   "non_heap_committed": 0,
24   "non_heap_total": 0,
25   "new_gen_max": 5606735872,
26   "new_gen_used": 22140536,
27   "old_gen_max": 11231297536,
28   "old_gen_used": 45225064,
29   "survivor_space_max": 4194304,
30   "survivor_space_used": 1048576,
31   "metaspace_max": 70672384,
32   "metaspace_used": 69416912,
33   "metaspace_total": 70672384,
34   "perm_gen_usage": 0.0
35 }
*/

type JVM struct {
	AppID                  string  `json:"appid"`                    // 应用id
	AppsysID               string  `json:"appsysid"`                 // 应用系统id
	AgentID                string  `json:"agentid"`                  // agent id
	AgentVersion           int     `json:"agent_version"`            // agent version
	TS                     int64   `json:"ts"`                       // 时间 毫秒单位
	JvmCpu                 int     `json:"jvm_cpu"`                  // jvm cpu占用率
	SysCpu                 int     `json:"sys_cpu"`                  // 系统cpu占用率
	GcType                 string  `json:"gc_type"`                  // gc 类型
	FgcCount               int     `json:"fgc_count"`                // full gc 次数
	FgcTime                int     `json:"fgc_time"`                 // full gc 时间累计值
	GcOldCount             int     `json:"gc_old_count"`             // gc 次数累计值
	GcOldTime              int     `json:"gc_old_time"`              // gc 时间累加值
	GcOldCountNew          int     `json:"gc_old_count_new"`         // 新探针的gc次数累加值
	GcOldTimeNew           int     `json:"gc_old_time_new"`          // 新探针的gc时间累加值
	TotalPhysicalMemory    int     `json:"total_physical_memory"`    // 物理内存总量
	JdbcConnNum            int     `json:"jdbc_conn_num"`            // jdbc连接数
	ThreadNum              int     `json:"thread_num"`               // 线程数
	HeapMax                int     `json:"heap_max"`                 // 最大堆内存
	HeapUsed               int     `json:"heap_used"`                // 堆内存的使用量
	NonHeapMax             int     `json:"non_heap_max"`             // 最大非堆内存
	NonHeapUsed            int     `json:"non_heap_used"`            // 非堆内存使用量
	NonHeapCommitted       int     `json:"non_heap_committed"`       // 已提交的非堆内存
	NonHeapTotal           int     `json:"non_heap_total"`           // 非堆内存总量
	NewGenMax              int     `json:"new_gen_max"`              // 最大年轻代内存
	NewGenUsed             int     `json:"new_gen_used"`             // 年轻代使用量
	NewGenCommitted        int     `json:"new_gen_committed"`        // 已提交的年轻代内存
	OldGenMax              int     `json:"old_gen_max"`              // 最大老年代内存
	OldGenUsed             int     `json:"old_gen_used"`             // 老年代使用量
	OldGenCommitted        int     `json:"old_gen_committed"`        // 已提交的老年代内存
	SurvivorSpaceMax       int     `json:"survivor_space_max"`       // 最大幸存区内存
	SurvivorSpaceUsed      int     `json:"survivor_space_used"`      // 幸存区内存使用量
	SurvivorSpaceCommitted int     `json:"survivor_space_committed"` // 幸存区内存使用量
	MetaspaceMax           int     `json:"metaspace_max"`            // 最大元空间的值
	MetaspaceUsed          int     `json:"metaspace_used"`           // 已用的元空间的值
	MetaspaceCommitted     int     `json:"metaspace_committed"`      // 已提交的元空间的值
	MetaspaceTotal         int     `json:"metaspace_total"`          // 元空间的总量
	PermGenMax             int     `json:"perm_gen_max"`             // 最大永久代内存
	PermGenUsed            int     `json:"perm_gen_used"`            // 已使用永久代内存
	PermGenUsage           float64 `json:"perm_gen_usage"`           // 可用的永久代内存
	PermGenCommitted       int     `json:"perm_gen_committed"`       // 已提交的永久代内存
	CodeCacheCommitted     int     `json:"code_cache_committed"`     // 已提交的代码缓存区内存
	IoIdle                 float64 `json:"io_idle"`                  // io空闲时间百分比
}

func (jvm *JVM) ToPoint() *point.Point {
	if jvm == nil {
		return nil
	}
	projectID := projectFilter(jvm.AppID)
	if projectID == "" {
		log.Warnf("filter: can find projectID for %s", jvm.AppID)
		return nil
	}

	opts := point.DefaultMetricOptions()
	opts = append(opts, point.WithTime(time.UnixMilli(jvm.TS)))
	var kvs point.KVs
	kvs = kvs.AddTag("appid", jvm.AppID).
		AddTag("appsysid", jvm.AppsysID).
		AddTag(ProjectKey, projectID).
		AddTag("agent_id", jvm.AgentID).
		AddTag("agent_version", strconv.Itoa(jvm.AgentVersion)).
		AddTag("gc_type", jvm.GcType).
		Add("jvm_cpu", jvm.JvmCpu, false, false).
		Add("sys_cpu", jvm.SysCpu, false, false).
		Add("fgc_count", jvm.FgcCount, false, false).
		Add("fgc_time", jvm.FgcTime, false, false).
		Add("gc_old_count", jvm.GcOldCount, false, false).
		Add("gc_old_time", jvm.GcOldTime, false, false).
		Add("gc_old_count_new", jvm.GcOldCountNew, false, false).
		Add("gc_old_time_new", jvm.GcOldTimeNew, false, false).
		Add("total_physical_memory", jvm.TotalPhysicalMemory, false, false).
		Add("jdbc_conn_num", jvm.JdbcConnNum, false, false).
		Add("thread_num", jvm.ThreadNum, false, false).
		Add("heap_max", jvm.HeapMax, false, false).
		Add("heap_used", jvm.HeapUsed, false, false).
		Add("non_heap_max", jvm.NonHeapMax, false, false).
		Add("non_heap_used", jvm.NonHeapUsed, false, false).
		Add("non_heap_committed", jvm.NonHeapCommitted, false, false).
		Add("non_heap_total", jvm.NonHeapTotal, false, false).
		Add("new_gen_max", jvm.NewGenMax, false, false).
		Add("new_gen_used", jvm.NewGenUsed, false, false).
		Add("new_gen_committed", jvm.NewGenCommitted, false, false).
		Add("old_gen_max", jvm.OldGenMax, false, false).
		Add("old_gen_used", jvm.OldGenUsed, false, false).
		Add("old_gen_committed", jvm.OldGenCommitted, false, false).
		Add("survivor_space_max", jvm.SurvivorSpaceMax, false, false).
		Add("survivor_space_used", jvm.SurvivorSpaceUsed, false, false).
		Add("survivor_space_committed", jvm.SurvivorSpaceCommitted, false, false).
		Add("metaspace_max", jvm.MetaspaceMax, false, false).
		Add("metaspace_used", jvm.MetaspaceUsed, false, false).
		Add("metaspace_committed", jvm.MetaspaceCommitted, false, false).
		Add("metaspace_total", jvm.MetaspaceTotal, false, false).
		Add("perm_gen_max", jvm.PermGenMax, false, false).
		Add("perm_gen_used", jvm.PermGenUsed, false, false).
		Add("perm_gen_usage", jvm.PermGenUsage, false, false).
		Add("perm_gen_committed", jvm.PermGenCommitted, false, false).
		Add("code_cache_committed", jvm.CodeCacheCommitted, false, false).
		Add("io_idle", jvm.IoIdle, false, false)

	pt := point.NewPointV2("bfy-jvm", kvs, opts...)

	return pt
}

func JVMParse(msg *sarama.ConsumerMessage) (pts []*point.Point, category point.Category) {
	jvm := &JVM{}

	err := json.Unmarshal(msg.Value, jvm)
	if err != nil {
		log.Errorf("unmarshal jvm err=%v", err)
		return
	}

	pt := jvm.ToPoint()
	if pt == nil {
		log.Warnf("point is nil")
		return
	}
	pt.AddTag("topic", msg.Topic)
	pts = append(pts, pt)
	return pts, point.Metric
}
