package sockets

import (
	"github.com/undefined7887/foundation/internal/utils"
)

type defaultHandler struct {
	topic string
	model interface{}
	fun   []func(interface{}) interface{}
}

func (s *defaultHandler) Topic() string {
	return s.topic
}

func (s *defaultHandler) Model() interface{} {
	return utils.CopyInterfaceValue(s.model)
}

func (s *defaultHandler) Serve(data interface{}) interface{} {
	for _, fun := range s.fun {
		result := fun(data)

		if result != nil {
			return result
		}
	}

	return nil
}

func NewDefaultHandler(
	topic string,
	model interface{},
	fun ...func(interface{}) interface{},
) Handler {
	return &defaultHandler{
		topic: topic,
		model: model,
		fun:   fun,
	}
}
