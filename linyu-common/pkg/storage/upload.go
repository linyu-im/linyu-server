package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
)

// 临时目录（本地存储合并用）
var tempDir = "./temp/chunks/%s"

const uploadChunkTTL = 60 * time.Minute

// CompletedPartInfo 分片完成信息
type CompletedPartInfo struct {
	PartNumber int32
	ETag       string
}

// MultipartStorage 对象存储分片直传：UploadPart 写入最终对象，合并仅 Complete
type MultipartStorage interface {
	InitMultipart(fileKey string) (uploadID string, err error)
	UploadPart(fileKey, uploadID string, partNumber int32, reader io.Reader, size int64) (etag string, err error)
	CompleteMultipart(fileKey, uploadID string, parts []CompletedPartInfo) (url string, err error)
	AbortMultipart(fileKey, uploadID string) error
}

type multipartUploadSession struct {
	UploadID string `json:"uploadId"`
	FileKey  string `json:"fileKey"`
}

// GetUploadChunkInfo 获取切片信息
func GetUploadChunkInfo(fileHash string) []string {
	key := fmt.Sprintf(constant.RedisKey.UploadChunk, fileHash)
	chunks, err := db.CacheDB.SMembers(key)
	if err != nil {
		return []string{}
	}
	return chunks
}

// UploadChunk 切片上传（S3/OSS 直传分片；本地落盘）
func UploadChunk(chunkFile *multipart.FileHeader, chunkIndex, fileHash, fileName, storagePrefix string) error {
	if ms, ok := S.(MultipartStorage); ok {
		return uploadChunkMultipart(ms, chunkFile, chunkIndex, fileHash, fileName, storagePrefix)
	}
	return uploadChunkToLocal(chunkFile, chunkIndex, fileHash)
}

func uploadChunkMultipart(ms MultipartStorage, chunkFile *multipart.FileHeader, chunkIndex, fileHash, fileName, storagePrefix string) error {
	if storagePrefix == "" {
		return fmt.Errorf("storagePrefix is required")
	}
	if fileName == "" {
		fileName = fileHash
	}

	session, err := getOrCreateMultipartSession(ms, fileHash, fileName, storagePrefix)
	if err != nil {
		return err
	}

	partKey := fmt.Sprintf(constant.RedisKey.UploadPartETag, fileHash, chunkIndex)
	if etag, err := db.CacheDB.Get(partKey); err == nil && etag != "" {
		key := fmt.Sprintf(constant.RedisKey.UploadChunk, fileHash)
		return db.CacheDB.SAdd(key, uploadChunkTTL, chunkIndex)
	}

	partIndex, err := strconv.ParseInt(chunkIndex, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid chunkIndex: %s", chunkIndex)
	}
	partNumber := int32(partIndex + 1)

	file, err := chunkFile.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	etag, err := ms.UploadPart(session.FileKey, session.UploadID, partNumber, file, chunkFile.Size)
	if err != nil {
		return err
	}
	if err = db.CacheDB.Set(partKey, etag, uploadChunkTTL); err != nil {
		return err
	}

	key := fmt.Sprintf(constant.RedisKey.UploadChunk, fileHash)
	return db.CacheDB.SAdd(key, uploadChunkTTL, chunkIndex)
}

func getOrCreateMultipartSession(ms MultipartStorage, fileHash, fileName, storagePrefix string) (*multipartUploadSession, error) {
	metaKey := fmt.Sprintf(constant.RedisKey.UploadMultipart, fileHash)
	session := &multipartUploadSession{}
	if err := db.CacheDB.GetObject(metaKey, session); err == nil && session.UploadID != "" && session.FileKey != "" {
		return session, nil
	}

	fileKey := buildUploadFileKey(fileHash, fileName, storagePrefix)
	uploadID, err := ms.InitMultipart(fileKey)
	if err != nil {
		return nil, err
	}
	session = &multipartUploadSession{
		UploadID: uploadID,
		FileKey:  fileKey,
	}
	if err = db.CacheDB.SetObject(metaKey, session, uploadChunkTTL); err != nil {
		_ = ms.AbortMultipart(fileKey, uploadID)
		return nil, err
	}
	return session, nil
}

func buildUploadFileKey(fileHash, fileName, storagePrefix string) string {
	fileBase := fileHash
	if strings.HasPrefix(storagePrefix, "msgfile/") {
		fileBase = fmt.Sprintf("%s-%d", fileHash, time.Now().UnixNano())
	}
	return fmt.Sprintf("%s/%s%s", storagePrefix, fileBase, filepath.Ext(fileName))
}

func uploadChunkToLocal(chunkFile *multipart.FileHeader, chunkIndex, fileHash string) error {
	dir := fmt.Sprintf(tempDir, fileHash)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	savePath := fmt.Sprintf("%s/%s.part", dir, chunkIndex)

	if _, err := os.Stat(savePath); err == nil {
		key := fmt.Sprintf(constant.RedisKey.UploadChunk, fileHash)
		return db.CacheDB.SAdd(key, uploadChunkTTL, chunkIndex)
	}
	if err := utils.SaveUploadedFile(chunkFile, savePath); err != nil {
		return err
	}
	key := fmt.Sprintf(constant.RedisKey.UploadChunk, fileHash)
	return db.CacheDB.SAdd(key, uploadChunkTTL, chunkIndex)
}

// MergeChunk 切片合并（S3/OSS 仅 CompleteMultipart，本地再拼盘）
func MergeChunk(fileHash string, totalChunks int, fileName string, storagePrefix string) (error, string) {
	if ms, ok := S.(MultipartStorage); ok {
		return mergeChunkMultipart(ms, fileHash, totalChunks)
	}

	fileKey := buildUploadFileKey(fileHash, fileName, storagePrefix)
	chunkDir := fmt.Sprintf(tempDir, fileHash)
	storagePath, err := S.Merge(fileKey, chunkDir, totalChunks)
	if err != nil {
		return err, ""
	}
	_ = os.RemoveAll(chunkDir)
	key := fmt.Sprintf(constant.RedisKey.UploadChunk, fileHash)
	if err = db.CacheDB.Del(key); err != nil {
		return err, ""
	}
	return nil, storagePath
}

func mergeChunkMultipart(ms MultipartStorage, fileHash string, totalChunks int) (error, string) {
	if totalChunks <= 0 {
		return fmt.Errorf("invalid totalChunks: %d", totalChunks), ""
	}

	metaKey := fmt.Sprintf(constant.RedisKey.UploadMultipart, fileHash)
	session := &multipartUploadSession{}
	if err := db.CacheDB.GetObject(metaKey, session); err != nil || session.UploadID == "" || session.FileKey == "" {
		return fmt.Errorf("multipart session not found, please re-upload chunks"), ""
	}

	parts := make([]CompletedPartInfo, totalChunks)
	for i := 0; i < totalChunks; i++ {
		partKey := fmt.Sprintf(constant.RedisKey.UploadPartETag, fileHash, strconv.Itoa(i))
		etag, err := db.CacheDB.Get(partKey)
		if err != nil || etag == "" {
			return fmt.Errorf("missing uploaded part: %d", i), ""
		}
		parts[i] = CompletedPartInfo{
			PartNumber: int32(i + 1),
			ETag:       etag,
		}
	}

	storagePath, err := ms.CompleteMultipart(session.FileKey, session.UploadID, parts)
	if err != nil {
		return err, ""
	}

	cleanupMultipartSession(fileHash, totalChunks)
	return nil, storagePath
}

func cleanupMultipartSession(fileHash string, totalChunks int) {
	_ = db.CacheDB.Del(fmt.Sprintf(constant.RedisKey.UploadChunk, fileHash))
	_ = db.CacheDB.Del(fmt.Sprintf(constant.RedisKey.UploadMultipart, fileHash))
	for i := 0; i < totalChunks; i++ {
		_ = db.CacheDB.Del(fmt.Sprintf(constant.RedisKey.UploadPartETag, fileHash, strconv.Itoa(i)))
	}
}
