package storage

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage struct {
	client *minio.Client
	bucket string
}

func NewStorage(config Config) (*Storage, error) {
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(config.Credentials.AccessKey, config.Credentials.SecretKey, ""),
	})
	if err != nil {
		return nil, err
	}

	return &Storage{client: client, bucket: config.Bucket}, nil
}
