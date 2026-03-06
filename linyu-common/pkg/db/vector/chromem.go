package vector

import (
	"context"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
	"github.com/philippgille/chromem-go"
	"log"
)

type ChromemClient struct {
	client *chromem.DB
}

func NewChromemClient() *ChromemClient {
	db, _ := chromem.NewPersistentDB(config.C.Vector.ChromemVector.Path, false)
	return &ChromemClient{
		client: db,
	}
}

func (c *ChromemClient) CreateCollection(collectionName string) {
	_, err := c.client.GetOrCreateCollection(collectionName, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *ChromemClient) Search(collectionName string, embeddings []float32, filter map[string]string, num int, similarity float32) ([]*VData, error) {
	collection, err := c.client.GetOrCreateCollection(collectionName, nil, nil)
	if err != nil {
		return []*VData{}, err
	}
	l := collection.Count()
	if l <= 0 {
		return nil, nil
	}
	if num > collection.Count() {
		num = collection.Count()
	}
	results, err := collection.QueryEmbedding(
		context.Background(),
		embeddings,
		num,
		filter,
		nil,
	)
	var data []*VData
	for _, v := range results {
		if v.Similarity < similarity {
			continue
		}
		data = append(data, &VData{
			ID:         v.ID,
			Content:    v.Metadata["content"],
			SessionId:  v.Metadata["sessionId"],
			Similarity: v.Similarity,
		})
	}
	return data, err
}

func (c *ChromemClient) Insert(collectionName string, embeddings []float32, metadata map[string]string) error {
	collection, _ := c.client.GetOrCreateCollection(collectionName, nil, nil)
	err := collection.AddDocument(context.Background(), chromem.Document{
		ID:        utils.GenerateSfIDString(),
		Embedding: embeddings,
		Metadata:  metadata,
	})
	return err
}
