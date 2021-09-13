package transport

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/gorilla/mux"
)

type http1Server struct {
	tlsConfig         *tls.Config
	certFile, keyFile string

	server *http.Server
	router *mux.Router
}

// NewHTTP1Server creates new RESTServer based only HTTP/1.x protocol.
// Works with both secured and unsecured connections.
//
// Listen method accepts address to bind to
func NewHTTP1Server(tlsConfig *tls.Config) RESTServer {
	return &http1Server{
		tlsConfig: tlsConfig,
		router:    mux.NewRouter(),
	}
}

// NewHTTP1ServerHelper is tls helper for NewHTTP1Server function
func NewHTTP1ServerHelper(certFile, keyFile string) RESTServer {
	return &http1Server{
		certFile: certFile,
		keyFile:  keyFile,
		router:   mux.NewRouter(),
	}
}
func (l *http1Server) Listen(addr string) error {
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

func (l *http1Server) Shutdown(ctx context.Context) error {
	return l.server.Shutdown(ctx)
}

func (l *http1Server) Router() *mux.Router {
	return l.router
}
