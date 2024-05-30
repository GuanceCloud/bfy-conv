# bfy-conv
目前已经集成在 DataKit 中，功能是 负责解析从kafka收到的消息转成行协议数据，其中包括：链路，jvm指标。

kafka 中的数据为 bfy 二开 Pinpoint 链路工具，再将span和SpanChunk进行序列化之后发送的。对比原始的 Pinpoint 结构，其中很多tag中的标签都提取到一级字段中。

## thrift
根据 thrift 文件生成 go 结构体。
```shell
thrift --gen go span.thrift
```

## trace

根据 agentid 和 transactionID 从redis中获取trace_id，或者从请求的头部信息中根据 B3 协议获取。再存在redis中。

1.3.3 版本以后，tSpan 中的 traceparent 为标准 W3C 协议格式，`trace_id` 和  `span_id` 都是从其中获取。

1.4.x 版本为新结构，[parseV2](./parseV2/readme.md)