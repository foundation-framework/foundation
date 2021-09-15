package session

import (
	"github.com/Workiva/go-datastructures/set"
	"github.com/intale-llc/foundation/rand"
	"github.com/intale-llc/foundation/transport/sockets"
)

type ClientPool struct {
	rooms          map[string]*set.Set
	accessLevelFun func(client *Client) AccessLevel
}

func NewClientPool(accessLevelFun func(client *Client) AccessLevel) *ClientPool {
	if accessLevelFun == nil {
		accessLevelFun = accessLevelFunDefault
	}

	return &ClientPool{
		rooms:          map[string]*set.Set{},
		accessLevelFun: accessLevelFun,
	}
}

func (p *ClientPool) NewClient(connection sockets.Conn) *Client {
	return &Client{
		Conn: connection,
		pool: p,
		id:   rand.UUIDv4(),
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
