package transport

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/gorilla/mux"
)

type http2Server struct {
	tlsConfig         *tls.Config
	certFile, keyFile string

	server *http.Server
	router *mux.Router
}

// NewHTTP2Server creates new RESTServer based on HTTP/2 and
// (for backward compatibility) HTTP/1.x protocol. HTTP/2 will
// only work with secured connections, while HTTP/1 will work
// with both secured and unsecured connections.
//
// Listen method accepts address to bind to
func NewHTTP2Server(tlsConfig *tls.Config) RESTServer {
	return &http2Server{
		tlsConfig: tlsConfig,
		router:    mux.NewRouter(),
	}
}

// NewHTTP2ServerHelper is tls helper for NewHTTP2Server function
func NewHTTP2ServerHelper(certFile, keyFile string) RESTServer {
	return &http2Server{
		certFile: certFile,
		keyFile:  keyFile,
		router:   mux.NewRouter(),
	}
}

func (l *http2Server) Listen(addr string) error {
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

func (l *http2Server) Shutdown(ctx context.Context) error {
	return l.server.Shutdown(ctx)
}

func (l *http2Server) Router() *mux.Router {
	return l.router
}
