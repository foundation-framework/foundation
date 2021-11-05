package session

// Adapter describes a mechanism to share rooms between servers
type Adapter interface {
	// Broadcast broadcasts message to other adapter members
	Broadcast(sessionID, room, topic string, data interface{})

	// HandleBroadcast handles incoming broadcast
	HandleBroadcast(func(sessionID, room, topic string, data interface{}))
}
