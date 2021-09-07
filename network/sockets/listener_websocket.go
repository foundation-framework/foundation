package sockets

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/intale-llc/foundation/network/rest"
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

type websocketListener struct {
	restListener rest.Listener
	upgrader     *Upgrader

	onConnection []func(Connection)
	onError      []func(error)
	onClose      []func(error)
	closeChan    chan error
}

func NewWebsocketListener(restListener rest.Listener, upgrader *Upgrader) Listener {
	result := &websocketListener{
		restListener: restListener,
		upgrader:     upgrader,

		onConnection: []func(Connection){},
		onError:      []func(error){},
		onClose:      []func(error){},
		closeChan:    make(chan error),
	}

	go result.handleRestListenerClose()
	return result
}

func (l *websocketListener) Listen(path string) error {
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

func (l *websocketListener) Close(ctx context.Context) error {
	return l.restListener.Close(ctx)
}

func (l *websocketListener) OnConnection(fun func(connection Connection)) {
	l.onConnection = append(l.onConnection, fun)
}

func (l *websocketListener) callConnectionCb(connection Connection) {
	for _, fun := range l.onConnection {
		fun(connection)
	}
}

func (l *websocketListener) OnError(fun func(err error)) {
	l.onError = append(l.onError, fun)
}

func (l *websocketListener) callErrorCb(err error) {
	for _, fun := range l.onError {
		fun(err)
	}
}

func (l *websocketListener) OnClose(fun func(err error)) {
	l.onClose = append(l.onClose, fun)
}

func (l *websocketListener) callCloseCb(err error) {
	for _, fun := range l.onClose {
		fun(err)
	}
}

func (l *websocketListener) CloseChan() <-chan error {
	return l.closeChan
}

func (l *websocketListener) handleRestListenerClose() {
	err := <-l.restListener.CloseChan()

	l.closeChan <- err
	l.callCloseCb(err)
}
