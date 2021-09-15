package sockets

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/intale-llc/foundation/internal/utils"
)

var (
	pingTimeout      = time.Second * 10
	pingCheckTimeout = time.Second * 4
)

type websocketConn struct {
	conn     *websocket.Conn
	listener *websocketListener
	encoder  Encoder

	pingTicker  *time.Ticker
	writerMutex sync.Mutex

	readCounter  utils.ReadCounter
	writeCounter utils.WriteCounter

	onMessage map[string]Handler
	onClose   []func(err error)
	onError   []func(err error)
	onFatal   func(topic string, data interface{}, msg interface{})
}

func newWebsocketConn(conn *websocket.Conn, listener *websocketListener) Conn {
	result := &websocketConn{
		conn:     conn,
		listener: listener,
		encoder:  NewMsgpackEncoder(),

		pingTicker: time.NewTicker(pingTimeout),

		onMessage: map[string]Handler{},
		onClose:   []func(err error){},
		onError:   []func(err error){},
		// onFatal must be nil to print default messages to terminal
	}

	go func() {
		// Read loop reads and decodes any data in connection
		result.readLoop()
	}()

	go func() {
		// Ping loop keeps the connection alive
		result.pingLoop()
	}()

	return result
}

func (c *websocketConn) SetEncoder(encoder Encoder) {
	c.encoder = encoder
}

func (c *websocketConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *websocketConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *websocketConn) Write(topic string, data interface{}) error {
	c.writerMutex.Lock()
	defer c.writerMutex.Unlock()

	writer, err := c.conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return err
	}

	c.writeCounter.Reset(writer)

	if err := c.writeMessage(topic, data); err != nil {
		// writeMessage can only return network errors that will be handled in readLoop
		c.callErrorCb(err)
	}

	if err := writer.Close(); err != nil {
		// Not critical error (any critical errors we handle inside read loop)
		c.callErrorCb(err)
	}

	return nil
}

func (c *websocketConn) writeMessage(topic string, data interface{}) error {
	c.encoder.ResetWriter(&c.writeCounter)

	if err := c.encoder.WriteTopic(topic); err != nil {
		return err
	}

	if err := c.encoder.WriteData(data); err != nil {
		return err
	}

	if err := c.encoder.Flush(); err != nil {
		return err
	}

	return nil
}

func (c *websocketConn) writePing() error {
	c.writerMutex.Lock()
	defer c.writerMutex.Unlock()

	return c.conn.WriteMessage(websocket.PingMessage, nil)
}

func (c *websocketConn) Close(context.Context) error {
	if err := c.conn.Close(); err != nil {
		return err
	}

	return nil
}

func (c *websocketConn) BytesSent() uint64 {
	return c.writeCounter.Count()
}

func (c *websocketConn) BytesReceived() uint64 {
	return c.readCounter.Count()
}

func (c *websocketConn) OnMessage(handler Handler) {
	c.onMessage[handler.Topic()] = handler
}

func (c *websocketConn) OnClose(fun func(err error)) {
	c.onClose = append(c.onClose, fun)
}

func (c *websocketConn) callCloseCb(err error) {
	for _, fun := range c.onClose {
		fun(err)
	}
}

func (c *websocketConn) OnError(fun func(err error)) {
	c.onError = append(c.onError, fun)
}

func (c *websocketConn) callErrorCb(err error) {
	for _, fun := range c.onError {
		fun(err)
	}
}

func (c *websocketConn) OnFatal(fun func(topic string, data interface{}, msg interface{})) {
	c.onFatal = fun
}

func (c *websocketConn) readLoop() {
	for {
		messageType, reader, err := c.conn.NextReader()
		if err != nil {
			c.callCloseCb(err)

			// We MUST explicitly close connection
			// Without this close, a connection file descriptor is sometimes leaked
			_ = c.conn.Close()

			return
		}

		// Resetting ping timer after any data received (ping messages are included)
		c.pingTicker.Reset(pingTimeout)

		if messageType == websocket.BinaryMessage {
			c.readCounter.ResetReader(reader)
			c.readMessage()
		}
	}
}

func (c *websocketConn) readMessage() {
	c.encoder.ResetReader(&c.readCounter)

	topic, err := c.encoder.ReadTopic()
	if err != nil {
		c.callErrorCb(err)
		return
	}

	handler := c.onMessage[topic]
	if handler == nil {
		c.callErrorCb(fmt.Errorf("no handler found for \"%s\" topic", topic))
		return
	}

	data := handler.Model()
	if err := c.encoder.ReadData(data); err != nil {
		c.callErrorCb(err)
		return
	}

	go func() {
		defer c.panicCatcher(topic, data)
		replyData := handler.Serve(context.TODO(), data)

		if replyData == nil {
			return
		}

		if err := c.Write(topic, replyData); err != nil {
			c.callErrorCb(err)
		}
	}()
}

func (c *websocketConn) panicCatcher(topic string, data interface{}) {
	msg := recover()
	if msg == nil {
		return
	}

	if c.onFatal != nil {
		c.onFatal(topic, data, msg)
		return
	}

	log.Printf("sockets: panic on '%s' handler: %s\n%s", topic, msg, utils.StackString())
}

func (c *websocketConn) pingLoop() {
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Time{})
	})

	for {
		<-c.pingTicker.C

		// Important:
		// If the connection is closed, we will detect it inside the read loop

		if err := c.writePing(); err != nil {
			return
		}

		if err := c.conn.SetReadDeadline(time.Now().Add(pingCheckTimeout)); err != nil {
			return
		}
	}
}

func DialWebsocket(addr string) (Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return nil, err
	}

	return newWebsocketConn(conn, nil), err
}
