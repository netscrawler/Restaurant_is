package minio

import (
	"context"
	"log/slog"

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

	stLogger := logger.With("func", "setup")

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
		client,
	}
}
