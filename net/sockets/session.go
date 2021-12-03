package sockets

import (
	"context"
	"sync"

	"github.com/intale-llc/foundation/rand"
	"github.com/intale-llc/foundation/session"
)

type Session struct {
	Conn

	hub   *session.Hub
	rooms map[string]struct{}

	id   string
	data map[string]interface{}
	rmux sync.RWMutex
}

func NewSession(hub *session.Hub, conn Conn) *Session {
	return &Session{
		Conn: conn,

		hub:   hub,
		rooms: map[string]struct{}{},

		id:   rand.UUID(),
		data: map[string]interface{}{},
	}
}

func (s *Session) ID() string {
	return s.id
}

func (s *Session) ServeBroadcast(topic string, data interface{}) {
	// Ignoring any errors
	_ = s.Write(topic, data)
}

func (s *Session) GetData(key string) interface{} {
	s.rmux.RLock()
	defer s.rmux.RUnlock()

	return s.data[key]
}

func (s *Session) SetData(key string, data interface{}) {
	s.rmux.Lock()
	defer s.rmux.Unlock()

	s.data[key] = data
}

func (s *Session) DropData(key string) {
	s.rmux.Lock()
	defer s.rmux.Unlock()

	delete(s.data, key)
}

func (s *Session) Join(room string) {
	if s.hub.Join(s, room) {
		s.rooms[room] = struct{}{}
	}
}

func (s *Session) Leave(room string) {
	if s.hub.Leave(s, room) {
		delete(s.rooms, room)
	}
}

func (s *Session) LeaveAll() {
	for room := range s.rooms {
		s.hub.Leave(s, room)
	}

	s.rooms = map[string]struct{}{}
}

func (s *Session) Broadcast(room string, topic string, data interface{}) {
	s.hub.Broadcast(s, room, topic, data)
}

func (s *Session) Context() context.Context {
	return PackSession(context.TODO(), s)
}
