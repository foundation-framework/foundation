package sockets

import (
	"context"
	"reflect"

	"github.com/foundation-framework/foundation/errors"
)

type simpleMessageHandler struct {
	ctx   context.Context
	topic string
	model interface{}
	fn    MessageHandlerFunc
}

func NewSimpleMessageHandler(
	ctx context.Context,
	topic string,
	model interface{},
	fn MessageHandlerFunc,
) MessageHandler {
	return &simpleMessageHandler{
		ctx:   ctx,
		topic: topic,
		model: ensurePointer(model),
		fn:    fn,
	}
}

func (h *simpleMessageHandler) Topic() string {
	return h.topic
}

func (h *simpleMessageHandler) Model() interface{} {
	return copyInterfaceValue(h.model)
}

func (h *simpleMessageHandler) Serve(data interface{}) interface{} {
	return h.fn(h.ctx, data)
}

type simpleReplyHandler struct {
	ctx   context.Context
	model interface{}
	fn    ReplyHandlerFunc
}

func NewSimpleReplyHandler(
	ctx context.Context,
	model interface{},
	fn ReplyHandlerFunc,
) ReplyHandler {
	return &simpleReplyHandler{
		ctx:   ctx,
		model: ensurePointer(model),
		fn:    fn,
	}
}

func (h *simpleReplyHandler) Model() interface{} {
	return copyInterfaceValue(h.model)
}

func (h *simpleReplyHandler) Serve(data interface{}) {
	h.fn(h.ctx, data)
}

type middlewareMessageHandler struct {
	ctx         context.Context
	topic       string
	model       interface{}
	middlewares []MessageHandlerMiddleware
}

func NewMiddlewareMessageHandler(
	ctx context.Context,
	topic string,
	model interface{},
	middlewares ...MessageHandlerMiddleware,
) MessageHandler {
	return &middlewareMessageHandler{
		ctx:         ctx,
		topic:       topic,
		model:       ensurePointer(model),
		middlewares: middlewares,
	}
}

func (h *middlewareMessageHandler) Topic() string {
	return h.topic
}

func (h *middlewareMessageHandler) Model() interface{} {
	return copyInterfaceValue(h.model)
}

func (h *middlewareMessageHandler) Serve(data interface{}) interface{} {
	handler := MessageHandlerFunc(
		func(ctx context.Context, data interface{}) interface{} {
			return data
		},
	)

	for i := len(h.middlewares) - 1; i >= 0; i -= 1 {
		handler = h.middlewares[i](handler)
	}

	return handler(h.ctx, data)
}

func ensurePointer(i interface{}) interface{} {
	if reflect.TypeOf(i).Kind() != reflect.Ptr {
		errors.Panicf("net/sockets: handler model must be a pointer")
	}

	return i
}

func copyInterfaceValue(i interface{}) interface{} {
	return reflect.New(reflect.ValueOf(i).Elem().Type()).Interface()
}
