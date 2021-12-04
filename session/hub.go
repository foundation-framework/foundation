package session

import (
	"sync"
)

type Hub struct {
	adapters []Adapter

	rooms      map[string]map[Session]struct{}
	roomsMutex sync.RWMutex
}

func NewHub(adapters ...Adapter) *Hub {
	result := &Hub{
		adapters: adapters,
		rooms:    map[string]map[Session]struct{}{},
	}

	for _, adapter := range adapters {
		adapter.HandleBroadcast(result.broadcast)
	}

	return result
}

func (p *Hub) Broadcast(session Session, room, topic string, data interface{}) {
	// Firstly broadcast to adapters
	p.broadcastAdapters(session.ID(), room, topic, data)

	// Then broadcast to current hub
	p.broadcast(session.ID(), room, topic, data)
}

func (p *Hub) Join(session Session, room string) bool {
	p.roomsMutex.Lock()
	defer p.roomsMutex.Unlock()

	if p.rooms[room] == nil {
		p.rooms[room] = map[Session]struct{}{}
	}

	if _, ok := p.rooms[room][session]; ok {
		return false
	}

	p.rooms[room][session] = struct{}{}
	return true
}

func (p *Hub) Leave(session Session, room string) bool {
	p.roomsMutex.Lock()
	defer p.roomsMutex.Unlock()

	if p.rooms[room] == nil {
		return false
	}

	// Removing session entry
	delete(p.rooms[room], session)

	// Removing session room
	if len(p.rooms[room]) == 0 {
		delete(p.rooms, room)
	}

	return true
}

func (p *Hub) iterateRoom(room string, fn func(session Session)) {
	p.roomsMutex.RLock()
	defer p.roomsMutex.RUnlock()

	if p.rooms[room] == nil {
		return
	}

	for session := range p.rooms[room] {
		fn(session)
	}

	return
}

func (p *Hub) broadcast(sessionID, room, topic string, data interface{}) {
	p.iterateRoom(room, func(session Session) {
		if session.ID() != sessionID {
			session.ServeBroadcast(topic, data)
		}
	})
}

func (p *Hub) broadcastAdapters(sessionID, room, topic string, data interface{}) {
	for _, adapter := range p.adapters {
		adapter.Broadcast(sessionID, room, topic, data)
	}
}
