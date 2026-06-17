package storage

import (
	"fmt"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 临时目录
var tempDir = "./temp/chunks/%s"

// GetUploadChunkInfo 获取切片信息
func GetUploadChunkInfo(fileHash string) []string {
	key := fmt.Sprintf(constant.RedisKey.UploadChunk, fileHash)
	chunks, err := db.CacheDB.SMembers(key)
	if err != nil {
		return []string{}
	}
	return chunks
}

// UploadChunk 切片上传
func UploadChunk(chunkFile *multipart.FileHeader, chunkIndex string, fileHash string) error {
	// 创建临时目录并保存文件
	dir := fmt.Sprintf(tempDir, fileHash)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	savePath := fmt.Sprintf("%s/%s.part", dir, chunkIndex)

	if _, err := os.Stat(savePath); err == nil {
		// 物理文件已存在
		key := fmt.Sprintf(constant.RedisKey.UploadChunk, fileHash)
		return db.CacheDB.SAdd(key, 60*time.Minute, chunkIndex)
	}
	if err := utils.SaveUploadedFile(chunkFile, savePath); err != nil {
		return err
	}
	key := fmt.Sprintf(constant.RedisKey.UploadChunk, fileHash)
	return db.CacheDB.SAdd(key, 60*time.Minute, chunkIndex)
}

// MergeChunk 切片合并
func MergeChunk(fileHash string, totalChunks int, fileName string, storagePrefix string) (error, string) {
	chunkDir := fmt.Sprintf(tempDir, fileHash)
	fileBase := fileHash
	if strings.HasPrefix(storagePrefix, "msgfile/") {
		fileBase = fmt.Sprintf("%s-%d", fileHash, time.Now().UnixNano())
	}
	fileKey := fmt.Sprintf("%s/%s%s", storagePrefix, fileBase, filepath.Ext(fileName))

	storagePath, err := S.Merge(fileKey, chunkDir, totalChunks)
	if err != nil {
		return err, ""
	}
	// 上传成功后清理分片
	key := fmt.Sprintf(constant.RedisKey.UploadChunk, fileHash)
	if err = db.CacheDB.Del(key); err != nil {
		return err, ""
	}
	if err = os.RemoveAll(chunkDir); err != nil {
		return err, ""
	}
	return nil, storagePath
}
