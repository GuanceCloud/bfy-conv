package parse

import (
	"fmt"
	"github.com/GuanceCloud/bfy-conv/mock"
	"testing"
)

func Test_parseTSpan(t *testing.T) {
	buf := mock.GetKafkaSpanByte()
	tspan, err := parseTSpan(buf[4:])
	if err != nil {
		t.Log(err)
	}

	t.Logf("tspan = %+v \n", tspan)
	t.Logf("tspan = %+v", tspan.ApiId)
}

func TestParseServiceType(t *testing.T) {
	ParseServiceType()

	fmt.Println(len(ServiceTypeMap))
}
