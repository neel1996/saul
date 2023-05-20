package service

//go:generate mockgen -destination=../mocks/mock_document_upload_service.go -package=mocks -source=document_upload_service.go

import (
	"bytes"
	"context"
	"core/clients"
	"core/configuration"
	"core/constants"
	"core/kafka"
	"core/log"
	"core/model"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

type DocumentUploadService interface {
	UploadDocument(ctx context.Context, pdfBytes []byte) (string, error)
}

type documentUploadService struct {
	config      configuration.Configuration
	minioClient clients.MinioClient
	producer    kafka.DocumentDetailsProducer
}

func (service documentUploadService) UploadDocument(ctx context.Context, pdfBytes []byte) (string, error) {
	logger := log.NewLogger(ctx)
	logger.Info("Uploading document to Minio")

	checksum, err := service.calculateChecksum(pdfBytes)
	if err != nil {
		logger.Errorf("Failed to calculate checksum: %v", err)
		return "", err
	}
	objectName := fmt.Sprintf("%s/%s", checksum, constants.DefaultDocumentName)

	_, err = service.minioClient.PutObject(
		ctx,
		service.config.Minio.Bucket,
		objectName,
		pdfBytes,
	)
	if err != nil {
		logger.Errorf("Failed to upload document to Minio: %v", err)
		return "", constants.DocumentUploadError
	}

	documentDetails := model.DocumentDetails{
		FileName: constants.DefaultDocumentName,
		Checksum: checksum,
	}

	err = service.producer.ProduceDocumentDetailsEvent(ctx, documentDetails)
	if err != nil {
		logger.Errorf("Failed to produce document details event: %v", err)
		return "", constants.DocumentUploadError
	}

	return checksum, nil
}

func (service documentUploadService) calculateChecksum(pdfBytes []byte) (string, error) {
	hash := sha256.New()
	_, err := io.Copy(hash, bytes.NewReader(pdfBytes))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func NewDocumentUploadService(config configuration.Configuration, minioClient clients.MinioClient, producer kafka.DocumentDetailsProducer) DocumentUploadService {
	return documentUploadService{
		config,
		minioClient,
		producer,
	}
}
