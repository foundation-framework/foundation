package sockets

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
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

type websocketServer struct {
	upgrader *Upgrader

	connHandlers  []func(Conn)
	errorHandlers []func(error)
	closeHandlers []func(error)
}

// NewWebsocketServer creates new Server based on Websocket protocol
//
// Listen method accepts path to handle incoming connections
// (path is used to create handler for restListener)
func NewWebsocketServer(upgrader *Upgrader) Server {
	return &websocketServer{
		upgrader: upgrader,

		connHandlers:  []func(Conn){},
		errorHandlers: []func(error){},
		closeHandlers: []func(error){},
	}
}

func (l *websocketServer) Handler() http.Handler {
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

		l.callConnHandlers(newWebsocketConn(conn, l))
	})
}

func (l *websocketServer) HandleConn(fun func(conn Conn)) {
	l.connHandlers = append(l.connHandlers, fun)
}

func (l *websocketServer) callConnHandlers(conn Conn) {
	for _, fun := range l.connHandlers {
		fun(conn)
	}
}

func (l *websocketServer) HandleError(fun func(err error)) {
	l.errorHandlers = append(l.errorHandlers, fun)
}

func (l *websocketServer) callErrorHandlers(err error) {
	for _, fun := range l.errorHandlers {
		fun(err)
	}
}

func (l *websocketServer) HandleClose(fun func(err error)) {
	l.closeHandlers = append(l.closeHandlers, fun)
}

func (l *websocketServer) callCloseHandlers(err error) {
	for _, fun := range l.closeHandlers {
		fun(err)
	}
}
