package redisDB

import (
	"context"
	"time"

	"github.com/go-redis/redis"
)

type TCredential struct {
	MasterName       string
	Password         string
	SentinelAddrs    []string
	SentinelPassword string
}

type TRedis struct {
	ctx     context.Context
	Session *redis.Client
}

func New(_cred TCredential) (tRedis *TRedis, err error) {
	defer func() {
		if e := recover(); e != nil {
			tRedis = nil
		}
		return
	}()

	tRedis = new(TRedis)
	tRedis.ctx = context.Background()

	tRedis.Session = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    _cred.MasterName,
		Password:      _cred.Password,
		SentinelAddrs: _cred.SentinelAddrs,
		// SentinelPassword: _cred.SentinelPassword,
	})
	return
}

func (r *TRedis) Disconnect() error {
	if r.Session == nil {
		return nil
	}
	return r.Session.Close()
}

func (r *TRedis) Set(_key string, _value interface{}, _expiration time.Duration) error {
	// return r.Session.Set(r.ctx, _key, _value, _expiration).Err()
	return r.Session.Set(_key, _value, _expiration).Err()
}

func (r *TRedis) Get(_key string) (string, error) {
	// return r.Session.Get(r.ctx, _key).Result()
	return r.Session.Get(_key).Result()
}
