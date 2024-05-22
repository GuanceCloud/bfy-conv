package parseV2

import (
	"github.com/GuanceCloud/cliutils/logger"
	"github.com/GuanceCloud/cliutils/point"
	"github.com/IBM/sarama"
)

var log = logger.DefaultSLogger("bfy")

func HandlerTopic(msg *sarama.ConsumerMessage) (pts []*point.Point, category point.Category) {
	switch msg.Topic {
	case "dwd_request":
		return request(msg)
	case "dwd_jvmstats":
		return JVMParse(msg)
	case "dwd_callevents":
		return parseCallTree(msg)
	case "dwd_sql":
		return parseSQL(msg)
	case "dwd_exception":
		return parseException(msg)
	case "ods_metadata":
		metadata(msg)
	default:
		log.Infof("unknown topic:%s", msg.Topic)
	}
	return
}

func SetLogging(l *logger.Logger) {
	log = l
}
