package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"github.com/maypok86/otter"
)

var LongPeriod = time.Hour * 24 * 365 * 12

type OtterClient struct {
	cache otter.CacheWithVariableTTL[string, string]
	mu    sync.Mutex
}

func NewOtterClient() *OtterClient {
	conf := config.C.Cache.OtterCache
	capacity := conf.Capacity
	if capacity <= 0 {
		capacity = 10_000
	}
	c, err := otter.MustBuilder[string, string](capacity).
		CollectStats().
		Cost(func(key string, value string) uint32 {
			return 1
		}).
		WithVariableTTL().
		Build()
	if err != nil {
		panic("failed to create otter cache: " + err.Error())
	}
	return &OtterClient{
		cache: c,
	}
}

func (o *OtterClient) Set(key string, value interface{}, expiration time.Duration) error {
	var valStr string
	switch v := value.(type) {
	case string:
		valStr = v
	case []byte:
		valStr = string(v)
	default:
		b, err := json.Marshal(value)
		if err != nil {
			return err
		}
		valStr = string(b)
	}
	if expiration <= 0 {
		expiration = LongPeriod
	}
	o.cache.Set(key, valStr, expiration)
	return nil
}

func (o *OtterClient) Get(key string) (string, error) {
	val, ok := o.cache.Get(key)
	if !ok {
		return "", errors.New("cache: key not found")
	}
	return val, nil
}

func (o *OtterClient) Del(key string) error {
	o.cache.Delete(key)
	return nil
}

func (o *OtterClient) Exists(key string) (bool, error) {
	return o.cache.Has(key), nil
}

// SetObject 存储对象
func (o *OtterClient) SetObject(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if expiration <= 0 {
		expiration = LongPeriod
	}

	o.cache.Set(key, string(data), expiration)
	return nil
}

func (o *OtterClient) GetObject(key string, dest interface{}) error {
	val, ok := o.cache.Get(key)
	if !ok {
		return errors.New("cache: key not found")
	}

	return json.Unmarshal([]byte(val), dest)
}

func (o *OtterClient) SAdd(key string, expiration time.Duration, members ...interface{}) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	var current map[string]struct{}
	val, ok := o.cache.Get(key)
	if ok {
		_ = json.Unmarshal([]byte(val), &current)
	}

	if current == nil {
		current = make(map[string]struct{})
	}

	for _, m := range members {
		str := fmt.Sprintf("%v", m)
		current[str] = struct{}{}
	}

	data, err := json.Marshal(current)
	if err != nil {
		return err
	}

	if expiration <= 0 {
		expiration = LongPeriod
	}

	o.cache.Set(key, string(data), expiration)

	return nil
}

func (o *OtterClient) SRem(key string, members ...interface{}) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	val, ok := o.cache.Get(key)
	if !ok {
		return nil
	}

	var current map[string]struct{}
	if err := json.Unmarshal([]byte(val), &current); err != nil {
		return err
	}
	for _, member := range members {
		delete(current, fmt.Sprintf("%v", member))
	}
	if len(current) == 0 {
		o.cache.Delete(key)
		return nil
	}

	data, err := json.Marshal(current)
	if err != nil {
		return err
	}
	o.cache.Set(key, string(data), LongPeriod)
	return nil
}

func (o *OtterClient) SMembers(key string) ([]string, error) {
	val, ok := o.cache.Get(key)
	if !ok {
		return []string{}, nil
	}
	var resMap map[string]struct{}
	_ = json.Unmarshal([]byte(val), &resMap)

	keys := make([]string, 0, len(resMap))
	for k := range resMap {
		keys = append(keys, k)
	}
	return keys, nil
}

// SetIfGreater 进程内加锁后比较再写入（Otter 本身非跨进程共享）
func (o *OtterClient) SetIfGreater(key string, newValue string, expiration time.Duration) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if cur, ok := o.cache.Get(key); ok && cur != "" {
		curN, curErr := strconv.ParseInt(cur, 10, 64)
		newN, newErr := strconv.ParseInt(newValue, 10, 64)
		if curErr == nil && newErr == nil && newN <= curN {
			return nil
		}
	}
	if expiration <= 0 {
		expiration = LongPeriod
	}
	o.cache.Set(key, newValue, expiration)
	return nil
}
