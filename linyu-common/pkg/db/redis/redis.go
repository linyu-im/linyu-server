package redis

import (
	"context"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

var Ctx = context.Background()

// CreateRedisClient 创建redis客户端
func CreateRedisClient() *RedisClient {
	c := config.C.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password,
		DB:       c.DB,
	})
	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		panic("failed to connect redis database: " + err.Error())
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
