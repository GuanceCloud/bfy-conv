package parseV2

import (
	"github.com/GuanceCloud/cliutils/logger"
	"github.com/GuanceCloud/cliutils/point"
	"github.com/IBM/sarama"
	"strings"
)

const (
	ProjectKey = "project_id"
)

var (
	log       = logger.DefaultSLogger("bfy")
	appFilter *AppFilter
)

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

type AppFilter struct {
	Projects map[string][]string
}

func InitAppFilter(apps map[string]string) {
	if apps == nil || len(apps) == 0 {
		return
	}
	af := &AppFilter{
		Projects: make(map[string][]string),
	}
	for pname, anames := range apps {
		ns := strings.Split(anames, ",")
		af.Projects[pname] = ns
	}
	appFilter = af
}

func projectFilter(appId string) (projectID string) {
	if appFilter != nil {
		filter := false
		// 过滤 app 名称， 通过之后增加tag：project="project_name"
		for pName, appNames := range appFilter.Projects {
			for _, name := range appNames {
				if name == appId {
					projectID = pName
					filter = true
					break
				}
			}
		}
		if !filter {
			log.Debugf("can not find appId: %s in filters", appId)
			return
		}
	}
	return
}
