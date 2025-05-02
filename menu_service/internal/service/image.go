package service

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/netscrawler/Restaurant_is/menu_service/internal/domain"
)

// Storage — интерфейс для работы с MinIO/S3-совместимым хранилищем.
type Storage interface {
	// BucketExists проверяет существование бакета.
	BucketExists(ctx context.Context, bucketName string) (bool, error)

	// PresignedPutObject генерирует pre-signed URL для загрузки файла.
	PresignedPutObject(
		ctx context.Context,
		bucketName, objectName string,
		expiry time.Duration,
	) (*url.URL, error)
}

type Image struct {
	storage    Storage
	bucketName string
	expiry     time.Duration
}

func NewImageService(repo Storage, bucket string, expiry time.Duration) *Image {
	return &Image{
		storage:    repo,
		bucketName: bucket,
		expiry:     expiry,
	}
}

// CreateURL generates a pre-signed URL to upload an image via PUT method.
func (i *Image) CreateURL(
	ctx context.Context,
	filename string,
	contentType string,
) (string, string, error) {
	objKey := path.Join("uploads", fmt.Sprintf("%s-%s", uuid.New().String(), filename))

	exists, err := i.storage.BucketExists(ctx, i.bucketName)
	if err != nil {
		return "", "", fmt.Errorf(
			"%w failed to check bucket existence: %w",
			domain.ErrInternal,
			err,
		)
	}
	if !exists {
		return "", "", fmt.Errorf("%w bucket %q does not exist", domain.ErrInternal, i.bucketName)
	}

	reqParams := make(url.Values)
	reqParams.Set("Content-Type", contentType)

	presignedURL, err := i.storage.PresignedPutObject(ctx, i.bucketName, objKey, i.expiry)
	if err != nil {
		return "", "", fmt.Errorf(
			"%w failed to generate presigned PUT URL: %w",
			domain.ErrInternal,
			err,
		)
	}

	return presignedURL.String(), objKey, nil
}
