package session

import (
	"sync"

	"github.com/Workiva/go-datastructures/set"
	"github.com/intale-llc/foundation/transport"
)

type Data interface {
	ID() string
}

type Client struct {
	transport.Connection

	pool  *ClientPool
	rooms *set.Set

	id   string
	data interface{}
	rmux sync.RWMutex
}

func (c *Client) ID() string {
	return c.id
}

func (c *Client) DataID() string {
	c.rmux.RLock()
	defer c.rmux.RUnlock()

	if c.data == nil {
		return ""
	}

	if d, ok := c.data.(Data); ok {
		return d.ID()
	}

	return ""
}

func (c *Client) Data() interface{} {
	c.rmux.RLock()
	defer c.rmux.RUnlock()

	return c.data
}

func (c *Client) SetData(data interface{}) {
	c.rmux.Lock()
	defer c.rmux.Unlock()

	c.data = data
}

func (c *Client) Access(level AccessLevel) bool {
	c.rmux.RLock()
	defer c.rmux.RUnlock()

	return c.pool.accessLevelFun(c) >= level
}

func (c *Client) Join(room string) {
	c.pool.join(c, room)
	c.rooms.Add(room)
}

func (c *Client) Leave(room string) {
	c.pool.leave(c, room)
	c.rooms.Remove(room)
}

func (c *Client) LeaveAll() {
	rooms := c.rooms.Flatten()
	c.rooms.Clear()

	for _, room := range rooms {
		c.pool.leave(c, room.(string))
	}
}

func (c *Client) Broadcast(room string, topic string, data interface{}) {
	c.pool.iterateRoom(room, func(client *Client) {
		// Ignore any errors
		_ = client.Write(topic, data)
	})
}
