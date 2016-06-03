package ratelimiter_test

import (
	"github.com/rxwen/rolling-rate-limiter"
	"testing"
)

func TestImplementLogInterface(t *testing.T) {
	// compilation should fail if RedisRollingRateLimiter doesn't implement RollingRateLimiter interface.
	var _ ratelimiter.RollingRateLimiter = ratelimiter.RedisRollingRateLimiter{}
	var _ ratelimiter.RollingRateLimiter = (*ratelimiter.RedisRollingRateLimiter)(nil)
}
