package vector

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
)

type VData struct {
	ID         string  `json:"id"`
	Content    string  `json:"content"`
	SessionId  string  `json:"sessionId"`
	Similarity float32 `json:"similarity"`
}

type Vector interface {
	CreateLtmCollection(collectionName string)
	Insert(collectionName string, embeddings []float32, metadata map[string]string) error
	Search(collectionName string, embeddings []float32, filter map[string]string, num int, similarity float32) ([]*VData, error)
}

func InitVector() Vector {
	var v Vector
	switch config.C.Vector.Type {
	case config.ChromemVectorType:
		v = NewChromemClient()
	case config.WeaviateVectorType:
		v = NewWeaviateClient()
	default:
		panic("vector type not supported: " + config.C.Vector.Type)
	}
	v.CreateLtmCollection(constant.VectorCollection.LongTermMemory)
	return v
}
