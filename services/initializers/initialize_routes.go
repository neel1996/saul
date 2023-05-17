package initializers

import (
	"github.com/gin-gonic/gin"
	"github.com/neel1996/saul/configuration"
	"github.com/neel1996/saul/log"
	"github.com/neel1996/saul/middleware"
)

func InitializeRoutes(config configuration.Configuration) *gin.Engine {
	r := gin.Default()
	r.Use(authMiddleware.Authenticate)
	setupCORS(config, r)

	group := r.Group("/api/saul/v1")
	{
		group.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "OK",
			})
		})
		group.POST("/login", loginController.Login)
	}

	return r
}

func setupCORS(config configuration.Configuration, r *gin.Engine) gin.IRoutes {
	return r.Use(func(context *gin.Context) {
		logger := log.NewLogger(context)
		for _, ignoreUrl := range config.CorsIgnoreUrls {
			if context.Request.URL.Path == ignoreUrl {
				logger.Infof("%s is not part of secure list. Ignoring CORS", ignoreUrl)
				context.Next()
				return
			}
		}

		logger.Infof("Setting up CORS for %s", context.Request.URL.Path)
		middleware.CorsMiddleware(r, config)
		context.Next()
		return
	})
}
