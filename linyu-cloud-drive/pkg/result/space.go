package result

type SpaceFileCategoryStat struct {
	FileCategory string `json:"fileCategory"`
	FileCount    int64  `json:"fileCount"`
	TotalSize    int64  `json:"totalSize"`
}
