package initializers

import (
	"context"
	"core/configuration"
	"core/middleware"
	"core/model/request"
	"encoding/json"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"log"
)

func InitializeRoutes(config configuration.Configuration) *gin.Engine {
	r := gin.Default()
	//r.Use(authMiddleware.Authenticate)
	middleware.CorsMiddleware(r, config)
	initializeSocket(r)

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

func initializeSocket(r *gin.Engine) {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())

		return nil
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	server.OnEvent("/", "join", func(s socketio.Conn, msg string) {
		log.Println("join:", msg)
		s.Join(msg)
		s.Emit("joined", msg)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("closed", reason)
	})

	server.OnEvent("/", "message", func(s socketio.Conn, msg string) {
		log.Println("message:", msg)

		var req request.DocumentQaSocketMessage
		err := json.Unmarshal([]byte(msg), &req)
		if err != nil {
			log.Println("error unmarshalling message:", err.Error())
			return
		}

		answer, err := documentAnalyzerService.AnalyzeDocument(context.Background(), req.DocumentId, req.Question)
		if err != nil {
			log.Println("error analyzing document:", err.Error())
			return
		}

		b, err := json.Marshal(answer)
		if err != nil {
			log.Println("error marshalling answer:", err.Error())
			return
		}

		s.Emit("answer", string(b))
	})

	go server.Serve()

	r.GET("/saul/socket.io/*any", gin.WrapH(server))
	r.POST("/saul/socket.io/*any", gin.WrapH(server))
}
