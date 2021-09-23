package session

import (
	"github.com/Workiva/go-datastructures/set"
	"github.com/intale-llc/foundation/net/sockets"
	"github.com/intale-llc/foundation/utils/rand"
)

type ClientPool struct {
	rooms map[string]*set.Set
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
		id:   randutil.UUIDv4(),
		data: map[string]interface{}{},
	}
}

func (p *ClientPool) join(client *Client, room string) {
	if p.rooms[room] == nil {
		p.rooms[room] = set.New()
	}

	p.rooms[room].Add(client)
}

func (p *ClientPool) leave(client *Client, room string) {
	p.rooms[room].Remove(client)

	if p.rooms[room].Len() == 0 {
		p.rooms[room] = nil
	}
}

func (p *ClientPool) iterateRoom(room string, fun func(client *Client)) {
	clients := p.rooms[room].Flatten()

	for _, client := range clients {
		fun(client.(*Client))
	}
}
