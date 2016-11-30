package ratelimiter

import (
	"encoding/json"

	"github.com/garyburd/redigo/redis"
)

type RateConfig struct {
	Rate     int `json:"rate"`
	Interval int `json:"interval"`
}

// ListenRateConfigForLimiter subscribes on a channel, and configure limiter on RateConfig message
func ListenRateConfigForLimiter(channel string, conn redis.Conn, limiter *RedisRollingRateLimiter) error {
	psc := redis.PubSubConn{conn}
	err := psc.Subscribe(channel)
	if err != nil {
		return err
	}
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			config := RateConfig{}
			err = json.Unmarshal(v.Data, &config)
			if err != nil {
			} else {
				limiter.interval = 0
				limiter.rate = 0
			}
		case redis.Subscription:
		case error:
			return v
		}
	}
}
