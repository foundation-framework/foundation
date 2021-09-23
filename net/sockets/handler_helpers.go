package sockets

import (
	"context"
)

type stopHandler struct {
	topic string
	model interface{}
	fns   []func(context.Context, interface{}) (string, interface{})
}

// NewStopHandler will create a Handler
//
// Any chained function return non-nil result will stop subsequent calls
func NewStopHandler(
	topic string,
	model interface{},
	fns ...func(context.Context, interface{}) (string, interface{}),
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

func (s *stopHandler) Serve(ctx context.Context, data interface{}) (string, interface{}) {
	for _, fn := range s.fns {
		replyTopic, replyData := fn(ctx, data)

		if replyData != nil {
			return replyTopic, replyData
		}
	}

	return "", nil
}

type stopLastHandler struct {
	topic string
	model interface{}
	fns   []func(context.Context, interface{}) (string, interface{})
}

// NewStopLastHandler will create a Handler
//
// Any chained function return non-nil result will stop subsequent
// calls and call the last function with that result
func NewStopLastHandler(
	topic string,
	model interface{},
	fns ...func(context.Context, interface{}) (string, interface{}),
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

func (s *stopLastHandler) Serve(ctx context.Context, data interface{}) (string, interface{}) {
	for i, fn := range s.fns {
		replyTopic, replyData := fn(ctx, data)

		if replyData == nil {
			continue
		}

		lastIndex := len(s.fns) - 1
		if i != lastIndex {
			return s.fns[lastIndex](ctx, replyData)
		}

		return replyTopic, replyData
	}

	// Impossible
	return "", nil
}
