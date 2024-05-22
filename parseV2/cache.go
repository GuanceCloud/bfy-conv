package parseV2

import (
	"context"
	redis "github.com/redis/go-redis/v9"
	"sync"
	"time"
)

//var c1 redis.Conn
var pool *redis.Client
var localCache *Cache

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

type Cache struct {
	data map[string]cacheItem
	lock sync.RWMutex
}

type cacheItem struct {
	value  string
	expiry time.Time
}

func NewCache() *Cache {
	cache := &Cache{
		data: make(map[string]cacheItem),
		lock: sync.RWMutex{},
	}

	go cache.CleanupRoutine()
	return cache
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

func (c *Cache) CleanupRoutine() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		}
	}
}
