package parse

import (
	"github.com/gomodule/redigo/redis"
	"sync"
	"time"
)

//var c1 redis.Conn
var pool *redis.Pool
var localCache *Cache

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

	localCache = NewCache()
}

func RedigoGet(key string) string {
	va, ok := localCache.Get(key)
	if ok {
		return va
	}

	c := pool.Get()
	defer func() { c.Close() }()

	cachedTraceID, err := redis.String(c.Do("GET", key))
	if err != nil || cachedTraceID == "" {
		log.Debugf("can not get %s form redis ,err=%v , or trace_id=%s", xid, err, cachedTraceID)

	}
	return cachedTraceID
}

func RedigoSet(key, val string) {
	c := pool.Get()
	defer func() { c.Close() }()
	c.Do("SET", key, val, "EX", 600)
	localCache.Set(key, val)
}

type Cache struct {
	data map[string]cacheItem
	lock sync.RWMutex
}

type cacheItem struct {
	value  string
	expiry time.Time
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]cacheItem),
		lock: sync.RWMutex{},
	}
}

func (c *Cache) Set(key, value string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data[key] = cacheItem{
		value:  value,
		expiry: time.Now().Add(time.Minute),
	}
}

func (c *Cache) Get(key string) (string, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	item, ok := c.data[key]
	if ok && time.Now().Before(item.expiry) {
		return item.value, true
	}

	return "", false
}

func (c *Cache) Cleanup() {
	c.lock.Lock()
	defer c.lock.Unlock()

	now := time.Now()
	for key, item := range c.data {
		if now.After(item.expiry) {
			delete(c.data, key)
		}
	}
}

func (c *Cache) CleanupRoutine() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.Cleanup()
		}
	}
}
