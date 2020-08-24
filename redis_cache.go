// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/23 - 9:09 PM - UTC/GMT+08:00

package session

import (
	"github.com/go-redis/redis"
	"sync"
	"time"
)

type RedisStore struct {
	mx     sync.Mutex
	client *redis.Client
}

func newRedisStore() *RedisStore {
	return &RedisStore{
		client: redis.NewClient(&redis.Options{
			Addr:     _Cfg.RedisAddr,
			Password: _Cfg.RedisPassword, // no password set
			DB:       _Cfg.RedisDB,       // use default DB
		}),
	}
}

func (r *RedisStore) Writer(id, key string, data interface{}) error {
	tmpKey := _Cfg.RedisKeyPrefix + id
	err := r.client.HSet(tmpKey, key, data).Err()
	if err != nil {
		return err
	}
	// redis auto del expire data
	r.client.Expire(tmpKey, time.Duration(_Cfg.MaxAge))
	return nil
}

func (r *RedisStore) Reader(id, key string) ([]byte, error) {
	tmpKey := _Cfg.RedisKeyPrefix + id
	result, err := r.client.HGet(tmpKey, key).Result()
	if err != nil {
		return nil, err
	}
	return Serialize(result)
}

func (r *RedisStore) Remove(id, key string) {
	tmpKey := _Cfg.RedisKeyPrefix + id
	r.client.HDel(tmpKey, key)
}

func (r *RedisStore) clean(id string) {
	tmpKey := _Cfg.RedisKeyPrefix + id
	r.client.Del(tmpKey)
}
