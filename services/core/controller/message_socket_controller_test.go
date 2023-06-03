package controller

import (
	"context"
	"core/mocks"
	"core/model/response"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type MessageSocketControllerTestSuite struct {
	suite.Suite
	mockController          *gomock.Controller
	mockSocket              *mocks.MockConn
	mockAnalyzerService     *mocks.MockDocumentAnalyzerService
	messageSocketController MessageSocketController
}

func TestMessageSocketControllerTestSuite(t *testing.T) {
	suite.Run(t, new(MessageSocketControllerTestSuite))
}

func (suite *MessageSocketControllerTestSuite) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.mockSocket = mocks.NewMockConn(suite.mockController)
	suite.mockAnalyzerService = mocks.NewMockDocumentAnalyzerService(suite.mockController)

	suite.messageSocketController = NewMessageSocketController(suite.mockAnalyzerService)
}

func (suite *MessageSocketControllerTestSuite) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *MessageSocketControllerTestSuite) TestHandleMessage_ShouldCallDocumentAnalyzerServiceAndEmitAnswer() {
	suite.mockAnalyzerService.EXPECT().
		AnalyzeDocument(context.Background(), "checksum", "question").
		Return(response.LayoutLMAnswer{
			Score:  0.9,
			Answer: "answer",
		}, nil).
		Times(1)

	suite.mockSocket.EXPECT().Emit("answer", `{"score":0.9,"answer":"answer"}`).Times(1)

	suite.messageSocketController.HandleMessage(suite.mockSocket, `{"documentId": "checksum", "question": "question"}`)
}

func (suite *MessageSocketControllerTestSuite) TestHandleMessage_WhenDocumentAnalyzerServiceFails_ShouldEmitError() {
	suite.mockAnalyzerService.EXPECT().
		AnalyzeDocument(context.Background(), "checksum", "question").
		Return(response.LayoutLMAnswer{}, errors.New("error")).
		Times(1)

	suite.mockSocket.EXPECT().Emit("error", "error").Times(1)

	suite.messageSocketController.HandleMessage(suite.mockSocket, `{"documentId": "checksum", "question": "question"}`)
}

func (suite *MessageSocketControllerTestSuite) TestHandleMessage_WhenMessageIsInvalid_ShouldEmitError() {
	suite.mockSocket.EXPECT().Emit("error", "invalid character 'a' looking for beginning of value").Times(1)

	suite.messageSocketController.HandleMessage(suite.mockSocket, `a`)
}
