package sockets

import (
	"context"
)

type stopHandler struct {
	topic string
	model interface{}
	fns   []func(context.Context, interface{}) interface{}
}

// NewStopHandler will create a Handler
//
// Any chained function return non-nil result will stop subsequent calls
func NewStopHandler(
	topic string,
	model interface{},
	fns ...func(context.Context, interface{}) interface{},
) Handler {
	return &stopHandler{
		topic: topic,
		model: model,
		fns:   fns,
	}
}

func (s *stopHandler) Topic() string {
	return s.topic
}

func (s *stopHandler) Model() interface{} {
	return copyInterfaceValue(s.model)
}

func (s *stopHandler) Serve(ctx context.Context, data interface{}) interface{} {
	for _, fn := range s.fns {
		result := fn(ctx, data)

		if result != nil {
			return result
		}
	}

	return nil
}

type stopLastHandler struct {
	topic string
	model interface{}
	fns   []func(context.Context, interface{}) interface{}
}

// NewStopLastHandler will create a Handler
//
// Any chained function return non-nil result will stop subsequent
// calls and call the last function with that result
func NewStopLastHandler(
	topic string,
	model interface{},
	fns ...func(context.Context, interface{}) interface{},
) Handler {
	return &stopLastHandler{
		topic: topic,
		model: model,
		fns:   fns,
	}
}

func (s *stopLastHandler) Topic() string {
	return s.topic
}

func (s *stopLastHandler) Model() interface{} {
	return copyInterfaceValue(s.model)
}

func (s *stopLastHandler) Serve(ctx context.Context, data interface{}) interface{} {
	for i, fn := range s.fns {
		result := fn(ctx, data)

		if result == nil {
			continue
		}

		lastIndex := len(s.fns) - 1
		if i != lastIndex {
			return s.fns[lastIndex](ctx, result)
		}

		return result
	}

	return nil
}