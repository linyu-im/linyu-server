package param

type UploadFileInfoParam struct {
	FileHash   string `json:"fileHash" binding:"required"`
	FileSize   int64  `json:"fileSize" binding:"required"`
	FileName   string `json:"fileName" binding:"required"`
	ParentID   string `json:"parentId" binding:"required"`
	TotalChunk int    `json:"totalChunk" binding:"required"`
	SpaceID    string `json:"spaceId" binding:"required"`
}
