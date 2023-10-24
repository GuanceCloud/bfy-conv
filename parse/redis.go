package parse

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

//var c1 redis.Conn
var pool *redis.Pool

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

	pool = &redis.Pool{
		MaxIdle:     16,
		MaxActive:   1024,
		IdleTimeout: 300,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", host+":"+port, setdb, setPasswd, timeout, readTimeout, writeTimeout)
		},
	}
}
