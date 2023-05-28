package controller

import "core/service"

type DocumentQaSocketController interface {
	AnswerQuestion(documentId string, question string) (string, error)
}

type documentQaSocketController struct {
	documentAnalyzerService service.DocumentAnalyzerService
}

func (controller documentQaSocketController) AnswerQuestion(documentId string, question string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func NewDocumentQaSocketController(documentAnalyzerService service.DocumentAnalyzerService) DocumentQaSocketController {
	return documentQaSocketController{
		documentAnalyzerService: documentAnalyzerService,
	}
}
