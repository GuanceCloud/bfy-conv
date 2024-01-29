# bfy-conv
负责解析从kafka收到的消息成行协议，其中包括：链路，jvm指标。

## thrift
根据 thrift 文件生成 go 结构体。
```shell
thrift --gen go span.thrift
```

## trace

tSpan 中的 traceparent 中