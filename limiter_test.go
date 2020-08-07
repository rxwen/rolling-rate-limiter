package ratelimiter_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/rxwen/resourcepool"
	"github.com/rxwen/resourcepool/redispool"
	ratelimiter "github.com/rxwen/rolling-rate-limiter"
	"github.com/stretchr/testify/assert"
)

func TestImplementLogInterface(t *testing.T) {
	// compilation should fail if RedisRollingRateLimiter doesn't implement RollingRateLimiter interface.
	var _ ratelimiter.RollingRateLimiter = ratelimiter.RedisRollingRateLimiter{}
	var _ ratelimiter.RollingRateLimiter = (*ratelimiter.RedisRollingRateLimiter)(nil)
}

func updateRate(pool *resourcepool.ResourcePool, channel string, cfg ratelimiter.RateConfig) {
	c, e := pool.Get()
	if e != nil {
		panic(e)
	}
	destroy := false
	defer func() { pool.Putback(c, destroy) }() // use a func to wrap the putback to avoid evalute destroy value now
	conn := c.(redis.Conn)
	data, err := json.Marshal(cfg)
	if err != nil {
		panic(err)
	}

	time.Sleep(100 * time.Millisecond)
	_, err = conn.Do("PUBLISH", channel, data)
	if err != nil {
		panic(err)
	}

	time.Sleep(300 * time.Millisecond)
}

func TestCheckRate(t *testing.T) {
	assert := assert.New(t)
	redisPool, _ := redispool.CreateRedisConnectionPool("192.168.2.175:6379", 10, 5)
	const interval = 1
	var l *ratelimiter.RedisRollingRateLimiter = ratelimiter.NewRedisRollingRateLimiter("demo", redisPool, interval, 3)
	assert.True(l.Check("aaa"))
	assert.True(l.Check("aaa"))
	assert.True(l.Check("aaa"))
	assert.False(l.Check("aaa"))
	assert.True(l.Check("bbb"))
	assert.False(l.Check("aaa"))
	time.Sleep(interval * time.Second)
	assert.True(l.Check("aaa"))
	time.Sleep(interval * time.Second)

	channel := "ratechannel"
	go func() {
		e := ratelimiter.ListenRateConfigForLimiter(channel, redisPool, l)
		if e != nil {
			panic(e)
		}
	}()
	updateRate(redisPool, channel, ratelimiter.RateConfig{Rate: 5, Interval: interval})
	assert.True(l.Check("aaa"))
	assert.True(l.Check("aaa"))
	assert.True(l.Check("aaa"))
	assert.True(l.Check("aaa"))
	assert.True(l.Check("aaa"))
	assert.False(l.Check("aaa"))
	assert.True(l.Check("bbb"))
	assert.False(l.Check("aaa"))
	updateRate(redisPool, channel, ratelimiter.RateConfig{Rate: 0, Interval: interval})
	for i := 0; i < 20; i++ {
		assert.True(l.Check("ccc"))
	}
	updateRate(redisPool, channel, ratelimiter.RateConfig{Rate: 5, Interval: 0})
	for i := 0; i < 20; i++ {
		assert.True(l.Check("ddd"))
	}
}

func TestCheckReset(t *testing.T) {
	assert := assert.New(t)
	redisPool, _ := redispool.CreateRedisConnectionPool("192.168.2.175:6379", 10, 5)
	const interval = 5
	var l *ratelimiter.RedisRollingRateLimiter = ratelimiter.NewRedisRollingRateLimiter("demo2", redisPool, interval, 3)
	assert.True(l.Check("bbb"))
	assert.True(l.Check("bbb"))
	assert.True(l.Check("bbb"))
	assert.False(l.Check("bbb"))
	l.Reset("bbb")
	assert.True(l.Check("bbb"))
	assert.True(l.Check("bbb"))
	assert.True(l.Check("bbb"))
	assert.False(l.Check("bbb"))
	l.Reset("bbb")
	assert.True(l.Check("bbb"))
	assert.True(l.Check("bbb"))
	assert.True(l.Check("bbb"))
	assert.False(l.Check("bbb"))
	l.Reset("bbb")
}
