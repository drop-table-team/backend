package storage

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
)

func (s Storage) UploadFile(file bytes.Buffer, objectName string, sourceName string, mimetype string) (*minio.UploadInfo, error) {
	object := bytes.NewReader(file.Bytes())
	metadata := map[string]string{
		"sourcename": sourceName,
	}

	info, err := s.client.PutObject(context.Background(), s.bucket, objectName, object, object.Size(), minio.PutObjectOptions{
		ContentType:  mimetype,
		UserMetadata: metadata,
	})
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (s Storage) DeleteFile(objectName string) error {
	err := s.client.RemoveObject(context.Background(), s.bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
