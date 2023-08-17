package parse

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

var c1 redis.Conn

func InitRedis(host string, port string, password string, db int) {
	if host == "" {
		host = "10.200.14.188"
	}

	if port == "" {
		port = "6379"
	}

	log.Infof("redis_host=%s redis_password=%s redis_port=%s redis_db=%d", host, password, port, db)
	setdb := redis.DialDatabase(db)                         //库
	setPasswd := redis.DialPassword(password)               //密码
	timeout := redis.DialConnectTimeout(5 * time.Second)    //连接超时时间
	readTimeout := redis.DialReadTimeout(5 * time.Second)   //读超时时间
	writeTimeout := redis.DialWriteTimeout(5 * time.Second) //写超时时间
	var err error
	c1, err = redis.Dial("tcp", host+":"+port, setdb, setPasswd, timeout, readTimeout, writeTimeout)
	if err != nil {
		log.Error(err)
	}
}

func StopRedis() {
	if c1 != nil {
		_ = c1.Close()
	}
}
