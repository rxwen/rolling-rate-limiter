package ratelimiter

import (
	"time"

	"github.com/garyburd/redigo/redis"

	"github.com/rxwen/resourcepool"
)

// RollingRateLimiter interface makes it easier to do unit testing.
type RollingRateLimiter interface {
	Check(key string) bool
	Reset(key string)
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
	if l.interval == 0 || l.rate == 0 {
		return true
	}
	now := time.Now().Unix()
	nowNano := time.Now().UnixNano()
	timeToClean := now - int64(l.interval)
	c, e := l.pool.Get()
	if e != nil {
		return false
	}
	destroy := false
	defer func() { l.pool.Putback(c, destroy) }() // use a func to wrap the putback to avoid evalute destroy value now

	key = l.prefix + key
	conn := c.(redis.Conn)
	conn.Send("MULTI")
	conn.Send("ZREMRANGEBYSCORE", key, 0, timeToClean)
	conn.Send("ZADD", key, nowNano, nowNano)
	conn.Send("EXPIRE", key, l.interval)
	_, err := conn.Do("EXEC")
	if err != nil {
		destroy = true
		return false
	}

	items, err := redis.Strings(conn.Do("ZRANGE", key, 0, -1))
	if err != nil || len(items) > l.rate {
		return false
	} else {
		return true
	}
}

func (l RedisRollingRateLimiter) Reset(key string) {
	c, e := l.pool.Get()
	if e != nil {
		return
	}
	destroy := false
	defer func() { l.pool.Putback(c, destroy) }() // use a func to wrap the putback to avoid evalute destroy value now

	key = l.prefix + key
	conn := c.(redis.Conn)
	_, err := conn.Do("DEL", key)
	if err != nil {
		destroy = true
		return
	}
}
