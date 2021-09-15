package http

import (
	"context"
	"crypto/tls"
	stdhttp "net/http"

	"github.com/gorilla/mux"
)

type http1Server struct {
	tlsConfig         *tls.Config
	certFile, keyFile string

	server *stdhttp.Server
	router *mux.Router
}

// NewHTTP1Listen creates new Listener based only HTTP/1.x protocol.
// Works with both secured and unsecured connections.
//
// Listen method accepts address to bind to
func NewHTTP1Listen(tlsConfig *tls.Config) Listener {
	return &http1Server{
		tlsConfig: tlsConfig,
		router:    mux.NewRouter(),
	}
}

// NewHTTP1ListenerHelper is tls helper for NewHTTP1Listen function
func NewHTTP1ListenerHelper(certFile, keyFile string) Listener {
	return &http1Server{
		certFile: certFile,
		keyFile:  keyFile,
		router:   mux.NewRouter(),
	}
}
func (l *http1Server) Listen(addr string) error {
	l.server = &stdhttp.Server{
		Addr:    addr,
		Handler: l.router,

		// Disabling http2 support
		TLSNextProto: map[string]func(*stdhttp.Server, *tls.Conn, stdhttp.Handler){},
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
