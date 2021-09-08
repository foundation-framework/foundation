package rest

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type httpListener struct {
	server   http.Server
	certFile string
	keyFile  string

	router *mux.Router

	onError []func(error)
	onClose []func(error)
}

// NewHttpListener created new Listener based on http protocol
//
// Listen method accepts address to bind to
func NewHttpListener(certFile, keyFile string) Listener {
	return &httpListener{
		certFile: certFile,
		keyFile:  keyFile,

		router: mux.NewRouter(),

		onError: []func(error){},
		onClose: []func(error){},
	}
}

func (l *httpListener) Listen(addr string) error {
	l.server.Addr = addr
	l.server.Handler = l.router

	if l.certFile == "" || l.keyFile == "" {
		return l.server.ListenAndServe()
	} else {
		return l.server.ListenAndServeTLS(l.certFile, l.keyFile)
	}
}

func (l *httpListener) Shutdown(ctx context.Context) error {
	return l.server.Shutdown(ctx)
}

func (l *httpListener) Router() *mux.Router {
	return l.router
}

func (l *httpListener) OnError(fun func(err error)) {
	l.onError = append(l.onError, fun)
}

func (l *httpListener) callErrorCb(err error) {
	for _, fun := range l.onError {
		fun(err)
	}
}

func (l *httpListener) OnClose(fun func(err error)) {
	l.onClose = append(l.onClose, fun)
}

func (l *httpListener) callCloseCb(err error) {
	for _, fun := range l.onClose {
		fun(err)
	}
}
