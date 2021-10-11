package session

// Adapter describes a mechanism to share rooms between servers
type Adapter interface {
	// Broadcast broadcasts message to other adapter members
	Broadcast(clientID, room, topic string, data interface{})

	// HandleBroadcast handles incoming broadcast
	HandleBroadcast(func(clientID, room, topic string, data interface{}))
}
