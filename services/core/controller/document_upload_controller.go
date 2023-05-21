package controller

import (
	"core/constants"
	"core/log"
	"core/model/response"
	"core/service"
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

type DocumentUploadController interface {
	UploadDocument(ctx *gin.Context)
}

type documentUploadController struct {
	documentUploadService service.DocumentUploadService
}

func (controller documentUploadController) UploadDocument(ctx *gin.Context) {
	logger := log.NewLogger(ctx)
	logger.Info("Handling uploaded document")

	file, err := ctx.FormFile("file")
	if err != nil {
		logger.Errorf("Failed to get file: %v", err)
		ctx.JSON(400, constants.RequestValidationError)
		return
	}

	fileBytes, err := controller.fileContentAsBytes(ctx, file)
	if err != nil {
		ctx.JSON(400, constants.RequestValidationError)
		return
	}

	checksum, err := controller.documentUploadService.UploadDocument(ctx, fileBytes)
	if err != nil {
		logger.Errorf("Failed to upload document: %v", err)
		e, ok := err.(constants.Error)
		if !ok {
			ctx.JSON(500, constants.InternalServerError)
			return
		}
		ctx.JSON(e.GetGinResponse())
		return
	}

	res := response.DocumentUploadResponse{
		Checksum: checksum,
	}

	ctx.JSON(200, res)
}

func (controller documentUploadController) fileContentAsBytes(ctx *gin.Context, file *multipart.FileHeader) ([]byte, error) {
	logger := log.NewLogger(ctx)

	fileContent, err := file.Open()
	if err != nil {
		logger.Errorf("Failed to open file: %v", err)
		return nil, err
	}

	var fileBytes = make([]byte, file.Size)
	_, err = fileContent.Read(fileBytes)
	if err != nil {
		logger.Errorf("Failed to read file: %v", err)
		return nil, err
	}

	logger.Infof("Returning file content with size %d as bytes", len(fileBytes))
	return fileBytes, nil
}

func NewDocumentUploadController(documentUploadService service.DocumentUploadService) DocumentUploadController {
	return documentUploadController{
		documentUploadService,
	}
}
