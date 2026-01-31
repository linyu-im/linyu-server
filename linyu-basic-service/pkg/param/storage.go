package param

type StorageDeleteParam struct {
	FileKey string `json:"fileKey" binding:"required"`
}
