package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
)

type OssStorage struct {
	BaseURL    string
	BucketName string
	Endpoint   string
	Client     *oss.Client
}

func NewOssStorage(c config.OssStorageConfig) *OssStorage {
	if c.Endpoint == "" || c.AccessKeyID == "" || c.AccessKeySecret == "" || c.BucketName == "" {
		panic("oss storage required config: Endpoint, AccessKeyID, AccessKeySecret, BucketName must be set")
	}

	client, err := oss.New(c.Endpoint, c.AccessKeyID, c.AccessKeySecret)
	if err != nil {
		panic("failed to create oss client: " + err.Error())
	}

	return &OssStorage{
		BaseURL:    c.BaseURL,
		BucketName: c.BucketName,
		Endpoint:   c.Endpoint,
		Client:     client,
	}
}

func (s *OssStorage) getBucket() (*oss.Bucket, error) {
	return s.Client.Bucket(s.BucketName)
}

func (s *OssStorage) Upload(fileKey string, reader io.Reader) (string, error) {
	bucket, err := s.getBucket()
	if err != nil {
		return "", err
	}

	// 上传文件流
	err = bucket.PutObject(fileKey, reader)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to oss: %v", err)
	}

	return s.GetURL(fileKey, 0)
}

func (s *OssStorage) Download(fileKey string, writer io.Writer) error {
	bucket, err := s.getBucket()
	if err != nil {
		return err
	}

	body, err := bucket.GetObject(fileKey)
	if err != nil {
		return err
	}
	defer body.Close()

	_, err = io.Copy(writer, body)
	return err
}

func (s *OssStorage) Delete(fileKey string) error {
	if fileKey == "" {
		return nil
	}

	bucket, err := s.getBucket()
	if err != nil {
		return err
	}

	err = bucket.DeleteObject(fileKey)
	if err != nil {
		return fmt.Errorf("failed to delete oss object: %v", err)
	}

	return nil
}

func (s *OssStorage) GetURL(fileKey string, expire int64) (string, error) {
	if fileKey == "" {
		return "", nil
	}

	if s.BaseURL != "" {
		u, err := url.Parse(s.BaseURL)
		if err != nil {
			return "", fmt.Errorf("invalid BaseURL: %v", err)
		}

		u.Path = path.Join(u.Path, fileKey)

		if strings.HasSuffix(fileKey, "/") && !strings.HasSuffix(u.Path, "/") {
			u.Path += "/"
		}
		return u.String(), nil
	}

	bucket, err := s.getBucket()
	if err != nil {
		return "", err
	}

	if expire > 0 {
		signedURL, err := bucket.SignURL(fileKey, oss.HTTPGet, expire)
		if err != nil {
			return "", fmt.Errorf("failed to sign OSS URL: %v", err)
		}
		return signedURL, nil
	}

	return fmt.Sprintf("https://%s.%s/%s", s.BucketName, s.Endpoint, fileKey), nil
}

func (s *OssStorage) Merge(fileKey string, chunkDir string, totalChunks int) (string, error) {
	bucket, err := s.getBucket()
	if err != nil {
		return "", err
	}

	imur, err := bucket.InitiateMultipartUpload(fileKey)
	if err != nil {
		return "", fmt.Errorf("failed to initiate multipart upload: %v", err)
	}

	parts := make([]oss.UploadPart, totalChunks)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		wg        sync.WaitGroup
		errOnce   sync.Once
		uploadErr error
		sem       = make(chan struct{}, mergeUploadConcurrency)
	)

	for i := 0; i < totalChunks; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				return
			}
			defer func() { <-sem }()

			if ctx.Err() != nil {
				return
			}

			partPath := fmt.Sprintf("%s/%d.part", chunkDir, i)
			fileInfo, err := os.Stat(partPath)
			if err != nil {
				errOnce.Do(func() {
					uploadErr = fmt.Errorf("failed to stat part %d: %v", i, err)
					cancel()
				})
				return
			}

			partNumber := i + 1
			part, err := bucket.UploadPartFromFile(imur, partPath, 0, fileInfo.Size(), partNumber)
			if err != nil {
				errOnce.Do(func() {
					uploadErr = fmt.Errorf("failed to upload part %d: %v", i, err)
					cancel()
				})
				return
			}
			parts[i] = part
		}(i)
	}
	wg.Wait()

	if uploadErr != nil {
		_ = bucket.AbortMultipartUpload(imur)
		return "", uploadErr
	}

	_, err = bucket.CompleteMultipartUpload(imur, parts)
	if err != nil {
		_ = bucket.AbortMultipartUpload(imur)
		return "", fmt.Errorf("failed to complete multipart upload: %v", err)
	}

	return s.GetURL(fileKey, 0)
}

func (s *OssStorage) InitMultipart(fileKey string) (string, error) {
	bucket, err := s.getBucket()
	if err != nil {
		return "", err
	}
	imur, err := bucket.InitiateMultipartUpload(fileKey)
	if err != nil {
		return "", fmt.Errorf("failed to init multipart upload: %v", err)
	}
	return imur.UploadID, nil
}

func (s *OssStorage) UploadPart(fileKey, uploadID string, partNumber int32, reader io.Reader, size int64) (string, error) {
	bucket, err := s.getBucket()
	if err != nil {
		return "", err
	}
	imur := oss.InitiateMultipartUploadResult{
		Bucket:   s.BucketName,
		Key:      fileKey,
		UploadID: uploadID,
	}
	part, err := bucket.UploadPart(imur, reader, size, int(partNumber))
	if err != nil {
		return "", fmt.Errorf("failed to upload part %d: %v", partNumber, err)
	}
	return part.ETag, nil
}

func (s *OssStorage) CompleteMultipart(fileKey, uploadID string, parts []CompletedPartInfo) (string, error) {
	bucket, err := s.getBucket()
	if err != nil {
		return "", err
	}
	imur := oss.InitiateMultipartUploadResult{
		Bucket:   s.BucketName,
		Key:      fileKey,
		UploadID: uploadID,
	}
	ossParts := make([]oss.UploadPart, len(parts))
	for i, part := range parts {
		ossParts[i] = oss.UploadPart{
			PartNumber: int(part.PartNumber),
			ETag:       part.ETag,
		}
	}
	_, err = bucket.CompleteMultipartUpload(imur, ossParts)
	if err != nil {
		_ = s.AbortMultipart(fileKey, uploadID)
		return "", fmt.Errorf("failed to complete multipart upload: %v", err)
	}
	return s.GetURL(fileKey, 0)
}

func (s *OssStorage) AbortMultipart(fileKey, uploadID string) error {
	bucket, err := s.getBucket()
	if err != nil {
		return err
	}
	imur := oss.InitiateMultipartUploadResult{
		Bucket:   s.BucketName,
		Key:      fileKey,
		UploadID: uploadID,
	}
	return bucket.AbortMultipartUpload(imur)
}
