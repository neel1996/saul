package initializers

import (
	"core/configuration"
	"core/middleware"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(config configuration.Configuration) *gin.Engine {
	r := gin.Default()
	//r.Use(authMiddleware.Authenticate)
	middleware.CorsMiddleware(r, config)
	initializeSockets(config, r)

	group := r.Group("/api/saul/v1")
	{
		group.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "OK",
			})
		})
		group.POST("/login", loginController.Login)
		group.POST("/upload", documentUploadController.UploadDocument)
	}

	return r
}
