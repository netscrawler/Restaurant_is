package minio

import (
	"context"
	"log/slog"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage struct {
	*minio.Client
}

// MustSetup подключается к MinIO, валидирует соединение и возвращает инстанс Storage.
// Паника при любой ошибке.
func MustSetup(
	ctx context.Context,
	endpoint, accessKey, secretKey string,
	useSSL bool,
	buckets []string,
	log *slog.Logger,
) *Storage {
	logger := log.With("storage", "mini_io")

	stLogger := logger

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		stLogger.ErrorContext(ctx, "error creating client", slog.String("error", err.Error()))
		panic(err)
	}

	// Проверка доступности
	existBuckets, err := client.ListBuckets(ctx)
	if err != nil {
		stLogger.ErrorContext(ctx, "error connect client", slog.String("error", err.Error()))
		panic(err)
	}

	stLogger.InfoContext(ctx, "Successfully connected to MinIO")

	stLogger.InfoContext(ctx, "Check exist buckets")

	existBucketsMap := make(map[string]struct{})

	for _, b := range existBuckets {
		existBucketsMap[b.Name] = struct{}{}
	}

	for _, b := range buckets {
		_, ok := existBucketsMap[b]
		if !ok {
			stLogger.InfoContext(ctx, "bucket not found creating new", slog.String("bucket", b))
			client.MakeBucket(ctx, b, minio.MakeBucketOptions{
				// TODO: add region creation from config
				Region:        "",
				ObjectLocking: true,
			})
		}

	}

	stLogger.InfoContext(ctx, "all buckets ready")

	return &Storage{
		Client: client,
	}
}

// PresignedGetObject генерирует pre-signed URL для скачивания файла.
func (s *Storage) PresignedGetObject(
	ctx context.Context,
	bucketName, objectName string,
	expiry time.Duration,
) (*url.URL, error) {
	url, err := s.Client.PresignedGetObject(ctx, bucketName, objectName, expiry, nil)
	return url, err
}

// PresignedPutObject генерирует pre-signed URL для загрузки файла.
func (s *Storage) PresignedPutObject(
	ctx context.Context,
	bucketName, objectName string,
	expiry time.Duration,
) (*url.URL, error) {
	return s.Client.PresignedPutObject(ctx, bucketName, objectName, expiry)
}

// BucketExists проверяет существование бакета.
func (s *Storage) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	return s.Client.BucketExists(ctx, bucketName)
}
