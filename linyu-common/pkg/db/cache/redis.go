package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"github.com/redis/go-redis/v9"
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

func (r *RedisClient) SRem(key string, members ...interface{}) error {
	return r.client.SRem(r.ctx, key, members...).Err()
}

func (r *RedisClient) SMembers(key string) ([]string, error) {
	return r.client.SMembers(r.ctx, key).Result()
}

// SetIfGreater 原子更新：仅当新值大于当前值时写入
func (r *RedisClient) SetIfGreater(key string, newValue string, expiration time.Duration) error {
	const script = `
local cur = redis.call('GET', KEYS[1])
if (not cur) or (tonumber(ARGV[1]) > tonumber(cur)) then
    if tonumber(ARGV[2]) > 0 then
        redis.call('SET', KEYS[1], ARGV[1], 'PX', ARGV[2])
    else
        redis.call('SET', KEYS[1], ARGV[1])
    end
end
return 1
`
	ms := int64(0)
	if expiration > 0 {
		ms = expiration.Milliseconds()
	}
	return r.client.Eval(r.ctx, script, []string{key}, newValue, ms).Err()
}
