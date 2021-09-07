package sockets

import (
	"github.com/intale-llc/foundation/internal/utils"
)

// TODO: Documentation

type stopHandler struct {
	topic string
	model interface{}
	fun   []func(interface{}) interface{}
}

func NewStopHandler(
	topic string,
	model interface{},
	fun ...func(interface{}) interface{},
) Handler {
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

type chainHandler struct {
	topic string
	model interface{}
	fun   []func(interface{}) interface{}
}

func NewChainHandler(
	topic string,
	model interface{},
	fun ...func(interface{}) interface{},
) Handler {
	return &chainHandler{
		topic: topic,
		model: model,
		fun:   fun,
	}
}

func (s *chainHandler) Topic() string {
	return s.topic
}

func (s *chainHandler) Model() interface{} {
	return utils.CopyInterfaceValue(s.model)
}

func (s *chainHandler) Serve(data interface{}) interface{} {
	result := data
	for _, fun := range s.fun {
		result = fun(data)
	}

	return result
}
