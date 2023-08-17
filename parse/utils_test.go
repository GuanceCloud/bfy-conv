package parse

import (
	"fmt"
	"strings"
	"testing"
)

func Test_xid(t *testing.T) {
	buf := []byte{0x00, 0x1a, 0x74, 0x79, 0x6a, 0x79, 0x67, 0x6c, 0x2d, 0x70, 0x74, 0x2d, 0x75, 0x61, 0x70, 0xfc, 0xf8, 0xf8, 0x92, 0xf1, 0x30, 0x54, 0x6d, 0x65, 0x73, 0x6f, 0x73, 0x2d, 0x62, 0x39, 0x39, 0x37, 0x38, 0x37, 0x33, 0x39, 0x2d, 0x30, 0x37, 0x0d, 0x2d, 0x34, 0x30, 0x35, 0x36, 0x2d, 0x38, 0x33, 0x65, 0x35, 0x2d, 0x61, 0x32, 0x37, 0x63, 0x34, 0x39, 0x64, 0x36, 0x39, 0x65, 0x38, 0x36, 0xaf, 0xae, 0xe9, 0x06}
	appid := "eurus-uw-pt"
	agent := "xytb-uw-pt-7d876c8d-7vjg6"

	result := xid(buf, appid, agent)

	fmt.Println(result)
	splitResult := strings.Split(result, "^")
	appID := splitResult[0]
	timestamp := splitResult[1]
	agentID := splitResult[2]
	sequence := splitResult[3]

	t.Logf("appID = %s , right is tyjygl-pt-uap \n timestamp =%s right is 1679640378492 \n agentID=%s right is mesos-b9978739-070d-4056-83e5-a27c49d69e86 \n sequence=%s right is 14309167", appID, timestamp, agentID, sequence)

}

//\x00\x1atyjygl-pt-uap\xfc\xf8\xf8\x92\xf10Tmesos-b9978739-070d-4056-83e5-a27c49d69e86\xaf\xae\xe9\x06

func TestXid(t *testing.T) {
	buf := make([]byte, 0)
	buf = append(buf, 0x00, 0x1a)
	buf = append(buf, []byte("tyjygl-pt-uap")...)
	buf = append(buf, 0xfc, 0xf8, 0xf8, 0x92, 0xf1)
	buf = append(buf, []byte("0Tmesos-b9978739-070d-4056-83e5-a27c49d69e86")...)
	buf = append(buf, 0xaf, 0xae, 0xe9, 0x06)
	appid := "eurus-uw-pt"
	agent := "xytb-uw-pt-7d876c8d-7vjg6"
	fmt.Println("--------------------")
	fmt.Println(string(buf))
	fmt.Println(buf)
	fmt.Println("--------------------")
	dst := "\x00\x1atyjygl-pt-uap\xfc\xf8\xf8\x92\xf10Tmesos-b9978739-070d-4056-83e5-a27c49d69e86\xaf\xae\xe9\x06"
	dstb := []byte(dst)
	fmt.Println("--------------------")
	fmt.Println(string(dstb))
	fmt.Println("--------------------")
	result := xid(buf, appid, agent)

	fmt.Println(result)
	splitResult := strings.Split(result, "^")
	appID := splitResult[0]
	timestamp := splitResult[1]
	agentID := splitResult[2]
	sequence := splitResult[3]

	t.Logf("appID = %s , right is tyjygl-pt-uap \n timestamp =%s right is 1679640378492 \n agentID=%s right is mesos-b9978739-070d-4056-83e5-a27c49d69e86 \n sequence=%s right is 14309167", appID, timestamp, agentID, sequence)
}
