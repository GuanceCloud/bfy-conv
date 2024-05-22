package parse

import "testing"

func TestInitRedis(t *testing.T) {
	InitRedis("10.200.14.188", "6379", "", 0)
	key := "bfy-cache-001"
	val := "span-test"
	RedigoSet(key, val)

	val_r := RedigoGet(key)
	if val_r != val {
		t.Errorf("not find val or val=%s", val)
	} else {
		t.Logf("ok val=%s", val)
	}
}
