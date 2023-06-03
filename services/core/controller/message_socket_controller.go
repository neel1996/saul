package controller

import (
	"context"
	"core/log"
	"core/model/request"
	"core/service"
	"encoding/json"
	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
)

type MessageSocketController interface {
	HandleMessage(s socketio.Conn, msg string)
}

type messageSocketController struct {
	documentAnalyzerService service.DocumentAnalyzerService
}

func (controller messageSocketController) HandleMessage(s socketio.Conn, msg string) {
	logger := log.NewLogger().WithFields(logrus.Fields{"message": msg})
	logger.Info("Handling message")

	var req request.DocumentQaSocketMessage
	err := json.Unmarshal([]byte(msg), &req)
	if err != nil {
		controller.emitError(s, err)
		return
	}

	answer, err := controller.documentAnalyzerService.AnalyzeDocument(context.Background(), req.DocumentId, req.Question)
	if err != nil {
		controller.emitError(s, err)
		return
	}

	b, _ := json.Marshal(answer)

	logger.Info("Emitting answer")
	s.Emit("answer", string(b))
}

func (controller messageSocketController) emitError(s socketio.Conn, err error) {
	logger := log.NewLogger().WithFields(logrus.Fields{"error": err.Error()})
	logger.Error("Emitting error")

	s.Emit("error", err.Error())
}

func NewMessageSocketController(documentAnalyzerService service.DocumentAnalyzerService) MessageSocketController {
	return messageSocketController{documentAnalyzerService: documentAnalyzerService}
}
