# 新结构对应关系梳理

- Request 对应之前pp结构的Tspan ， `Topic：dwd_request` json 格式
- CallTree 对应 tspanChunk ， `Topic：dwd_callevents` proto 格式
- sql数据 是对callTree中event补充 ， `Topic：dwd_sql`  json
- exception数据 也是对Event的补充 ， `Topic：dwd_exception` json
- metaData 对应的是pp中的meta，是sql和api字符串 ， `Topic：dwd_metadata` thrift 格式
- jvm 生成指标数据 ， `Topic：dwd_jvmstats` json

## 对应关系

| OT 链路 中字段名称                                   | bfy字段             | 对应关系说明                                                          |
|:----------------------------------------------|:------------------|:----------------------------------------------------------------|
| Resource.Attributes 中 `service.name`          | appid             | request 中appid作为服务名                                             |
| trace_id                                      | trace_id          | request中的trace_id, event中的trace_id                              |
| span_id                                       | span_id 或者eventId | request中使用spanid。但是 event中spanid对应的是request中的spanid，需要使用eventId |
| parent_span_id                                | pspanid 或者 无      | request中使用pspanid，event中需要根据depth确认上下级关系后使用上一级的spanid作为父spanid  |
| name                                          | apiid             | event根据apiid，request中使用URL                                      |
| kind                                          | 待定                | -                                                               |
| start_time_unix_nano                          | start             | 起始时间是毫秒，需要 *1e6                                                 |
| end_time_unix_nano                            | dur               | dur加起始时间就是endtime                                               |
| Status                                        | status            | 0是正常，非0就是异常                                                     |
| Attributes                                    | 很多字段              | 不一一列举                                                           |


> 补充：上面的spanid 在request和event中存在且同一链路中相同，只能作为顶层spanid使用，子span（event list）中需要随机生成一个。

## 流程梳理

四个Topic：

### 1 dwd_request
Request 数据。 在接收到之后填充 `apiid` 对应的`meta` 数据。

### 2 dwd_sql
sql 数据用来生成日志， 通过 `traceId` 与链路关联。还有一个叫 event_cid 与span_id trace_id性质一致 这里不采用。

在生成日志之前，通过group字段与`sqlMeta`中的hash关联查找对用的sql语句，并填充到日志`sql_template`中。

### 3 dwd_exception

错误信息。同上 与 sql逻辑一致。

### 4 dwd_callevents 
CallTree 数据。由 `CallEvent` 数组组成。 spanId作为第一个Event的`parent_span_id`，并根据depth判断event的父子级关系。

每个`event` 通过 `trace_id` 与 `Request` 关联。也可以通过event_cid与request关联（等同于spanId和traceId）。

通过 event中id 又与sql或者异常的id关联。

### 5 meta
元数据。 元数据在Agent启动时后发送`api meta`信息到收集端，而后 会断断续续发送个别的元数据出来。

消费和处理逻辑：

多台消费者共同消费（一个group）元数据，将数据发送到缓存（暂定Redis），并在本地内存中存一份。
这样的好处的，不浪费消费资源情况下做到消息同步。存储的格式为：key："Agentid + string/sql + id" 对应的val为 一个string（未定，消息数据太冗余）。

在处理 `Request` 或者 `Event` 时,遇到 apiid 可从本地内存查找，如果没有从reids查找并放本地内存一份，如没有，则将id上传。

## 对应关系
|           | Request           | CallTree                 | sql                                     | Exception | sql Meta | func Meta |
|-----------|-------------------|--------------------------|-----------------------------------------|-----------|----------|-----------|
| Request   |                   | spanid,traceId,event_cid | spanId traceId agentId userId sessionId |           | apiID    | apiID     |
| CallTree  | spanId            |                          | eventId                                 |           | group    | group     |
| Sql       | spanID traceId    | EventId                  |                                         |           | group    | group     |
| Exception | userId,session_id | EventId                  |                                         |           |          |           |
| sql Meta  | apiId             | EventId                  | group                                   |           |          |           |
| func Meta | apiId             |                          |                                         |           |          |           |


