package main

import (
	"fmt"
	"time"

	"github.com/rxwen/resourcepool/redispool"
	"github.com/rxwen/rolling-rate-limiter"
)

func main() {
	redisPool, _ := redispool.CreateRedisConnectionPool("192.168.2.175:6379", 10, 5)
	r := ratelimiter.NewRedisRollingRateLimiter("test", redisPool, 20, 10)
	for i := 0; i < 4; i++ {
		time.Sleep(1 * time.Second)
		fmt.Println(r.Check("ip_address_of_client"))
	}
}
