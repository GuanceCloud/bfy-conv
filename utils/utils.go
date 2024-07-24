package utils

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"
)

func Xid(buf []byte, appid string, agent string) string {
	i := 1
	appid_size, offset := binary.Varint(buf[i:])

	i += offset
	app_id := ""
	if appid_size == -1 {
		app_id = appid
	} else {
		if utf8.Valid(buf[i:(i + int(appid_size))]) {
			app_id = string(buf[i:(i + int(appid_size))])
		}
	}

	if appid_size == -1 {
		i += 0
	} else {
		i += int(appid_size)
	}

	timestamp, offset := binary.Uvarint(buf[i:])
	i += offset
	agent_size, offset := binary.Varint(buf[i:])
	i += offset
	agent_id := ""
	if agent_size == -1 || agent_size == 0 {
		agent_id = agent
	} else {

		if utf8.Valid(buf[i:(i + int(agent_size))]) {
			agent_id = string(buf[i:(i + int(agent_size))])
		} else {
			agent_id = agent
		}
	}

	if agent_size == -1 {
		i += 0
	} else {
		i += int(agent_size)
	}

	seq, _ := binary.Uvarint(buf[i:])
	appid = app_id

	return strings.Join([]string{
		appid,
		fmt.Sprintf("%d", timestamp),
		agent_id,
		fmt.Sprintf("%d", seq),
	}, "^")
}

func Code(buf []byte) (int, error) {
	// 检查协议长度必须为4
	if len(buf) <= 4 {
		return -1, errors.New("invalid Protocol Length")
	}

	// 检查协议签名必须为0xEF
	signature := buf[0]
	if signature != 0xEF {
		return -1, errors.New("invalid Protocol Signature")
	}

	code := 0
	if buf[2] == 0 && buf[3] == 0 {
		code = 0
	} else if buf[2] == 0 && buf[3] == 40 {
		code = 40
	} else if buf[2] == 0 && buf[3] == 70 {
		code = 70
	} else if buf[2] == 0 && buf[3] == 56 {
		code = 56
	} else if buf[2] == 0 && buf[3] == 57 {
		code = 57
	} else if buf[2] == 0 && buf[3] == 50 {
		code = 50
	} else {
		return -1, errors.New("invalid Protocol Code")
	}

	return code, nil
}

/*
func getTidFromHeader(header string, key string, xid string) string {
	if header == "" {
		return ""
	}
	headers := strings.Split(header, ";")

	for _, h := range headers {
		if strings.HasPrefix(h, key) {
			vals := strings.Split(h, ",")
			if len(vals) >= 2 {
				RedigoSet(xid, vals[1])
				return vals[1]
			}
		}
	}
	return ""
}

func getTidFromRedis(xid string) string {
	traceID := ""
	cachedTraceID := RedigoGet(xid)
	if utf8.Valid([]byte(cachedTraceID)) {
		traceID = cachedTraceID
	}

	return traceID
}
*/

/*func ServiceName(serviceType int16) string {
	if sts, ok := ServiceMaps[serviceType]; ok {
		return sts[0]
	}
	return "Unknown"
}

func SourceType(serviceType int16) string {
	if sts, ok := ServiceMaps[serviceType]; ok {
		return sts[1]
	}
	return "Unknown"
}*/

func GetRandomWithAll() int64 {
	rand.Seed(time.Now().UnixNano())
	return int64(rand.Intn(math.MaxInt))
}
