package websocket

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

	connHandlers  []func(sockets.Conn)
	errorHandlers []func(error)
	closeHandlers []func(error)
}

// NewServer creates new Server based on Websocket protocol
//
// Listen method accepts path to handle incoming connections
// (path is used to create handler for restListener)
func NewServer(upgrader *Upgrader) sockets.Server {
	return &server{
		upgrader: upgrader,

		connHandlers:  []func(sockets.Conn){},
		errorHandlers: []func(error){},
		closeHandlers: []func(error){},
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
			l.callErrorHandlers(err)
			return
		}

		l.callConnHandlers(newConn(conn, l))
	})
}

func (l *server) HandleConn(fun func(conn sockets.Conn)) {
	l.connHandlers = append(l.connHandlers, fun)
}

func (l *server) callConnHandlers(conn sockets.Conn) {
	for _, fun := range l.connHandlers {
		fun(conn)
	}
}

func (l *server) HandleError(fun func(err error)) {
	l.errorHandlers = append(l.errorHandlers, fun)
}

func (l *server) callErrorHandlers(err error) {
	for _, fun := range l.errorHandlers {
		fun(err)
	}
}

func (l *server) HandleClose(fun func(err error)) {
	l.closeHandlers = append(l.closeHandlers, fun)
}

func (l *server) callCloseHandlers(err error) {
	for _, fun := range l.closeHandlers {
		fun(err)
	}
}
