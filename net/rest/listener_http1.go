package rest

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/gorilla/mux"
)

type http1Listener struct {
	tlsConfig         *tls.Config
	certFile, keyFile string

	server *http.Server
	router *mux.Router
}

// NewHTTP1Listener creates new Listener based only HTTP/1.x protocol.
// Works with both secured and unsecured connections.
//
// Listen method accepts address to bind to
func NewHTTP1Listener(tlsConfig *tls.Config) Listener {
	return &http1Listener{
		tlsConfig: tlsConfig,
		router:    mux.NewRouter(),
	}
}

// NewHTTP1ListenerHelper is tls helper for NewHTTP1Listener function
func NewHTTP1ListenerHelper(certFile, keyFile string) Listener {
	return &http1Listener{
		certFile: certFile,
		keyFile:  keyFile,
		router:   mux.NewRouter(),
	}
}
func (l *http1Listener) Listen(addr string) error {
	l.server = &http.Server{
		Addr:    addr,
		Handler: l.router,

		// Disabling http2 support
		TLSNextProto: map[string]func(*http.Server, *tls.Conn, http.Handler){},
	}

	if l.tlsConfig == nil && (l.certFile == "" || l.keyFile == "") {
		return l.server.ListenAndServe()
	} else {
		return l.server.ListenAndServeTLS(l.certFile, l.keyFile)
	}
}

func (l *http1Listener) Shutdown(ctx context.Context) error {
	return l.server.Shutdown(ctx)
}

func (l *http1Listener) Router() *mux.Router {
	return l.router
}
