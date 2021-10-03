package sockets

import (
	"context"
)

type stopHandler struct {
	ctxFn func() context.Context
	topic string
	model interface{}
	fns   []HandlerFunc
}

// NewStopHandler will create a Handler
//
// Any chained function return non-nil result will
// stop subsequent calls
//
// The return value is the first non-nil value returned
// by chained functions
func NewStopHandler(
	ctxFn func() context.Context,
	topic string,
	model interface{},
	fns ...HandlerFunc,
) Handler {
	return &stopHandler{
		ctxFn: ctxFn,
		topic: topic,
		model: model,
		fns:   fns,
	}
}

func (s *stopHandler) Context() context.Context {
	return s.ctxFn()
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
	ctxFn func() context.Context
	topic string
	model interface{}
	fns   []HandlerFunc
}

// NewStopLastHandler will create a Handler
//
// Any chained function return non-nil result will stop subsequent
// calls and call the last function with that result
//
// The return value is the first non-nil value returned
// by chained functions
func NewStopLastHandler(
	ctxFn func() context.Context,
	topic string,
	model interface{},
	fns ...HandlerFunc,
) Handler {
	return &stopLastHandler{
		ctxFn: ctxFn,
		topic: topic,
		model: model,
		fns:   fns,
	}
}

func (s *stopLastHandler) Context() context.Context {
	return s.ctxFn()
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
			s.fns[lastIndex](ctx, replyData)
		}

		return replyTopic, replyData
	}

	// Impossible
	return "", nil
}