package repository

import (
	"bytes"
	"context"
	"webby/room-service/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type FileStorage struct {
	region  string
	bucket  string
	storage *s3.Client
}

func NewFileStorage(cfg *config.Config) *FileStorage {
	options := s3.Options{
		Region:      cfg.Aws.Region,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(cfg.Aws.AccessKey, cfg.Aws.SecretKey, "")),
	}

	client := s3.New(options)

	return &FileStorage{
		region:  cfg.Aws.Region,
		bucket:  cfg.Aws.Bucket,
		storage: client,
	}
}

func (fs *FileStorage) Save(ctx context.Context, key string, data []byte) (string, error) {
	if _, err := fs.storage.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	}); err != nil {
		return "", err
	}

	fileURL := "https://" + fs.bucket + ".s3." + fs.region + ".amazonaws.com/" + key

	return fileURL, nil
}

func (fs *FileStorage) Remove(ctx context.Context, key string) error {
	if _, err := fs.storage.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(key),
	}); err != nil {
		return err
	}

	return nil
}
