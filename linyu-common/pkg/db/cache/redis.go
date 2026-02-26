package cache

import (
	"context"
	"encoding/json"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

var Ctx = context.Background()

func NewRedisClient() *RedisClient {
	c := config.C.Cache.RedisCache
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password,
		DB:       c.DB,
	})
	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		panic("failed to connect cache database: " + err.Error())
	}
	return &RedisClient{client: rdb, ctx: context.Background()}
}

// Set 设置key-value，带过期时间
func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

// Get 获取key
func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

// Del 删除key
func (r *RedisClient) Del(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// Exists 判断key是否存在
func (r *RedisClient) Exists(key string) (bool, error) {
	count, err := r.client.Exists(r.ctx, key).Result()
	return count > 0, err
}

// SetObject 存储对象
func (r *RedisClient) SetObject(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(r.ctx, key, data, expiration).Err()
}

func (r *RedisClient) GetObject(key string, dest interface{}) error {
	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (r *RedisClient) SAdd(key string, ttl time.Duration, members ...interface{}) error {
	err := r.client.SAdd(r.ctx, key, members...).Err()
	if err != nil {
		return err
	}
	if ttl > 0 {
		return r.client.Expire(r.ctx, key, ttl).Err()
	}
	return nil
}

func (r *RedisClient) SMembers(key string) ([]string, error) {
	return r.client.SMembers(r.ctx, key).Result()
}
