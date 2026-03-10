package vector

import (
	"context"
	"github.com/google/uuid"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
)

type WeaviateClient struct {
	client *weaviate.Client
}

func NewWeaviateClient() *WeaviateClient {
	cfg := weaviate.Config{
		Host:   config.C.Vector.WeaviateVector.Host,
		Scheme: config.C.Vector.WeaviateVector.Scheme,
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &WeaviateClient{
		client: client,
	}
}

func (c *WeaviateClient) CollectionIsExist(collectionName string) bool {
	classes, err := c.client.Schema().Getter().Do(context.Background())
	if err != nil {
		return false
	}
	for _, cls := range classes.Classes {
		if cls.Class == collectionName {
			return true
		}
	}
	return false
}

func (c *WeaviateClient) CreateLtmCollection(collectionName string) {
	if !c.CollectionIsExist(collectionName) {
		return
	}
	class := &models.Class{
		Class:      collectionName,
		Vectorizer: "none",
		Properties: []*models.Property{
			{
				Name:     "content",
				DataType: []string{"text"},
			},
			{
				Name:     "sessionId",
				DataType: []string{"text"},
			},
		},
	}
	c.client.Schema()
	err := c.client.Schema().
		ClassCreator().
		WithClass(class).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func (c *WeaviateClient) Insert(collectionName string,
	embeddings []float32,
	metadata map[string]string,
) error {
	ctx := context.Background()

	props := map[string]interface{}{
		"content":   metadata["content"],
		"sessionId": metadata["sessionId"],
	}

	_, err := c.client.Data().Creator().
		WithClassName(collectionName).
		WithID(uuid.New().String()).
		WithProperties(props).
		WithVector(embeddings).
		Do(ctx)

	return err
}

func (c *WeaviateClient) Search(
	collectionName string,
	embeddings []float32,
	filter map[string]string,
	num int,
	similarity float32,
) ([]*VData, error) {

	ctx := context.Background()

	fields := []graphql.Field{
		{Name: "content"},
		{Name: "sessionId"},
		{
			Name: "_additional",
			Fields: []graphql.Field{
				{Name: "id"},
				{Name: "distance"},
			},
		},
	}

	classKey := cases.Title(language.English).String(collectionName)
	builder := c.client.GraphQL().Get().
		WithClassName(classKey).
		WithFields(fields...).
		WithNearVector(
			c.client.GraphQL().NearVectorArgBuilder().WithVector(embeddings),
		)

	if num > 0 {
		builder = builder.WithLimit(num)
	}

	// sessionId 过滤
	if filter != nil {
		if v, ok := filter["sessionId"]; ok {
			where := filters.Where().
				WithPath([]string{"sessionId"}).
				WithOperator(filters.Equal).
				WithValueText(v)
			builder = builder.WithWhere(where)
		}
	}

	// 执行查询
	result, err := builder.Do(ctx)
	if err != nil {
		return nil, err
	}

	var data []*VData

	get := result.Data["Get"].(map[string]interface{})
	objsInterface, ok := get[classKey]
	if !ok {
		return data, nil
	}

	objs := objsInterface.([]interface{})

	for _, obj := range objs {
		m := obj.(map[string]interface{})
		add := m["_additional"].(map[string]interface{})

		distance := add["distance"].(float64)
		score := float32(1 - distance)
		if score < similarity {
			continue
		}

		data = append(data, &VData{
			ID:         add["id"].(string),
			Content:    m["content"].(string),
			SessionId:  m["sessionId"].(string),
			Similarity: score,
		})
	}

	return data, nil
}
