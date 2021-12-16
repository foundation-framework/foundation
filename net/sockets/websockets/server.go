package websockets

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/intale-llc/foundation/net/sockets"
	"github.com/intale-llc/foundation/units"
)

type Upgrader websocket.Upgrader

var DefaultUpgrader = &Upgrader{
	HandshakeTimeout: time.Second * 5,
	ReadBufferSize:   units.Kilobyte * 65,
	WriteBufferSize:  units.Kilobyte * 65,
	CheckOrigin:      func(r *http.Request) bool { return true },

	// Other settings are nil & false
}

type server struct {
	upgrader *Upgrader

	connCb  func(sockets.Conn, http.Header)
	errorCb func(error)
}

// NewServer creates new Server based on WebSockets protocol
//
// Listen method accepts path to handle incoming connections
// (path is used to create handler for restListener)
func NewServer(upgrader *Upgrader) sockets.Server {
	return &server{
		upgrader: upgrader,

		connCb:  func(sockets.Conn, http.Header) {},
		errorCb: func(error) {},
	}
}

func (l *server) Handler() http.Handler {
	var upgrader = (*websocket.Upgrader)(l.upgrader)

	// Use DefaultUpgrader if nil
	if upgrader == nil {
		upgrader = (*websocket.Upgrader)(DefaultUpgrader)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			l.errorCb(err)
			return
		}

		l.connCb(newConn(conn, l), r.Header)
	})
}

func (l *server) OnConn(fn func(conn sockets.Conn, header http.Header)) {
	l.connCb = fn
}

func (l *server) OnError(fn func(err error)) {
	l.errorCb = fn
}
