package storage

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"strings"

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
	//初始化
	imur, err := bucket.InitiateMultipartUpload(fileKey)
	if err != nil {
		return "", fmt.Errorf("failed to initiate multipart upload: %v", err)
	}
	var parts []oss.UploadPart
	//逐个上传
	for i := 0; i < totalChunks; i++ {
		partPath := fmt.Sprintf("%s/%d.part", chunkDir, i)

		fileInfo, err := os.Stat(partPath)
		if err != nil {
			bucket.AbortMultipartUpload(imur)
			return "", fmt.Errorf("failed to stat part %d: %v", i, err)
		}
		partSize := fileInfo.Size()

		part, err := bucket.UploadPartFromFile(imur, partPath, 0, partSize, i)
		if err != nil {
			bucket.AbortMultipartUpload(imur)
			return "", fmt.Errorf("failed to upload part %d: %v", i, err)
		}
		parts = append(parts, part)
	}

	_, err = bucket.CompleteMultipartUpload(imur, parts)
	if err != nil {
		bucket.AbortMultipartUpload(imur)
		return "", fmt.Errorf("failed to complete multipart upload: %v", err)
	}

	return s.GetURL(fileKey, 0)
}
