package service

import (
	"context"
	"core/configuration"
	"core/constants"
	"core/mocks"
	"core/model"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/suite"
	"testing"
)

type DocumentUploadServiceTestSuite struct {
	suite.Suite
	context                     context.Context
	mockController              *gomock.Controller
	config                      configuration.Configuration
	mockMinioClient             *mocks.MockMinioClient
	mockDocumentDetailsProducer *mocks.MockDocumentDetailsProducer
	service                     DocumentUploadService
}

func TestDocumentUploadServiceTestSuite(t *testing.T) {
	suite.Run(t, new(DocumentUploadServiceTestSuite))
}

func (suite *DocumentUploadServiceTestSuite) SetupTest() {
	suite.context = context.Background()
	suite.mockController = gomock.NewController(suite.T())
	suite.config = configuration.Configuration{
		Minio: configuration.Minio{
			Bucket: "test-bucket",
		},
	}
	suite.mockMinioClient = mocks.NewMockMinioClient(suite.mockController)
	suite.mockDocumentDetailsProducer = mocks.NewMockDocumentDetailsProducer(suite.mockController)

	suite.service = NewDocumentUploadService(suite.config, suite.mockMinioClient, suite.mockDocumentDetailsProducer)
}

func (suite *DocumentUploadServiceTestSuite) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *DocumentUploadServiceTestSuite) TestUploadDocument_ShouldUploadDocumentToMinioAndProduceEvent() {
	pdfBytes := []byte("test")
	checksum := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	objectName := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08/document.pdf"

	documentDetails := model.DocumentDetails{
		FileName: "document.pdf",
		Checksum: checksum,
	}
	suite.mockMinioClient.EXPECT().
		PutObject(suite.context, suite.config.Minio.Bucket, objectName, pdfBytes).
		Return(minio.UploadInfo{}, nil).
		Times(1)

	suite.mockDocumentDetailsProducer.EXPECT().
		ProduceDocumentDetailsEvent(suite.context, documentDetails).
		Return(nil).
		Times(1)

	fileChecksum, err := suite.service.UploadDocument(suite.context, pdfBytes)

	suite.Nil(err)
	suite.Equal(checksum, fileChecksum)
}

func (suite *DocumentUploadServiceTestSuite) TestUploadDocument_WhenMinioUploadFails_ShouldReturnError() {
	pdfBytes := []byte("test")
	objectName := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08/document.pdf"

	suite.mockMinioClient.EXPECT().
		PutObject(suite.context, suite.config.Minio.Bucket, objectName, pdfBytes).
		Return(minio.UploadInfo{}, errors.New("bucket not found")).
		Times(1)

	_, err := suite.service.UploadDocument(suite.context, pdfBytes)

	suite.NotNil(err)
	suite.Equal(constants.DocumentUploadError, err)
}

func (suite *DocumentUploadServiceTestSuite) TestUploadDocument_WhenProducingEventFails_ShouldReturnError() {
	pdfBytes := []byte("test")
	checksum := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	objectName := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08/document.pdf"

	documentDetails := model.DocumentDetails{
		FileName: "document.pdf",
		Checksum: checksum,
	}
	suite.mockMinioClient.EXPECT().
		PutObject(suite.context, suite.config.Minio.Bucket, objectName, pdfBytes).
		Return(minio.UploadInfo{
			ChecksumSHA256: checksum,
		}, nil).
		Times(1)

	suite.mockDocumentDetailsProducer.EXPECT().
		ProduceDocumentDetailsEvent(suite.context, documentDetails).
		Return(errors.New("failed to produce event")).
		Times(1)

	_, err := suite.service.UploadDocument(suite.context, pdfBytes)

	suite.NotNil(err)
	suite.Equal(constants.DocumentUploadError, err)
}
