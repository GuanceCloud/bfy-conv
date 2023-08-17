package parse

import (
	"github.com/gomodule/redigo/redis"
	"os"
	"strconv"
	"time"
)

var c1 redis.Conn

func InitRedis() {
	redisHost := os.Getenv("PMC_REDIS_HOST")
	if redisHost == "" {
		redisHost = "10.200.14.188"
	}
	redis_password := os.Getenv("PMC_REDIS_PASSWORD")

	redis_port := os.Getenv("PMC_REDIS_PORT")
	if redis_port == "" {
		redis_port = "6379"
	}
	redis_db := os.Getenv("PMC_REDIS_DB")
	if redis_db == "" {
		redis_db = "0"
	}

	db, err := strconv.Atoi(redis_db)
	if err != nil {
		db = 0
	}
	log.Infof("redis_host=%s redis_password=%s redis_port=%s redis_db=%s", redisHost, redis_password, redis_port, redis_db)
	setdb := redis.DialDatabase(db)                         //库
	setPasswd := redis.DialPassword(redis_password)         //密码
	timeout := redis.DialConnectTimeout(5 * time.Second)    //连接超时时间
	readTimeout := redis.DialReadTimeout(5 * time.Second)   //读超时时间
	writeTimeout := redis.DialWriteTimeout(5 * time.Second) //写超时时间

	c1, err = redis.Dial("tcp", redisHost+":"+redis_port, setdb, setPasswd, timeout, readTimeout, writeTimeout)
	if err != nil {
		log.Error(err)
	}
}

func StopRedis() {
	if c1 != nil {
		_ = c1.Close()
	}
}
