package sockets

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/undefined7887/foundation/units"
)

type Upgrader websocket.Upgrader

var DefaultUpgrader = Upgrader{
	HandshakeTimeout: time.Second * 5,
	ReadBufferSize:   units.Kilobyte * 65,
	WriteBufferSize:  units.Kilobyte * 65,
	CheckOrigin:      func(r *http.Request) bool { return true },

	// Other settings are nil & false
}

type listener struct {
	handler  *http.ServeMux
	upgrader *Upgrader

	handleConnection func(connection Connection)
	handleError      func(err error)
}

func (l *listener) Listen(addr string) error {
	var (
		handler  = l.handler
		upgrader = (*websocket.Upgrader)(l.upgrader)
	)

	if l.handler == nil {
		handler = http.DefaultServeMux
	}

	if upgrader == nil {
		upgrader = (*websocket.Upgrader)(&DefaultUpgrader)
	}

	handler.HandleFunc(addr, func(w http.ResponseWriter, r *http.Request) {
		inner, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			l.handleError(err)
			return
		}

		l.handleConnection(newWsConnection(inner))
	})

	return nil
}

func (l *listener) Close() error {
	// Only for interface realization
	return nil
}

func (l *listener) OnConnection(h func(connection Connection)) {
	l.handleConnection = h
}

func (l *listener) OnError(h func(err error)) {
	l.handleError = h
}

func (l *listener) OnClose(h func(err error)) {
	// Only for interface realization
}

func NewWsListener(handler *http.ServeMux, upgrader *Upgrader) Listener {
	return &listener{
		handler:  handler,
		upgrader: upgrader,

		handleConnection: func(connection Connection) {},
		handleError:      func(err error) {},
	}
}
