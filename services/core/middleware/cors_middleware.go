package middleware

import (
	"core/configuration"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware(r *gin.Engine, config configuration.Configuration) {
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  false,
		AllowOrigins:     config.Cors.AllowedOrigins,
		AllowMethods:     config.Cors.AllowedMethods,
		AllowHeaders:     config.Cors.AllowedHeaders,
		AllowCredentials: true,
		AllowWildcard:    false,
	}))
}
