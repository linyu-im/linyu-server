package service

import (
	"context"
	"github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"sync"
)

var EmbedderPool sync.Map

func GetEmbedder() (*openai.Embedder, error) {
	if m, ok := EmbedderPool.Load("default"); ok {
		return m.(*openai.Embedder), nil
	}
	embedder, err := openai.NewEmbedder(context.Background(), &openai.EmbeddingConfig{
		BaseURL:    config.C.AI.Embedding.BaseUrl,
		APIKey:     config.C.AI.Embedding.ApiKey,
		Model:      config.C.AI.Embedding.Model,
		Dimensions: &config.C.AI.Embedding.Dimensions,
		Timeout:    0,
	})
	EmbedderPool.Store("default", embedder)
	return embedder, err
}
