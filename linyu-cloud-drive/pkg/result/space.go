package result

type SpaceFileCategoryStat struct {
	FileCategory string `json:"fileCategory"`
	FileCount    int64  `json:"fileCount"`
	TotalSize    int64  `json:"totalSize"`
}

// SpaceDirTreeNode 目录树节点
type SpaceDirTreeNode struct {
	ID       string              `json:"id"`
	FileName string              `json:"fileName"`
	ParentID string              `json:"parentId"`
	Path     string              `json:"path"`
	Level    int                 `json:"level"`
	Children []*SpaceDirTreeNode `json:"children"`
}
