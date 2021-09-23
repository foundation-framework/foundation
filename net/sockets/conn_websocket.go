package sockets

import (
	"context"
	"fmt"
	"log"
	"net"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pingTimeout      = time.Second * 10
	pingCheckTimeout = time.Second * 4
)

type websocketConn struct {
	conn     *websocket.Conn
	listener *websocketServer
	encoder  Encoder

	pingTicker  *time.Ticker
	writerMutex sync.Mutex

	readCounter  readCounter
	writeCounter writeCounter

	msgHandlers   map[string]Handler
	closeHandlers []func(err error)
	errorHandlers []func(err error)
	fatalHandler  func(topic string, data interface{}, msg interface{})
}

func newWebsocketConn(conn *websocket.Conn, listener *websocketServer) Conn {
	result := &websocketConn{
		conn:     conn,
		listener: listener,
		encoder:  NewMsgpackEncoder(),

		pingTicker: time.NewTicker(pingTimeout),

		msgHandlers:   map[string]Handler{},
		closeHandlers: []func(err error){},
		errorHandlers: []func(err error){},
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
		c.callErrorHandlers(err)
	}

	if err := writer.Close(); err != nil {
		// Not critical error (any critical errors we handle inside read loop)
		c.callErrorHandlers(err)
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

func (c *websocketConn) HandleMsg(handlers ...Handler) {
	for _, handler := range handlers {
		c.msgHandlers[handler.Topic()] = handler
	}
}

func (c *websocketConn) HandleClose(fun func(err error)) {
	c.closeHandlers = append(c.closeHandlers, fun)
}

func (c *websocketConn) callCloseHandlers(err error) {
	for _, fun := range c.closeHandlers {
		fun(err)
	}
}

func (c *websocketConn) HandleError(fun func(err error)) {
	c.errorHandlers = append(c.errorHandlers, fun)
}

func (c *websocketConn) callErrorHandlers(err error) {
	for _, fun := range c.errorHandlers {
		fun(err)
	}
}

func (c *websocketConn) HandleFatal(fun func(topic string, data interface{}, msg interface{})) {
	c.fatalHandler = fun
}

func (c *websocketConn) readLoop() {
	for {
		messageType, reader, err := c.conn.NextReader()
		if err != nil {
			c.callCloseHandlers(err)

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
		c.callErrorHandlers(err)
		return
	}

	handler := c.msgHandlers[topic]
	if handler == nil {
		c.callErrorHandlers(fmt.Errorf("no handler found for \"%s\" topic", topic))
		return
	}

	data := handler.Model()
	if err := c.encoder.ReadData(data); err != nil {
		c.callErrorHandlers(err)
		return
	}

	go func() {
		defer c.panicCatcher(topic, data)
		replyTopic, replyData := handler.Serve(handler.Context(), data)

		if replyData == nil {
			return
		}

		if replyTopic == "" {
			replyTopic = topic
		}

		if err := c.Write(topic, replyData); err != nil {
			c.callErrorHandlers(err)
		}
	}()
}

func (c *websocketConn) panicCatcher(topic string, data interface{}) {
	msg := recover()
	if msg == nil {
		return
	}

	if c.fatalHandler != nil {
		c.fatalHandler(topic, data, msg)
		return
	}

	log.Printf("websocket: panic on '%s' handler: %s\n%s", topic, msg, string(debug.Stack()))
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
