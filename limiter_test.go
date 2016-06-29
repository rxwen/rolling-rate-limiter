package ratelimiter_test

import (
	"testing"
	"time"

	"github.com/rxwen/rolling-rate-limiter"
	"github.com/stretchr/testify/assert"
)

func TestImplementLogInterface(t *testing.T) {
	// compilation should fail if RedisRollingRateLimiter doesn't implement RollingRateLimiter interface.
	var _ ratelimiter.RollingRateLimiter = ratelimiter.RedisRollingRateLimiter{}
	var _ ratelimiter.RollingRateLimiter = (*ratelimiter.RedisRollingRateLimiter)(nil)
}

func TestCheckRate(t *testing.T) {
	assert := assert.New(t)
	const interval = 3
	var l ratelimiter.RollingRateLimiter = ratelimiter.NewRedisRollingRateLimiter("172.16.154.128:6379", interval, 3)
	assert.True(l.Check("aaa"))
	assert.True(l.Check("aaa"))
	assert.True(l.Check("aaa"))
	assert.False(l.Check("aaa"))
	assert.True(l.Check("bbb"))
	assert.False(l.Check("aaa"))
	time.Sleep(interval * time.Second)
	assert.True(l.Check("aaa"))
}
