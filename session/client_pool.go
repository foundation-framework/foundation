package session

import (
	"sync"

	"github.com/Workiva/go-datastructures/set"
	"github.com/intale-llc/foundation/net/sockets"
	"github.com/intale-llc/foundation/rand"
)

type ClientPool struct {
	rooms map[string]*set.Set
	mux   sync.RWMutex
}

func NewClientPool() *ClientPool {
	return &ClientPool{
		rooms: map[string]*set.Set{},
	}
}

func (p *ClientPool) NewClient(connection sockets.Conn) *Client {
	return &Client{
		Conn: connection,
		pool: p,
		id:   rand.UUIDv4(),
		data: map[string]interface{}{},
	}
}

func (p *ClientPool) join(client *Client, room string) {
	p.createRoom(room)
	p.rooms[room].Add(client)
}

func (p *ClientPool) leave(client *Client, room string) {
	if !p.roomExists(room) {
		return
	}

	p.rooms[room].Remove(client)
	p.removeRoom(room)
}

func (p *ClientPool) iterateRoom(room string, fun func(client *Client)) {
	clients := p.getRoom(room)
	if clients == nil {
		return
	}

	for _, client := range clients {
		fun(client.(*Client))
	}
}

func (p *ClientPool) roomExists(room string) bool {
	p.mux.RLock()
	defer p.mux.RUnlock()

	return p.rooms[room] != nil
}

func (p *ClientPool) createRoom(room string) {
	p.mux.Lock()
	defer p.mux.Unlock()

	if p.rooms[room] == nil {
		p.rooms[room] = set.New()
	}
}

func (p *ClientPool) getRoom(room string) []interface{} {
	p.mux.RLock()
	defer p.mux.RUnlock()

	if p.rooms[room] == nil {
		return nil
	}

	return p.rooms[room].Flatten()
}

func (p *ClientPool) removeRoom(room string) {
	p.mux.Lock()
	defer p.mux.Unlock()

	if p.rooms[room].Len() == 0 {
		delete(p.rooms, room)
	}
}
