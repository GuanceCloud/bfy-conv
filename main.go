package main

import (
	"github.com/GuanceCloud/bfy-conv/mock"
	"github.com/GuanceCloud/bfy-conv/parse"
	"github.com/GuanceCloud/cliutils/logger"
)

func main() {
	slog := logger.SLogger("hander")
	parse.SetLogger(slog)
	parse.InitRedis("", "", "", 0)
	defer func() { parse.StopRedis() }()

	pts := parse.Handle(mock.GetKafkaSpanByte())
	for _, pt := range pts {
		bts, err := pt.MarshalJSON()
		if err == nil {
			slog.Debugf("pts = %s", string(bts))
			slog.Debugf("=========")
		}
	}
}
