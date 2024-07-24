package utils

import (
	"encoding/json"
	"fmt"
)

var ServiceTypeMap = map[int16]*serviceType{}

type serviceType struct {
	ID                     int    `json:"_id"`
	TypeID                 int    `json:"typeId"`
	TypeDesc               string `json:"typeDesc"`
	Name                   string `json:"name"`
	IsQueue                bool   `json:"isQueue"`
	IsIncludeDestinationID int    `json:"isIncludeDestinationId"`
	IsInternalMethod       int    `json:"isInternalMethod"`
	IsRecordStatistics     int    `json:"isRecordStatistics"`
	IsRpcClient            int    `json:"isRpcClient"`
	IsTerminal             int    `json:"isTerminal"`
	IsUnknown              int    `json:"isUnknown"`
	IsUser                 int    `json:"isUser"`
}

func ParseServiceType() {
	sts := make(map[int16]*serviceType)
	err := json.Unmarshal([]byte(serviceMap), &sts)
	if err != nil {
		fmt.Println(err)
	} else {
		ServiceTypeMap = sts
	}
}

func GetSourceType(id int16) string {
	st, ok := ServiceTypeMap[id]
	if !ok {
		return "unknown"
	}
	switch {
	case st.IsQueue:
		return "message_queue"
	case st.IsIncludeDestinationID == 1:
		return "db"
	case st.IsRecordStatistics == 1:
		return "custom"
	case st.IsInternalMethod == 1:
		return "custom"
	case st.IsRpcClient == 1:
		return "http"
	case st.IsUser == 1:
		return "custom"
	case st.IsTerminal == 1:
		return "db"
	case st.IsUnknown == 1:
		return "unknown"
	default:

		return "unknown"
	}
}

func GetServiceName(id int16) string {
	st, ok := ServiceTypeMap[id]
	if !ok {
		return "unknown"
	}
	return st.Name
}
