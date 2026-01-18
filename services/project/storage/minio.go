package storage

import (
	"bytes"
	"context"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	client     *minio.Client
	bucketName string
}

func NewMinioStorage(endpoint, accessKey, secretKey string, useSSL bool, bucketName string) (*MinioStorage, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize minio client: %w", err)
	}

	storage := &MinioStorage{
		client:     minioClient,
		bucketName: bucketName,
	}

	if err := storage.createBucketIfNotExists(context.Background()); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *MinioStorage) createBucketIfNotExists(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = s.client.MakeBucket(ctx, s.bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Printf("Bucket %s created successfully", s.bucketName)
	}

	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}
		]
	}`, s.bucketName)

	if err := s.client.SetBucketPolicy(ctx, s.bucketName, policy); err != nil {
		log.Printf("Warning: failed to set bucket policy: %v", err)
	}

	return nil
}

func (s *MinioStorage) PutObject(ctx context.Context, objectName string, data *bytes.Buffer, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, s.bucketName, objectName, data, int64(data.Len()), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to put object: %w", err)
	}
	return fmt.Sprintf("%s/%s", s.bucketName, objectName), nil
}

func (s *MinioStorage) DeletePicture(ctx context.Context, picturePath string) error {
	parts := strings.SplitN(picturePath, "/", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid picture path format")
	}
	objectName := parts[1]

	err := s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}
