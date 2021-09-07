package sockets

import (
	"fmt"
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

type websocketConnection struct {
	conn    *websocket.Conn
	encoder Encoder

	pingTicker  *time.Ticker
	writerMutex sync.Mutex

	readCounter  utils.ReadCounter
	writeCounter utils.WriteCounter

	onMessage map[string]Handler
	onClose   func(err error)
	onError   func(err error)
}

func newWebsocketConnection(innerConn *websocket.Conn) Connection {
	result := &websocketConnection{
		conn:    innerConn,
		encoder: NewMsgpackEncoder(),

		pingTicker: time.NewTicker(pingTimeout),

		onMessage: map[string]Handler{},
		onClose:   func(err error) {},
		onError:   func(err error) {},
	}

	go result.readLoop()
	go result.pingLoop()

	return result
}

func (c *websocketConnection) SetEncoder(encoder Encoder) {
	c.encoder = encoder
}

func (c *websocketConnection) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *websocketConnection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *websocketConnection) Write(topic string, data interface{}) error {
	c.writerMutex.Lock()
	defer c.writerMutex.Unlock()

	writer, err := c.conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return err
	}

	c.writeCounter.Reset(writer)

	if err := c.writeMessage(topic, data); err != nil {
		// writeMessage can only return network errors that will be handled in readLoop
		c.onError(err)
	}

	if err := writer.Close(); err != nil {
		// Not critical error (any critical errors we handle inside read loop)
		c.onError(err)
	}

	return nil
}

func (c *websocketConnection) writeMessage(topic string, data interface{}) error {
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

func (c *websocketConnection) writePing() error {
	c.writerMutex.Lock()
	defer c.writerMutex.Unlock()

	return c.conn.WriteMessage(websocket.PingMessage, nil)
}

func (c *websocketConnection) Close() error {
	if err := c.conn.Close(); err != nil {
		return err
	}

	c.onClose(nil)
	return nil
}

func (c *websocketConnection) BytesSent() uint64 {
	return c.writeCounter.Count()
}

func (c *websocketConnection) BytesReceived() uint64 {
	return c.readCounter.Count()
}

func (c *websocketConnection) OnMessage(handler Handler) {
	c.onMessage[handler.Topic()] = handler
}

func (c *websocketConnection) OnClose(handler func(err error)) {
	c.onClose = handler
}

func (c *websocketConnection) OnError(handler func(err error)) {
	c.onError = handler
}

func (c *websocketConnection) readLoop() {
	for {
		messageType, reader, err := c.conn.NextReader()
		if err != nil {
			c.onClose(err)

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

func (c *websocketConnection) readMessage() {
	c.encoder.ResetReader(&c.readCounter)

	topic, err := c.encoder.ReadTopic()
	if err != nil {
		c.onError(err)
		return
	}

	handler := c.onMessage[topic]
	if handler == nil {
		c.onError(fmt.Errorf("no handler found for \"%s\" topic", topic))
		return
	}

	data := handler.Model()
	if err := c.encoder.ReadData(data); err != nil {
		c.onError(err)
		return
	}

	go func() {
		replyData := handler.Serve(data)
		if replyData == nil {
			return
		}

		if err := c.Write(topic, replyData); err != nil {
			c.onError(err)
		}
	}()
}

func (c *websocketConnection) pingLoop() {
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

func DialWebsocket(addr string) (Connection, error) {
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return nil, err
	}

	return newWebsocketConnection(conn), err
}
