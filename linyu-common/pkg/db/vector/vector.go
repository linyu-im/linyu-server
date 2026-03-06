package vector

type VData struct {
	ID         string  `json:"id"`
	Content    string  `json:"content"`
	SessionId  string  `json:"sessionId"`
	Similarity float32 `json:"similarity"`
}

type Vector interface {
	CreateCollection(collectionName string)
	Insert(collectionName string, embeddings []float32, metadata map[string]string) error
	Search(collectionName string, embeddings []float32, filter map[string]string, num int, similarity float32) ([]*VData, error)
}

func InitVector() Vector {
	return NewChromemClient()
}
