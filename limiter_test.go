package ratelimiter_test

import (
	"github.com/rxwen/rolling-rate-limiter"
	"testing"
)

func TestImplementLogInterface(t *testing.T) {
	// compilation should fail if MySqlUserManagement doesn't implement UserManagement interface.
	var _ ratelimiter.RollingRateLimiter = ratelimiter.RedisRollingRateLimiter{}
	var _ ratelimiter.RollingRateLimiter = (*ratelimiter.RedisRollingRateLimiter)(nil)
}
