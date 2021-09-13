package rest

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/gorilla/mux"
)

type http2Listener struct {
	tlsConfig         *tls.Config
	certFile, keyFile string

	server *http.Server
	router *mux.Router
}

// NewHTTP2Listener creates new Listener based on HTTP/2 and
// (for backward compatibility) HTTP/1.x protocol. HTTP/2 will
// only work with secured connections, while HTTP/1 will work
// with both secured and unsecured connections.
//
// Listen method accepts address to bind to
func NewHTTP2Listener(tlsConfig *tls.Config) Listener {
	return &http2Listener{
		tlsConfig: tlsConfig,
		router:    mux.NewRouter(),
	}
}

// NewHTTP2ListenerHelper is tls helper for NewHTTP2Listener function
func NewHTTP2ListenerHelper(certFile, keyFile string) Listener {
	return &http2Listener{
		certFile: certFile,
		keyFile:  keyFile,
		router:   mux.NewRouter(),
	}
}

func (l *http2Listener) Listen(addr string) error {
	l.server = &http.Server{
		Addr:      addr,
		Handler:   l.router,
		TLSConfig: l.tlsConfig,
	}

	if l.tlsConfig == nil && (l.certFile == "" || l.keyFile == "") {
		return l.server.ListenAndServe()
	} else {
		return l.server.ListenAndServeTLS(l.certFile, l.keyFile)
	}
}

func (l *http2Listener) Shutdown(ctx context.Context) error {
	return l.server.Shutdown(ctx)
}

func (l *http2Listener) Router() *mux.Router {
	return l.router
}
