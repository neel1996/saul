package mocks

//go:generate mockgen -destination=./mock_socket_stub.go -package=mocks -source=socket_stub.go

import (
	socketio "github.com/googollee/go-socket.io"
	"io"
	"net"
	"net/http"
	"net/url"
)

type Conn interface {
	io.Closer
	socketio.Namespace

	ID() string
	URL() url.URL
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	RemoteHeader() http.Header
}
