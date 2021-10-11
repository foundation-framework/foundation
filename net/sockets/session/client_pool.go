package session

import (
	"sync"

	"github.com/intale-llc/foundation/net/sockets"
	"github.com/intale-llc/foundation/rand"
)

type ClientPool struct {
	adapters []Adapter

	rooms map[string]map[*Client]struct{}
	mux   sync.RWMutex
}

func NewClientPool(adapters ...Adapter) *ClientPool {
	result := &ClientPool{
		adapters: adapters,
		rooms:    map[string]map[*Client]struct{}{},
	}

	result.setupAdapters(adapters)
	return result
}

func (p *ClientPool) NewClient(connection sockets.Conn) *Client {
	return &Client{
		Conn: connection,
		pool: p,
		id:   rand.UUID(),
		data: map[string]interface{}{},
	}
}

func (p *ClientPool) Broadcast(client *Client, room, topic string, data interface{}) {
	// Firstly broadcast to adapters
	p.broadcastAdapters(client.id, room, topic, data)

	// Then broadcast to current pool
	p.broadcast(client.id, room, topic, data)
}

func (p *ClientPool) join(client *Client, room string) bool {
	p.mux.Lock()
	defer p.mux.Unlock()

	if p.rooms[room] == nil {
		p.rooms[room] = map[*Client]struct{}{}
	}

	if _, ok := p.rooms[room][client]; ok {
		return false
	}

	p.rooms[room][client] = struct{}{}
	return true
}

func (p *ClientPool) leave(client *Client, room string) bool {
	p.mux.Lock()
	defer p.mux.Unlock()

	if p.rooms[room] == nil {
		return false
	}

	// Removing client entry
	delete(p.rooms[room], client)

	// Removing client room
	if len(p.rooms[room]) == 0 {
		delete(p.rooms, room)
	}

	return true
}

func (p *ClientPool) iterateRoom(room string, fn func(client *Client)) {
	p.mux.RLock()
	defer p.mux.RUnlock()

	if p.rooms[room] == nil {
		return
	}

	for client := range p.rooms[room] {
		fn(client)
	}

	return
}

func (p *ClientPool) broadcast(clientID, room, topic string, data interface{}) {
	p.iterateRoom(room, func(client *Client) {
		if client.id != clientID {
			// Ignoring any errors
			_ = client.Write(topic, data)
		}
	})
}

func (p *ClientPool) setupAdapters(adapters []Adapter) {
	for _, adapter := range adapters {
		adapter.HandleBroadcast(p.broadcast)
	}
}

func (p *ClientPool) broadcastAdapters(clientID, room, topic string, data interface{}) {
	for _, adapter := range p.adapters {
		adapter.Broadcast(clientID, room, topic, data)
	}
}
