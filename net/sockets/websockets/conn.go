package websockets

import (
	"context"
	"log"
	"net"
	"net/http"
	"reflect"
	"runtime/debug"
	"sync"
	"time"

	"github.com/foundation-framework/foundation/errors"
	"github.com/foundation-framework/foundation/net/sockets"
	"github.com/foundation-framework/foundation/rand"
	"github.com/gorilla/websocket"
)

var (
	pingTimeout      = time.Second * 10
	pingCheckTimeout = time.Second * 4
)

type conn struct {
	inner    *websocket.Conn
	server   *server
	acceptWg sync.WaitGroup

	encoder sockets.Encoder
	pinger  *time.Ticker

	reader      readerCounter
	readerMutex sync.Mutex

	writer      writerCounter
	writerMutex sync.Mutex

	closeCb []func(err error)
	errorCb func(err error)
	fatalCb func(topic string, data interface{}, msg interface{})

	messageHandlers map[string]sockets.MessageHandler
	replyHandlers   map[string]sockets.ReplyHandler
	handlersMutex   sync.Mutex
}

func newConn(inner *websocket.Conn, server *server) sockets.Conn {
	result := &conn{
		inner:   inner,
		server:  server,
		encoder: sockets.NewMsgpackEncoder(),

		pinger: time.NewTicker(pingTimeout),

		closeCb: []func(err error){},
		errorCb: func(err error) {},
		// fatalCb must be nil to print default messages to terminal
		messageHandlers: map[string]sockets.MessageHandler{},
		replyHandlers:   map[string]sockets.ReplyHandler{},
	}

	// For accept function
	result.acceptWg.Add(1)

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

func (c *conn) Accept() {
	c.acceptWg.Done()
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
	c.encoder = encoder
}

func (c *conn) readMessageLoop() {
	c.acceptWg.Wait()

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
		c.pinger.Reset(pingTimeout)

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
		<-c.pinger.C

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

	c.encoder.ResetReader(&c.reader)

	id, err := c.encoder.ReadString()
	if err != nil {
		c.errorCb(err)
		return
	}

	topic, err := c.encoder.ReadString()
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
	if err := c.encoder.ReadData(data); err != nil {
		c.errorCb(err)
		return
	}

	go func() {
		defer c.panicCatcher(topic, data)

		if handler, ok := handler.(sockets.ReplyHandler); ok {
			handler.Serve(data)
			return
		}

		if handler, ok := handler.(sockets.MessageHandler); ok {
			replyData := handler.Serve(data)
			if replyData == nil {
				return
			}

			if err := c.write(id, topic, replyData); err != nil {
				c.errorCb(err)
			}
		}
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

	log.Printf(
		"websockets: panic on \"%s\" handler: %s\n%s",

		topic,
		msg,
		string(debug.Stack()),
	)
}

func (c *conn) Write(topic string, data interface{}) error {
	id := rand.UUID()
	if err := c.write(id, topic, data); err != nil {
		return err
	}

	return nil
}

func (c *conn) WriteWithReply(topic string, data interface{}, handler sockets.ReplyHandler) error {
	id := rand.UUID()
	if err := c.write(id, topic, data); err != nil {
		return err
	}

	c.setReplyHandler(id, handler)
	return nil
}

func (c *conn) write(id, topic string, data interface{}) error {
	c.writerMutex.Lock()
	defer c.writerMutex.Unlock()

	writer, err := c.inner.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return err
	}

	c.writer.Reset(writer)

	if err := c.writeMessage(id, topic, data); err != nil {
		// writeMessage can only return network errors that will be handled in readMessageLoop
		return err
	}

	if err := writer.Close(); err != nil {
		// Not a critical error (any critical errors we handle inside read loop)
		return err
	}

	return nil
}

func (c *conn) writeMessage(id, topic string, data interface{}) error {
	c.encoder.ResetWriter(&c.writer)

	if err := c.encoder.WriteString(id); err != nil {
		return err
	}
	if err := c.encoder.WriteString(topic); err != nil {
		return err
	}
	if err := c.encoder.WriteData(data); err != nil {
		return err
	}

	return c.encoder.Flush()
}

func (c *conn) writePing() error {
	c.writerMutex.Lock()
	defer c.writerMutex.Unlock()

	return c.inner.WriteMessage(websocket.PingMessage, nil)
}

func (c *conn) SetMessageHandlers(handlers ...sockets.MessageHandler) {
	c.handlersMutex.Lock()
	defer c.handlersMutex.Unlock()

	for _, handler := range handlers {
		c.messageHandlers[handler.Topic()] = handler
	}
}

func (c *conn) RemoveMessageHandlers(handlers ...sockets.MessageHandler) {
	c.handlersMutex.Lock()
	defer c.handlersMutex.Unlock()

	for _, handler := range handlers {
		delete(c.messageHandlers, handler.Topic())
	}
}

func (c *conn) setReplyHandler(id string, handler sockets.ReplyHandler) {
	c.handlersMutex.Lock()
	defer c.handlersMutex.Unlock()

	c.replyHandlers[id] = handler
}

func (c *conn) findHandler(id, topic string) sockets.HandlerBase {
	c.handlersMutex.Lock()
	defer c.handlersMutex.Unlock()

	replyHandler := c.replyHandlers[id]
	if replyHandler != nil {
		// Automatically remove reply replyHandler
		delete(c.replyHandlers, id)

		return replyHandler
	}

	messageHandler := c.messageHandlers[topic]
	if messageHandler != nil {
		return messageHandler
	}

	return nil
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

func Dial(addr string, headers http.Header) (sockets.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(addr, headers)
	if err != nil {
		return nil, err
	}

	return newConn(conn, nil), err
}

func isPointer(i interface{}) bool {
	return reflect.TypeOf(i).Kind() == reflect.Ptr
}
