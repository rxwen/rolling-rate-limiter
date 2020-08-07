package ratelimiter

import (
	"encoding/json"

	"github.com/garyburd/redigo/redis"
	"github.com/rxwen/resourcepool"
)

type RateConfig struct {
	Rate     int `json:"rate"`
	Interval int `json:"interval"`
}

// ListenRateConfigForLimiter subscribes on a channel, and configure limiter on RateConfig message
func ListenRateConfigForLimiter(channel string, pool *resourcepool.ResourcePool,
	limiter *RedisRollingRateLimiter) error {
	c, e := pool.Get()
	if e != nil {
		return e
	}
	destroy := false
	defer func() { pool.Putback(c, destroy) }() // use a func to wrap the putback to avoid evalute destroy value now
	conn := c.(redis.Conn)
	psc := redis.PubSubConn{conn}
	err := psc.Subscribe(channel)
	if err != nil {
		destroy = true
		return err
	}
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			config := RateConfig{}
			err = json.Unmarshal(v.Data, &config)
			if err == nil {
				limiter.interval = config.Interval
				limiter.rate = config.Rate
			}
		case redis.Subscription:
		case error:
			destroy = true
			return v
		}
	}
}
