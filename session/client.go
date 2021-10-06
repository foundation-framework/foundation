package session

import (
	"context"
	"sync"

	"github.com/Workiva/go-datastructures/set"
	"github.com/intale-llc/foundation/net/sockets"
)

type Client struct {
	sockets.Conn

	pool  *ClientPool
	rooms *set.Set

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

func (c *Client) Context() context.Context {
	return PackClient(context.TODO(), c)
}
