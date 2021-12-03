package sockets

import (
	"context"
)

type simpleHandler struct {
	topic     string
	model     interface{}
	handlerFn HandlerFunc
}

// NewSimpleHandler creates simple message Handler
//
// Always uses context.Background() as context and
// has only one chained function
func NewSimpleHandler(
	topic string,
	model interface{},
	handlerFn HandlerFunc,
) Handler {
	return &simpleHandler{
		topic:     topic,
		model:     model,
		handlerFn: handlerFn,
	}
}

func (h *simpleHandler) Context(caller func(ctx context.Context)) {
	caller(context.Background())
}

func (h *simpleHandler) Topic() string {
	return h.topic
}

func (h *simpleHandler) Model() interface{} {
	return copyInterfaceValue(h.model)
}

func (h *simpleHandler) Serve(ctx context.Context, data interface{}) interface{} {
	return h.handlerFn(ctx, data)
}

type simpleReplyHandler struct {
	model     interface{}
	handlerFn ReplyHandlerFunc
}

// NewSimpleReplyHandler creates simple message reply Handler
//
// Always uses context.Background() as context and
// has only one chained function
func NewSimpleReplyHandler(
	model interface{},
	handlerFn ReplyHandlerFunc,
) Handler {
	return &simpleReplyHandler{
		model:     model,
		handlerFn: handlerFn,
	}
}

func (h *simpleReplyHandler) Context(caller func(ctx context.Context)) {
	caller(context.Background())
}

func (h *simpleReplyHandler) Topic() string {
	return ""
}

func (h *simpleReplyHandler) Model() interface{} {
	return copyInterfaceValue(h.model)
}

func (h *simpleReplyHandler) Serve(ctx context.Context, data interface{}) interface{} {
	h.handlerFn(ctx, data)
	return nil
}

type defaultHandler struct {
	contextFn  func() context.Context
	topic      string
	model      interface{}
	handlerFns []HandlerFunc
}

// NewDefaultHandler creates default message Handler
//
// Any chained function return non-nil result will stop subsequent
// calls and call the last function with that result
//
// The return value - is the first non-nil value returned
// by chained functions
func NewDefaultHandler(
	contextFn func() context.Context,
	topic string,
	model interface{},
	handlerFns ...HandlerFunc,
) Handler {
	return &defaultHandler{
		contextFn:  contextFn,
		topic:      topic,
		model:      model,
		handlerFns: handlerFns,
	}
}

func (h *defaultHandler) Context(caller func(ctx context.Context)) {
	caller(h.contextFn())
}

func (h *defaultHandler) Topic() string {
	return h.topic
}

func (h *defaultHandler) Model() interface{} {
	return copyInterfaceValue(h.model)
}

func (h *defaultHandler) Serve(ctx context.Context, data interface{}) interface{} {
	for i, handlerFn := range h.handlerFns {
		replyData := handlerFn(ctx, data)

		if replyData == nil {
			continue
		}

		lastIndex := len(h.handlerFns) - 1
		if i != lastIndex {
			h.handlerFns[lastIndex](ctx, replyData)
		}

		return replyData
	}

	// Impossible
	return nil
}

type complexHandler struct {
	contextFn  func(caller func(ctx context.Context))
	topic      string
	model      interface{}
	handlerFns []HandlerFunc
}

// NewComplexHandler will create a complex message Handler
//
// Complex handler is like default message handler
// but with another context function
func NewComplexHandler(
	contextFn func(caller func(ctx context.Context)),
	topic string,
	model interface{},
	handlerFns ...HandlerFunc,
) Handler {
	return &complexHandler{
		contextFn:  contextFn,
		topic:      topic,
		model:      model,
		handlerFns: handlerFns,
	}
}

func (h *complexHandler) Context(caller func(ctx context.Context)) {
	h.contextFn(caller)
}

func (h *complexHandler) Topic() string {
	return h.topic
}

func (h *complexHandler) Model() interface{} {
	return copyInterfaceValue(h.model)
}

func (h *complexHandler) Serve(ctx context.Context, data interface{}) interface{} {
	for i, handlerFn := range h.handlerFns {
		replyData := handlerFn(ctx, data)

		if replyData == nil {
			continue
		}

		lastIndex := len(h.handlerFns) - 1
		if i != lastIndex {
			h.handlerFns[lastIndex](ctx, replyData)
		}

		return replyData
	}

	// Impossible
	return nil
}
