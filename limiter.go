package ratelimiter

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

type RedisRollingRateLimiter struct {
	conn     redis.Conn
	interval int
	rate     int
}

func NewRedisRollingRateLimiter(endpoint string, interval, rate int) *RedisRollingRateLimiter {
	c, err := redis.Dial("tcp", endpoint)
	if err != nil {
		panic("failed to connect redis")
	}
	return &RedisRollingRateLimiter{
		conn:     c,
		interval: interval,
		rate:     rate,
	}
}

func (l RedisRollingRateLimiter) Check(key string) bool {
	now := time.Now().Unix()
	timeToClean := now - int64(l.interval)
	l.conn.Send("MULTI")
	l.conn.Send("ZREMRANGEBYSCORE", key, 0, timeToClean)
	l.conn.Send("ZADD", key, now, now)
	l.conn.Send("EXPIRE", key, l.interval)
	status, err := l.conn.Do("EXEC")
	fmt.Println(status, err)

	items, err := redis.Strings(l.conn.Do("ZRANGE", key, 0, -1))
	if len(items) >= l.rate {
		return false
	} else {
		return true
	}
}
