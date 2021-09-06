package sockets

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/undefined7887/foundation/internal/utils"
)

var (
	pingTimeout      = time.Second * 10
	pingCheckTimeout = time.Second * 4
)

type wsConnection struct {
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

func (c *wsConnection) SetEncoder(encoder Encoder) {
	c.encoder = encoder
}

func (c *wsConnection) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *wsConnection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *wsConnection) Write(topic string, data interface{}) error {
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

func (c *wsConnection) writeMessage(topic string, data interface{}) error {
	c.encoder.Writer(&c.writeCounter)

	if err := c.encoder.WriteTopic(topic); err != nil {
		return err
	}

	if err := c.encoder.WriteData(data); err != nil {
		return err
	}

	return nil
}

func (c *wsConnection) writePing() error {
	c.writerMutex.Lock()
	defer c.writerMutex.Unlock()

	return c.conn.WriteMessage(websocket.PingMessage, nil)
}

func (c *wsConnection) Close() error {
	if err := c.conn.Close(); err != nil {
		return err
	}

	c.onClose(nil)
	return nil
}

func (c *wsConnection) BytesSent() uint64 {
	return c.writeCounter.Count()
}

func (c *wsConnection) BytesReceived() uint64 {
	return c.readCounter.Count()
}

func (c *wsConnection) OnMessage(handler Handler) {
	c.onMessage[handler.Topic()] = handler
}

func (c *wsConnection) OnClose(handler func(err error)) {
	c.onClose = handler
}

func (c *wsConnection) OnError(handler func(err error)) {
	c.onError = handler
}

func (c *wsConnection) readLoop() {
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
			c.readCounter.Reset(reader)
			c.readMessage()
		}
	}
}

func (c *wsConnection) readMessage() {
	c.encoder.Reader(&c.readCounter)

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

func (c *wsConnection) pingLoop() {
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Time{})
	})

	for {
		<-c.pingTicker.C

		var (
			writeErr    = c.writePing()
			deadlineErr = c.conn.SetReadDeadline(time.Now().Add(pingCheckTimeout))
		)

		if writeErr != nil || deadlineErr != nil {
			// If the connection is closed, we will detect it inside the read loop
			return
		}
	}
}

func DialWs(addr string) (Connection, error) {
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return nil, err
	}

	return newWsConnection(conn), err
}

func newWsConnection(innerConn *websocket.Conn) Connection {
	result := &wsConnection{
		conn:       innerConn,
		pingTicker: time.NewTicker(pingTimeout),

		onMessage: map[string]Handler{},
		onClose:   func(err error) {},
		onError:   func(err error) {},
	}

	go result.readLoop()
	go result.pingLoop()

	return result
}
