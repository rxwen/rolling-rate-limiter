package ratelimiter

import (
	"log"
	"time"

	"github.com/garyburd/redigo/redis"

	"github.com/rxwen/resourcepool"
)

// RollingRateLimiter interface makes it easier to do unit testing.
type RollingRateLimiter interface {
	Check(key string) bool
}

// RedisRollingRateLimiter
type RedisRollingRateLimiter struct {
	pool     *resourcepool.ResourcePool
	interval int
	rate     int
	prefix   string
}

func NewRedisRollingRateLimiter(prefix string, redisPool *resourcepool.ResourcePool, interval, rate int) *RedisRollingRateLimiter {
	return &RedisRollingRateLimiter{
		pool:     redisPool,
		interval: interval,
		rate:     rate,
	}
}

func (l RedisRollingRateLimiter) Check(key string) bool {
	now := time.Now().Unix()
	nowNano := time.Now().UnixNano()
	timeToClean := now - int64(l.interval)
	c, e := l.pool.Get()
	if e != nil {
		return false
	}
	defer l.pool.Release(c)

	key = l.prefix + key
	conn := c.(redis.Conn)
	conn.Send("MULTI")
	conn.Send("ZREMRANGEBYSCORE", key, 0, timeToClean)
	conn.Send("ZADD", key, nowNano, nowNano)
	conn.Send("EXPIRE", key, l.interval)
	status, err := conn.Do("EXEC")
	if err != nil {
		log.Println(status, err)
		return false
	}

	items, err := redis.Strings(conn.Do("ZRANGE", key, 0, -1))
	if err != nil || len(items) > l.rate {
		log.Println(err)
		return false
	} else {
		return true
	}
}
