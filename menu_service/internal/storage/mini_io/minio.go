package minio

import (
	"context"
	"log/slog"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage struct {
	log   *slog.Logger
	Minio *minio.Client
}

// MustSetup подключается к MinIO, валидирует соединение и возвращает инстанс Storage.
// Паника при любой ошибке.
func MustSetup(
	ctx context.Context,
	endpoint, accessKey, secretKey string,
	useSSL bool,
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
	_, err = client.ListBuckets(ctx)
	if err != nil {
		stLogger.ErrorContext(ctx, "error connect client", slog.String("error", err.Error()))
		panic(err)
	}

	stLogger.InfoContext(ctx, "Successfully connected to MinIO")

	return &Storage{
		log:   logger,
		Minio: client,
	}
}
