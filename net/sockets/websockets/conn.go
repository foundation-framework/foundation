package websockets

import (
	"context"
	"log"
	"net"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/intale-llc/foundation/errors"
	"github.com/intale-llc/foundation/net/sockets"
	"github.com/intale-llc/foundation/rand"
)

var (
	pingTimeout      = time.Second * 10
	pingCheckTimeout = time.Second * 4
)

type conn struct {
	inner  *websocket.Conn
	server *server

	pingTicker *time.Ticker

	encoder      sockets.Encoder
	encoderMutex sync.Mutex

	reader      readerCounter
	readerMutex sync.Mutex

	writer      writerCounter
	writerMutex sync.Mutex

	closeCb []func(err error)
	errorCb func(err error)
	fatalCb func(topic string, data interface{}, msg interface{})

	messageHandlers map[string]sockets.Handler
	replyHandlers   map[string]sockets.Handler
	handlersMutex   sync.Mutex
}

func newConn(inner *websocket.Conn, server *server) sockets.Conn {
	result := &conn{
		inner:   inner,
		server:  server,
		encoder: sockets.NewMsgpackEncoder(),

		pingTicker: time.NewTicker(pingTimeout),

		closeCb: []func(err error){},
		errorCb: func(err error) {},
		// fatalCb must be nil to print default messages to terminal
		messageHandlers: map[string]sockets.Handler{},
		replyHandlers:   map[string]sockets.Handler{},
	}

	go func() {
		// Read loop reads and decodes any data in connection
		result.readMessageLoop()
	}()

	go func() {
		// Ping loop keeps the connection alive
		result.readPingLoop()
	}()

	return result
}

func (c *conn) LocalAddr() net.Addr {
	return c.inner.LocalAddr()
}

func (c *conn) RemoteAddr() net.Addr {
	return c.inner.RemoteAddr()
}

func (c *conn) BytesSent() uint64 {
	c.writerMutex.Lock()
	defer c.writerMutex.Unlock()

	return c.writer.Count()
}

func (c *conn) BytesReceived() uint64 {
	c.readerMutex.Lock()
	defer c.readerMutex.Unlock()

	return c.reader.Count()
}

func (c *conn) SetEncoder(encoder sockets.Encoder) {
	c.encoderMutex.Lock()
	defer c.encoderMutex.Unlock()

	c.encoder = encoder
}

func (c *conn) getEncoder() sockets.Encoder {
	c.encoderMutex.Lock()
	defer c.encoderMutex.Unlock()

	return c.encoder
}

func (c *conn) readMessageLoop() {
	for {
		messageType, reader, err := c.inner.NextReader()
		if err != nil {
			c.callCloseCb(err)

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

func (c *conn) readPingLoop() {
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

func (c *conn) readMessage() {
	c.readerMutex.Lock()
	defer c.readerMutex.Unlock()

	encoder := c.getEncoder()
	encoder.ResetReader(&c.reader)

	id, err := encoder.ReadString()
	if err != nil {
		c.errorCb(err)
		return
	}

	topic, err := encoder.ReadString()
	if err != nil {
		c.errorCb(err)
		return
	}

	handler := c.findHandler(id, topic)
	if handler == nil {
		c.errorCb(errors.Newf("no handler found for \"%s\" topic", topic))
		return
	}

	data := handler.Model()
	if err := encoder.ReadData(data); err != nil {
		c.errorCb(err)
		return
	}

	go func() {
		defer c.panicCatcher(topic, data)
		handler.Context(func(ctx context.Context) {
			replyData := handler.Serve(ctx, data)

			if isReplyHandler(handler) || replyData == nil {
				return
			}

			if err := c.writeMessage(id, topic, replyData); err != nil {
				c.errorCb(err)
			}
		})
	}()
}

func (c *conn) panicCatcher(topic string, data interface{}) {
	msg := recover()
	if msg == nil {
		return
	}

	if c.fatalCb != nil {
		c.fatalCb(topic, data, msg)
		return
	}

	log.Printf("websockets: panic on \"%s\" handler: %s\n%s", topic, msg, string(debug.Stack()))
}

func (c *conn) Write(topic string, data interface{}, handler ...sockets.Handler) error {
	if len(handler) > 1 {
		return errors.New("more than one reply handler passed")
	}

	var replyHandler sockets.Handler
	if len(handler) > 0 && handler[0] != nil {
		replyHandler = handler[0]
	}

	if replyHandler != nil && !isReplyHandler(replyHandler) {
		return errors.New("reply handler must have empty topic")
	}

	id := rand.UUID()
	if err := c.writeMessage(id, topic, data); err != nil {
		return err
	}

	if replyHandler != nil {
		c.setReplyHandler(id, replyHandler)
	}

	return nil
}

func (c *conn) writeMessage(id, topic string, data interface{}) error {
	c.writerMutex.Lock()
	defer c.writerMutex.Unlock()

	writer, err := c.inner.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return err
	}

	c.writer.Reset(writer)

	if err := c.encodeMessage(id, topic, data); err != nil {
		// encodeMessage can only return network errors that will be handled in readMessageLoop
		return err
	}

	if err := writer.Close(); err != nil {
		// Not a critical error (any critical errors we handle inside read loop)
		return err
	}

	return nil
}

func (c *conn) encodeMessage(id, topic string, data interface{}) error {
	encoder := c.getEncoder()
	encoder.ResetWriter(&c.writer)

	if err := encoder.WriteString(id); err != nil {
		return err
	}
	if err := encoder.WriteString(topic); err != nil {
		return err
	}
	if err := encoder.WriteData(data); err != nil {
		return err
	}

	return encoder.Flush()
}

func (c *conn) writePing() error {
	c.writerMutex.Lock()
	defer c.writerMutex.Unlock()

	return c.inner.WriteMessage(websocket.PingMessage, nil)
}

func (c *conn) SetMessageHandlers(handlers ...sockets.Handler) {
	c.handlersMutex.Lock()
	defer c.handlersMutex.Unlock()

	for _, handler := range handlers {
		if !isPointer(handler.Model()) {
			errors.Panicf("\"%s\" handler: model must be a pointer", handler.Topic())
		}

		if isReplyHandler(handler) {
			errors.Panicf("reply handler cannot be used as message handler")
		}

		c.messageHandlers[handler.Topic()] = handler
	}
}

func (c *conn) RemoveMessageHandlers(handlers ...sockets.Handler) {
	c.handlersMutex.Lock()
	defer c.handlersMutex.Unlock()

	for _, handler := range handlers {
		delete(c.messageHandlers, handler.Topic())
	}
}

func (c *conn) setReplyHandler(id string, handler sockets.Handler) {
	c.handlersMutex.Lock()
	defer c.handlersMutex.Unlock()

	c.replyHandlers[id] = handler
}

func (c *conn) findHandler(id, topic string) sockets.Handler {
	c.handlersMutex.Lock()
	defer c.handlersMutex.Unlock()

	handler := c.replyHandlers[id]

	if handler == nil {
		handler = c.messageHandlers[topic]
	} else {
		// Automatically remove reply handler
		delete(c.replyHandlers, id)
	}

	return handler
}

func (c *conn) OnError(fn func(err error)) {
	c.errorCb = fn
}

func (c *conn) OnFatal(fn func(topic string, data interface{}, panicMsg interface{})) {
	c.fatalCb = fn
}

func (c *conn) OnClose(fn func(err error)) {
	c.closeCb = append(c.closeCb, fn)
}

func (c *conn) callCloseCb(err error) {
	for _, fn := range c.closeCb {
		fn(err)
	}
}

func (c *conn) Close(context.Context) error {
	if err := c.inner.Close(); err != nil {
		return err
	}

	return nil
}

func Dial(addr string) (sockets.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return nil, err
	}

	return newConn(conn, nil), err
}

func isReplyHandler(handler sockets.Handler) bool {
	return handler.Topic() == ""
}
