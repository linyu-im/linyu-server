package param

type SpaceFileListParam struct {
	ParentID string `json:"parentId" binding:"required"`
}

type SpaceCreateDirParam struct {
	ParentID string `json:"parentId" binding:"required"`
	DirName  string `json:"dirName" binding:"required"`
}

type SpaceFileDeleteParam struct {
	SpaceFileIDs []string `json:"spaceFileIds" binding:"required,min=1"`
}

type SpaceRecycleRestoreParam struct {
	SpaceRecycleIDs []string `json:"spaceRecycleIds" binding:"required,min=1"`
}

type SpaceRecycleDeleteParam struct {
	SpaceRecycleIDs []string `json:"spaceRecycleIds" binding:"required,min=1"`
}
