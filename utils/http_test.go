package utils

import (
	"io"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServeRequest(t *testing.T) {
	var serveErr error

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		err := ServeHttpRequest("localhost:8080", "/test", func(writer http.ResponseWriter, request *http.Request) {
			_, serveErr = writer.Write([]byte("test"))
		})

		if err != nil && serveErr != nil {
			serveErr = err
		}

		wg.Done()
	}()

	req, err := http.NewRequest("GET", "http://localhost:8080/test", nil)
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	result, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	require.Equal(t, "test", string(result))

	wg.Wait()
	require.NoError(t, serveErr)
}
