package initializers

import (
	"core/configuration"
	"core/kafka/consumers"
	"core/log"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

func initializeSockets(config configuration.Configuration, r *gin.Engine) {
	logger := log.NewLogger()
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		logger.Infof("connected: %v", s.ID())
		go consumers.NewDocumentStatusConsumer(config, documentAnalyzerService).ConsumeDocumentStatus(s)

		return nil
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		logger.Errorf("socket error: %v", e)
	})

	server.OnEvent("/", "join", func(s socketio.Conn, msg string) {
		logger.Infof("Joining room %s", msg)

		s.Join(msg)
		s.Emit("joined", msg)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		logger.Warnf("Socket closed: %v", reason)
	})

	// sockets
	server.OnEvent("/", "message", messageSocketController.HandleMessage)

	go server.Serve()

	r.GET("/saul/socket.io/*any", gin.WrapH(server))
	r.POST("/saul/socket.io/*any", gin.WrapH(server))
}
