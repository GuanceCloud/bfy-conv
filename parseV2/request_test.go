package parseV2

import (
	"encoding/json"
	"testing"
)

var requestTestData = `
{
  "model": "/oaNotify-cont",
  "group": "451DC0F25922FD15",
  "trxid": "shop_web^1715024901017^0kQfJK23Iu^2351",
  "trace_id": "000000000000092f0d8ca3e554b3a735",
  "span_id": "174155503488541",
  "pspan_id": "522207364658462",
  "ip_addr": "172.17.0.128",
  "host": "172.17.0.120:38080",
  "path": "/admin/oa/oaNotify/self/count",
  "agent_id": "shop38080_ot",
  "agent_ip": "172.17.0.120",
  "root_appid": "shop_web",
  "pappsysid": "dc6b396d-cbb8-44bd-859f-93ac4267225b",
  "pappid": "shop_web",
  "papp_type": 10,
  "ret_code": 200,
  "method": "GET",
  "modelid": "451DC0F25922FD15",
  "url": "updateSession\u003d0\u0026t\u003d1715024901017",
  "service_type": 1010,
  "appid": "shop_38080_ot",
  "appsysid": "dc6b396d-cbb8-44bd-859f-93ac4267225b",
  "pagent_id": "",
  "pagent_ip": "",
  "status": 0,
  "err_4xx": 0,
  "err_5xx": 0,
  "tag": null,
  "code": "⚊NULL⚊",
  "dur": 7,
  "header": "host,172.17.0.120:38080;connection,close;apptrace-pappname,shop_web;apptrace-traceid,shop_web^1715024901017^0kQfJK23Iu^2351;apptrace-sid,8171f6da-ed6c-4664-8c5a-9c1dec8c7e91;user-agent,Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36;accept,*/*;apptrace-spanid,174155503488541;apptrace-papptype,10;apptrace-uid,b8ac9134-e69b-4b82-8bac-d1e033238eb0;pagentid,websdk;x-requested-with,XMLHttpRequest;apptrace-pspanid,522207364658462;referer,http://172.17.0.128:18884/admin?login;accept-encoding,gzip, deflate;accept-language,zh-CN,zh;q\u003d0.9;cookie,rememberMe\u003dSN1hXxMwrUUAJ3s5AYMvSOaTVwKk/Ag4+9D7V4qWtiMEmb74n2teaOwihrvTJmNEmIT9PABAd1o3GnIDOvXzNdJUoo4/8Hq27xPwZhblzkH0QyZm4CL6glBSLh2C8thbkwftXKC4E2dSJOPFFwHyCX7YqKGVRjZzKeJIv7dX1+pUxw9MnYvf/GcTnpKPKGuMwBj6tOEIPAzDuSxPAPYd2V8bPnLZjvNoo012m8kTcCejUhTFsMzWmSM7UbhVbX5/OywMWoPYSVNPcXyrUo4UjcsglHLwiP2LpE3dF4m68ZSDqvtIi6pwkx6T0Flv/h1KfQz+8AWwO0DVWGR7hi2IaA65nbY4mq2lJCy82junaOS/IHr0JnFf3UyuwXf+g53ONu7q634xZLVl2EnCDL7Muz2hfc9vk6w0yx2NQgliD0BZD6cL0mHjrmInFJeU0poQEB7kFsXuYSV0sK10/TDBXsSzrLodjlLzUZ8QRnPYIoIbU8SuJ5gNqzdu/7ewD8mOwqYoqR/vf7zvc5Qzgbrlxp/QEhnXeKzJCH7yV9r0+B1X/ixPu4N9544EfeKdQYu9vANkY2ubK0BkqbocqHOmnMyrPBmu8aW7OdaVBWXEt9Fw7pqM6tHyhpmstO6vL3QJJtWK24sJK8yQMqIq8pgF0XYJe/25OSZco0nBBBXFExECqGQA1z+y5GxQaGEkmZeYzV+PWy8O0QsGXn4i1JpNftIGy/pYkMw2B7S9IhcZ36ZjO+WsbP6CBhQYskA8/oWPww2knWfSv/4XabYSpz6NDfIVv6daiWkiQkbnI2dsv/E\u003d; pdomain\u003d172.17.0.128; JSESSIONID\u003d5C96AF9DC1CAD319C9A2EE4D94E9CEBA; ppageid\u003d874b42175a5f68a8517cdaa88145fa74; puri\u003d/admin/sys/user/info; shop.session.id\u003deb6ab80292364bcc9220b84ca80fcef2",
  "body": "",
  "res_header": "Accept-Charset,[Accept-Charset, Content-Type, Content-Length, Date, Connection];Connection,[Accept-Charset, Content-Type, Content-Length, Date, Connection];Content-Length,[Accept-Charset, Content-Type, Content-Length, Date, Connection];Date,[Accept-Charset, Content-Type, Content-Length, Date, Connection];Content-Type,[Accept-Charset, Content-Type, Content-Length, Date, Connection]",
  "res_body": "",
  "uevent_model": null,
  "uevent_id": null,
  "user_id": "b8ac9134-e69b-4b82-8bac-d1e033238eb0",
  "session_id": "8171f6da-ed6c-4664-8c5a-9c1dec8c7e91",
  "province": "Unknown",
  "city": "Unknown",
  "biz_data": {
    "cip": "127.0.0.1",
    "cid": "00",
    "cname": "未知"
  },
  "page_id": "7d92448a21049f10ffa28719f08426cb",
  "page_group": "0611b3ff18d459bf30ba82c89278d9e3a961d994",
  "api_id": 4,
  "exception": 0,
  "type": "server",
  "ts": 1715024901050,
  "browser": "chrome/124.0.0.0",
  "device": "Windows x64",
  "os": "Win10",
  "os_version": "10",
  "os_version_number": 10000000,
  "is_otel": true
}
`

func TestHandle(t *testing.T) {
	rd := &RequestData{}
	err := json.Unmarshal([]byte(requestTestData), rd)
	if err != nil {
		t.Errorf("unmarshal err%v", err)
		return
	}

	t.Logf("request data=%+v", rd)
	//{Model:/oaNotify-cont
	//Group:451DC0F25922FD15
	//Trxid:shop_web^1715024901017^0kQfJK23Iu^2351 TraceID:000000000000092f0d8ca3e554b3a735 SpanID:174155503488541
	//PSpanID:522207364658462 IPAddr:172.17.0.128 Host:172.17.0.120:38080 Path:/admin/oa/oaNotify/self/count
	//AgentID:shop38080_ot AgentIP:172.17.0.120 RootAppID:shop_web
	//PAppSysID:dc6b396d-cbb8-44bd-859f-93ac4267225b PAppID:shop_web PAppType:10 RetCode:200 Method:GET ModelID:451DC0F25922FD15 URL:updateSession=0&t=1715024901017
	//ServiceType:1010 AppID:shop_38080_ot AppSysID:dc6b396d-cbb8-44bd-859f-93ac4267225b PAgentID: PAgentIP: Status:0 Err4xx:0 Err5xx:0 Tag: Code:⚊NULL⚊ Dur:7
	//Header:host,172.17.0.120:38080;connection,close;apptrace-pappname,shop_web;apptrace-traceid,shop_web^1715024901017^0kQfJK23Iu^2351;apptrace-sid,8171f6da-ed6c-4664-8c5a-9c1dec8c7e91;user-agent,Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36;accept,*/*;apptrace-spanid,174155503488541;apptrace-papptype,10;apptrace-uid,b8ac9134-e69b-4b82-8bac-d1e033238eb0;pagentid,websdk;x-requested-with,XMLHttpRequest;apptrace-pspanid,522207364658462;referer,http://172.17.0.128:18884/admin?login;accept-encoding,gzip, deflate;accept-language,zh-CN,zh;q=0.9;cookie,rememberMe=SN1hXxMwrUUAJ3s5AYMvSOaTVwKk/Ag4+9D7V4qWtiMEmb74n2teaOwihrvTJmNEmIT9PABAd1o3GnIDOvXzNdJUoo4/8Hq27xPwZhblzkH0QyZm4CL6glBSLh2C8thbkwftXKC4E2dSJOPFFwHyCX7YqKGVRjZzKeJIv7dX1+pUxw9MnYvf/GcTnpKPKGuMwBj6tOEIPAzDuSxPAPYd2V8bPnLZjvNoo012m8kTcCejUhTFsMzWmSM7UbhVbX5/OywMWoPYSVNPcXyrUo4UjcsglHLwiP2LpE3dF4m68ZSDqvtIi6pwkx6T0Flv/h1KfQz+8AWwO0DVWGR7hi2IaA65nbY4mq2lJCy82junaOS/IHr0JnFf3UyuwXf+g53ONu7q634xZLVl2EnCDL7Muz2hfc9vk6w0yx2NQgliD0BZD6cL0mHjrmInFJeU0poQEB7kFsXuYSV0sK10/TDBXsSzrLodjlLzUZ8QRnPYIoIbU8SuJ5gNqzdu/7ewD8mOwqYoqR/vf7zvc5Qzgbrlxp/QEhnXeKzJCH7yV9r0+B1X/ixPu4N9544EfeKdQYu9vANkY2ubK0BkqbocqHOmnMyrPBmu8aW7OdaVBWXEt9Fw7pqM6tHyhpmstO6vL3QJJtWK24sJK8yQMqIq8pgF0XYJe/25OSZco0nBBBXFExECqGQA1z+y5GxQaGEkmZeYzV+PWy8O0QsGXn4i1JpNftIGy/pYkMw2B7S9IhcZ36ZjO+WsbP6CBhQYskA8/oWPww2knWfSv/4XabYSpz6NDfIVv6daiWkiQkbnI2dsv/E=; pdomain=172.17.0.128; JSESSIONID=5C96AF9DC1CAD319C9A2EE4D94E9CEBA; ppageid=874b42175a5f68a8517cdaa88145fa74; puri=/admin/sys/user/info; shop.session.id=eb6ab80292364bcc9220b84ca80fcef2 Body: ResHeader:Accept-Charset,[Accept-Charset, Content-Type, Content-Length, Date, Connection];Connection,[Accept-Charset, Content-Type, Content-Length, Date, Connection];Content-Length,[Accept-Charset, Content-Type, Content-Length, Date, Connection];Date,[Accept-Charset, Content-Type, Content-Length, Date, Connection];Content-Type,[Accept-Charset, Content-Type, Content-Length, Date, Connection]
	//ResBody: UEventModel: UEventID: UserID:b8ac9134-e69b-4b82-8bac-d1e033238eb0 SessionID:8171f6da-ed6c-4664-8c5a-9c1dec8c7e91 Province:Unknown City:Unknown BizData:{Cip:127.0.0.1 Cid:00 Cname:未知}
	//PageID:7d92448a21049f10ffa28719f08426cb PageGroup:0611b3ff18d459bf30ba82c89278d9e3a961d994 APIID:4 Exception:0 Type:server Ts:1715024901050 Browser:chrome/124.0.0.0 Device:Windows x64 OS:Win10 OSVersion:10 OSVersionNumber:10000000 IsOtel:true}
}
