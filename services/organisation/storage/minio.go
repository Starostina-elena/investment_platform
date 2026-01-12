package storage

import (
	"bytes"
	"context"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	client       *minio.Client
	avatarBucket string
	docsBucket   string
}

func NewMinioStorage(endpoint, accessKey, secretKey string, useSSL bool, avatarBucket string, docsBucket string) (*MinioStorage, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize minio client: %w", err)
	}

	storage := &MinioStorage{
		client:       minioClient,
		avatarBucket: avatarBucket,
		docsBucket:   docsBucket,
	}

	if err := storage.createBucketIfNotExists(context.Background(), storage.avatarBucket, true); err != nil {
		return nil, err
	}
	if err := storage.createBucketIfNotExists(context.Background(), storage.docsBucket, false); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *MinioStorage) createBucketIfNotExists(ctx context.Context, bucket string, makePublic bool) error {
	exists, err := s.client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = s.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Printf("Bucket %s created successfully", bucket)
	}

	if makePublic {
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
	}`, bucket)

		if err := s.client.SetBucketPolicy(ctx, bucket, policy); err != nil {
			log.Printf("Warning: failed to set bucket policy: %v", err)
		}
	}

	return nil
}

func (s *MinioStorage) PutAvatar(ctx context.Context, objectName string, data *bytes.Buffer, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, s.avatarBucket, objectName, data, int64(data.Len()), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to put object: %w", err)
	}
	return fmt.Sprintf("%s/%s", s.avatarBucket, objectName), nil
}

func (s *MinioStorage) DeleteAvatar(ctx context.Context, avatarPath string) error {
	return s.deleteObject(ctx, avatarPath)
}

func (s *MinioStorage) PutDoc(ctx context.Context, objectName string, data *bytes.Buffer, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, s.docsBucket, objectName, data, int64(data.Len()), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to put object: %w", err)
	}
	return fmt.Sprintf("%s/%s", s.docsBucket, objectName), nil
}

func (s *MinioStorage) GetObject(ctx context.Context, objectPath string) ([]byte, string, error) {
	bucket, objectName, err := splitPath(objectPath)
	if err != nil {
		return nil, "", err
	}
	obj, err := s.client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, "", fmt.Errorf("failed to get object: %w", err)
	}
	defer obj.Close()

	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read object: %w", err)
	}
	info, err := obj.Stat()
	if err != nil {
		return nil, "", fmt.Errorf("failed to stat object: %w", err)
	}

	return data, info.ContentType, nil
}

func (s *MinioStorage) DeleteDoc(ctx context.Context, objectPath string) error {
	return s.deleteObject(ctx, objectPath)
}

func (s *MinioStorage) deleteObject(ctx context.Context, objectPath string) error {
	bucket, objectName, err := splitPath(objectPath)
	if err != nil {
		return err
	}

	err = s.client.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

func splitPath(path string) (string, string, error) {
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid object path format")
	}
	return parts[0], parts[1], nil
}
