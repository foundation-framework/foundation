package websocket

import (
	"context"
	"fmt"
	"log"
	"net"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/intale-llc/foundation/net/sockets"
)

var (
	pingTimeout      = time.Second * 10
	pingCheckTimeout = time.Second * 4
)

type conn struct {
	inner    *websocket.Conn
	listener *server
	encoder  sockets.Encoder

	pingTicker  *time.Ticker
	writerMutex sync.Mutex

	reader readerCounter
	writer writerCounter

	msgHandlers   map[string]sockets.Handler
	closeHandlers []func(err error)
	errorHandlers []func(err error)
	fatalHandler  func(topic string, data interface{}, msg interface{})
}

func newConn(inner *websocket.Conn, listener *server) sockets.Conn {
	result := &conn{
		inner:    inner,
		listener: listener,
		encoder:  sockets.NewMsgpackEncoder(),

		pingTicker: time.NewTicker(pingTimeout),

		msgHandlers:   map[string]sockets.Handler{},
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

func (c *conn) SetEncoder(encoder sockets.Encoder) {
	c.encoder = encoder
}

func (c *conn) LocalAddr() net.Addr {
	return c.inner.LocalAddr()
}

func (c *conn) RemoteAddr() net.Addr {
	return c.inner.RemoteAddr()
}

func (c *conn) Write(topic string, data interface{}) error {
	c.writerMutex.Lock()
	defer c.writerMutex.Unlock()

	writer, err := c.inner.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return err
	}

	c.writer.Reset(writer)

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

func (c *conn) writeMessage(topic string, data interface{}) error {
	c.encoder.ResetWriter(&c.writer)

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

func (c *conn) writePing() error {
	c.writerMutex.Lock()
	defer c.writerMutex.Unlock()

	return c.inner.WriteMessage(websocket.PingMessage, nil)
}

func (c *conn) Close(context.Context) error {
	if err := c.inner.Close(); err != nil {
		return err
	}

	return nil
}

func (c *conn) BytesSent() uint64 {
	return c.writer.Count()
}

func (c *conn) BytesReceived() uint64 {
	return c.reader.Count()
}

func (c *conn) HandleMsg(handlers ...sockets.Handler) {
	for _, handler := range handlers {
		c.msgHandlers[handler.Topic()] = handler
	}
}

func (c *conn) HandleClose(fun func(err error)) {
	c.closeHandlers = append(c.closeHandlers, fun)
}

func (c *conn) callCloseHandlers(err error) {
	for _, fun := range c.closeHandlers {
		fun(err)
	}
}

func (c *conn) HandleError(fun func(err error)) {
	c.errorHandlers = append(c.errorHandlers, fun)
}

func (c *conn) callErrorHandlers(err error) {
	for _, fun := range c.errorHandlers {
		fun(err)
	}
}

func (c *conn) HandleFatal(fun func(topic string, data interface{}, msg interface{})) {
	c.fatalHandler = fun
}

func (c *conn) readLoop() {
	for {
		messageType, reader, err := c.inner.NextReader()
		if err != nil {
			c.callCloseHandlers(err)

			// We MUST explicitly close connection
			// Without this close, a connection file descriptor is sometimes leaked
			_ = c.inner.Close()

			return
		}

		// Resetting ping timer after any data received (ping messages are included)
		c.pingTicker.Reset(pingTimeout)

		if messageType == websocket.BinaryMessage {
			c.reader.ResetReader(reader)
			c.readMessage()
		}
	}
}

func (c *conn) readMessage() {
	c.encoder.ResetReader(&c.reader)

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

func (c *conn) panicCatcher(topic string, data interface{}) {
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

func (c *conn) pingLoop() {
	c.inner.SetPongHandler(func(string) error {
		return c.inner.SetReadDeadline(time.Time{})
	})

	for {
		<-c.pingTicker.C

		// Important:
		// If the connection is closed, we will detect it inside the read loop

		if err := c.writePing(); err != nil {
			return
		}

		if err := c.inner.SetReadDeadline(time.Now().Add(pingCheckTimeout)); err != nil {
			return
		}
	}
}

func Dial(addr string) (sockets.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return nil, err
	}

	return newConn(conn, nil), err
}
