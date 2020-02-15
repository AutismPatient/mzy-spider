package stock

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"time"
)

type MyRedis struct {
	pool *redis.Pool
}

var (
	Redis = &MyRedis{}
)

func init() {
	Redis.newRedisPool()
}
func (r *MyRedis) newRedisPool() {
	dial := func() (conn redis.Conn, e error) {
		return redis.Dial("tcp", "127.0.0.1:6379",
			redis.DialConnectTimeout(10*time.Second),
			redis.DialReadTimeout(10*time.Second),
			redis.DialWriteTimeout(20*time.Second),
			redis.DialDatabase(2),
			redis.DialPassword("123"),
			redis.DialKeepAlive(30*time.Second),
		)
	}
	_, err := dial()
	if err != nil {
		panic(err)
	}
	r.pool = &redis.Pool{
		Dial:            dial,
		Wait:            true,
		MaxConnLifetime: 5 >> 10 * time.Second,
		MaxIdle:         100,
		IdleTimeout:     5 >> 10 * time.Second,
		MaxActive:       120,
	}
}
func (r *MyRedis) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if r == nil {
		return nil, errors.New("error: pool is null")
	}
	conn := r.pool.Get()
	defer conn.Close()
	return conn.Do(commandName, args...)
}
