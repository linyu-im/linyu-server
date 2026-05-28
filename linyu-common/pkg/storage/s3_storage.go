package storage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	s3config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
	"io"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

type S3Storage struct {
	BaseURL    string
	BucketName string
	Endpoint   string
	Client     *s3.Client
}

func NewS3Storage(c config.S3StorageConfig) *S3Storage {

	if c.Endpoint == "" || c.AccessKeyID == "" || c.AccessKeySecret == "" || c.BucketName == "" {
		panic("s3 storage required config: Endpoint, AccessKeyID, AccessKeySecret, BucketName must be set")
	}

	cfg, err := s3config.LoadDefaultConfig(context.TODO(),
		s3config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(c.AccessKeyID, c.AccessKeySecret, ""),
		),
		s3config.WithRegion(c.Region),
	)

	if err != nil {
		panic("failed to load aws config: " + err.Error())
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if c.Endpoint != "" {
			o.BaseEndpoint = aws.String(c.Endpoint)
		}
		o.UsePathStyle = true
	})

	return &S3Storage{
		BaseURL:    c.BaseURL,
		BucketName: c.BucketName,
		Endpoint:   c.Endpoint,
		Client:     client,
	}
}

func (s *S3Storage) Upload(fileKey string, reader io.Reader) (string, error) {

	_, err := s.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(fileKey),
		Body:   reader,
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file to s3: %v", err)
	}

	return s.GetURL(fileKey, 0)
}

func (s *S3Storage) Download(fileKey string, writer io.Writer) error {

	resp, err := s.Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(fileKey),
	})

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(writer, resp.Body)
	return err
}

func (s *S3Storage) Delete(fileKey string) error {

	if fileKey == "" {
		return nil
	}

	_, err := s.Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(fileKey),
	})

	if err != nil {
		return fmt.Errorf("failed to delete s3 object: %v", err)
	}

	return nil
}

func (s *S3Storage) GetURL(fileKey string, expire int64) (string, error) {

	if fileKey == "" {
		return "", nil
	}

	if s.BaseURL != "" {
		u, err := url.Parse(s.BaseURL)
		if err != nil {
			return "", err
		}
		u.Path = path.Join(u.Path, fileKey)
		if strings.HasSuffix(fileKey, "/") && !strings.HasSuffix(u.Path, "/") {
			u.Path += "/"
		}
		return u.String(), nil
	}

	if expire == 0 {
		return fmt.Sprintf("%s/%s/%s", s.Endpoint, s.BucketName, fileKey), nil
	}

	presignClient := s3.NewPresignClient(s.Client)
	req, err := presignClient.PresignGetObject(
		context.TODO(),
		&s3.GetObjectInput{
			Bucket: aws.String(s.BucketName),
			Key:    aws.String(fileKey),
		},
		s3.WithPresignExpires(time.Duration(expire)*time.Second),
	)

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned url: %v", err)
	}

	return req.URL, nil
}

func (s *S3Storage) Merge(fileKey string, chunkDir string, totalChunks int) (string, error) {

	ctx := context.TODO()

	createResp, err := s.Client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(fileKey),
	})
	if err != nil {
		return "", err
	}

	uploadID := *createResp.UploadId

	var completedParts []types.CompletedPart

	for i := 0; i < totalChunks; i++ {

		partPath := fmt.Sprintf("%s/%d.part", chunkDir, i)

		file, err := os.Open(partPath)
		if err != nil {
			return "", err
		}

		resp, err := s.Client.UploadPart(ctx, &s3.UploadPartInput{
			Bucket:     aws.String(s.BucketName),
			Key:        aws.String(fileKey),
			UploadId:   aws.String(uploadID),
			PartNumber: utils.Int32Ptr(int32(i)),
			Body:       file,
		})

		file.Close()

		if err != nil {
			return "", err
		}

		completedParts = append(completedParts, types.CompletedPart{
			ETag:       resp.ETag,
			PartNumber: utils.Int32Ptr(int32(i)),
		})
	}

	_, err = s.Client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(s.BucketName),
		Key:      aws.String(fileKey),
		UploadId: aws.String(uploadID),
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})

	if err != nil {
		return "", err
	}

	return s.GetURL(fileKey, 0)
}
