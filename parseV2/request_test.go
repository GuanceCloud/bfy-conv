package parseV2

import "testing"

func TestRequest(t *testing.T) {
	sendToKafka("dwd_request", []byte(requestTestData), t)
}
