package httputils

import (
	"context"
	"net/http"
	"sync"
)

func ServeRequest(addr, path string, fn http.HandlerFunc) error {
	wg := sync.WaitGroup{}
	wg.Add(1)

	mux := http.NewServeMux()
	httpServer := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	mux.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		fn(writer, request)
		wg.Done()
	})

	var serverErr error
	go func() {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			serverErr = err
			wg.Done()
		}
	}()

	wg.Wait()
	if err := httpServer.Shutdown(context.TODO()); err != nil {
		return err
	}

	return serverErr
}
