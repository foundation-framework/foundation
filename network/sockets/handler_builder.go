package sockets

import (
	"github.com/intale-llc/foundation/internal/utils"
)

type stopHandler struct {
	topic string
	model interface{}
	fun   []func(interface{}) interface{}
}

// NewStopHandler will create a MessageHandler
//
// Any chained function return non-nil result will stop subsequent calls
func NewStopHandler(
	topic string,
	model interface{},
	fun ...func(interface{}) interface{},
) MessageHandler {
	return &stopHandler{
		topic: topic,
		model: model,
		fun:   fun,
	}
}

func (s *stopHandler) Topic() string {
	return s.topic
}

func (s *stopHandler) Model() interface{} {
	return utils.CopyInterfaceValue(s.model)
}

func (s *stopHandler) Serve(data interface{}) interface{} {
	for _, fun := range s.fun {
		result := fun(data)

		if result != nil {
			return result
		}
	}

	return nil
}

type stopLastHandler struct {
	topic string
	model interface{}
	fun   []func(interface{}) interface{}
}

// NewStopLastHandler will create a MessageHandler
//
// Any chained function return non-nil result will stop subsequent
// calls and call the last function with that result
func NewStopLastHandler(
	topic string,
	model interface{},
	fun ...func(interface{}) interface{},
) MessageHandler {
	return &stopLastHandler{
		topic: topic,
		model: model,
		fun:   fun,
	}
}

func (s *stopLastHandler) Topic() string {
	return s.topic
}

func (s *stopLastHandler) Model() interface{} {
	return utils.CopyInterfaceValue(s.model)
}

func (s *stopLastHandler) Serve(data interface{}) interface{} {
	for i, fun := range s.fun {
		result := fun(data)

		if result == nil {
			continue
		}

		lastIndex := len(s.fun) - 1
		if i != lastIndex {
			return s.fun[lastIndex](result)
		}

		return result
	}

	return nil
}
