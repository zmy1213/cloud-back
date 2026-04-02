package minio

import (
	"context"
	"time"

	minioV7 "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func Check(endpoint, accessKey, secretKey, bucket string, useSSL bool, timeout time.Duration) error {
	client, err := minioV7.New(endpoint, &minioV7.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err = client.BucketExists(ctx, bucket)
	return err
}
