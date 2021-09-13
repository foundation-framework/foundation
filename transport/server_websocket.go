package transport

import (
	"context"
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
	restListener RESTServer
	upgrader     *Upgrader

	onConnection []func(Connection)
	onError      []func(error)
	onClose      []func(error)
}

// NewWebsocketServer creates new Server based on Websocket protocol
//
// Listen method accepts path to handle incoming connections
// (path is used to create handler for restListener)
func NewWebsocketServer(restListener RESTServer, upgrader *Upgrader) Server {
	return &websocketServer{
		restListener: restListener,
		upgrader:     upgrader,

		onConnection: []func(Connection){},
		onError:      []func(error){},
		onClose:      []func(error){},
	}
}

func (l *websocketServer) Listen(path string) error {
	var upgrader = (*websocket.Upgrader)(l.upgrader)

	// Use DefaultUpgrader if nil
	if upgrader == nil {
		upgrader = (*websocket.Upgrader)(DefaultUpgrader)
	}

	l.restListener.Router().HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			l.callErrorCb(err)
			return
		}

		l.callConnectionCb(newWebsocketConnection(conn, l))
	})

	return nil
}

func (l *websocketServer) Shutdown(ctx context.Context) error {
	return l.restListener.Shutdown(ctx)
}

func (l *websocketServer) OnConnection(fun func(connection Connection)) {
	l.onConnection = append(l.onConnection, fun)
}

func (l *websocketServer) callConnectionCb(connection Connection) {
	for _, fun := range l.onConnection {
		fun(connection)
	}
}

func (l *websocketServer) OnError(fun func(err error)) {
	l.onError = append(l.onError, fun)
}

func (l *websocketServer) callErrorCb(err error) {
	for _, fun := range l.onError {
		fun(err)
	}
}

func (l *websocketServer) OnClose(fun func(err error)) {
	l.onClose = append(l.onClose, fun)
}

func (l *websocketServer) callCloseCb(err error) {
	for _, fun := range l.onClose {
		fun(err)
	}
}
