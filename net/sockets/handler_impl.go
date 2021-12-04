package sockets

import (
	"context"
	"reflect"
)

type simpleMessageHandler struct {
	topic string
	model interface{}
	fn    MessageHandlerFunc
}

func NewSimpleMessageHandler(
	topic string,
	model interface{},
	fn MessageHandlerFunc,
) MessageHandler {
	return &simpleMessageHandler{
		topic: topic,
		model: model,
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
	return h.fn(data)
}

type simpleReplyHandler struct {
	model interface{}
	fn    ReplyHandlerFunc
}

func NewSimpleReplyHandler(
	model interface{},
	fn ReplyHandlerFunc,
) ReplyHandler {
	return &simpleReplyHandler{
		model: model,
		fn:    fn,
	}
}

func (h *simpleReplyHandler) Model() interface{} {
	return copyInterfaceValue(h.model)
}

func (h *simpleReplyHandler) Serve(data interface{}) {
	h.fn(data)
}

type middlewareMessageHandler struct {
	topic       string
	model       interface{}
	middlewares []func(next MessageHandlerFuncCtx) MessageHandlerFuncCtx
}

func NewMiddlewareMessageHandler(
	topic string,
	model interface{},
	middlewares ...func(next MessageHandlerFuncCtx) MessageHandlerFuncCtx,
) MessageHandler {
	return &middlewareMessageHandler{
		topic:       topic,
		model:       model,
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
	handler := MessageHandlerFuncCtx(
		func(ctx context.Context, data interface{}) interface{} {
			return data
		},
	)

	for i := len(h.middlewares) - 1; i >= 0; i -= 1 {
		handler = h.middlewares[i](handler)
	}

	return handler(context.Background(), data)
}

func copyInterfaceValue(i interface{}) interface{} {
	if reflect.TypeOf(i).Kind() == reflect.Ptr {
		return reflect.New(reflect.ValueOf(i).Elem().Type()).Interface()
	} else {
		return reflect.New(reflect.TypeOf(i)).Interface()
	}
}
