package cache

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"github.com/maypok86/otter"
)

type OtterClient struct {
	cache otter.CacheWithVariableTTL[string, string]
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
		expiration = time.Hour * 24 * 365 * 12
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
