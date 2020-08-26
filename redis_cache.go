// Copyright (c) 2020 HigKer
// Open Source: MIT License
// Author: SDing <deen.job@qq.com>
// Date: 2020/8/23 - 9:09 PM - UTC/GMT+08:00

package session

import (
	"context"
	"github.com/go-redis/redis"
	"sync"
	"time"
)

type RedisStore struct {
	mx     sync.Mutex
	client *redis.Client
}

func newRedisStore() (*RedisStore, error) {
	pool := redis.NewClient(&redis.Options{
		Addr:     _Cfg.RedisAddr,
		Password: _Cfg.RedisPassword, // no password set
		DB:       _Cfg.RedisDB,       // use default DB
	})
	err := pool.Ping().Err()
	if err != nil {
		return nil, err
	}
	return &RedisStore{
		client: pool,
	}, nil
}

func (r *RedisStore) Writer(ctx context.Context) error {
	cv := ctx.Value(contextValue).(map[string]interface{})
	serialize, err := Serialize(cv[contextValueData])
	if err != nil {
		return err
	}
	tmpKey := _Cfg.RedisKeyPrefix + cv[contextValueID].(string)
	_, err = r.client.HSet(tmpKey, cv[contextValueKey].(string), serialize).Result()
	// redis auto del expire data
	if err != nil {
		return ErrorSetValue
	}
	err = r.client.Expire(tmpKey, time.Duration(_Cfg.MaxAge)*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisStore) Reader(ctx context.Context) ([]byte, error) {
	cv := ctx.Value(contextValue).(map[string]interface{})
	tmpKey := _Cfg.RedisKeyPrefix + cv[contextValueID].(string)
	result, err := r.client.HGet(tmpKey, cv[contextValueKey].(string)).Result()
	if err != nil {
		return nil, err
	}
	// json.Unmarshal(obj) 是对象类型转换   []byte(str) 这个string类型 这里result是string类型所以用这个
	// https://www.jianshu.com/p/f778206ac54c
	return []byte(result), err
}

func (r *RedisStore) Remove(ctx context.Context) {
	cv := ctx.Value(contextValue).(map[string]interface{})
	tmpKey := _Cfg.RedisKeyPrefix + cv[contextValueID].(string)
	r.client.HDel(tmpKey, cv[contextValueKey].(string))
}

func (r *RedisStore) Clean(ctx context.Context) {
	cv := ctx.Value(contextValue).(map[string]interface{})
	tmpKey := _Cfg.RedisKeyPrefix + cv[contextValueID].(string)
	r.client.Del(tmpKey)
}
