package param

type SpaceFileListParam struct {
	ParentID string `json:"parentId" binding:"required"`
}

type SpaceFileListAllParam struct {
	ParentID string `json:"parentId" binding:"required"`
}

type SpaceCreateDirParam struct {
	ParentID string `json:"parentId" binding:"required"`
	DirName  string `json:"dirName" binding:"required"`
}

type SpaceFileDeleteParam struct {
	SpaceFileIDs []string `json:"spaceFileIds" binding:"required,min=1"`
}

type SpaceFileMoveParam struct {
	SpaceFileIDs   []string `json:"spaceFileIds" binding:"required,min=1"`
	TargetParentID string   `json:"targetParentId" binding:"required"`
}

type SpaceFileRenameParam struct {
	SpaceFileID string `json:"spaceFileId" binding:"required"`
	NewName     string `json:"newName" binding:"required"`
}

type SpaceRecycleRestoreParam struct {
	SpaceRecycleIDs []string `json:"spaceRecycleIds" binding:"required,min=1"`
}

type SpaceRecycleDeleteParam struct {
	SpaceRecycleIDs []string `json:"spaceRecycleIds" binding:"required,min=1"`
}
