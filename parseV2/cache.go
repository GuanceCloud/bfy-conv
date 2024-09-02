package parseV2

import (
	"context"
	"fmt"
	redis "github.com/redis/go-redis/v9"
	"sync"
	"time"
)

//var c1 redis.Conn
var pool *redis.Client
var localCache *Cache
var getFromOld bool

var DefaultCacheTime = int64(24 * 7)

func InitRedis(host string, port string, password string, db int) {
	if host == "" {
		host = "10.200.14.188"
	}

	if port == "" {
		port = "6379"
	}

	log.Infof("redis_host=%s redis_password=%s redis_port=%s redis_db=%d", host, password, port, db)

	pool = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       db,
	})

	localCache = NewCache()
}

func SetCacheTime(hour int64, isFromOld bool) { // 1.5.1 add func
	if hour > 0 {
		DefaultCacheTime = hour
	}

	getFromOld = isFromOld
}

func Close() {
	_ = pool.Close()
}

func RedigoGet(key string) string {
	va, ok := localCache.Get(key)
	if ok {
		return va
	}
	if pool != nil {
		val, err := pool.Do(context.Background(), "GET", key).Result()
		if err == redis.Nil {
			log.Debugf("can not get %s form redis ,err=%v ", key, err)
			return va
		}
		return val.(string)
	}

	return va
}

func RedigoSet(key, val string) {
	if pool != nil {
		pool.Do(context.Background(), "SET", key, val).Result()
	}

	localCache.Set(key, val)
}

func RedigoHGet(key, name string) string {
	defer localCache.flush(key)
	lcacheKey := fmt.Sprintf("%s-%s", key, name)
	va, ok := localCache.Get(lcacheKey)
	if ok {
		return va
	}
	if pool != nil {
		val, err := pool.HGet(context.Background(), key, name).Result()
		if err == redis.Nil {
			log.Debugf("can not get %s form redis ,err=%v ", key, err)
			return va
		}
		return val
	}
	return va
}

func RedigoHSet(key, name, val string) {
	if pool != nil {
		pool.HSet(context.Background(), key, name, val).Result()
	}
	lcacheKey := fmt.Sprintf("%s-%s", key, name)
	localCache.Set(lcacheKey, val)
	localCache.flush(key)
}

type Cache struct {
	data         map[string]cacheItem
	redisHashKey map[string]time.Time
	lock         sync.RWMutex
}

type cacheItem struct {
	value  string
	expiry time.Time
}

func NewCache() *Cache {
	cache := &Cache{
		data:         make(map[string]cacheItem),
		redisHashKey: make(map[string]time.Time),
		lock:         sync.RWMutex{},
	}

	go cache.CleanupRoutine()
	return cache
}

func (c *Cache) Set(key, value string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data[key] = cacheItem{
		value:  value,
		expiry: time.Now().Add(time.Minute * 10),
	}
}

func (c *Cache) flush(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.redisHashKey[key] = time.Now()
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

func (c *Cache) cleanup() {
	c.lock.Lock()
	defer c.lock.Unlock()

	now := time.Now()
	for key, item := range c.data {
		if now.After(item.expiry) {
			delete(c.data, key)
		}
	}
}

func (c *Cache) flushRedis() {
	// todo 如果6个小时使用过，就重置一次时间
	c.lock.Lock()
	defer c.lock.Unlock()

	now := time.Now()
	for key, t := range c.redisHashKey {
		d := now.Sub(t)
		if d < time.Hour*6 {
			if pool != nil {
				pool.Expire(context.Background(), key, time.Hour*time.Duration(DefaultCacheTime))
			}
		}
	}
}

func (c *Cache) CleanupRoutine() {
	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
			c.flushRedis()
		}
	}
}
