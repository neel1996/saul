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
	GetObject(ctx context.Context, bucketName string, objectName string) ([]byte, error)
}

type minioClient struct {
	minio *minio.Client
}

func (client minioClient) GetObject(ctx context.Context, bucketName string, objectName string) ([]byte, error) {
	object, err := client.minio.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		logger.Errorf("Error getting image from Minio: %v", err)
		return nil, err
	}
	defer object.Close()

	objectInfo, err := object.Stat()
	if err != nil || objectInfo.Size == 0 {
		logger.Errorf("Error getting image info from Minio: %v", err)
		return nil, err
	}

	logger.Infof("Reading %d bytes from Minio", objectInfo.Size)
	b := make([]byte, objectInfo.Size)
	for {
		_, err = object.Read(b)
		if err != nil {
			break
		}
	}

	return b, nil
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
