package sockets

import (
	"context"
)

type simpleHandler struct {
	topic string
	model interface{}
	fn    HandlerFunc
}

// NewSimpleHandler creates simple Handler
//
// Always returns context.Background() as context and
// has only one chained function
func NewSimpleHandler(
	topic string,
	model interface{},
	fn HandlerFunc,
) Handler {
	return &simpleHandler{
		topic: topic,
		model: model,
		fn:    fn,
	}
}

func (h *simpleHandler) Context() context.Context {
	return context.Background()
}

func (h *simpleHandler) Topic() string {
	return h.topic
}

func (h *simpleHandler) Model() interface{} {
	return copyInterfaceValue(h.model)
}

func (h *simpleHandler) Serve(ctx context.Context, data interface{}) interface{} {
	return h.fn(ctx, data)
}

type simpleReplyHandler struct {
	model interface{}
	fn    ReplyHandlerFunc
}

// NewSimpleReplyHandler creates simple Handler for handle message reply
//
// Always returns context.Background() as context and
// has only one chained function
func NewSimpleReplyHandler(
	model interface{},
	fn ReplyHandlerFunc,
) Handler {
	return &simpleReplyHandler{
		model: model,
		fn:    fn,
	}
}

func (h *simpleReplyHandler) Context() context.Context {
	return context.Background()
}

func (h *simpleReplyHandler) Topic() string {
	return ""
}

func (h *simpleReplyHandler) Model() interface{} {
	return copyInterfaceValue(h.model)
}

func (h *simpleReplyHandler) Serve(ctx context.Context, data interface{}) interface{} {
	h.fn(ctx, data)
	return nil
}

type defaultHandler struct {
	ctxFn func() context.Context
	topic string
	model interface{}
	fns   []HandlerFunc
}

// NewDefaultHandler will create a Handler
//
// Any chained function return non-nil result will stop subsequent
// calls and call the last function with that result
//
// The return value - is the first non-nil value returned
// by chained functions
func NewDefaultHandler(
	ctxFn func() context.Context,
	topic string,
	model interface{},
	fns ...HandlerFunc,
) Handler {
	return &defaultHandler{
		ctxFn: ctxFn,
		topic: topic,
		model: model,
		fns:   fns,
	}
}

func (h *defaultHandler) Context() context.Context {
	return h.ctxFn()
}

func (h *defaultHandler) Topic() string {
	return h.topic
}

func (h *defaultHandler) Model() interface{} {
	return copyInterfaceValue(h.model)
}

func (h *defaultHandler) Serve(ctx context.Context, data interface{}) interface{} {
	for i, fn := range h.fns {
		replyData := fn(ctx, data)

		if replyData == nil {
			continue
		}

		lastIndex := len(h.fns) - 1
		if i != lastIndex {
			h.fns[lastIndex](ctx, replyData)
		}

		return replyData
	}

	// Impossible
	return nil
}
