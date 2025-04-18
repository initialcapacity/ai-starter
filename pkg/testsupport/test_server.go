package testsupport

import (
	"fmt"
	"github.com/initialcapacity/ai-starter/pkg/websupport"
	"net/http"
	"testing"
	"time"
)

func StartTestServer(t *testing.T, handlers websupport.Handlers) string {
	server := websupport.NewServer(func(mux *http.ServeMux) {
		handlers(mux)
		mux.HandleFunc("GET /health-just-for-tests", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	port, _ := server.Start("localhost", 0)
	AssertHealthy(t, port, "/health-just-for-tests")

	t.Cleanup(func() { stopTestServer(t, server) })

	return fmt.Sprintf("http://localhost:%d", port)
}

func stopTestServer(t *testing.T, server *websupport.Server) {
	err := server.Stop()
	if err != nil {
		t.Errorf("unable to stop server: %s", err)
	}
}

func AssertHealthy(t *testing.T, port int, path string) {
	statusCode := make(chan int)

	go func() {
		for {
			resp, err := http.Get(fmt.Sprintf("http://localhost:%d%s", port, path))
			if err == nil {
				statusCode <- resp.StatusCode
				return
			}
		}
	}()

	select {
	case code := <-statusCode:
		if code != http.StatusOK {
			t.Errorf("server responded with a non 200 code: %d", code)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("server did not respond in 100 milliseconds")
	}
}

func Handle(mux *http.ServeMux, path string, responseBody string) {
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(responseBody))
	})
}
