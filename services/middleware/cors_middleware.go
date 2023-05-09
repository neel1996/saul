package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/neel1996/saul/configuration"
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
