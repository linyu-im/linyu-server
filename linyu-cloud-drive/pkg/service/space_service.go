package service

import (
	"errors"
	"path/filepath"
	"strings"
	"time"

	basicService "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/service"
	driveDao "github.com/linyu-im/linyu-server/linyu-cloud-drive/internal/dao"
	driveModel "github.com/linyu-im/linyu-server/linyu-cloud-drive/pkg/model"
	driveResult "github.com/linyu-im/linyu-server/linyu-cloud-drive/pkg/result"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
	"gorm.io/gorm"
)

var SpaceService = newSpaceService()

func newSpaceService() *spaceService {
	return &spaceService{}
}

type spaceService struct{}

// CheckUserSpaceQuota 校验用户空间容量是否足够容纳本次上传文件
func (s *spaceService) CheckUserSpaceQuota(userId string, fileSize int64) error {
	space, err := s.GetOrCreateUserSpace(userId)
	if err != nil {
		return err
	}
	return s.CheckSpaceQuota(space, fileSize)
}

// CheckSpaceQuota 校验空间容量是否足够
func (s *spaceService) CheckSpaceQuota(space *driveModel.Space, fileSize int64) error {
	// QuotaBytes 为 0 表示不限制
	if space.QuotaBytes > 0 && space.UsedBytes+fileSize > space.QuotaBytes {
		return errors.New("cloud-drive.space.quota-exceeded")
	}
	return nil
}

// CalcSpaceFilesUsedBytes 计算文件/文件夹占用容量（文件夹通过 path like 汇总全部子文件）
func (s *spaceService) CalcSpaceFilesUsedBytes(tx *gorm.DB, spaceId string, files []*driveModel.SpaceFile) (int64, int64, error) {
	topFiles := filterTopLevelSpaceFiles(files)
	var totalSize, fileCount int64
	for _, file := range topFiles {
		if file.IsDir {
			size, count, err := driveDao.SpaceFileDao.SumFileSizeSelfAndDescendants(tx, spaceId, file.ID, file.Path)
			if err != nil {
				return 0, 0, err
			}
			totalSize += size
			fileCount += count
			continue
		}
		totalSize += file.FileSize
		fileCount++
	}
	return totalSize, fileCount, nil
}

// filterTopLevelSpaceFiles 过滤掉已被其他选中文件夹覆盖的项，避免重复统计
func filterTopLevelSpaceFiles(files []*driveModel.SpaceFile) []*driveModel.SpaceFile {
	result := make([]*driveModel.SpaceFile, 0, len(files))
	for _, file := range files {
		covered := false
		for _, other := range files {
			if other.ID == file.ID || !other.IsDir {
				continue
			}
			if strings.HasPrefix(file.Path, other.Path+"/") {
				covered = true
				break
			}
		}
		if !covered {
			result = append(result, file)
		}
	}
	return result
}

func (s *spaceService) GetOrCreateUserSpace(userId string) (*driveModel.Space, error) {
	if space := driveDao.SpaceDao.GetByTypeAndTargetID(db.RDB, constant.SpaceType.User, userId); space != nil {
		return space, nil
	}

	spaceID := utils.GenerateSfIDString()
	space := &driveModel.Space{
		ID:         spaceID,
		SpaceType:  constant.SpaceType.User,
		TargetID:   userId,
		OwnerID:    userId,
		SpaceName:  "我的云盘",
		QuotaBytes: config.C.Storage.SpaceQuota,
		Status:     constant.SpaceStatus.Active,
	}

	err := db.RDB.Transaction(func(tx *gorm.DB) error {
		if err := driveDao.SpaceDao.Create(tx, space); err != nil {
			return err
		}
		return driveDao.SpaceMemberDao.Create(tx, &driveModel.SpaceMember{
			ID:         utils.GenerateSfIDString(),
			SpaceID:    spaceID,
			UserID:     userId,
			MemberRole: constant.SpaceMemberRole.Owner,
			Status:     constant.SpaceMemberStatus.Active,
		})
	})
	if err != nil {
		if exist := driveDao.SpaceDao.GetByTypeAndTargetID(db.RDB, constant.SpaceType.User, userId); exist != nil {
			return exist, nil
		}
		return nil, err
	}
	return space, nil
}

func (s *spaceService) ListUserSpaceFiles(userId string, parentId string) ([]*driveModel.SpaceFile, error) {
	space, err := s.GetOrCreateUserSpace(userId)
	if err != nil {
		return nil, err
	}
	if _, err := SpaceFileService.VerifyDirPermission(parentId, space.ID, userId); err != nil {
		return nil, err
	}
	return driveDao.SpaceFileDao.ListBySpaceIdAndParentId(db.RDB, space.ID, parentId)
}

// ListUserSpaceAllFiles 查询目录下全部文件（含子目录，仅文件不含目录）
func (s *spaceService) ListUserSpaceAllFiles(userId string, parentId string) ([]*driveModel.SpaceFile, error) {
	space, err := s.GetOrCreateUserSpace(userId)
	if err != nil {
		return nil, err
	}
	parentFile, err := SpaceFileService.VerifyDirPermission(parentId, space.ID, userId)
	if err != nil {
		return nil, err
	}

	pathPrefix := ""
	if parentFile != nil {
		pathPrefix = parentFile.Path
	}
	return driveDao.SpaceFileDao.ListFilesUnderPath(db.RDB, space.ID, pathPrefix)
}

// GetUserSpaceFileDetail 查询文件或文件夹详情
func (s *spaceService) GetUserSpaceFileDetail(userId, spaceFileId string) (*driveResult.SpaceFileDetail, error) {
	if spaceFileId == "" {
		return nil, errors.New("param.error")
	}
	space, err := s.GetOrCreateUserSpace(userId)
	if err != nil {
		return nil, err
	}

	file := driveDao.SpaceFileDao.GetById(db.RDB, spaceFileId)
	if file == nil || file.SpaceID != space.ID {
		return nil, errors.New("param.error")
	}

	location, err := s.buildLocation(file)
	if err != nil {
		return nil, err
	}

	detail := &driveResult.SpaceFileDetail{
		ID:           file.ID,
		FileName:     file.FileName,
		IsDir:        file.IsDir,
		FileType:     file.FileType,
		FileCategory: file.FileCategory,
		Location:     location,
		Size:         file.FileSize,
		Contains:     nil,
		UpdatedAt:    file.UpdatedAt,
	}

	if file.IsDir {
		totalSize, _, err := driveDao.SpaceFileDao.SumFileSizeSelfAndDescendants(db.RDB, space.ID, file.ID, file.Path)
		if err != nil {
			return nil, err
		}
		fileCount, folderCount, err := driveDao.SpaceFileDao.CountDescendants(db.RDB, space.ID, file.Path)
		if err != nil {
			return nil, err
		}
		detail.Size = totalSize
		detail.Contains = &driveResult.SpaceFileContains{
			FileCount:   fileCount,
			FolderCount: folderCount,
		}
	}
	return detail, nil
}

// buildLocation 根据 path 组装可读位置，如 /文档/工作
func (s *spaceService) buildLocation(file *driveModel.SpaceFile) (string, error) {
	if file.ParentID == "root" || file.Path == "" {
		return "/", nil
	}

	parts := strings.Split(strings.Trim(file.Path, "/"), "/")
	if len(parts) <= 1 {
		return "/", nil
	}
	// path 含自身 id，位置取父级名称路径
	ancestorIds := parts[:len(parts)-1]
	ancestors, err := driveDao.SpaceFileDao.ListByIds(db.RDB, ancestorIds)
	if err != nil {
		return "", err
	}
	nameMap := make(map[string]string, len(ancestors))
	for _, item := range ancestors {
		nameMap[item.ID] = item.FileName
	}

	names := make([]string, 0, len(ancestorIds))
	for _, id := range ancestorIds {
		name, ok := nameMap[id]
		if !ok || name == "" {
			continue
		}
		names = append(names, name)
	}
	if len(names) == 0 {
		return "/", nil
	}
	return "/" + strings.Join(names, "/"), nil
}

// ListUserSpaceDirTree 查询用户空间目录树（仅目录）
func (s *spaceService) ListUserSpaceDirTree(userId string) ([]*driveResult.SpaceDirTreeNode, error) {
	space, err := s.GetOrCreateUserSpace(userId)
	if err != nil {
		return nil, err
	}
	dirs, err := driveDao.SpaceFileDao.ListDirsBySpaceId(db.RDB, space.ID)
	if err != nil {
		return nil, err
	}

	nodeMap := make(map[string]*driveResult.SpaceDirTreeNode, len(dirs))
	roots := make([]*driveResult.SpaceDirTreeNode, 0)
	for _, dir := range dirs {
		node := &driveResult.SpaceDirTreeNode{
			ID:       dir.ID,
			FileName: dir.FileName,
			ParentID: dir.ParentID,
			Path:     dir.Path,
			Level:    dir.Level,
			Children: make([]*driveResult.SpaceDirTreeNode, 0),
		}
		nodeMap[dir.ID] = node
	}
	for _, dir := range dirs {
		node := nodeMap[dir.ID]
		if dir.ParentID == "root" {
			roots = append(roots, node)
			continue
		}
		if parent, ok := nodeMap[dir.ParentID]; ok {
			parent.Children = append(parent.Children, node)
			continue
		}
		// 父目录缺失时兜底挂到根，避免丢节点
		roots = append(roots, node)
	}
	return roots, nil
}

func (s *spaceService) CreateUserSpaceDir(userId string, parentId string, dirName string) (*driveModel.SpaceFile, error) {
	space, err := s.GetOrCreateUserSpace(userId)
	if err != nil {
		return nil, err
	}
	parentFile, err := SpaceFileService.VerifyDirPermission(parentId, space.ID, userId)
	if err != nil {
		return nil, err
	}

	path := ""
	level := 0
	resolvedParentId := "root"
	if parentFile != nil {
		path = parentFile.Path
		level = parentFile.Level
		resolvedParentId = parentFile.ID
	}
	id := utils.GenerateSfIDString()
	dir := &driveModel.SpaceFile{
		ID:                  id,
		SpaceID:             space.ID,
		UserID:              userId,
		PhysicalID:          "",
		PhysicalStoragePath: "",
		ParentID:            resolvedParentId,
		Path:                path + "/" + id,
		Level:               level + 1,
		FileName:            dirName,
		IsDir:               true,
		FileType:            constant.FileType.Folder,
		FileSize:            0,
	}
	if err := driveDao.SpaceFileDao.Create(db.RDB, dir); err != nil {
		return nil, err
	}
	return dir, nil
}

func (s *spaceService) DeleteUserSpaceFiles(userId string, spaceFileIds []string) error {
	space, err := s.GetOrCreateUserSpace(userId)
	if err != nil {
		return err
	}

	idSet := make(map[string]struct{}, len(spaceFileIds))
	uniqueIds := make([]string, 0, len(spaceFileIds))
	for _, id := range spaceFileIds {
		if id == "" {
			return errors.New("param.error")
		}
		if _, ok := idSet[id]; ok {
			continue
		}
		idSet[id] = struct{}{}
		uniqueIds = append(uniqueIds, id)
	}

	files, err := driveDao.SpaceFileDao.ListByIds(db.RDB, uniqueIds)
	if err != nil {
		return err
	}
	if len(files) != len(uniqueIds) {
		return errors.New("param.error")
	}
	for _, file := range files {
		if file.SpaceID != space.ID {
			return errors.New("param.error")
		}
	}

	totalSize, fileCount, err := s.CalcSpaceFilesUsedBytes(db.RDB, space.ID, files)
	if err != nil {
		return err
	}

	expireAt := localtime.LocalTime(time.Now().In(localtime.Location).
		AddDate(0, 0, constant.DefaultSpaceRecycleExpireDays))
	ids := make([]string, 0, len(files))
	recycles := make([]*driveModel.SpaceRecycle, 0, len(files))
	for _, file := range files {
		ids = append(ids, file.ID)
		recycles = append(recycles, &driveModel.SpaceRecycle{
			ID:          utils.GenerateSfIDString(),
			UserID:      userId,
			SpaceID:     space.ID,
			SpaceFileID: file.ID,
			DeletedBy:   userId,
			ExpireAt:    &expireAt,
		})
	}

	return db.RDB.Transaction(func(tx *gorm.DB) error {
		if err := driveDao.SpaceFileDao.DeleteByIds(tx, ids); err != nil {
			return err
		}
		if err := driveDao.SpaceRecycleDao.CreateBatch(tx, recycles); err != nil {
			return err
		}
		if totalSize == 0 && fileCount == 0 {
			return nil
		}
		return driveDao.SpaceDao.DecUsedBytesById(tx, space.ID, totalSize, fileCount)
	})
}

// TransferSaveUserSpaceFiles 将源文件/目录转存（复制）到当前用户空间指定目录
func (s *spaceService) TransferSaveUserSpaceFiles(userId string, spaceFileIds []string, targetDirId string) error {
	destSpace, err := s.GetOrCreateUserSpace(userId)
	if err != nil {
		return err
	}

	targetParent, err := SpaceFileService.VerifyDirPermission(targetDirId, destSpace.ID, userId)
	if err != nil {
		return errors.New("cloud-drive.space.target-dir-not-exist")
	}

	targetPathPrefix := ""
	targetLevel := 0
	resolvedParentId := "root"
	if targetParent != nil {
		targetPathPrefix = targetParent.Path
		targetLevel = targetParent.Level
		resolvedParentId = targetParent.ID
	}

	idSet := make(map[string]struct{}, len(spaceFileIds))
	uniqueIds := make([]string, 0, len(spaceFileIds))
	for _, id := range spaceFileIds {
		if id == "" {
			return errors.New("param.error")
		}
		if _, ok := idSet[id]; ok {
			continue
		}
		idSet[id] = struct{}{}
		uniqueIds = append(uniqueIds, id)
	}

	files, err := driveDao.SpaceFileDao.ListByIds(db.RDB, uniqueIds)
	if err != nil {
		return err
	}
	// ListByIds 默认排除软删除；查不到视为已过期/已删除
	if len(files) == 0 || len(files) != len(uniqueIds) {
		return errors.New("cloud-drive.space.transfer-file-expired")
	}

	spaceCache := make(map[string]*driveModel.Space)
	for _, file := range files {
		space, ok := spaceCache[file.SpaceID]
		if !ok {
			space = driveDao.SpaceDao.GetById(db.RDB, file.SpaceID)
			if space == nil || space.SpaceType != constant.SpaceType.User {
				return errors.New("cloud-drive.space.transfer-source-denied")
			}
			spaceCache[file.SpaceID] = space
		}
		ownerId := space.OwnerID
		if ownerId == "" {
			ownerId = space.TargetID
		}
		if ownerId != userId && !basicService.ContactsService.IsFriend(userId, ownerId) {
			return errors.New("cloud-drive.space.transfer-source-denied")
		}
	}

	topFiles := filterTopLevelSpaceFiles(files)

	var totalSize, fileCount int64
	for _, file := range topFiles {
		size, count, calcErr := s.CalcSpaceFilesUsedBytes(db.RDB, file.SpaceID, []*driveModel.SpaceFile{file})
		if calcErr != nil {
			return calcErr
		}
		totalSize += size
		fileCount += count
	}
	if err := s.CheckSpaceQuota(destSpace, totalSize); err != nil {
		return err
	}

	return db.RDB.Transaction(func(tx *gorm.DB) error {
		var copiedSize, copiedCount int64
		for _, file := range topFiles {
			size, count, copyErr := s.copySpaceFileTree(tx, file, destSpace.ID, userId, resolvedParentId, targetPathPrefix, targetLevel)
			if copyErr != nil {
				return copyErr
			}
			copiedSize += size
			copiedCount += count
		}
		if copiedSize == 0 && copiedCount == 0 {
			return nil
		}
		return driveDao.SpaceDao.IncUsedBytesById(tx, destSpace.ID, copiedSize, copiedCount)
	})
}

// copySpaceFileTree 将单个文件或整棵目录树复制到目标空间目录下，复用物理文件引用
func (s *spaceService) copySpaceFileTree(
	tx *gorm.DB,
	src *driveModel.SpaceFile,
	destSpaceId, userId, destParentId, destParentPath string,
	destParentLevel int,
) (int64, int64, error) {
	nodes := []*driveModel.SpaceFile{src}
	if src.IsDir {
		list, err := driveDao.SpaceFileDao.ListSelfAndDescendants(tx, src.SpaceID, src.ID, src.Path)
		if err != nil {
			return 0, 0, err
		}
		nodes = list
	}

	idMap := make(map[string]string, len(nodes))
	pathMap := make(map[string]string, len(nodes))
	levelMap := make(map[string]int, len(nodes))
	var totalSize, fileCount int64

	for _, node := range nodes {
		newId := utils.GenerateSfIDString()
		idMap[node.ID] = newId

		var newParentId, newPath string
		var newLevel int
		if node.ID == src.ID {
			newParentId = destParentId
			newPath = destParentPath + "/" + newId
			newLevel = destParentLevel + 1
		} else {
			mappedParent, ok := idMap[node.ParentID]
			if !ok {
				return 0, 0, errors.New("param.error")
			}
			newParentId = mappedParent
			newPath = pathMap[node.ParentID] + "/" + newId
			newLevel = levelMap[node.ParentID] + 1
		}
		pathMap[node.ID] = newPath
		levelMap[node.ID] = newLevel

		clone := &driveModel.SpaceFile{
			ID:                  newId,
			SpaceID:             destSpaceId,
			UserID:              userId,
			PhysicalID:          node.PhysicalID,
			PhysicalStoragePath: node.PhysicalStoragePath,
			ParentID:            newParentId,
			Path:                newPath,
			Level:               newLevel,
			FileName:            node.FileName,
			IsDir:               node.IsDir,
			FileType:            node.FileType,
			FileCategory:        node.FileCategory,
			FileSize:            node.FileSize,
			Status:              node.Status,
		}
		if err := driveDao.SpaceFileDao.Create(tx, clone); err != nil {
			return 0, 0, err
		}
		if !node.IsDir && node.PhysicalID != "" {
			if err := driveDao.PhysicalFileDao.FileRefIncById(tx, node.PhysicalID); err != nil {
				return 0, 0, err
			}
			totalSize += node.FileSize
			fileCount++
		}
	}
	return totalSize, fileCount, nil
}

func (s *spaceService) MoveUserSpaceFiles(userId string, spaceFileIds []string, targetParentId string) error {
	space, err := s.GetOrCreateUserSpace(userId)
	if err != nil {
		return err
	}

	targetParent, err := SpaceFileService.VerifyDirPermission(targetParentId, space.ID, userId)
	if err != nil {
		return err
	}

	targetPathPrefix := ""
	targetLevel := 0
	resolvedParentId := "root"
	if targetParent != nil {
		targetPathPrefix = targetParent.Path
		targetLevel = targetParent.Level
		resolvedParentId = targetParent.ID
	}

	idSet := make(map[string]struct{}, len(spaceFileIds))
	uniqueIds := make([]string, 0, len(spaceFileIds))
	for _, id := range spaceFileIds {
		if id == "" {
			return errors.New("param.error")
		}
		if _, ok := idSet[id]; ok {
			continue
		}
		idSet[id] = struct{}{}
		uniqueIds = append(uniqueIds, id)
	}

	files, err := driveDao.SpaceFileDao.ListByIds(db.RDB, uniqueIds)
	if err != nil {
		return err
	}
	if len(files) != len(uniqueIds) {
		return errors.New("param.error")
	}
	for _, file := range files {
		if file.SpaceID != space.ID {
			return errors.New("param.error")
		}
		// 不能移动到自身
		if file.ID == resolvedParentId {
			return errors.New("param.error")
		}
		// 不能把目录移动到自己的子孙目录下
		if file.IsDir && targetPathPrefix != "" &&
			(targetPathPrefix == file.Path || strings.HasPrefix(targetPathPrefix, file.Path+"/")) {
			return errors.New("param.error")
		}
	}

	// 只移动顶层项，避免父子同时选中时重复处理
	topFiles := filterTopLevelSpaceFiles(files)

	return db.RDB.Transaction(func(tx *gorm.DB) error {
		for _, file := range topFiles {
			if file.ParentID == resolvedParentId {
				continue
			}
			oldPath := file.Path
			oldLevel := file.Level
			newPath := targetPathPrefix + "/" + file.ID
			newLevel := targetLevel + 1

			if err := driveDao.SpaceFileDao.UpdateMove(tx, file.ID, resolvedParentId, newPath, newLevel); err != nil {
				return err
			}
			if file.IsDir {
				levelDelta := newLevel - oldLevel
				if err := driveDao.SpaceFileDao.UpdateDescendantsPathAndLevel(tx, space.ID, oldPath, newPath, levelDelta); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// RenameUserSpaceFile 重命名文件或目录
func (s *spaceService) RenameUserSpaceFile(userId, spaceFileId, newName string) (*driveModel.SpaceFile, error) {
	newName = strings.TrimSpace(newName)
	if spaceFileId == "" || newName == "" {
		return nil, errors.New("param.error")
	}

	space, err := s.GetOrCreateUserSpace(userId)
	if err != nil {
		return nil, err
	}

	file := driveDao.SpaceFileDao.GetById(db.RDB, spaceFileId)
	if file == nil || file.SpaceID != space.ID {
		return nil, errors.New("param.error")
	}
	if file.FileName == newName {
		return file, nil
	}

	fileType := file.FileType
	fileCategory := file.FileCategory
	updateTypeMeta := false
	if !file.IsDir {
		fileType = strings.TrimPrefix(filepath.Ext(newName), ".")
		fileCategory = utils.FileCategoryFromExt(fileType)
		updateTypeMeta = true
	}

	if err := driveDao.SpaceFileDao.UpdateName(db.RDB, file.ID, newName, fileType, fileCategory, updateTypeMeta); err != nil {
		return nil, err
	}

	file.FileName = newName
	if updateTypeMeta {
		file.FileType = fileType
		file.FileCategory = fileCategory
	}
	return file, nil
}

func (s *spaceService) ListUserSpaceCategoryStats(userId string) ([]*driveResult.SpaceFileCategoryStat, error) {
	space, err := s.GetOrCreateUserSpace(userId)
	if err != nil {
		return nil, err
	}
	stats, err := driveDao.SpaceFileDao.StatByCategory(db.RDB, space.ID)
	if err != nil {
		return nil, err
	}

	statMap := make(map[string]*driveResult.SpaceFileCategoryStat, len(stats))
	for _, stat := range stats {
		statMap[stat.FileCategory] = stat
	}

	categories := []string{
		constant.FileCategory.Image,
		constant.FileCategory.Video,
		constant.FileCategory.Document,
		constant.FileCategory.Audio,
		constant.FileCategory.Archive,
		constant.FileCategory.Other,
	}
	result := make([]*driveResult.SpaceFileCategoryStat, 0, len(categories))
	for _, category := range categories {
		if stat, ok := statMap[category]; ok {
			result = append(result, stat)
			continue
		}
		result = append(result, &driveResult.SpaceFileCategoryStat{
			FileCategory: category,
			FileCount:    0,
			TotalSize:    0,
		})
	}
	return result, nil
}
