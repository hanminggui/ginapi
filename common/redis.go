package common

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
	"time"
)

// Redis redis连接池实例
var Redis *redisTool

func InitRedis() {
	if Redis != nil {
		return
	}
	config := viper.GetStringMapString("redis")
	Redis = &redisTool{
		redisPool: &redis.Pool{
			MaxIdle:     viper.GetInt("redis.maxIdle"),                        // 最大空闲连接数
			MaxActive:   viper.GetInt("redis.maxActive"),                      // 最大连接数
			IdleTimeout: viper.GetDuration("redis.idleTimeout") * time.Second, // 空闲连接存活时长
			Wait:        true,                                                 // 连接数满后新的请求排队等待
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", config["addr"],
					redis.DialPassword(config["password"]),
					redis.DialDatabase(viper.GetInt("redis.database")),
					redis.DialConnectTimeout(viper.GetDuration("redis.connectTimeout")*time.Second),
					redis.DialReadTimeout(viper.GetDuration("redis.readTimeout")*time.Second),
					redis.DialWriteTimeout(viper.GetDuration("redis.writeTimeout")*time.Second),
				)
			},
		},
	}
	go Redis.heartbeatRedis()
}

type redisTool struct {
	redisPool *redis.Pool
}

func (rds *redisTool) close() {
	rds.redisPool.Close()
}

// 心跳检测
func (rds *redisTool) heartbeatRedis() {
	rds.ping()
	t := time.Tick(time.Second * 10)
	for range t {
		rds.ping()
	}
}

func (rds *redisTool) ping() {
	rc := rds.redisPool.Get()
	defer rc.Close()
	FatalError(rc.Err())
	PONG, err := redis.String(rc.Do("PING"))
	FatalError(err)
	if PONG != "PONG" {
		FatalError(errors.New("redis 连接失败"))
	}
	Log.Debug("ping redis success")
}

func (rds *redisTool) Get(key string, dest interface{}) (exists bool) {
	rc := rds.redisPool.Get()
	defer rc.Close()
	exists, err := redis.Bool(rc.Do("EXISTS", key))
	FatalError(err)
	if !exists {
		return
	}
	v, err := rc.Do("GET", key)
	FatalError(err)
	tmp, err := redis.Bytes(v, nil)
	PanicError(err)
	json.Unmarshal(tmp, dest)
	return
}

func (rds *redisTool) Set(key string, v interface{}, expire int64) {
	rc := rds.redisPool.Get()
	defer rc.Close()
	jsonv, err := json.Marshal(v)
	PanicError(err)
	if expire > 0 {
		_, err := rc.Do("SET", key, jsonv, "EX", expire)
		FatalError(err)
	} else {
		_, err := rc.Do("SET", key, jsonv)
		FatalError(err)
	}
}

func (rds *redisTool) Exists(key string) (exists bool) {
	rc := rds.redisPool.Get()
	defer rc.Close()
	exists, err := redis.Bool(rc.Do("EXISTS", key))
	FatalError(err)
	return
}

func (rds *redisTool) Expire(key string, expire int64) {
	rc := rds.redisPool.Get()
	defer rc.Close()
	_, err := rc.Do("EXPIRE", key, expire)
	FatalError(err)
}

func (rds *redisTool) Del(key string) {
	rc := rds.redisPool.Get()
	defer rc.Close()
	_, err := rc.Do("DEL", key)
	FatalError(err)
}

func (rds *redisTool) Incr(key string, step int8) {
	rc := rds.redisPool.Get()
	defer rc.Close()
	_, err := rc.Do("INCR", key, step)
	FatalError(err)
}

func (rds *redisTool) Do(commandName string, args ...interface{}) (v interface{}) {
	rc := rds.redisPool.Get()
	defer rc.Close()
	v, err := rc.Do(commandName, args)
	// 命令使用错误不能把服务搞崩了
	PanicError(err)
	return
}

func (rds *redisTool) Obtain(key string, expire int64, getValue func() interface{}, v interface{}) {
	exists := rds.Get(key, v)
	if exists {
		return
	}
	x := getValue()
	rds.Set(key, x, expire)
	rds.Get(key, v)
}
