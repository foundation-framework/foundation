package sockets

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/undefined7887/foundation/units"
)

type Upgrader websocket.Upgrader

var DefaultUpgrader = &Upgrader{
	HandshakeTimeout: time.Second * 5,
	ReadBufferSize:   units.Kilobyte * 65,
	WriteBufferSize:  units.Kilobyte * 65,
	CheckOrigin:      func(r *http.Request) bool { return true },

	// Other settings are nil & false
}

type websocketListener struct {
	handler  *http.ServeMux
	upgrader *Upgrader

	onConnection func(connection Connection)
	onError      func(err error)
}

func NewWebsocketListener(handler *http.ServeMux, upgrader *Upgrader) Listener {
	return &websocketListener{
		handler:  handler,
		upgrader: upgrader,

		onConnection: func(connection Connection) {},
		onError:      func(err error) {},
	}
}

func (l *websocketListener) Listen(addr string) error {
	var (
		handler  = l.handler
		upgrader = (*websocket.Upgrader)(l.upgrader)
	)

	if l.handler == nil {
		handler = http.DefaultServeMux
	}

	if upgrader == nil {
		upgrader = (*websocket.Upgrader)(DefaultUpgrader)
	}

	handler.HandleFunc(addr, func(w http.ResponseWriter, r *http.Request) {
		inner, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			l.onError(err)
			return
		}

		l.onConnection(newWebsocketConnection(inner))
	})

	return nil
}

func (l *websocketListener) Close() error {
	// Only for interface realization
	return nil
}

func (l *websocketListener) OnConnection(h func(connection Connection)) {
	l.onConnection = h
}

func (l *websocketListener) OnError(h func(err error)) {
	l.onError = h
}

func (l *websocketListener) OnClose(h func(err error)) {
	// Only for interface realization
}
