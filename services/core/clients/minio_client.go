package clients

//go:generate mockgen -destination=../mocks/mock_minio_client.go -package=mocks -source=minio_client.go

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
)

type MinioClient interface {
	PutObject(ctx context.Context, bucketName string, objectName string, fileBytes []byte) (minio.UploadInfo, error)
	ListObjects(ctx context.Context, bucketName string, opts minio.ListObjectsOptions) <-chan minio.ObjectInfo
	GetObject(ctx context.Context, bucketName string, objectName string) (*minio.Object, error)
}

type minioClient struct {
	minio *minio.Client
}

func (client minioClient) GetObject(ctx context.Context, bucketName string, objectName string) (*minio.Object, error) {
	return client.minio.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
}

func (client minioClient) ListObjects(ctx context.Context, bucketName string, opts minio.ListObjectsOptions) <-chan minio.ObjectInfo {
	return client.minio.ListObjects(ctx, bucketName, opts)
}

func (client minioClient) PutObject(ctx context.Context, bucketName string, objectName string, fileBytes []byte) (minio.UploadInfo, error) {
	return client.minio.PutObject(
		ctx,
		bucketName,
		objectName,
		bytes.NewBuffer(fileBytes),
		int64(len(fileBytes)),
		minio.PutObjectOptions{ContentType: "application/pdf"},
	)
}

func NewMinioClient(minio *minio.Client) MinioClient {
	return &minioClient{
		minio: minio,
	}
}
