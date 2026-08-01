package result

import "github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"

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

// SpaceFileContains 文件夹包含统计
type SpaceFileContains struct {
	FileCount   int64 `json:"fileCount"`
	FolderCount int64 `json:"folderCount"`
}

// SpaceFileDetail 文件/文件夹详情
type SpaceFileDetail struct {
	ID           string              `json:"id"`
	FileName     string              `json:"fileName"`
	IsDir        bool                `json:"isDir"`
	FileType     string              `json:"fileType"`     // 类型：文件后缀；文件夹为 folder
	FileCategory string              `json:"fileCategory"` // 文件分类；文件夹可为空
	Location     string              `json:"location"`     // 位置，如 /文档/工作
	Size         int64               `json:"size"`         // 大小；文件夹为子文件总和
	Contains     *SpaceFileContains  `json:"contains"`     // 包含；文件为 null
	UpdatedAt    localtime.LocalTime `json:"updatedAt"`    // 修改时间
}
