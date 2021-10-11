package session

import (
	"context"
	"sync"

	"github.com/intale-llc/foundation/net/sockets"
)

type Client struct {
	sockets.Conn

	pool  *ClientPool
	rooms map[string]struct{}

	id   string
	data map[string]interface{}
	rmux sync.RWMutex
}

func (c *Client) ID() string {
	return c.id
}

func (c *Client) Data(key string) interface{} {
	c.rmux.RLock()
	defer c.rmux.RUnlock()

	return c.data[key]
}

func (c *Client) SetData(key string, data interface{}) {
	c.rmux.Lock()
	defer c.rmux.Unlock()

	c.data[key] = data
}

func (c *Client) DropData(key string) {
	c.rmux.Lock()
	defer c.rmux.Unlock()

	delete(c.data, key)
}

func (c *Client) Check(fns ...func(c *Client) bool) bool {
	c.rmux.RLock()
	defer c.rmux.RUnlock()

	for _, fn := range fns {
		if !fn(c) {
			return false
		}
	}

	return true
}

func (c *Client) Join(room string) {
	if c.pool.join(c, room) {
		c.rooms[room] = struct{}{}
	}
}

func (c *Client) Leave(room string) {
	if c.pool.leave(c, room) {
		delete(c.rooms, room)
	}
}

func (c *Client) LeaveAll() {
	for room := range c.rooms {
		c.pool.leave(c, room)
	}

	c.rooms = map[string]struct{}{}
}

func (c *Client) Broadcast(room string, topic string, data interface{}) {
	c.pool.Broadcast(c, room, topic, data)
}

func (c *Client) Context() context.Context {
	return PackClient(context.TODO(), c)
}
