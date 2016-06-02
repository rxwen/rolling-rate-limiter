package main

import (
	"fmt"
	"github.com/rxwen/rolling-rate-limiter"
	"time"
)

func main() {
	r := ratelimiter.NewRedisRollingRateLimiter("172.16.154.128:6379", 20, 10)
	for i := 0; i < 4; i++ {
		time.Sleep(1 * time.Second)
		fmt.Println(r.Check("ip_address_of_client"))
	}
}
