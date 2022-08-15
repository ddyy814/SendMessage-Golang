package main

import (
	"time"
	"github.com/garyburd/redigo/redis"
)

// define a pool
var pool *redis.Pool

// when start application, init pool
func initPool(address string, maxIdle int, maxActive int, idleTimeout time.Duration) {
	pool = &redis.Pool{
		MaxIdle: maxIdle,
		MaxActive: maxActive,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address)
		},
	}
}