package session

type Session interface {
	ID() string
	ServeBroadcast(topic string, data interface{})
}
