package controller

import (
	"bytes"
	"core/constants"
	"core/mocks"
	"core/model/response"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"mime/multipart"
	"net/http/httptest"
	"testing"
)

type DocumentUploadControllerTestSuite struct {
	suite.Suite
	mockContext               *gin.Context
	recorder                  *httptest.ResponseRecorder
	mockController            *gomock.Controller
	mockDocumentUploadService *mocks.MockDocumentUploadService
	controller                DocumentUploadController
}

func TestDocumentUploadControllerTestSuite(t *testing.T) {
	suite.Run(t, new(DocumentUploadControllerTestSuite))
}

func (suite *DocumentUploadControllerTestSuite) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.recorder = httptest.NewRecorder()
	suite.mockContext, _ = gin.CreateTestContext(suite.recorder)
	suite.mockDocumentUploadService = mocks.NewMockDocumentUploadService(suite.mockController)

	suite.controller = NewDocumentUploadController(suite.mockDocumentUploadService)
}

func (suite *DocumentUploadControllerTestSuite) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *DocumentUploadControllerTestSuite) TestUploadDocument_ShouldRespondWithDocumentChecksumAndStatus200() {
	fileContent := []byte("test")
	fileName := "test.pdf"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	file, err := writer.CreateFormFile("file", fileName)
	suite.NoError(err)

	_, err = file.Write(fileContent)
	suite.NoError(err)
	suite.NoError(writer.Close())

	expectedChecksum := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"

	suite.mockContext.Request = httptest.NewRequest("POST", "/document", body)
	suite.mockContext.Request.Header.Set("Content-Type", writer.FormDataContentType())
	suite.mockDocumentUploadService.EXPECT().UploadDocument(suite.mockContext, fileContent).Return(expectedChecksum, nil)

	suite.controller.UploadDocument(suite.mockContext)

	var res response.DocumentUploadResponse
	responseBody := suite.recorder.Body.Bytes()
	err = json.Unmarshal(responseBody, &res)
	suite.NoError(err)

	suite.Equal(200, suite.recorder.Code)
	suite.Equal(expectedChecksum, res.Checksum)
}

func (suite *DocumentUploadControllerTestSuite) TestUploadDocument_WhenDocumentUploadServiceReturnsError_ShouldRespondWithStatus500() {
	fileContent := []byte("test")
	fileName := "test.pdf"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	file, err := writer.CreateFormFile("file", fileName)
	suite.NoError(err)

	_, err = file.Write(fileContent)
	suite.NoError(err)
	suite.NoError(writer.Close())

	suite.mockContext.Request = httptest.NewRequest("POST", "/document", body)
	suite.mockContext.Request.Header.Set("Content-Type", writer.FormDataContentType())
	suite.mockDocumentUploadService.EXPECT().
		UploadDocument(suite.mockContext, fileContent).
		Return("", constants.DocumentUploadError).
		Times(1)

	suite.controller.UploadDocument(suite.mockContext)

	suite.Equal(500, suite.recorder.Code)
}

func (suite *DocumentUploadControllerTestSuite) TestUploadDocument_WhenFileIsNotProvided_ShouldRespondWithStatus400() {
	suite.mockContext.Request = httptest.NewRequest("POST", "/document", nil)

	suite.controller.UploadDocument(suite.mockContext)

	suite.Equal(400, suite.recorder.Code)
}
